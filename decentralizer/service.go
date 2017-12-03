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
	lane "gopkg.in/oleiade/lane.v1"
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
	seen map[string]*time.Time
	queryQueue *lane.Deque
}

type services map[string]*service

func newService(name, hash string, peer *Peer) (*service, error) {
	if peer == nil {
		return nil, errors.New("Peer should not be nil")
	}
	instance := &service{
		name:   name,
		hash: hash,
		running: true,
		self:	peer,
		Cron:   gocron.NewScheduler(),
		peers: map[string]*Peer{},
		seen: map[string]*time.Time{},
		queryQueue: lane.NewDeque(),
	}
	instance.Cron.Every(5).Seconds().Do(instance.manage)
	instance.Cron.Every(1).Seconds().Do(instance.query)
	instance.Cron.Start()
	return instance, nil
}

func (s *service) GetPeers() []*pb.Peer {
	var peers []*pb.Peer
	for _, peer := range s.peers {
		//TODO make the limit a parameter.
		if len(peers) >= 100 {
			break
		}
		peers = append(peers, peer.Peer)
	}
	return peers
}

func (s *service) SetDetail(name string, value string) {
	s.self.Details[name] = value;
}

func (s *service) PeerDiscovered(pbPeer *pb.Peer) {
	ip := net.ParseIP(pbPeer.Ip)
	if ip == nil {
		return
	}
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

func (s *service) DiscoveredAddress(ip string, port int, priority bool) {
	address := ip + ":" + strconv.Itoa(int(port))
	if s.seen[address] != nil {
		return
	}
	last := time.Now()
	s.seen[address] = &last
	if priority {
		s.queryQueue.Prepend(address)
	} else {
		s.queryQueue.Append(address)
	}
}

//Called every 15 seconds to do maintenance work
func (s *service) manage() {
	for key, peer := range s.peers {
		diff := time.Now().Sub(peer.seen)
		//Delete if expired expired
		if diff >= 10*time.Second {
			logger.Infof("Peer %s has expired", key)
			delete(s.peers, key)
			continue
		}

		if diff >= 5*time.Second {
			//logger.Infof("Querying peer %s", key)
			s.queryQueue.Prepend(peer.GetAddress())
			continue
		}
	}
	for key, _ := range s.seen {
		diff := time.Now().Sub(*s.seen[key])
		if diff > time.Duration(30*time.Second) {
			delete(s.seen, key)
		}
	}
}

func (s *service) query() {
	i := 0
	for {
		address := s.queryQueue.Shift()
		if address == nil || i > 100 {
			break
		}

		go func(address string) {
			res, err := s.getServiceRequest(address)
			if err != nil {
				//logger.Warn(err)
				return
			}
			s.PeerDiscovered(res.Result)
			for _, peer := range res.Peers {
				s.DiscoveredAddress(peer.Ip, int(peer.Port), true)
			}
			logger.Infof("Querying address %s", address)
		}(address.(string))
		i++
	}
}

func (s *service) stop() {
	s.running = false
	s.Cron.Clear()
	logger.Infof("Stopping %s", s.name)
	if s.Announcement != nil {
		s.Announcement.Close()
	}
}