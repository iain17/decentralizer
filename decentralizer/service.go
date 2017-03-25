package decentralizer

import (
	"github.com/jasonlvhit/gocron"
	logger "github.com/Sirupsen/logrus"
	"github.com/anacrolix/dht"
	"github.com/pkg/errors"
	"net"
	"time"
	"strconv"
	"github.com/iain17/decentralizer/decentralizer/pb"
)

type service struct {
	name string
	hash string
	running bool
	Details interface{}
	Cron *gocron.Scheduler
	Announcement *dht.Announce
	self *Peer
	peers map[string]*Peer
	//Initial map. Used for map a address to a time we queried it.
	//Only used at the very first stage.
	seen map[string]time.Time
}

type services map[string]*service

func newService(name, hash string, peer *Peer) (*service, error) {
	if peer == nil {
		return nil, errors.New("Peer should not be nil")
	}
	return &service{
		name:   name,
		hash: hash,
		running: true,
		self:	peer,
		Cron:   gocron.NewScheduler(),
		peers: map[string]*Peer{},
		seen: map[string]time.Time{},
	}, nil
}

func (s *service) GetPeers() []*pb.Peer {
	var peers []*pb.Peer
	for key, peer := range s.peers {
		//TODO make this a parameter.
		if len(peers) >= 100 {
			break
		}
		//expired
		diff := time.Now().Sub(peer.seen)
		if diff > time.Duration(45 * time.Second) {
			delete(s.peers, key)
			continue
		}

		peers = append(peers, peer.Peer)
	}
	return peers
}

func (s *service) SetDetail(name string, value string) {
	s.self.Details[name] = value;
}

func (s *service) PeerDiscovered(pbPeer *pb.Peer) {
	peer := Peer{
		Peer: pbPeer,
		seen: time.Now(),
	}
	key := peer.GetKey()
	if s.peers[key] == nil {
		logger.Infof("Discovered %s", key)
		s.peers[key] = &peer
	}
	s.peers[key].seen = peer.seen

	for key, value := range s.peers[key].Details {
		logger.Infof("Detail: %s = %s", key, value)
	}
}

func (s *service) DiscoveredAddress(IP net.IP, Port uint32) {
	address := IP.String() + ":" + strconv.Itoa(int(Port))
	diff := time.Now().Sub(s.seen[address])
	if diff < time.Duration(15 * time.Second) {
		return
	}
	s.seen[address] = time.Now()
	go s.introduce(IP, address)
}

//TODO: queue? Only x amount of outgoing connections?
func (s *service) introduce(IP net.IP, address string) {
	res, err := s.getService(address)
	if err != nil {
		logger.Warn(err)
		return
	}

	res.Result.Ip = IP.String()
	s.PeerDiscovered(res.Result)
	//Discover peers of this peer.
	for _, peer := range res.Peers {
		ip := net.ParseIP(peer.Ip)
		if ip == nil {
			continue
		}
		s.DiscoveredAddress(ip, peer.Port)
	}
}

func (s *service) stop() {
	s.running = false
	logger.Infof("Stopping %s", s.name)
	if s.Announcement != nil {
		s.Announcement.Close()
	}
}