package search

import (
	"github.com/anacrolix/dht"
)

//This method gets called every 30 seconds
func (s *search) validate() {
	for s.Announcement.NumContacted() == 0 {

	}
	result := <- s.Announcement.Peers
	for _, peer := range result.Peers {
		s.addClient(&peer)
	}
}

func (s *search) addClient(peer *dht.Peer) {
	s.Clients[peer.String()] = &Client{
		IP: peer.IP,
		Port: peer.Port,
		Valid: false,
	}
}