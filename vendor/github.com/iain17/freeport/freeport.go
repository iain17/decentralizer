package freeport
// This asks the kernel for a free open port that is ready to use

import (
	"net"
)

/*
 Get an open TCP port.
 returns a int of the open port.
*/
func GetTCPPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

/*
 Get an open UDP port.
 returns a int of the open port.
*/
func GetUDPPort() int {
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.LocalAddr().(*net.UDPAddr).Port
}