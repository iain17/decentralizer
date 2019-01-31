// Copyright 2018 github.com/ucirello
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

// Package worker implements the build worker.
package worker // import "cirello.io/cci/pkg/worker"

import (
	"context"
	"fmt"
	"log"
	"os"

	"cirello.io/cci/pkg/grpc/api"
	"cirello.io/cci/pkg/grpc/client"
	"cirello.io/errors"
	"google.golang.org/grpc"
)

// Start the builders.
func Start(ctx context.Context, grpcServerAddr, buildsDir string) error {
	cc, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.E(err, "cannot dial to GRPC server")
	}
	cl := api.NewRunnerClient(cc)
	resp, err := cl.Configuration(context.Background(), &api.ConfigurationRequest{})
	if err != nil {
		return errors.E(err, "cannot load configuration")
	}
	for repoFullName, recipe := range resp.Configuration {
		total := int(recipe.Concurrency)
		for i := 0; i < total; i++ {
			buildsDir := fmt.Sprintf(buildsDir, i)
			if err := os.MkdirAll(buildsDir,
				os.ModePerm&0700); err != nil {
				return errors.E(err, "cannot create .cci build directory")
			}
			go worker(ctx, cc, buildsDir, repoFullName, i)
		}
	}
	return nil
}

func worker(ctx context.Context, cc *grpc.ClientConn, buildsDir, repoFullName string, i int) {
	c := client.New(cc)
	log.Println("starting worker for", repoFullName, i)
	defer log.Println("done with ", repoFullName, i)
	err := c.Run(ctx, buildsDir, repoFullName)
	if err != nil {
		log.Println("cannot run worker", repoFullName, i, ":", err)
	}
}
