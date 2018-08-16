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
	//logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"

	"github.com/spf13/cobra"
	"github.com/iain17/logger"
	"time"
	"os"
	"github.com/iain17/decentralizer/api"
	"context"
	"os/signal"
	"syscall"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/decentralizer/app"
	"gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
)
var verbose, daemon, isPrivateKey, isLimited, removeLock bool
var logPath, networkKey string
var port int

func init() {
	RootCmd.AddCommand(apiCmd)

	apiCmd.Flags().IntVarP(&port,"port", "p", 50010, "Port to run api on. +1 for http.")
	apiCmd.Flags().BoolVarP(&daemon,"daemon", "d", false, "Run daemon mode. Meaning it won't close")
	apiCmd.Flags().BoolVarP(&verbose,"verbose", "v", false, "Verbose will enable verbose logging")
	apiCmd.Flags().StringVarP(&logPath, "logPath", "l", "./adna.log", "Path of log file to output to")
	apiCmd.Flags().StringVarP(&networkKey, "network", "n", "", "Network key we should initialize with")
	apiCmd.Flags().BoolVar(&isPrivateKey, "isPrivate", false, "Is network key a private key or not (not used if network key not set)")
	apiCmd.Flags().BoolVar(&isLimited, "limited", false, "If we are on a limited (slower) connection (not used if network key not set)")
	apiCmd.Flags().BoolVar(&removeLock, "removeLock", false, "If set to true. It will remove to lock file")
}

const MAX_IDLE_TIME = 5 * time.Minute//Ignored in daemon mode
// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Runs the GRPC and HTTP api",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		logLvl := logger.INFO
		if verbose {
			logLvl = logger.DEBUG
			logging.InitForTesting(logging.DEBUG)
		} else {
			//Set ipfs logging
			logging.InitForTesting(logging.ERROR)
		}
		logger.AddOutput(logger.Stdout{
			MinLevel: logLvl,
			Colored:  true,
		})
		if logPath != "" {
			os.Remove(logPath)
			fileOut, err := logger.NewFileOut(logPath, logLvl)
			if err != nil {
				logger.Fatal(err)
			}
			logger.AddOutput(fileOut)
			//if verbose {
			//	ipfsLogOption := logging.Output(fileOut)
			//	logging.Configure(ipfsLogOption)
			//}
		}
		if removeLock {
			err := os.Remove(app.Base.Path+"/ipfs/repo.lock")
			if err != nil {
				logger.Warning(err)
			}
		}

		logger.Infof("Version: %s - %s", BRANCH, COMMIT_HASH)

		ctx, cancel := context.WithCancel(context.Background())
		s, err := api.New(ctx, port)
		if err != nil {
			logger.Fatal(err)
		}

		if networkKey != "" {
			err = s.SetNetwork(pb.VERSION.String(), networkKey, isPrivateKey, isLimited)
			if err != nil {
				logger.Fatal(err)
			}
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
		if s.App != nil {
			s.App.WaitTilEnoughPeers()
		}
		if s.App != nil && free.Add(MAX_IDLE_TIME).After(time.Now()) {
			logger.Warning("Idle. Closing process.")
			cancel()
		}
	}
}