# Decentralizer

This project aims to be the tool to decentralize any existing or in-development service. By running this project in the background you are able to find other peers. No servers needed, no NAT getting in the way. You'll be able to discover peers, share details and send messages up and down to them.
On top of that this project provides an easy to use API to find other peers around a topic and save files to the network.

## Features

- Matchmaking
- File storage

## Why

Often times I've wondered what it would take to create a simple piece of software that will take care of all the hard work to decentralize a service.

## How

- DHT
- IRC
- [IPFS](https://github.com/ipfs/go-ipfs)
- [UDP hole punching (Stun)](https://github.com/ccding/go-stun)
- UpNp

## Disclaimer

This project is a work in progress. Please do not use it unless you know what you're doing. Don't expect any support and do NOT use this in production!