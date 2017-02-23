package upnp

import (
	"github.com/prestonTao/upnp"
)

var mapping *upnp.Upnp

func init() {
	mapping = new(upnp.Upnp)
}

func Open(localPort, remotePort int, protocol string) error {
	return mapping.AddPortMapping(localPort, remotePort, protocol)
}