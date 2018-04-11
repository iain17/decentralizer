package discovery

import (
	"context"
	"time"
	"github.com/thoj/go-ircevent"
	"github.com/Pallinder/go-randomdata"
	"crypto/tls"
	"encoding/hex"
	"strings"
	"net"
	"strconv"
	"github.com/iain17/logger"
	"io/ioutil"
	ttlru "github.com/iain17/kvcache/lttlru"
	"fmt"
	"hash"
	"hash/crc32"
)

type DiscoveryIRC struct {
	connection      *irc.Connection
	localNode *LocalNode
	context context.Context
	channel string
	logger *logger.Logger
	//A fallback way of sharing data.
	messages *ttlru.LruWithTTL
	message string
	crcTable hash.Hash32
	retries int
}

func (d *DiscoveryIRC) String() string {
	return "DiscoveryIRC"
}

func (d *DiscoveryIRC) init(ctx context.Context) (err error) {
	defer func() {
		if d.localNode.wg != nil {
			d.localNode.wg.Done()
		}
	}()
	d.logger = logger.New(d.String())
	d.context = ctx
	infoHash := d.localNode.discovery.network.InfoHash()
	d.channel = "#"+hex.EncodeToString(infoHash[:])
	d.crcTable = crc32.NewIEEE()

	d.messages, err = ttlru.NewTTL(10)
	if err != nil {
		return err
	}

	name := randomdata.SillyName()
	d.connection = irc.IRC(name, name)
	d.connection.Log.SetOutput(ioutil.Discard)
	//d.connection.Debug = true
	d.connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.connection.UseTLS = true

	d.connection.AddCallback("001", func(e *irc.Event) {
		logger.Debugf("Joined IRC: %s", d.channel)
		d.connection.Join(d.channel)
	})
	d.connection.AddCallback("366", func(e *irc.Event) {
		d.Advertise()
	})
	d.connection.AddCallback("PRIVMSG", func(event *irc.Event) {
		if d.localNode.netTableService.isEnoughPeers() {
			return
		}
		message := event.Message()
		parts := strings.Split(message, IRC_SEPERATOR)
		if len(parts) != 2 {
			//legacy
			if strings.HasPrefix(message,"JOIN US:") {
				d.onReceiveJoin(message[len("JOIN US:"):])
				return
			}

			d.logger.Debugf("Malformed message: '%s'", message)
			return
		}
		switch parts[0] {
		case IRC_JOIN_HANDLE:
			d.onReceiveJoin(parts[1])
			break
		case IRC_MESSAGE_HANDLE:
			d.onReceiveNetworkMessage(parts[1])
			break
		default:
			d.logger.Warningf("Unknown handle '%s'.", parts[0])
		}
	})
	err = d.connection.Connect(IRC_SERVER)
	return err
}

func (d *DiscoveryIRC) onReceiveJoin(data string) {
	logger.Debugf("Received IRC join message: %s", data)
	ipPort := strings.Split(data, ":")
	if len(ipPort) != 2 {
		d.logger.Warningf("Received a weird IRC message: %s", data)
		return
	}
	port, err := strconv.Atoi(ipPort[1])
	if err != nil {
		d.logger.Warning(err)
		return
	}
	d.localNode.netTableService.Discovered(&net.UDPAddr{
		IP:   net.ParseIP(ipPort[0]),
		Port: port,
	})
}

func (d *DiscoveryIRC) onReceiveNetworkMessage(data string) {
	logger.Debugf("Received IRC network message: %s", data)
	d.crcTable.Reset()
	d.crcTable.Write([]byte(data))
	d.messages.AddWithTTL(d.crcTable.Sum32(), data, 30 * time.Minute)
	logger.Debugf("Saved network message: %s", data)
}

func (d *DiscoveryIRC) Serve(ctx context.Context) {
	defer d.Stop()
	d.localNode.waitTilCoreReady()

	if err := d.init(ctx); err != nil {
		d.localNode.lastError = err
		panic(err)
	}
	d.localNode.waitTilReady()
	advertiseTicker := time.Tick(30 * time.Second)
	messageTicker := time.Tick(120 * time.Second)
	for {
		select {
		case <-d.context.Done():
			return
		case <-advertiseTicker:
			d.check()
			d.Advertise()
			break
		case <-messageTicker:
			d.check()
			d.Message()
			break
		}
	}
}

func (d *DiscoveryIRC) check() {
	if !d.connection.Connected() {
		if d.retries > 10 {
			panic("Too many retries to IRC")
		}
		time.Sleep(5 * time.Second)
		d.logger.Warning("Reconnecting...")
		err := d.connection.Connect(IRC_SERVER)
		if err != nil {
			d.logger.Error(err)
		}
		d.retries++
	}
}

func (d *DiscoveryIRC) Stop() {
	if d.connection != nil && d.connection.Connected() {
		d.connection.Disconnect()
	}
}

func (d *DiscoveryIRC) Advertise() {
	if !d.connection.Connected() {
		return
	}

	if d.localNode.netTableService.isEnoughPeers() {
		return
	}

	if d.localNode.ip == "" {
		d.logger.Warning("Not sending advertise message. No ip set.")
		return
	}
	d.Send(IRC_JOIN_HANDLE, fmt.Sprintf("%s:%d", d.localNode.ip, d.localNode.port))
}

func (d *DiscoveryIRC) Message() {
	if !d.connection.Connected() {
		return
	}

	if d.message == "" {
		d.logger.Debug("Not sending network message. No message set.")
		return
	}
	d.Send(IRC_MESSAGE_HANDLE, d.message)
}

func (d *DiscoveryIRC) Send(handle string, data string) {
	payload := handle + IRC_SEPERATOR + data
	d.connection.Privmsg(d.channel, payload)
	d.logger.Debugf("Sent IRC network: %s", payload)
}
