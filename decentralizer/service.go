package decentralizer

import (
	"github.com/jasonlvhit/gocron"
	//logger "github.com/Sirupsen/logrus"
	"github.com/anacrolix/dht"
	"github.com/pkg/errors"
	"net"
)

type service struct {
	name string
	Details interface{}
	Cron *gocron.Scheduler
	Announcement *dht.Announce
	self *Peer
	peers []*Peer
}

type services map[string]*service

func newService(name string, peer *Peer) (*service, error) {
	if peer == nil {
		return nil, errors.New("Peer should not be nil")
	}
	return &service{
		name:   name,
		self:	peer,
		Cron:   gocron.NewScheduler(),
		peers: []*Peer{},
	}, nil
}

func (s *service) GetPeers() []*Peer {
	return s.peers
}

func (s *service) SetDetail(name string, value interface{}) {
	s.self.Details[name] = value;
}

func (s *service) DiscoveredAddress(IP net.IP, Port int) {

}