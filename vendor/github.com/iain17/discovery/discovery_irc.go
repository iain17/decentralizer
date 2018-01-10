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
)

type DiscoveryIRC struct {
	connection      *irc.Connection
	localNode *LocalNode
	context context.Context
	channel string
	logger *logger.Logger
}

func (d *DiscoveryIRC) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logger.New("DiscoveryIRC")
	d.localNode = ln
	d.context = ctx
	infoHash := d.localNode.discovery.network.InfoHash()
	d.channel = "#"+hex.EncodeToString(infoHash[:])

	name := randomdata.SillyName()
	d.connection = irc.IRC(name, name)
	//d.connection.Debug = true
	d.connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.connection.UseTLS = true

	d.connection.AddCallback("001", func(e *irc.Event) {
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
		d.logger.Debugf("Received message: %s", message)
		if strings.HasPrefix(message, IRC_PREFIX) {
			ipPort := strings.Split(message[len(IRC_PREFIX):], ":")
			if len(ipPort) != 2 {
				d.logger.Warningf("Received a weird IRC message: %s", message)
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
		} else {
			d.logger.Debug("Message hasn't got the IRC prefix.")
		}
	})
	err = d.connection.Connect(IRC_SERVER)
	go d.Run()
	return err
}

func (d *DiscoveryIRC) Stop() {
	if d.connection != nil && d.connection.Connected() {
		d.connection.Disconnect()
	}
}

func (d *DiscoveryIRC) Run() {
	defer func () {
		d.logger.Info("Stopping...")
		d.Stop()
	}()
	retries := 0
	for {
		select {
		case <-d.context.Done():
			return
		default:
			if !d.connection.Connected() {
				if retries > 10 {
					return
				}
				time.Sleep(1 * time.Second)
				d.logger.Warning("Reconnecting...")
				err := d.connection.Connect(IRC_SERVER)
				if err != nil {
					d.logger.Error(err)
				}
				retries++
				continue
			}
			d.Advertise()
			time.Sleep(30 * time.Second)
		}
	}
}

func (d *DiscoveryIRC) Advertise() {
	if d.localNode.ip != "" {
		d.connection.Privmsgf(d.channel, "%s%s:%d", IRC_PREFIX, d.localNode.ip, d.localNode.port)
		d.logger.Debug("Sent IRC message")
	}
}