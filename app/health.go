package app

import "errors"

func (d *Decentralizer) Health() (bool, error) {
	peers := d.d.WaitForPeers(MIN_DISCOVERED_PEERS, 30)
	if len(peers) == 0 {
		return false, errors.New("Could not find any peers in 30 seconds. Something is wrong.")
	}
	return true, nil
}