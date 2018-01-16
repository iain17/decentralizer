package discovery

import (
	"context"
	"github.com/ccding/go-stun/stun"
	"fmt"
	"github.com/iain17/logger"
)

type StunService struct {
	initialized bool
	client *stun.Client
	localNode *LocalNode
	logger  *logger.Logger
	context context.Context
}

func (d *StunService) String() string {
	return "Stun"
}

func (s *StunService) init(ctx context.Context) error {
	if s.initialized {
		return nil
	}
	s.initialized = true
	s.logger = logger.New(s.String())
	s.context = ctx
	s.client = stun.NewClientWithConnection(s.localNode.listenerService.socket)
	return nil
}

func (s *StunService) Stop() {

}

func (s *StunService) Serve(ctx context.Context) (err error) {
	err = s.init(ctx)
	if err != nil {
		return err
	}
	s.logger.Info("Running")

	nat, host, err := s.client.Discover()
	if err != nil {
		return err
	}

	if host != nil {
		s.logger.Debugf("processed, family %d, host %q, port %d", host.Family(), host.IP(), host.Port())
		s.localNode.ip = host.IP()
	}

	switch nat {
	case stun.NATError:
		return fmt.Errorf("NAT error")
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

