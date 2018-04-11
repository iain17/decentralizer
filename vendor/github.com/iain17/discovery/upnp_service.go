package discovery

import (
	"github.com/miolini/upnp"
	"time"
	"context"
	"github.com/iain17/logger"
)

type UPnPService struct {
	mapping *upnp.Upnp
	localNode *LocalNode
	logger  *logger.Logger
	context context.Context
}

func (d *UPnPService) String() string {
	return "UpNp"
}

func (s *UPnPService) init(ctx context.Context) error {
	defer func() {
		if s.localNode.wg != nil {
			s.localNode.wg.Done()
		}
	}()
	s.mapping = new(upnp.Upnp)
	s.logger = logger.New(s.String())
	s.context = ctx
	return nil
}

func (s *UPnPService) Stop() {
	s.mapping.DelPortMapping(s.localNode.port, "UDP")
}

func (s *UPnPService) Serve(ctx context.Context) {
	s.localNode.waitTilCoreReady()
	defer s.Stop()

	if err := s.init(ctx); err != nil {
		s.localNode.lastError = err
		panic(err)
	}
	s.localNode.WaitTilReady()
	ticker := time.Tick(1 * time.Minute)
	for {
		select {
		case <-s.context.Done():
			return
		case <-ticker:
			err := s.process(s.localNode.port)
			if err != nil {
				s.logger.Error("error on forwarding process, %v", err)
			}
		}
	}
}

func (s *UPnPService) process(port int) (err error) {
	s.logger.Debugf("trying to map port %d...", port)
	if err := s.mapping.AddPortMapping(port, port, "UDP"); err == nil {
		if s.mapping.GatewayOutsideIP != "" {
			s.localNode.ip = s.mapping.GatewayOutsideIP
			//Disabled: Seems empty?
			//s.localNode.outgoingPort = s.mapping.OutsideMappingPort
		}
		s.logger.Debug("port mapping passed")
	} else {
		s.logger.Warningf("port mapping fail, %v", err)
	}
	return nil
}
