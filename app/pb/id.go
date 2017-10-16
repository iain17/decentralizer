package pb

import (
	"hash/fnv"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

func GetPeer(id libp2pPeer.ID) DPeer {
	h := fnv.New64a()
	sId := id.Pretty()
	h.Write([]byte(sId))
	return DPeer{
		DId: h.Sum64(),
		PId: sId,
	}
}