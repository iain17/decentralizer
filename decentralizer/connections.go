package decentralizer

import (
	"github.com/ccding/go-stun/stun"
	"net"
	"errors"
	"github.com/iain17/dht-hello/decentralizer/upnp"
	"fmt"
)

//Returns a forwarded udp connection.
func getUdpConn() (*net.UDPConn, *stun.Host, error) {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, nil, err
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	err = upnp.Open(port, port, "UDP")
	if err != nil {
		return nil, nil, err
	}
	nat, host, err := stun.NewClientWithConnection(conn).Discover()
	if nat != stun.NATFull {
		return nil, nil, errors.New(nat.String())
	}
	conn.Close()
	conn, err = net.ListenUDP("udp", conn.LocalAddr().(*net.UDPAddr))

	return conn, host, nil
}

func getTcpConn() (net.Listener, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	//Forward the port
	port := lis.Addr().(*net.TCPAddr).Port
	err = upnp.Open(port, port, "TCP")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not forward TCP RPC server. %v", err))
	}
	return lis, nil
}

