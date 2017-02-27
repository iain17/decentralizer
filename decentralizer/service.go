package decentralizer

import (
	"github.com/jasonlvhit/gocron"
	logger "github.com/Sirupsen/logrus"
	"github.com/anacrolix/dht"
	"github.com/pkg/errors"
	"net"
	"time"
	"strconv"
	"github.com/iain17/dht-hello/decentralizer/pb"
	"github.com/gogo/protobuf/proto"
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

func (s *service) SetDetail(name string, value interface{}) {
	s.self.Details[name] = value;
}

func (s *service) DiscoveredAddress(IP net.IP, Port int, localConn *net.UDPConn) {
	address := IP.String() + ":" + strconv.Itoa(Port)
	diff := time.Now().Sub(s.seen[address])
	if diff < time.Duration(30 * time.Second) {
		return
	}
	s.seen[address] = time.Now()
	err := s.introduce(address, localConn)
	if err != nil {
		logger.Warning(err)
	}
}

func (s *service) introduce(address string, localConn *net.UDPConn) error {
	remoteAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	request := &pb.IntroductionRequest{
		Hash: s.hash,
	}
	payload, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	logger.Infof("Introducing to %s", address)

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(payload)
	//res, err := localConn.WriteTo(payload, remoteAddr)
	return err
}