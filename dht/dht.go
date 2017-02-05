package dht

import (
	"github.com/anacrolix/dht"
	"github.com/go-openapi/errors"
	"github.com/iain17/dht-hello/dht/search"
)

var server *dht.Server

func init() {
	var err error
	server, err = dht.NewServer(&dht.ServerConfig{
		OnQuery: onQuery,
		OnAnnouncePeer: onAnnouncePeer,
	})
	if err != nil {
		panic(err)
	}
}

func NewSearch(identifier string, port int32, impliedPort bool) error {
	return search.New(server, identifier, port, impliedPort)
}

func GetClients(identifier string) (map[string]*search.Client, errors.Error) {
	search := search.Get(identifier)
	if search == nil || search.Announcement == nil {
		return nil, errors.New(404, "Could not find existing search for identifier: %s", identifier)
	}
	return search.Clients, nil
}

func StopSearch(identifier string) errors.Error {
	search := search.Get(identifier)
	if search == nil {
		return errors.New(404, "Could not find existing search for identifier: %s", identifier)
	}
	search.Stop()
	return nil
}