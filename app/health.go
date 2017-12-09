package app

import (
	"errors"
	"time"
)

func (d *Decentralizer) Health() (bool, error) {
	peers := d.d.WaitForPeers(MIN_DISCOVERED_PEERS, 1 * time.Second)
	if len(peers) == 0 {
		return false, errors.New("could not find any peers yet... Check your internet connection")
	}
	return true, nil
}