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
	peers []*Peer
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
		peers: []*Peer{},
		seen: map[string]time.Time{},
	}, nil
}

func (s *service) GetPeers() []*Peer {
	return s.peers
}

func (s *service) SetDetail(name string, value string) {
	s.self.Details[name] = value;
}

func (s *service) DiscoveredAddress(IP net.IP, Port int) {
	address := IP.String() + ":" + strconv.Itoa(Port)
	diff := time.Now().Sub(s.seen[address])
	if diff < time.Duration(30 * time.Second) {
		return
	}
	s.seen[address] = time.Now()
	go s.introduce(address)
}

func (s *service) introduce(address string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(1 * time.Second))
	if err != nil {
		logger.Warning(err)
		return
	}
	defer conn.Close()
	c := pb.NewDecentralizerClient(conn)
	res, err := c.RPCGetService(context.Background(), &pb.GetServiceRequest{
		Hash: s.hash,
	})
	if err != nil {
		logger.Warning(err)
		return
	}
	logger.Info(res)
}