# Discovery
This package is used to discover fellow peers around a network. You can create a new network with the network package. A network is a unique public key.

Use this package in combination with a p2p project like IPFS or libp2p. This is just used to dynamically find the initial peers.

## Features
- Discovery through public BitTorrent DHT.
- Discovery through public IRC channels.
- Discovery through other peers.
- Saves previously connected peers to quickly find them again.
- Attempts to forward UDP connection with UPNP.
- Attempts to forward UDP connection with stun udp hole punching.

## TODO
1. (optional) add callback for data received