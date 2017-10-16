package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"strconv"
)

func (d *Decentralizer) CreateSession(sessionType int32, ) {
	ipfs.Publish(d.i, strconv.Itoa(int(sessionType)), []byte{'O', 'K'})
}

func (d *Decentralizer) GetSession(sessionType int32, ) {
	//ipfs.Receive(topic, func(peer peer.ID, message string) {
	//	logger.Infof("Received: %s: %s\n", peer.String(), message)
	//})
}