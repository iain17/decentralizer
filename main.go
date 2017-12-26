package main

import (
	"context"
	"github.com/iain17/decentralizer/api"
	"github.com/iain17/logger"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"os"
	"os/signal"
)

func init() {
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.INFO, //logger.DEBUG,
		Colored:  true,
	})
	logging.Configure(logging.LevelError)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	s, err := api.New(ctx, 50010)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		s.Stop()
		cancel()
	}()

	select {
	case <-c:
		cancel()
	case <-ctx.Done():
	}
}
