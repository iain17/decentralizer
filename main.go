package main

import (
	"github.com/iain17/decentralizer/cmd"
	"github.com/getsentry/raven-go"
)

func init() {
	raven.SetDSN("https://0cf522b0b3d841d1b601296ed41e9b5c:4317615fa9ab47b28718b33fd843e497@sentry.io/306393")
	raven.SetRelease("v1")
}

func main() {
	raven.CapturePanicAndWait(func() {
		cmd.Execute()
	},nil)
}
