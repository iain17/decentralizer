package app

import (
	"github.com/iain17/discovery"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"strings"
	"time"
	"net"
	"errors"
	"gx/ipfs/QmX3U3YXCQ6UYBxq2LVWF8dARS1hPUTEYLrSx654Qyxyw6/go-multiaddr-net"
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

func ping(host string) bool {
	conn, err := net.DialTimeout("tcp", host, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	// log success
	return true
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
		if len(addr.Protocols()) == 0 {
			continue
		}
		netAddr, err := manet.ToNetAddr(addr)
		if err != nil {
			logger.Warning(err)
			continue
		}
		if !ping(netAddr.String()) {
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
