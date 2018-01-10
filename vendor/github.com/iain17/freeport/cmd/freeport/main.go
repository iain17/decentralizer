package main

import (
	"github.com/iain17/freeport"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	protocol = kingpin.Arg("protocol", "Protocol of the free port. TCP or UDP.").Required().String()
)

func main() {
	kingpin.Parse()
	switch(*protocol) {
		case "TCP":
			println(freeport.GetTCPPort())
			break
		case "UDP":
			println(freeport.GetUDPPort())
			break
		default:
			println("Invalid protocol specified. Either specify UDP or TCP")
	}
}
