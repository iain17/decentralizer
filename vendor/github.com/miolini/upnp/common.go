package upnp

import (
	// "log"
	"errors"
	"net"
	"strings"
)

//Get this performance networking ip address
func GetLocalIntenetIp() (string, error) {
	/*
	  Get all the local address
	  Analyzing energy networking ip address
	*/

	conn, err := net.Dial("udp4", "google.com:80")
	if err != nil {
		return "", errors.New("You can not connect to the network")
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

// This returns the list of local ip addresses which other hosts can connect
// to (NOTE: Loopback ip is ignored).
func GetLocalIPs() ([]*net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]*net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}

		ips = append(ips, &ipnet.IP)
	}

	return ips, nil
}
