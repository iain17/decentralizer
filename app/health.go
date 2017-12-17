package app

import (
	"errors"
	"time"
	timeout "github.com/iain17/timeout"
	"context"
)

func (d *Decentralizer) Health() (bool, error) {
	timeout.Do(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for len(d.i.Peerstore.Peers()) >= MIN_CONNECTED_PEERS {
					break
				}
			}
		}
	}, 5 * time.Second)
	peers := d.i.Peerstore.Peers()
	if len(peers) == 0 {
		return false, errors.New("could not find any peers yet... Check your internet connection")
	}
	return true, nil
}