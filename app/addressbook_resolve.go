package app

import (
	"errors"
	Peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"fmt"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/logger"
)


func getDecentralizedIdKey(decentralizedId uint64) string {
	return fmt.Sprintf("DecentralizedId_%d", decentralizedId)
}

//Try and find peer in DHT
func (d *Decentralizer) resolveDecentralizedId(decentralizedId uint64) (Peer.ID, error) {
	values, err := d.b.GetValues(d.i.Context(), DHT_DECENTRALIZED_ID_KEY_TYPE, getDecentralizedIdKey(decentralizedId), 1024)
	if err != nil {
		return "", err
	}
	seen := make(map[string]bool)
	for _, value := range values {
		id := string(value.Val)
		if seen[id] {
			continue
		}
		peerId, err := libp2pPeer.IDB58Decode(id)
		if err != nil {
			continue
		}
		_, possibleId := peerstore.PeerToDnId(peerId)
		if possibleId == decentralizedId {
			logger.Infof("Resolved %d == %s", id)
			return peerId, nil
		}
	}
	return "", errors.New("could not resolve id")
}