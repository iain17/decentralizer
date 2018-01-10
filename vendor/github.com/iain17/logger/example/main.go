package main

import "github.com/iain17/logger"

func main() {
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.INFO, //logger.DEBUG,
		Colored:  true,
	})
	fileOut, err := logger.NewFileOut("/tmp/test.log", logger.INFO)
	if err != nil {
		panic(err)
	}
	logger.AddOutput(fileOut)

	logger.Info("info message test!")
	logger.Debug("debug message test")
	logger.Warning("warning message test")
	logger.Error("error message test")
}
