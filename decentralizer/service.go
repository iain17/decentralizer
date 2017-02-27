package decentralizer

import (
	"github.com/jasonlvhit/gocron"
	logger "github.com/Sirupsen/logrus"
	"github.com/anacrolix/dht"
	"github.com/pkg/errors"
	"net"
	"time"
	"strconv"
	"google.golang.org/grpc"
	"github.com/iain17/dht-hello/decentralizer/pb"
	"golang.org/x/net/context"
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
	i := 0
	for _, peer := range s.peers {
		if i >= 100 {
			break
		}
		peers = append(peers, peer.Peer)
		i++
	}
	return nil
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
		s.peers[key] = &peer
	}
	s.peers[key].seen = peer.seen
	logger.Infof("Discovered %s", key)
	for key, value := range s.peers[key].Details {
		logger.Infof("Detail: %s = %s", key, value)
	}
}

func (s *service) DiscoveredAddress(IP net.IP, Port uint32) {
	address := IP.String() + ":" + strconv.Itoa(int(Port))
	diff := time.Now().Sub(s.seen[address])
	if diff < time.Duration(30 * time.Second) {
		return
	}
	s.seen[address] = time.Now()
	go s.introduce(IP, address)
}

//TODO: queue? Only x amount of outgoing connections?
func (s *service) introduce(IP net.IP, address string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(1 * time.Second))
	if err != nil {
		//logger.Warning(err)
		return
	}
	defer conn.Close()
	c := pb.NewDecentralizerClient(conn)
	res, err := c.RPCGetService(context.Background(), &pb.GetServiceRequest{
		Hash: s.hash,
	})
	if err != nil {
		//logger.Debug(err)
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