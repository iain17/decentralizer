#Decentralizer

This project aims to be the tool to decentralize any existing or in-development service. By running this project in the background you are able to find other peers. No servers needed, no NAT getting in the way. You'll be able to discover peers, share details and send messages up and down to them.

This project is inspired by [danoctavian's bluntly project](https://github.com/danoctavian/bluntly).

## Why

Often times I've wondered what it would take to create a simple piece of software that will take care of all the hard work to decentralize a service.

## How

- DHT
- Bonjour
- UDP ([kcp](https://github.com/xtaci/kcp-go))
- [UDP hole punching (Stun)](https://github.com/ccding/go-stun)

## Implementation

Follow the instructions below to implement decentralizer in your app.

## Client generation
Between the decentralizer tool and your actual service we communicate using [GRPC](grpc.io). Generate the client files for the language of your choice. This can be C++, Java, node or even another Go project.

## Usage

1. Bundle this application with your application and run it in the background. See the C++ example for an example.
Specify a port like so ```decentralizer serve --listen :8080``` so the host app knows where to connect to.

2. Implement the client in whatever language you're writing the app in.

3. Start off with a search request. As long as the request (stream) is open, we expect the host application to be still running and interested.
This search request should be ran on another thread of your application.

4. Request peers. You can pass along limits, minimum amounts etc. The service will then respond with peers.

##Todo

This project is a work in progress. I discourage you to use it in production.

- Instead of TCP and upnp. Use kcp and stun.
- Write a C++ example.
- Do local discovery with Bonjour.
- Enable to option to send messages betwen nodes.

## Worthy mentions

I looked into [IPFS](https://github.com/ipfs/go-ipfs/issues/3686) to power this project. However it came with a lot of extra junk and didn't allow to find common peers around a certain hash. Developing my own protocol and building it from scratch required less time in the end than the IPFS prototype.
