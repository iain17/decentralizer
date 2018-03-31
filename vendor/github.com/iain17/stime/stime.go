package stime

import (
	"github.com/beevik/ntp"
	"time"
)

var response = ntpResponse()

func ntpResponse() ntp.Response {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return ntp.Response{}
	}
	return *response
}

func Now() time.Time {
	return time.Now().Add(response.ClockOffset).UTC()
}