package freeport
// This asks the kernel for a free open port that is ready to use

import (
	"net"
	"fmt"
)

func IsFreeTCP(port int) bool {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}

func IsFreeUDP(port int) bool {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}

func IsFree(netType string, port int) bool {
	if netType == "udp" {
		return IsFreeUDP(port)
	}
	if netType == "tcp" {
		return IsFreeTCP(port)
	}
	return false
}

/*
 Get an open TCP port.
 returns a int of the open port.
*/
func GetTCPPort() int {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
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
 Get an open TCP port that has ports free going upwards.
 returns a int of the open port. You can safely say that all the ports above it till numRange are also free
*/
func GetPortRange(netType string, numRange int) int {
	var port int
	free := false
	for !free {
		free = true
		if netType == "tcp" {
			port = GetTCPPort()
		} else {
			port = GetUDPPort()
		}
		for i := 0; i < numRange; i++ {
			if !IsFree(netType, port+i) {
				free = false
				break
			}
		}
	}
	return port
}

/*
 Get an open UDP port.
 returns a int of the open port.
*/
func GetUDPPort() int {
	addr, err := net.ResolveUDPAddr("udp", ":0")
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