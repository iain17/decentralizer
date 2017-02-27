package upnp

import (
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
	"github.com/scottjg/upnp"
	logger "github.com/Sirupsen/logrus"
	"time"
	"strings"
)

var natpmpClient *natpmp.Client
var upnpClient *upnp.Upnp
var initErr error

func init() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		initErr = err
		return
	}
	logger.Infof("Gateway found %s", gatewayIP.String())
	natpmpClient = natpmp.NewClientWithTimeout(gatewayIP, 3 * time.Second)

	//upnp
	upnpClient = new(upnp.Upnp)
}

func Open(localPort, remotePort int, protocol string) error {
	err := openUpnp(localPort, remotePort, protocol)
	if err != nil {
		err = openNatpmp(localPort, remotePort, protocol)
	}
	if err == nil {
		logger.Infof("Forwarded %d -> %d", localPort, remotePort)
	}
	return nil
}

func openUpnp(localPort, remotePort int, protocol string) error {
	return upnpClient.AddPortMapping(localPort, remotePort, 0, strings.ToUpper(protocol), "decentralizer")
}

func openNatpmp(localPort, remotePort int, protocol string) error {
	if natpmpClient == nil {
		return initErr
	}
	_, err := natpmpClient.AddPortMapping(protocol, localPort, remotePort, 0)
	return err
}