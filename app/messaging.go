package app

import (
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

func (d *Decentralizer) initMessaging() {
	d.i.PeerHost.SetStreamHandler(DIRECT_MESSAGE, d.directMessageReceived)
}

func (d *Decentralizer) SendMessage(peer peer.ID, data []byte) (bool, error) {
	return false, nil
}

func (d *Decentralizer) directMessageReceived(stream inet.Stream) {

}
