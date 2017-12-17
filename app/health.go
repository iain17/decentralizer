package app

import (
	"errors"
	"time"
	timeout "github.com/iain17/timeout"
	"context"
	"fmt"
)

func (d *Decentralizer) Health() (bool, error) {
	//MIN_CONNECTED_PEERS
	numPeers := len(d.i.Peerstore.Peers())
	if numPeers < 3 {
		timeout.Do(func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					numPeers := len(d.i.Peerstore.Peers())
					if numPeers >= 3 {
						return
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
		}, 5*time.Second)
	}
	numPeers = len(d.i.Peerstore.Peers())
	if numPeers < 3 {
		return false, errors.New(fmt.Sprintf("only connected to %d peers of the minimum of %d", numPeers, 3))
	}
	return true, nil
}