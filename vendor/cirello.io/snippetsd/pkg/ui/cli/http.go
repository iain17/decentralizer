// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"net"
	"net/http"

	"cirello.io/errors"
	"cirello.io/snippetsd/pkg/ui/web"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) httpMode() cli.Command {
	return cli.Command{
		Name:        "http",
		Aliases:     []string{"serve"},
		Usage:       "http mode",
		Description: "starts snippets web server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "host",
				Value:  "localhost",
				EnvVar: "HOST",
			},
			cli.StringFlag{
				Name:   "port",
				Value:  "8080",
				EnvVar: "PORT",
			},
		},
		Action: func(ctx *cli.Context) error {
			addr := net.JoinHostPort(ctx.String("host"), ctx.String("port"))
			l, err := net.Listen("tcp", addr)
			if err != nil {
				return errors.E(ctx, err, "cannot bind port")
			}
			err = http.Serve(l, web.New(c.db))
			return errors.E(ctx, err)
		},
	}
}
