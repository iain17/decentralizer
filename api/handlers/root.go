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
	var err error
	dService, err = decentralizer.New()
	if err != nil {
		panic(err)
	}

	err = dService.AddService("iain", 0)
	if err != nil {
		panic(err)
	}
	service := dService.GetService("iain")
	service.SetDetail("cool", "1")
}
