// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"context"
	"github.com/iain17/decentralizer/api"
	"github.com/iain17/logger"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"os"
	"github.com/kardianos/service"
)
const PROD = true
func init() {
	if !PROD {
		logger.AddOutput(logger.Stdout{
			MinLevel: logger.INFO, //logger.DEBUG,
			Colored:  true,
		})
	}
	logging.Configure(logging.LevelError)
}

type sLogger struct {}
var serviceLogger service.Logger
func (s sLogger) Print(level int, message string) error {
	switch level {
	case logger.ERROR:
		serviceLogger.Error(message)
		break
	case logger.WARNING:
		serviceLogger.Warning(message)
	default:
		serviceLogger.Info(message)
	}
	return nil
}

type program struct{
	ctx context.Context
	cancel context.CancelFunc
	api *api.Server
}

func (p *program) Start(s service.Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go p.run()
	return nil
}
func (p *program) run() {
	var err error
	p.api, err = api.New(p.ctx, 50010)
	if err != nil {
		serviceLogger.Error(err)
		os.Exit(0)
	}
}
func (p *program) Stop(s service.Service) error {
	p.cancel()
	s.Stop()
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "Decentralizer",
		DisplayName: "Adna",
		Description: "Takes care of all the hard parts ;)",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}
	serviceLogger, err = s.Logger(nil)
	if PROD {
		logger.AddOutput(sLogger{})
	}

	if err != nil {
		logger.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}