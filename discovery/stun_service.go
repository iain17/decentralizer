package discovery

import (
	"github.com/op/go-logging"
	"context"
	"github.com/ccding/go-stun/stun"
	"time"
	"fmt"
)

type StunService struct {
	client *stun.Client
	localNode *LocalNode
	logger  *logging.Logger
	context context.Context
}

func (s *StunService) Init(ctx context.Context, ln *LocalNode) error {
	s.logger = logging.MustGetLogger("Stun")
	s.localNode = ln
	s.context = ctx
	s.client = stun.NewClientWithConnection(s.localNode.listenerService.socket)
	go s.Run()
	return nil
}

func (s *StunService) Stop() {

}

func (s *StunService) Run() {
	defer s.Stop()

	for {
		select {
		case <-s.context.Done():
			return
		default:
			err := s.process()
			if err != nil {
				s.logger.Debugf("error on forwarding process, %v", err)
			}
			time.Sleep(time.Minute)
		}
	}
}

func (s *StunService) process() (err error) {
	nat, host, err := s.client.Discover()
	if err != nil {
		return err
	}

	if host != nil {
		s.logger.Infof("processed, family %d, host %q, port %d", host.Family(), host.IP(), host.Port())
		s.localNode.ip = host.IP()
	}

	switch nat {
	case stun.NATError:
		return fmt.Errorf("test failed")
	case stun.NATUnknown:
		return fmt.Errorf("unexpected response from the STUN server")
	case stun.NATBlocked:
		return fmt.Errorf("UDP is blocked")
	case stun.NATFull:
		return fmt.Errorf("full cone NAT")
	case stun.NATSymetric:
		return fmt.Errorf("symetric NAT")
	case stun.NATRestricted:
		return fmt.Errorf("restricted NAT")
	case stun.NATPortRestricted:
		return fmt.Errorf("port restricted NAT")
	case stun.NATNone:
		return fmt.Errorf("not behind a NAT")
	case stun.NATSymetricUDPFirewall:
		return fmt.Errorf("symetric UDP firewall")
	}
	s.logger.Info("NAT open!")
	return nil
}

