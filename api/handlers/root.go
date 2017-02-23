package handlers

import (
	logger "github.com/Sirupsen/logrus"
	"os"
	"github.com/iain17/dht-hello/decentralizer"
)

var dService decentralizer.Decentralizer

func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.DebugLevel)
	dService = decentralizer.New()
}
