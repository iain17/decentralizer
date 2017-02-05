package dht

import (
	"github.com/anacrolix/dht"
	"github.com/jasonlvhit/gocron"
	"crypto/sha1"
	logger "github.com/Sirupsen/logrus"
	"github.com/go-openapi/errors"
)

type search struct {
	hash string//SHA1 from identifier
	port int32
	impliedPort bool
	Announcement *dht.Announce
	cron *gocron.Scheduler
}

var hash string
var server *dht.Server
var searches map[string]*search

func init() {
	var err error
	server, err = dht.NewServer(&dht.ServerConfig{})
	if err != nil {
		panic(err)
	}
	searches = map[string]*search{}
	//Search("test", 8080, true)
}

func get(identifier string) *search {
	return searches[identifier]
}

func (s *search) stop() {
	if s.cron != nil {
		s.cron.Clear()
	}
	if s.Announcement != nil {
		s.Announcement.Close()
	}
	logger.Infof("Search with %x stopped", s.hash)
}

func (s *search) announce() error {
	if s.Announcement != nil {
		s.Announcement.Close()
	}
	logger.Infof("Announcing %x", s.hash)
	announcement, err := server.Announce(s.hash, int(s.port), s.impliedPort)
	if err != nil {
		logger.Debug(err)
		return err
	}
	s.Announcement = announcement
	return nil
}

func Search(identifier string, port int32, impliedPort bool) error {
	logger.Infof("Start searching for other DHT clients with identifier '%s' with port %d and impliedPort %t", identifier, port, impliedPort)
	existing := get(identifier)
	if existing != nil {
		existing.stop()
	}

	hash, err := getHash(identifier)
	if err != nil {
		return err
	}

	searches[identifier] = &search{
		hash: hash,
		port: port,
		impliedPort: impliedPort,
		cron: gocron.NewScheduler(),
	}
	err = searches[identifier].announce()
	if err != nil {
		return err
	}
	//searches[identifier].cron.Every(30).Seconds().Do(searches[identifier].announce)
	//searches[identifier].cron.Start()
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

func GetPeers(identifier string) ([]dht.Peer, errors.Error) {
	search := get(identifier)
	if search == nil || search.Announcement == nil {
		return nil, errors.New(404, "Could not find existing search for identifier: %s", identifier)
	}
	result := <- search.Announcement.Peers
	return result.Peers, nil
}

func StopSearch(identifier string) errors.Error {
	search := get(identifier)
	if search == nil {
		return errors.New(404, "Could not find existing search for identifier: %s", identifier)
	}
	search.stop()
	return nil
}