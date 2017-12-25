package app

import (
	"github.com/iain17/discovery"
	"github.com/iain17/logger"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	ma "gx/ipfs/QmW8s4zTsUoX1Q6CeYxVKPyqSKbF7H1YDUyTostBtZ8DaG/go-multiaddr"
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	pstore "gx/ipfs/QmYijbtjCxFEjSXaudaQAUz3LN5VKLssm8WCUsRoqzXmQR/go-libp2p-peerstore"
	"strings"
	"time"
)

func init() {
	if USE_OWN_BOOTSTRAPPING {
		core.DefaultBootstrapConfig = core.BootstrapConfig{
			MinPeerThreshold:  4,
			Period:            30 * time.Second,
			ConnectionTimeout: (30 * time.Second) / 3, // Period / 3
			BootstrapPeers: func() []pstore.PeerInfo {
				return nil
			},
		}
	}
}

func (d *Decentralizer) bootstrap() error {
	d.setInfo()
	bs := core.DefaultBootstrapConfig
	if USE_OWN_BOOTSTRAPPING {
		bs.BootstrapPeers = d.discover
	} else {
		bs.BootstrapPeers = nil
	}
	bs.MinPeerThreshold = MIN_CONNECTED_PEERS
	return d.i.Bootstrap(bs)
}

func (d *Decentralizer) discover() []pstore.PeerInfo {
	if d.d == nil {
		return nil
	}
	logger.Info("Bootstrapping")
	d.setInfo()
	var peers []pstore.PeerInfo
	for _, peer := range d.d.WaitForPeers(MIN_CONNECTED_PEERS, 5*time.Minute) {
		peerInfo, err := getInfo(peer)
		if err != nil {
			logger.Warning(err)
			continue
		}
		peers = append(peers, *peerInfo)
	}
	logger.Infof("Bootstrapped %d peers", len(peers))
	return peers
}

func (d *Decentralizer) setInfo() {
	if d.d == nil {
		return
	}
	ln := d.d.LocalNode
	addrs := ""
	for _, addr := range d.i.PeerHost.Addrs() {
		addrs += addr.String() + DELIMITER_ADDR
	}

	ln.SetInfo("peerId", d.i.Identity.Pretty())
	ln.SetInfo("addr", addrs)
}

func getInfo(remoteNode *discovery.RemoteNode) (*pstore.PeerInfo, error) {
	sPeerId := remoteNode.GetInfo("peerId")
	peerId, err := peer.IDB58Decode(sPeerId)
	if err != nil {
		return nil, err
	}
	var addrs []ma.Multiaddr
	for _, strAddr := range strings.Split(remoteNode.GetInfo("addr"), DELIMITER_ADDR) {
		addr, err := ma.NewMultiaddr(strAddr)
		if err != nil {
			logger.Warning(err)
			continue
		}
		addrs = append(addrs, addr)
	}
	return &pstore.PeerInfo{
		ID:    peerId,
		Addrs: addrs,
	}, nil
}
