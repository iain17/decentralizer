package app

import (
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"time"
	"github.com/ipfs/go-ipfs/core"
	"github.com/op/go-logging"
	"github.com/iain17/decentralizer/discovery"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	"context"
)

var logger = logging.MustGetLogger("bootstrap")

func init() {
	core.DefaultBootstrapConfig = core.BootstrapConfig{
		MinPeerThreshold:  4,
		Period:            30 * time.Second,
		ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
		BootstrapPeers: func() []pstore.PeerInfo {
			return nil
		},
	}
}

func (d *Decentralizer) bootstrap() []pstore.PeerInfo {
	logger.Info("Bootstrapping")
	d.setInfo()
	var peers []pstore.PeerInfo
	for _, peer := range d.d.WaitForPeers(MIN_DISCOVERED_PEERS, 300) {
		peerInfo, err := getInfo(peer)
		if err != nil {
			logger.Warning(err)
			continue
		}
		peers = append(peers, *peerInfo)
		d.i.PeerHost.Connect(context.Background(), *peerInfo)
	}
	return peers
}

func (d *Decentralizer) setInfo() {
	ln := d.d.LocalNode
	addrs := ma.Join(d.i.PeerHost.Addrs()...).String()
	ln.SetInfo("peerId", d.i.Identity.Pretty())
	ln.SetInfo("addr", addrs)
}

func getInfo(remoteNode *discovery.RemoteNode) (*pstore.PeerInfo, error) {
	sPeerId := remoteNode.GetInfo("peerId")
	peerId, err := peer.IDB58Decode(sPeerId)
	if err != nil {
		return nil, err
	}
	addrs, err := ma.NewMultiaddr(remoteNode.GetInfo("addr"))
	if err != nil {
		return nil, err
	}
	return &pstore.PeerInfo{
		ID: peerId,
		Addrs: ma.Split(addrs),
	}, nil
}