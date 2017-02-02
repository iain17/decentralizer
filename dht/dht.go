package dht

import (
	"github.com/anacrolix/dht"
	"time"
	"github.com/jasonlvhit/gocron"
)

var hash string
var server *dht.ServerConfig

func init() {
	var err error
	server, err = dht.NewServer(&dht.ServerConfig{})
	if err != nil {
		panic(err)
	}
	gocron.Every(30).Second().Do(announce, 1, "hello")
}

func announce() {
	println("Announcing...")
	announcement, err := server.Announce("60406582292516720455", 1231, false)

	if err != nil {
		panic(err)
	}

	peers := <-announcement.Peers
	println("Ips...")
	for _, peer := range peers.Peers {
		println(peer.String())
	}
	time.Sleep(5 * time.Second)
}

func Register(service string, port int) {
	hash = hash
}