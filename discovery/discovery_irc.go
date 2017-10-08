package discovery

import (
	"github.com/op/go-logging"
	"context"
	"time"
	"github.com/thoj/go-ircevent"
	"github.com/Pallinder/go-randomdata"
	"crypto/tls"
	"encoding/hex"
	"strings"
	"net"
	"strconv"
)

type DiscoveryIRC struct {
	connection      *irc.Connection
	localNode *LocalNode
	context context.Context
	channel string
	logger *logging.Logger
}

func (d *DiscoveryIRC) Init(ctx context.Context, ln *LocalNode) (err error) {
	d.logger = logging.MustGetLogger("DiscoveryIRC")
	d.localNode = ln
	d.context = ctx
	infoHash := d.localNode.discovery.network.InfoHash()
	d.channel = "#"+hex.EncodeToString(infoHash[:])

	name := randomdata.SillyName()
	d.connection = irc.IRC(name, name)
	//d.connection.Debug = true
	d.connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.connection.UseTLS = true

	d.connection.AddCallback("001", func(e *irc.Event) { d.connection.Join(d.channel) })
	d.connection.AddCallback("366", func(e *irc.Event) {  })

	d.connection.AddCallback("PRIVMSG", func(event *irc.Event) {
		message := event.Message()
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
		}
	})

	go d.Run()
	return d.connection.Connect(IRC_SERVER)
}

func (d *DiscoveryIRC) Stop() {
	if d.connection != nil {
		d.connection.Disconnect()
	}
}

func (d *DiscoveryIRC) Run() {
	t := time.NewTimer(30 * time.Second)
	defer func () {
		d.Stop()
		t.Stop()
	}()

	for {
		select {
		case <-d.context.Done():
			return
		case <-t.C:
			if !d.connection.Connected() {
				err := d.connection.Connect(IRC_SERVER)
				d.logger.Error(err)
				continue
			}

			if d.localNode.ip == "" {
				d.logger.Debug("Not sending a message because we don't know our ip yet.")
			}
			d.connection.Privmsgf(d.channel, "%s%s:%d", IRC_PREFIX, d.localNode.ip, d.localNode.port)
		}
	}
}