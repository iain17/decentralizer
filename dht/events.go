package dht

import (
	"github.com/anacrolix/dht/krpc"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/dht"
	"net"
	logger "github.com/Sirupsen/logrus"
)

func onQuery(query *krpc.Msg, source net.Addr) (propagate bool) {
	logger.Debug("onQuery: %s", query.String())
	return true
}
// Called when a peer successfully announces to us.
func onAnnouncePeer(infoHash metainfo.Hash, peer dht.Peer) {
	logger.Debug("onAnnouncePeer: %s", infoHash.AsString())
}