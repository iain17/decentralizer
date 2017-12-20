package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
)

func (d *Decentralizer) initPublisherFiles() {
	ipfs.Subscribe(d.i, PUBLISHER_TOPIC_FILES, d.publisherDefinitionChange)
}

func (d *Decentralizer) publisherDefinitionChange(peer peer.ID, data []byte) {

}