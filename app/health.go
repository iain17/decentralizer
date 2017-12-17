package app

import (
	"errors"
	"time"
	timeout "github.com/iain17/timeout"
	"context"
	"fmt"
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
	numPeers := len(d.i.Peerstore.Peers())
	if numPeers < MIN_CONNECTED_PEERS {
		return false, errors.New(fmt.Sprintf("only connected to %d peers of the minimum of %d", numPeers, MIN_CONNECTED_PEERS))
	}
	return true, nil
}