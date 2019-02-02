package app

import (
	"github.com/iain17/logger"
	"github.com/pkg/errors"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/config"
	"io/ioutil"
	"strings"
	"github.com/iain17/decentralizer/vars"
)

func init() {
	bs := core.DefaultBootstrapConfig
	bs.BootstrapPeers = func() []pstore.PeerInfo {
		return nil
	}
	core.DefaultBootstrapConfig = bs
}

func (d *Decentralizer) initBootstrap() error {
	bs := core.DefaultBootstrapConfig
	bs.BootstrapPeers = d.bootstrapPeers
	core.DefaultBootstrapConfig = bs
	go d.i.Bootstrap(bs)
	d.cron.Every(10).Seconds().Do(func() {
		d.shareOurBootstrap()
	})
	//Instantly start discovering and telling about ourselves if we are not on a limited connection.
	if !d.limitedConnection {
		d.startDiscovering()
	}
	return nil
}

func (d *Decentralizer) shareOurBootstrap() {
	if d.d == nil {return}
	peers, err := d.getBootstrapAddrs()
	if err != nil {
		logger.Warning(err)
		return
	}
	bootstrapNodes := serializeBootstrapAddrs(peers)
	d.d.LocalNode.SetInfo("bootstrap", bootstrapNodes)
	d.d.SetNetworkMessage(bootstrapNodes)
}

func (d *Decentralizer) saveBootstrapState() {
	peers, err := d.getBootstrapAddrs()
	if err != nil {
		logger.Debugf("Could not save bootstrap state: %s", err)
		return
	}
	data := serializeBootstrapAddrs(peers)
	file, err := d.fs.Create(vars.BOOTSTRAP_FILE)
	if err != nil {
		logger.Debugf("Could not save bootstrap state: %s", err)
		return
	}
	file.WriteString(data)
	logger.Debug("Saved bootstrap state")
}

func serializeBootstrapAddrs(bootstrapAddrs []config.BootstrapPeer) string {
	if bootstrapAddrs == nil {
		return ""
	}
	addrs := ""
	for _, addr := range bootstrapAddrs {
		addrs += addr.String() + vars.DELIMITER_ADDR
	}
	return addrs
}

/*
unserializes the ;; way of sharing bootstrap peers. the add parameter is what extra string should be appended to each entry. This can be helpful if you want to just append the
ipfs id so a multiaddr can qualifies for a bootstrap addr
 */
func unSerializeBootstrapAddrs(addrText string, add string) ([]config.BootstrapPeer, error) {
	addrs := strings.Split(addrText, vars.DELIMITER_ADDR)
	if len(addrs) == 0 {
		return nil, errors.New("no addressed given")
	}
	for i, _ := range addrs {
		addrs[i] += add
	}
	return config.ParseBootstrapPeers(addrs[:len(addrs)-1])
}

func (d *Decentralizer) getBootstrapAddrs() ([]config.BootstrapPeer, error) {
	connections := d.i.PeerHost.Network().Conns()
	var result []string
	for _, conn := range connections {
		if len(result) > vars.MAX_BOOTSTRAP_SHARE {
			break
		}
		addr := conn.RemoteMultiaddr().String() + "/ipfs/" + conn.RemotePeer().Pretty()
		result = append(result, addr)
	}
	return config.ParseBootstrapPeers(result)
}

func (d *Decentralizer) readBootstrapState() ([]config.BootstrapPeer, error) {
	file, err := d.fs.Open(vars.BOOTSTRAP_FILE)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return unSerializeBootstrapAddrs(string(data), "")
}

func (d *Decentralizer) bootstrapPeers() []pstore.PeerInfo {
	var result []config.BootstrapPeer
	//We have no connections at all yet?
	if len(d.i.PeerHost.Network().Peers()) == 0 {
		restoredAddrs, err := d.readBootstrapState()
		if err != nil {
			logger.Warning(err)
		} else {
			logger.Debugf("Bootstrapping with %d previous addresses", len(restoredAddrs))
			result = append(result, restoredAddrs...)
		}
	}
	if len(result) == 0 {
		d.startDiscovering()
	}

	if d.d != nil {
		peers := d.d.WaitForPeers(1, 0)
		for _, peer := range peers {
			if len(result) > vars.MAX_BOOTSTRAP_SHARE {
				break
			}

			//Peer sharing their bootstrap peers
			peer2peerBootstrap, err := unSerializeBootstrapAddrs(peer.GetInfo("bootstrap"), "")
			if err != nil {
				logger.Warning(err)
				continue
			}
			logger.Debugf("Discovered using: %s", peer2peerBootstrap)
			result = append(result, peer2peerBootstrap...)
		}
		for _, message := range d.d.GetNetworkMessages() {
			peerBootstrap, err := unSerializeBootstrapAddrs(message, "")
			if err != nil {
				logger.Warning(err)
				continue
			}
			logger.Debugf("Discovered using: %s", message)
			result = append(result, peerBootstrap...)
		}

		//Throw the peers that are sharing their bootstrap list in as bootstrap peers.
		if len(result) == 0 {
			for _, peer := range peers {
				peerBootstrap, err := unSerializeBootstrapAddrs(peer.GetInfo("addr"), "/ipfs/"+peer.GetInfo("peerId"))
				if err != nil {
					logger.Warning(err)
				} else {
					result = append(result, peerBootstrap...)
					logger.Debugf("Discovered using: %s", peerBootstrap)
				}
			}
		}
	}
	logger.Infof("Bootstrapping with %d addresses.", len(result))
	d.displayConnected()
	return toPeerInfos(result)
}

func (d *Decentralizer) displayConnected() {
	logger.Info("Connected table list:")
	for i, peer := range d.i.PeerHost.Network().Peers() {
		conns := d.i.PeerHost.Network().ConnsToPeer(peer)
		logger.Infof("%d. %s via %s", i, peer.Pretty(), conns[0].RemoteMultiaddr().String())
	}
	d.saveBootstrapState()
	d.onConnected()
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

func (d *Decentralizer) onConnected() {
	if d.connected { return }
	d.connected = true
	go d.updatePublisherDefinition()
}