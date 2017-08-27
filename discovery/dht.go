package discovery

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/whyrusleeping/mdns"
	"github.com/nictuku/dht"
	"crypto/sha1"
	"fmt"
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
	node      *dht.DHT
	ih        dht.InfoHash
	host      host.Host

	lk       sync.Mutex
	notifees []Notifee
	interval time.Duration
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

func NewDhtService(ctx context.Context, peerhost host.Host, interval time.Duration) (Service, error) {
	//hash, err := getHash(ServiceTag)
	//if err != nil {
	//	return nil, err
	//}
	ih, err := dht.DecodeInfoHash("deca7a89a1dbdc4b213de1c0d5351e92582f31fb")
	config := dht.NewConfig()

	addrs, err := getDialableListenAddrs(peerhost)
	if err != nil {
		return nil, fmt.Errorf("Could not get a dialable listen address: %s", err)
	}
	config.Port = addrs[0].Port
	node, err := dht.New(config)
	if err != nil {
		return nil, fmt.Errorf("new dht init err: %s", err)
	}

	s := &dhtService{
		host:     peerhost,
		node: 	  node,
		ih: 	  ih,
		interval: interval,
	}

	go s.pollForEntries(ctx)
	go s.awaitPeers(ctx)

	return s, nil
}

func (m *dhtService) Close() error {
	m.node.Stop()
	return nil
}

func (m *dhtService) pollForEntries(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	for {
		select {
		case <-ticker.C:
			m.request()
			log.Debug("dht query complete")
		case <-ctx.Done():
			log.Debug("dht service halting")
			return
		}
	}
}

func (m *dhtService) awaitPeers(ctx context.Context) {
	log.Debug("awaitPeers")
	for {
		select {
		case r := <-m.node.PeersRequestResults:
			log.Debug("We've got results")
			for _, peers := range r {
				for _, x := range peers {
					host := dht.DecodePeerAddress(x)
					m.addPeer(host)
				}
			}
		case <-ctx.Done():
			log.Debug("dht service halting")
			return
		}
	}
}

func (m *dhtService) request() {
	log.Debugf("sending request %s...", m.ih)
	m.node.PeersRequest(string(m.ih), true)
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

func (m *dhtService) addPeer(peer string) {
	m.lk.Lock()
	exists := m.lastPeers[peer] == 1
	if !exists {
		m.lastPeers[peer] = 1
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
	//maddr, err := manet.FromNetAddr(&net.TCPAddr{
	//	IP:   e.AddrV4,
	//	Port: e.Port,
	//})
	//if err != nil {
	//	log.Warning("Error parsing multiaddr from mdns entry: ", err)
	//	return
	//}
	//
	//pi := pstore.PeerInfo{
	//	ID:    mpeer,
	//	Addrs: []ma.Multiaddr{maddr},
	//}
	//
	//m.lk.Lock()
	//for _, n := range m.notifees {
	//	go n.HandlePeerFound(pi)
	//}
	//m.lk.Unlock()
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