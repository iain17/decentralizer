package search

import (
	"github.com/anacrolix/dht"
	"github.com/jasonlvhit/gocron"
	logger "github.com/Sirupsen/logrus"
	"crypto/sha1"
)

type search struct {
	server *dht.Server
	identifier string
	hash string//SHA1 from identifier
	port int32
	impliedPort bool
	Announcement *dht.Announce
	cron *gocron.Scheduler
	Clients map[string]*Client
}

var searches map[string]*search

func init() {
	searches = map[string]*search{}
}

func New(server *dht.Server, identifier string, port int32, impliedPort bool) error {
	logger.Infof("Start searching for other DHT clients with identifier '%s' with port %d and impliedPort %t", identifier, port, impliedPort)
	existing := Get(identifier)
	if existing != nil {
		existing.Stop()
	}

	hash, err := getHash(identifier)
	if err != nil {
		return err
	}

	searches[identifier] = &search{
		server: server,
		identifier: identifier,
		hash: hash,
		port: port,
		impliedPort: impliedPort,
		cron: gocron.NewScheduler(),
		Clients: map[string]*Client{},
	}
	err = searches[identifier].announce()
	if err != nil {
		return err
	}
	searches[identifier].cron.Every(30).Seconds().Do(searches[identifier].validate)
	searches[identifier].cron.Start()
	return nil
}

func getHash(identifier string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(identifier))
	if err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

func Get(identifier string) *search {
	return searches[identifier]
}

func (s *search) Stop() {
	if s.cron != nil {
		s.cron.Clear()
	}
	if s.Announcement != nil {
		s.Announcement.Close()
	}
	logger.Infof("Search with %x stopped", s.hash)
	searches[s.identifier] = nil
}

func (s *search) announce() error {
	if s.Announcement != nil {
		s.Announcement.Close()
	}
	logger.Infof("Announcing %x", s.hash)
	announcement, err := s.server.Announce(s.hash, int(s.port), s.impliedPort)
	if err != nil {
		logger.Debug(err)
		return err
	}
	s.Announcement = announcement
	s.validate()
	return nil
}