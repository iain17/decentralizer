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

- Write a C++ example.