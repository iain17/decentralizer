package app

import (
	"errors"
	Peer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"fmt"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/logger"
)


func getDecentralizedIdKey(decentralizedId uint64) string {
	return fmt.Sprintf("DecentralizedId_%d", decentralizedId)
}

//Try and find peer in DHT
func (d *Decentralizer) resolveDecentralizedId(decentralizedId uint64) (Peer.ID, error) {
	peers := d.b.Find(getDecentralizedIdKey(decentralizedId), 1024)
	seen := make(map[string]bool)
	for peer := range peers {
		id := peer.Pretty()
		if seen[id] {
			continue
		}
		seen[id] = true
		_, possibleId := peerstore.PeerToDnId(peer)
		if possibleId == decentralizedId {
			logger.Infof("Resolved %d == %s", id)
			return peer, nil
		}
	}
	return "", errors.New("could not resolve id.")
}