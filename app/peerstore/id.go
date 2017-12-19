package peerstore

import (
	libp2pPeer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"hash/fnv"
)

//libp2p peer id to uint64. Some apps expect some identification in the form of an integer. This will make it so we are compatible.
func PeerToDnId(id libp2pPeer.ID) (pId string, dID uint64) {
	h := fnv.New32a()
	pId = id.Pretty()
	h.Write([]byte(pId))
	dID = 0x110000100000000 + uint64(h.Sum32())
	return
}
