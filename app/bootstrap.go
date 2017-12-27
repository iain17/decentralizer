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
	"net"
	"errors"
	"gx/ipfs/QmSGL5Uoa6gKHgBBwQG8u1CWKUC8ZnwaZiLgFVTFBR2bxr/go-multiaddr-net"
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
	for _, peer := range d.d.WaitForPeers(MIN_CONNECTED_PEERS, 10*time.Second) {
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

func isPrivate(IP net.IP) error {
	var err error
	private := false
	if IP == nil {
		err = errors.New("invalid ip")
	} else {
		if !IP.IsGlobalUnicast() {
			return errors.New("multicast or loopback")
		}
		_, privateIPV61BitBlock, _ := net.ParseCIDR("2aa1::1/32")
		_, privateIPV62BitBlock, _ := net.ParseCIDR("fc00::/7")
		_, privateNineBitBlock, _ := net.ParseCIDR("9.0.0.0/8")
		_, privateLocalBitBlock, _ := net.ParseCIDR("127.0.0.0/8")
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = privateIPV61BitBlock.Contains(IP) || privateIPV62BitBlock.Contains(IP) || privateNineBitBlock.Contains(IP) || privateLocalBitBlock.Contains(IP) || private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	}
	if private {
		err = errors.New("private ip")
	}
	return err
}

func (d *Decentralizer) setInfo() {
	if d.d == nil {
		return
	}
	ln := d.d.LocalNode
	addrs := ""
	for _, addr := range d.i.PeerHost.Addrs() {
		addrText := addr.String()
		netAddr, err := manet.ToNetAddr(addr)
		if err != nil {
			continue
		}
		netAddrText := netAddr.String()
		var ip net.IP
		if value, ok := netAddr.(*net.TCPAddr); ok {
			ip = value.IP
		}
		if value, ok := netAddr.(*net.UDPAddr); ok {
			ip = value.IP
		}
		err = isPrivate(ip)
		if err != nil {
			logger.Debugf("Ignored: '%s': %s", netAddrText, err)
			continue
		}
		addrs += addrText + DELIMITER_ADDR
	}

	ln.SetInfo("peerId", d.i.Identity.Pretty())
	logger.Infof("Broadcasting: %s", addrs)
	ln.SetInfo("addr", addrs)
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
		if err != nil {
			logger.Warning(err)
			continue
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		return nil, errors.New("no addr set")
	}
	return &pstore.PeerInfo{
		ID:    peerId,
		Addrs: addrs,
	}, nil
}
