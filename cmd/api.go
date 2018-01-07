// Copyright © 2018 Iain Munro <iain@imunro.nl>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"

	"github.com/spf13/cobra"
	"github.com/iain17/logger"
	"time"
	"os"
	"github.com/iain17/decentralizer/api"
	"context"
	"os/signal"
	"syscall"
)
var verbose, daemon bool
var logPath string
var port int

func init() {
	RootCmd.AddCommand(apiCmd)

	apiCmd.Flags().IntVarP(&port,"port", "p", 50010, "Port to run api on. +1 for http.")
	apiCmd.Flags().BoolVarP(&daemon,"daemon", "d", false, "Run daemon mode. Meaning it won't close")
	apiCmd.Flags().BoolVarP(&verbose,"verbose", "v", false, "Verbose will enable verbose logging")
	apiCmd.Flags().StringVarP(&logPath, "logPath", "l", "adna.log", "Path of log file to output to")
}


const MAX_IDLE_TIME = 60 * time.Second//Ignored in daemon mode
// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Runs the GRPC and HTTP api",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		logLvl := logger.INFO
		if verbose {
			logLvl = logger.DEBUG
			logging.Configure(logging.LevelDebug)
		} else {
			//Set ipfs logging
			logging.Configure(logging.LevelError)
		}
		logger.AddOutput(logger.Stdout{
			MinLevel: logLvl, //logger.DEBUG,
			Colored:  true,
		})
		if logPath != "" {
			os.Remove(logPath)
			fileOut, err := logger.NewFileOut(logPath)
			if err != nil {
				logger.Fatal(err)
			}
			logger.AddOutput(fileOut)
		}

		ctx, cancel := context.WithCancel(context.Background())
		s, err := api.New(ctx, port)
		if err != nil {
			logger.Fatal(err)
		}

		if !daemon {
			go KillOnIdle(s, cancel)
		}

		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c,    syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT)
			select {
			case <-c:
				logger.Info("Stopping")
				cancel()
				s.Stop()
			}
		}()

		select {
		case <-ctx.Done():
			break
		}
	},
}

func KillOnIdle(s *api.Server, cancel context.CancelFunc) {
	logger.Warning("Killing on idle")
	var free time.Time
	for {
		time.Sleep(MAX_IDLE_TIME)
		s.Wg.Wait()
		free = time.Now()
		s.Wg.Wait()
		if free.Add(MAX_IDLE_TIME).After(time.Now()) {
			logger.Warning("Idle. Closing process.")
			cancel()
		}
	}
}