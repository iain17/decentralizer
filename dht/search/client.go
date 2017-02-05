package search

import (
	"net"
)

type Client struct{
	IP   net.IP
	Port int
	Valid bool
}