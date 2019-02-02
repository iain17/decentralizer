# Decentralizer

[![pipeline status](https://gitlab.com/atlascorporation/adna/badges/master/pipeline.svg)](https://gitlab.com/atlascorporation/adna/commits/master)
[![coverage report](https://gitlab.com/atlascorporation/adna/badges/master/coverage.svg)](https://gitlab.com/atlascorporation/adna/commits/master)

This project aims to be the tool to decentralize any existing or in-development service. By running this project in the background you are able to find other peers. No servers needed, no NAT getting in the way. You'll be able to discover peers, share details and send messages up and down to them.
On top of that this project provides an easy to use API to find other peers around a topic and save files to the network.

## Features

- Matchmaking
- File storage
- Peer to Peer communication
- Private and public key system. A way to control and sign updates to the network to control it.

## Why

Often times I've wondered what it would take to create a simple piece of software that will take care of all the hard work to decentralize a service.
Use cases for this piece of software are endless. From games to commercial enterprise software where a business wants a low cost highly scalable backend that will still work even after support has ended for it!

## How

- DHT
- IRC
- mDNS (bonjour)
- [IPFS](https://gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs)
- [UDP hole punching (Stun)](https://github.com/ccding/go-stun)
- UpNp

## Disclaimer

This project is a work in progress. Please do not use it unless you know what you're doing. I discourage anyone from using this in production.

#Generate
Build GRPC https://github.com/grpc/grpc/blob/v1.11.1/BUILDING.md
go build GRPC: go build ./vendor/google.golang.org/grpc

#Versions
protobufs version 3.5.0 rev e09c5db296004fbe3f74490e84dcd62c3c5ddb1b
grpc version v1.11.3 https://github.com/grpc/grpc/tree/v1.11.1

#UPX
UPX is used to pack the binary so its less big. It is not there to protect it.