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
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/repo/config"
)

var DefaultBootstrapAddresses = []string{
	"/dnsaddr/bootstrap.libp2p.io/ipfs/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/ipfs/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/ipfs/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/ipfs/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",            // mars.i.ipfs.io
	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",           // pluto.i.ipfs.io
	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",           // saturn.i.ipfs.io
	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",             // venus.i.ipfs.io
	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",            // earth.i.ipfs.io
	"/ip6/2604:a880:1:20::203:d001/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",  // pluto.i.ipfs.io
	"/ip6/2400:6180:0:d0::151:6001/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",  // saturn.i.ipfs.io
	"/ip6/2604:a880:800:10::4a:5001/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64", // venus.i.ipfs.io
	"/ip6/2a03:b0c0:0:1010::23:1001/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd", // earth.i.ipfs.io
}

func init() {
	if USE_OWN_BOOTSTRAPPING {
		bs := core.DefaultBootstrapConfig
		bs.BootstrapPeers = func() []pstore.PeerInfo {
			return nil
		}
	} else {
		bs := core.DefaultBootstrapConfig
		bs.BootstrapPeers = func() []pstore.PeerInfo {
			a, _ := config.DefaultBootstrapPeers()
			return toPeerInfos(a)
		}
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
	//bs := core.DefaultBootstrapConfig
	//if USE_OWN_BOOTSTRAPPING {
	//	bs.BootstrapPeers = d.discover
	//}
	//bs.MinPeerThreshold = MIN_CONNECTED_PEERS
	//return d.i.Bootstrap(bs)
	return nil
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
