package app

import (
	"github.com/iain17/logger"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/repo/config"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"io/ioutil"
	"strings"
	"time"
	"github.com/pkg/errors"
)

func init() {
	bs := core.DefaultBootstrapConfig
	bs.BootstrapPeers = func() []pstore.PeerInfo {
		logger.Info("Nulled bootstrap")
		return nil
	}
	core.DefaultBootstrapConfig = bs
}

func (d *Decentralizer) initBootstrap() error {
	bs := core.DefaultBootstrapConfig
	bs.BootstrapPeers = d.bootstrapPeers
	core.DefaultBootstrapConfig = bs
	d.i.Bootstrap(bs)
	d.cron.Every(10).Seconds().Do(func() {
		d.shareOurBootstrap()
		d.saveBootstrapState()
	})
	return nil
}

func (d *Decentralizer) shareOurBootstrap() {
	peers, err := d.getBootstrapAddrs()
	if err != nil {
		logger.Warning(err)
		return
	}
	d.d.LocalNode.SetInfo("bootstrap", serializeBootstrapAddrs(peers))
}

func (d *Decentralizer) saveBootstrapState() {
	peers, err := d.getBootstrapAddrs()
	if err != nil {
		logger.Warning(err)
		return
	}
	if len(peers) == 0 {
		return
	}
	data := serializeBootstrapAddrs(peers)
	file, err := d.fs.Create(BOOTSTRAP_FILE)
	if err != nil {
		logger.Warning(err)
		return
	}
	file.WriteString(data)
}

func serializeBootstrapAddrs(bootstrapAddrs []config.BootstrapPeer) string {
	if bootstrapAddrs == nil {
		return ""
	}
	addrs := ""
	for _, addr := range bootstrapAddrs {
		addrs += addr.String() + DELIMITER_ADDR
	}
	return addrs
}

func unSerializeBootstrapAddrs(addrText string) ([]config.BootstrapPeer, error) {
	addrs := strings.Split(addrText, DELIMITER_ADDR)
	if len(addrs) == 0 {
		return nil, errors.New("no addressed given")
	}
	return config.ParseBootstrapPeers(addrs[:len(addrs)-1])
}

func (d *Decentralizer) getBootstrapAddrs() ([]config.BootstrapPeer, error) {
	connections := d.i.PeerHost.Network().Conns()
	var result []string
	for _, conn := range connections {
		addr := conn.RemoteMultiaddr().String() + "/ipfs/" + conn.RemotePeer().Pretty()
		result = append(result, addr)
	}
	return config.ParseBootstrapPeers(result)
}

func (d *Decentralizer) readBootstrapState() ([]config.BootstrapPeer, error) {
	file, err := d.fs.Open(BOOTSTRAP_FILE)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return unSerializeBootstrapAddrs(string(data))
}

func (d *Decentralizer) bootstrapPeers() []pstore.PeerInfo {
	var result []config.BootstrapPeer
	//We have no connections at all yet?
	if len(d.i.PeerHost.Network().Peers()) == 0 {
		restoredAddrs, err := d.readBootstrapState()
		if err != nil {
			logger.Warning(err)
		} else {
			result = append(result, restoredAddrs...)
		}
	}
	peers := d.d.WaitForPeers(1, 0 * time.Second)
	for _, peer := range peers {
		peerBootstrap, err := unSerializeBootstrapAddrs(peer.GetInfo("bootstrap"))
		if err != nil {
			logger.Warning(err)
			continue
		}
		result = append(result, peerBootstrap...)
	}
	logger.Infof("Bootstrapping with %d addresses.", len(result))
	return toPeerInfos(result)
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