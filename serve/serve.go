package serve

import (
	logger "github.com/Sirupsen/logrus"
	desc "github.com/iain17/decentralizer/decentralizer"
	"net"
)

var service desc.Decentralizer

func setup() {
	var err error
	service, err = desc.New()
	if err != nil {
		panic(err)
	}
}

func Serve(addr string, http bool) {
	if service == nil {
		setup()
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	if http {
		go serveHttp(lis)
	} else {
		go serveGrpc(lis)
	}
	logger.Infof("Protobuf server listening at %s", addr)
	//TODO: Apart from protobuf, grpc. Could we add a simple http api?
}