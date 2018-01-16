package app

import (
	"errors"
	"fmt"
	"github.com/iain17/timeout"
	"time"
	"context"
)

func (d *Decentralizer) Health(WaitForMinConnections bool) (bool, error) {
	if WaitForMinConnections {
		numPeers := len(d.i.PeerHost.Network().Peers())
		if numPeers < MIN_CONNECTED_PEERS {
			timeout.Do(func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						numPeers := len(d.i.PeerHost.Network().Peers())
						if numPeers >= MIN_CONNECTED_PEERS {
							return
						}
						time.Sleep(100 * time.Millisecond)
					}
				}
			}, 5*time.Second)
		}
		numPeers = len(d.i.PeerHost.Network().Peers())
		if numPeers < MIN_CONNECTED_PEERS {
			addrs := ""
			for _, addr := range d.i.PeerHost.Network().ListenAddresses() {
				protocols := addr.Protocols()
				if protocols[0].Name != "ip4" && protocols[0].Name != "ip6" {
					continue
				}
				addrs += ", " + addr.String()
			}
			err := d.bootstrap()
			if err != nil {
				return false, err
			}
			percentage := 0.0
			if numPeers > 0 {
				total := float64(MIN_CONNECTED_PEERS)
				percentage = float64(numPeers) / total * 100
			}
			return false, errors.New(fmt.Sprintf("%.2f %% ready (Try portforwarding %s)", percentage, addrs))
		}
	}

	if d.publisherRecord == nil {
		return false, errors.New(fmt.Sprintf("Not ready yet. Waiting for publisher file..."))
	}
	if !d.publisherDefinition.Status {
		return false, errors.New("closed")
	}
	return true, nil
}