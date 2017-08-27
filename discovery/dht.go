package discovery

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/whyrusleeping/mdns"
	"github.com/anacrolix/dht"
	"crypto/sha1"
	"fmt"
	"time"
)

var log = logging.Logger("dht")

const ServiceTag = "decentralizer"

type Service interface {
	io.Closer
	RegisterNotifee(Notifee)
	UnregisterNotifee(Notifee)
}

type Notifee interface {
	HandlePeerFound(pstore.PeerInfo)
}

type dhtService struct {
	node      *dht.Server
	host      host.Host
	hash 	  string

	lk       sync.Mutex
	notifees []Notifee
	lastPeers map[string]byte
}

func getHash(value string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(value))
	if err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

func getDialableListenAddrs(ph host.Host) ([]*net.TCPAddr, error) {
	var out []*net.TCPAddr
	for _, addr := range ph.Addrs() {
		na, err := manet.ToNetAddr(addr)
		if err != nil {
			continue
		}
		tcp, ok := na.(*net.TCPAddr)
		if ok {
			out = append(out, tcp)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("failed to find good external addr from peerhost")
	}
	return out, nil
}

func NewDhtService(ctx context.Context, peerhost host.Host) (Service, error) {
	hash, err := getHash("abcyouandme")
	if err != nil {
		return nil, err
	}
	node, err := dht.NewServer(nil)
	if err != nil {
		return nil, fmt.Errorf("new dht init err: %s", err)
	}

	s := &dhtService{
		host:     peerhost,
		node: 	  node,
		hash:	  hash,
		lastPeers: make(map[string]byte),
	}

	go s.receive(ctx)

	return s, nil
}

func (m *dhtService) Close() error {
	m.node.Close()
	return nil
}

func (m *dhtService) receive(ctx context.Context) {
	log.Debug("awaitPeers")
	annoucement, err := m.request()
	if err != nil {
		log.Warning(err)
		return
	}
	for {
		select {
		case result := <- annoucement.Peers:
			for _, peer := range result.Peers {
				log.Info(peer)
				m.addPeer(peer)
			}
			time.Sleep(5 * time.Second)
			m.receive(ctx)
			return
		case <-ctx.Done():
			log.Debug("dht service halting")
			return
		}
	}
}

func (m *dhtService) request() (*dht.Announce, error) {
	addrs, err := getDialableListenAddrs(m.host)
	if err != nil {
		return nil, fmt.Errorf("Could not get a dialable listen address: %s", err)
	}
	log.Debugf("sending request %s...", m.hash)
	return m.node.Announce(m.hash, addrs[0].Port, false)
}

func (m *dhtService) RegisterNotifee(n Notifee) {
	m.lk.Lock()
	m.notifees = append(m.notifees, n)
	m.lk.Unlock()
}

func (m *dhtService) UnregisterNotifee(n Notifee) {
	m.lk.Lock()
	found := -1
	for i, notif := range m.notifees {
		if notif == n {
			found = i
			break
		}
	}
	if found != -1 {
		m.notifees = append(m.notifees[:found], m.notifees[found+1:]...)
	}
	m.lk.Unlock()
}

func (m *dhtService) addPeer(peer dht.Peer) {
	m.lk.Lock()
	exists := m.lastPeers[peer.String()] == 1
	if !exists {
		m.lastPeers[peer.String()] = 1
		if len(m.lastPeers) > 1000 {
			m.lastPeers = make(map[string]byte)
		}
	}
	m.lk.Unlock()
	if exists {
		return
	}

	//TODO: Self checker
	log.Debug("new peer %q received", peer)
	maddr, err := manet.FromNetAddr(&net.TCPAddr{
		IP:   peer.IP,
		Port: peer.Port,
	})
	if err != nil {
		log.Warning("Error parsing multiaddr from mdns entry: ", err)
		return
	}

	pi := pstore.PeerInfo{
		//ID:    mpeer,
		Addrs: []ma.Multiaddr{maddr},
	}

	m.lk.Lock()
	for _, n := range m.notifees {
		go n.HandlePeerFound(pi)
	}
	m.lk.Unlock()
}

func (m *dhtService) handleEntry(e *mdns.ServiceEntry) {
	log.Debugf("Handling MDNS entry: %s:%d %s", e.AddrV4, e.Port, e.Info)
	mpeer, err := peer.IDB58Decode(e.Info)
	if err != nil {
		log.Warning("Error parsing peer ID from mdns entry: ", err)
		return
	}

	if mpeer == m.host.ID() {
		log.Debug("got our own mdns entry, skipping")
		return
	}

	maddr, err := manet.FromNetAddr(&net.TCPAddr{
		IP:   e.AddrV4,
		Port: e.Port,
	})
	if err != nil {
		log.Warning("Error parsing multiaddr from mdns entry: ", err)
		return
	}

	pi := pstore.PeerInfo{
		ID:    mpeer,
		Addrs: []ma.Multiaddr{maddr},
	}

	m.lk.Lock()
	for _, n := range m.notifees {
		go n.HandlePeerFound(pi)
	}
	m.lk.Unlock()
}