package decentralizer

import (
	"github.com/anacrolix/dht"
	"time"
	logger "github.com/Sirupsen/logrus"
)

func (d *decentralizer) setupDht() error {
	var err error
	conn, host, err := getUdpConn()
	if err != nil {
		return err
	}
	d.ip = host.IP()
	d.dht, err = dht.NewServer(&dht.ServerConfig{
		Conn: conn,//Use the forwarded udp connection.
	})
	return err
}

func (s *decentralizer) discoveryUsingDht(hash string, service *service) {
	if service.Announcement != nil {
		service.Announcement.Close()
	}
	var err error
	service.Announcement, err = s.dht.Announce(hash, int(s.rpcPort), false)
	if err != nil {
		logger.Warn(err)
		return
	}
	for {
		peers, ok := <-service.Announcement.Peers
		if !ok || !service.running {
			break
		}
		for _, peer := range peers.Peers {
			service.DiscoveredAddress(peer.IP.String(), peer.Port, false)
		}

	}
	if service.running {
		time.Sleep(5 * time.Second)
		s.discoveryUsingDht(hash, service)
	}
}