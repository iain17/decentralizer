package discovery
//Discover peers around a certain name using mainline DHT bootstrap
import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/anacrolix/dht"
	"crypto/sha1"
	"fmt"
	"time"
	peer "github.com/libp2p/go-libp2p-peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"crypto/rand"
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

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func getDialableListenAddrs(ph host.Host) ([]*net.TCPAddr, error) {
	var out []*net.TCPAddr
	for _, addr := range ph.Addrs() {
		na, err := manet.ToNetAddr(addr)
		if err != nil {
			continue
		}
		tcp, ok := na.(*net.TCPAddr)
		if !IsPublicIP(tcp.IP) {
			continue
		}
		if ok {
			out = append(out, tcp)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("failed to find good external addr from peerhost")
	}
	return out, nil
}

func NewDhtService(ctx context.Context, peerhost host.Host, name string) (Service, error) {
	hash, err := getHash(name)
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
	//No external ips. Just discovering others.
	if err != nil {
		log.Warning(err)
		log.Debugf("Requesting peers...")
		return m.node.Announce(m.hash, 0, true)
	} else {
		log.Debugf("Requesting peers and announcing us...")
		return m.node.Announce(m.hash, addrs[0].Port, false)
	}
}

func (m *dhtService) RegisterNotifee(n Notifee) {
	m.lk.Lock()
	n.HandlePeerFound(pstore.PeerInfo{})
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

func (m *dhtService) addPeer(dhtPeer dht.Peer) {
	m.lk.Lock()
	exists := m.lastPeers[dhtPeer.String()] == 1
	if !exists {
		m.lastPeers[dhtPeer.String()] = 1
		if len(m.lastPeers) > 1000 {
			m.lastPeers = make(map[string]byte)
		}
	}
	m.lk.Unlock()
	if exists {
		return
	}

	//TODO: Self checker
	log.Debug("new peer %q received", dhtPeer)
	maddr, err := manet.FromNetAddr(&net.TCPAddr{
		IP:   dhtPeer.IP,
		Port: dhtPeer.Port,
	})
	if err != nil {
		log.Warning("Error parsing multiaddr from mdns entry: ", err)
		return
	}

	_, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	// A peers ID is the hash of its public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}

	pi := pstore.PeerInfo{
		ID:    pid,
		Addrs: []ma.Multiaddr{maddr},
	}

	m.lk.Lock()
	for _, n := range m.notifees {
		go n.HandlePeerFound(pi)
	}
	m.lk.Unlock()
}