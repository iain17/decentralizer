package utils

import (
	"strings"
	"net"
	"strconv"
	"github.com/iain17/logger"
)

// Convert uint to net.IP
func Inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0])
}

// Convert net.IP to int64
func Inet_aton(ipnr net.IP) uint64 {


	var sum uint64
	bits := strings.Split(ipnr.String(), ".")
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])

		sum += uint64(b0) << 24
		sum += uint64(b1) << 16
		sum += uint64(b2) << 8
		sum += uint64(b3)
		return sum
	}
	//if its an ipv6 address
	logger.Error("cant convert ip")

	return 2130706433 //127.0.0.1
}
