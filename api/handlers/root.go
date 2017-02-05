package handlers

import (
	logger "github.com/Sirupsen/logrus"
	"os"
)

func init() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.DebugLevel)
}
