package peerstore

import (
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"hash/fnv"
)

//libp2p peer id to uint64. Some apps expect some user id kind of thing. This will make it so we are compatible.
func PeerToDnId(id libp2pPeer.ID) (pId string, dID uint64) {
	h := fnv.New64a()
	pId = id.Pretty()
	h.Write([]byte(pId))
	dID = h.Sum64()
	return
}
