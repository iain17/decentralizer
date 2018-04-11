package stime

import (
	"github.com/beevik/ntp"
	"time"
	"github.com/iain17/logger"
)

var response = ntpResponse()

var servers = []string{
	"time.google.com",
	"pool.ntp.org",
	"0.pool.ntp.org",
	"1.pool.ntp.org",
}

func ntpResponse() ntp.Response {
	for _, server := range servers {
		response, err := ntp.Query(server)
		if err != nil {
			logger.Warning(err)
			continue
		}
		return *response
	}
	return ntp.Response{}
}

func Now() time.Time {
	return time.Now().Add(response.ClockOffset).UTC()
}

//Bit off topic for this package. But since we're doing a external request here to a NTP server we can also report back if the network is slow.
func IsBadNetwork() bool {
	return response.RTT > 1 * time.Second
}