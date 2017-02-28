package decentralizer

import (
	"github.com/anacrolix/dht"
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