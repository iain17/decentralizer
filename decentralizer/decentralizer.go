package decentralizer

import (
	"github.com/pkg/errors"
	"github.com/anacrolix/dht"
	"net"
	"crypto/sha1"
	logger "github.com/Sirupsen/logrus"
)

type Decentralizer interface {
	AddService(name string, port int32) error
	GetService(name string) *service
}

type decentralizer struct {
	services services
	rpcPort uint16
	introPort uint16
	ip string
	dht *dht.Server

	introConn *net.UDPConn
}

func New() (Decentralizer, error) {
	instance := &decentralizer{
		services: services{},

	}
	//Setup RPC server
	err := instance.listenRpcServer()
	if err != nil {
		logger.Warn(err)
	}

	//Setup intro server
	err = instance.setupIntroServer()
	if err != nil {
		return nil, err
	}

	//Setup Dht server
	err = instance.setupDht()
	if err != nil {
		panic(err)
	}

	return instance, nil
}

func (d *decentralizer) AddService(name string, port int32) error {
	hash, err := getHash(name)
	if err != nil {
		return err
	}
	if d.services[hash] != nil {
		return errors.New("A service with that name already exists.")
	}

	self := NewPeer(d.ip, int32(d.rpcPort), port, map[string]interface{}{})
	d.services[hash], err = newService(name, self)
	if err != nil {
		return err
	}
	d.setupService(hash, d.services[hash])
	return err
}

//TODO: Improve this whole situation.
func (s *decentralizer) setupService(hash string, service *service) {

	if service.Announcement != nil {
		service.Announcement.Close()
	}
	logger.Infof("Announcing %x", hash)
	var err error
	service.Announcement, err = s.dht.Announce(hash, int(s.introPort), true)
	if err != nil {
		logger.Warn(err)
	}

	go func() {
		for {
			_, ok := <-service.Announcement.Peers
			if !ok {
				logger.Debug("done!")
				break
			}

		}
	}()
}

//TODO: Cache that hash.
func getHash(value string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(value))
	if err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

func (d *decentralizer) GetService(name string) *service {
	hash, err := getHash(name)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return d.services[hash]
}