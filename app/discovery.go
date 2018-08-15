package app

import (
	"github.com/iain17/discovery"
	"github.com/iain17/logger"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"net"
	"errors"
	"gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	manet "gx/ipfs/QmV6FjemM1K8oXjrvuq3wuVWWoU2TLDPmNnKrxHzY3v6Ai/go-multiaddr-net"
	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pnet "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	"strings"
	"time"
)

func (d *Decentralizer) initDiscovery() error {
	d.cron.Every(10).Seconds().Do(func() {
		d.setSelfAddrs()
		d.setReachableAddrs()
	})
	return nil
}

func (d *Decentralizer) startDiscovering() error {
	if d.d == nil {
		return nil
	}
	addrs, err := getAddrs(d.i.PeerHost)
	if err != nil {
		return err
	}
	d.d, err = discovery.New(d.ctx, d.n, MAX_DISCOVERED_PEERS, d.peerDiscovered, d.limitedConnection, map[string]string{
		"peerId": d.i.Identity.Pretty(),
		"addr": addrs,
	})
	return err
}

func (d *Decentralizer) setSelfAddrs() {
	if d.d == nil {
		return
	}
	addrs, err := getAddrs(d.i.PeerHost)
	if err != nil {
		logger.Warning(err)
		return
	}
	d.d.LocalNode.SetInfo("addr", addrs)
}

func (d *Decentralizer) setReachableAddrs() {
	if d.d == nil {
		return
	}
	for _, peer := range d.d.WaitForPeers(MIN_CONNECTED_PEERS, 10*time.Second) {
		peerInfo, err := remoteNodeToPeerInfo(peer)
		if err != nil {
			//logger.Warning(err)
			peer.Close()
			continue
		}

		if d.i.PeerHost.Network().Connectedness(peerInfo.ID) == libp2pnet.Connected {
			peer.SetInfo("addr", serializeAddrs(d.i.PeerHost.Peerstore().Addrs(peerInfo.ID)))
		}
	}
}

func (d *Decentralizer) peerDiscovered(peer *discovery.RemoteNode) {
	info, err := remoteNodeToPeerInfo(peer)
	if err != nil {
		logger.Warning(err)
		return
	}
	d.i.HandlePeerFound(*info)
}

func getDialableListenAddrs(ph host.Host) ([]ma.Multiaddr, error) {
	var out []ma.Multiaddr
	for _, addr := range ph.Addrs() {
		na, err := manet.ToNetAddr(addr)
		if err != nil {
			continue
		}
		if _, ok := na.(*net.TCPAddr); ok {
			out = append(out, addr)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("failed to find good external addr from peerhost")
	}
	return out, nil
}

func getAddrs(ph host.Host) (string, error) {
	maAddrs, err := getDialableListenAddrs(ph)
	if err != nil {
		return "", err
	}
	return serializeAddrs(maAddrs), nil
}

func serializeAddrs(multiAddrs []ma.Multiaddr) string {
	if multiAddrs == nil {
		return ""
	}
	addrs := ""
	for _, addr := range multiAddrs {
		addrs += addr.String() + DELIMITER_ADDR
	}
	return addrs
}

func unSerializeAddrs(addrText string) []ma.Multiaddr {
	var addrs []ma.Multiaddr
	rawAddr := strings.Split(addrText, DELIMITER_ADDR)
	for _, strAddr := range rawAddr {
		addr, err := ma.NewMultiaddr(strAddr)
		if err != nil && addr != nil {
			logger.Warning(err)
			continue
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

func remoteNodeToPeerInfo(remoteNode *discovery.RemoteNode) (*pstore.PeerInfo, error) {
	sPeerId := remoteNode.GetInfo("peerId")
	peerId, err := peer.IDB58Decode(sPeerId)
	if err != nil {
		return nil, err
	}
	addrText := remoteNode.GetInfo("addr")
	addrs := unSerializeAddrs(addrText)
	if len(addrs) == 0 {
		return nil, errors.New("no addr set")
	}
	return &pstore.PeerInfo{
		ID:    peerId,
		Addrs: addrs,
	}, nil
}