package decentralizer

import (
	"github.com/anacrolix/dht"
	"github.com/anacrolix/dht/krpc"
	"github.com/anacrolix/torrent/metainfo"
	"net"
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
		OnQuery: d.onQuery,
		OnAnnouncePeer: d.onAnnouncePeer,
	})
	return err
}

func (d *decentralizer) onQuery(query *krpc.Msg, source net.Addr) (propagate bool) {
	//logger.Info("onQuery: %s", query.String())
	return true
}
// Called when a peer successfully announces to us.
func (d *decentralizer) onAnnouncePeer(infoHash metainfo.Hash, peer dht.Peer) {
	logger.Info("onAnnouncePeer", peer.IP, peer.Port)
}