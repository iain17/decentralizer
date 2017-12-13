package app

import (
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	inet "gx/ipfs/QmahYsGWry85Y7WUe2SX5G4JkH2zifEQAUtJVLZ24aC9DF/go-libp2p-net"
)

func (d *Decentralizer) initMessaging() {
	d.i.PeerHost.SetStreamHandler(DIRECT_MESSAGE, d.directMessageReceived)
}

func (d *Decentralizer) SendMessage(peer peer.ID, data []byte) (bool, error) {
	return false, nil
}

func (d *Decentralizer) directMessageReceived(stream inet.Stream)  {

}