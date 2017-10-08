package discovery

import (
	"fmt"
	"github.com/miolini/upnp"
	"time"
	"github.com/op/go-logging"
	"context"
)

type UPnPService struct {
	mapping *upnp.Upnp
	localNode *LocalNode
	logger  *logging.Logger
	context context.Context
}

func (s *UPnPService) Init(ctx context.Context, ln *LocalNode) error {
	s.mapping = new(upnp.Upnp)
	s.logger = logging.MustGetLogger("UpNp")
	s.localNode = ln
	s.context = ctx
	go s.Run()
	return nil
}

func (s *UPnPService) Stop() {
	s.mapping.DelPortMapping(s.localNode.port, "UDP")
}

func (s *UPnPService) Run() {
	defer s.Stop()
	for {
		select {
		case <-s.context.Done():
			return
		default:
			err := s.process(s.localNode.port)
			if err != nil {
				s.logger.Error("error on forwarding process, %v", err)
			}
			time.Sleep(time.Minute)
		}
	}
}

func (s *UPnPService) process(port int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	s.logger.Debugf("trying to map port %d...", port)
	if err := s.mapping.AddPortMapping(port, port, "UDP"); err == nil {
		if s.mapping.GatewayOutsideIP != "" {
			s.localNode.ip = s.mapping.GatewayOutsideIP
			//println(s.mapping.OutsideMappingPort)
		}
		s.logger.Debug("port mapping passed")
	} else {
		s.logger.Warningf("port mapping fail, %v", err)
	}
	return nil
}
