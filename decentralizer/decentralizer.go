package decentralizer

import (
	"github.com/pkg/errors"
	"github.com/anacrolix/dht"
	"crypto/sha1"
	logger "github.com/Sirupsen/logrus"
)

type Decentralizer interface {
	AddService(name string, port uint32) error
	GetService(name string) *service
	StopService(name string) error
}

type decentralizer struct {
	services services
	rpcPort uint16
	introPort uint16
	ip string
	dht *dht.Server
}

func New() (Decentralizer, error) {
	instance := &decentralizer{
		services: services{},

	}

	//Setup RPC server
	err := instance.listenRpcServer()
	if err != nil {
		logger.Error("Could not setup rpc server. This means you will not show up as a peer. You can only read!")
		logger.Warn(err)
	}

	//Setup Dht server
	err = instance.setupDht()
	if err != nil {
		return nil, err
	}

	logger.Info("Setup process finished.")

	return instance, nil
}

func (d *decentralizer) AddService(name string, port uint32) error {
	hash, err := getHash(name)
	if err != nil {
		return err
	}
	if d.services[hash] != nil {
		return errors.New("A service with that name already exists.")
	}

	self := NewPeer(d.ip, uint32(d.rpcPort), port, map[string]string{})
	d.services[hash], err = newService(name, hash, self)
	if err != nil {
		return err
	}
	go d.discoveryUsingDht(hash, d.services[hash])
	return err
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

func (d *decentralizer) StopService(name string) error {
	hash, err := getHash(name)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if d.services[hash] == nil {
		return errors.New("No service found")
	}
	d.services[hash].stop()
	delete(d.services, hash)
	return nil
}