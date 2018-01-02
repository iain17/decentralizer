package app

import (
	"github.com/iain17/discovery"
	"github.com/iain17/logger"
	"github.com/ipfs/go-ipfs/core"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"strings"
	"time"
	"errors"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/iain17/decentralizer/app/ipfs"
	"gx/ipfs/QmdQFrFnPrKRQtpeHKjZ3cVNwxmGKKS2TvhJTuN9C9yduh/go-libp2p-swarm"
	"gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

func init() {
	if USE_OWN_BOOTSTRAPPING {
		bs := core.DefaultBootstrapConfig
		bs.BootstrapPeers = func() []pstore.PeerInfo {
			return nil
		}
		core.DefaultBootstrapConfig = bs
	} else {
		bs := core.DefaultBootstrapConfig
		bs.BootstrapPeers = func() []pstore.PeerInfo {
			a, _ := config.DefaultBootstrapPeers()
			return toPeerInfos(a)
		}
		core.DefaultBootstrapConfig = bs
	}
}

func toPeerInfos(bpeers []config.BootstrapPeer) []pstore.PeerInfo {
	pinfos := make(map[peer.ID]*pstore.PeerInfo)
	for _, bootstrap := range bpeers {
		pinfo, ok := pinfos[bootstrap.ID()]
		if !ok {
			pinfo = new(pstore.PeerInfo)
			pinfos[bootstrap.ID()] = pinfo
			pinfo.ID = bootstrap.ID()
		}

		pinfo.Addrs = append(pinfo.Addrs, bootstrap.Transport())
	}

	var peers []pstore.PeerInfo
	for _, pinfo := range pinfos {
		peers = append(peers, *pinfo)
	}

	return peers
}

func (d *Decentralizer) bootstrap() error {
	if USE_OWN_BOOTSTRAPPING {
		bs := core.DefaultBootstrapConfig
		bs.BootstrapPeers = d.discover
		bs.MinPeerThreshold = MIN_CONNECTED_PEERS
		core.DefaultBootstrapConfig = bs
		return d.i.Bootstrap(bs)
	}
	return nil
}

func (d *Decentralizer) discover() []pstore.PeerInfo {
	if d.d == nil {
		return nil
	}
	logger.Info("Bootstrapping")
	d.setInfo()
	var peers []pstore.PeerInfo
	connected := 0
	for _, peer := range d.d.WaitForPeers(MIN_CONNECTED_PEERS, 10*time.Second) {
		peerInfo, err := getInfo(peer)
		if err != nil {
			logger.Warningf("Could not bootstrap %s: %s", peer.String(), err)
			continue
		}
		logger.Infof("Bootstrapping: %v", peerInfo)
		peers = append(peers, *peerInfo)

		if d.i.PeerHost.Network().Connectedness(peerInfo.ID) != net.Connected {
			d.clearBackOff(peerInfo.ID)
			err = d.i.PeerHost.Connect(d.i.Context(), *peerInfo)
			if err != nil {
				logger.Warning(err)
			} else {
				connected++
			}
		}
	}
	logger.Infof("Bootstrapped %d peers. We are connected to %d of those", len(peers), connected)
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
	logger.Infof("Broadcasting: %s", addrs)
	ln.SetInfo("addr", addrs)
}

func (d *Decentralizer) clearBackOff(id libp2pPeer.ID) {
	snet, ok := d.i.PeerHost.Network().(*swarm.Network)
	if ok {
		snet.Swarm().Backoff().Clear(id)
	}
}

func getInfo(remoteNode *discovery.RemoteNode) (*pstore.PeerInfo, error) {
	sPeerId := remoteNode.GetInfo("peerId")
	peerId, err := peer.IDB58Decode(sPeerId)
	if err != nil {
		return nil, err
	}
	var addrs []ma.Multiaddr
	addrText := remoteNode.GetInfo("addr")
	rawAddr := strings.Split(addrText, DELIMITER_ADDR)
	for _, strAddr := range rawAddr {
		addr, err := ma.NewMultiaddr(strAddr)
		if err != nil && addr != nil {
			logger.Warning(err)
			continue
		}
		if ipfs.IsAddrReachable(addr, false, true, false) {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return nil, errors.New("no addr set")
	}
	return &pstore.PeerInfo{
		ID:    peerId,
		Addrs: addrs,
	}, nil
}
