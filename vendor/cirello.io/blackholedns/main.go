// Copyright 2018 github.com/ucirello/blackholedns
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

// Command blackholedns implements a simple DNS daemon that diverts undesired
// requests to blackhole.
package main // import "cirello.io/blackholedns"

import (
	"bufio"
	"log"
	"os"
	"strings"

	dns "cirello.io/blackholedns/dns"
	"github.com/urfave/cli"
)

func init() {

}

func main() {
	app := cli.NewApp()
	app.Name = "blackholedns"
	app.Usage = "initiates blackholedns server"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bind",
			Value: "127.0.0.1:53",
		},
		cli.StringFlag{
			Name:  "forward",
			Value: "8.8.8.8:53",
		},
		cli.StringFlag{
			Name:  "blackhole-file,f",
			Value: "hosts",
		},
	}
	app.Action = func(c *cli.Context) error {
		fd, err := os.Open(c.String("blackhole-file"))
		if err != nil {
			return err
		}

		var blackholeList []string
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" || line[0] == '#' {
				continue
			}
			f := strings.Fields(line)
			blackholeList = append(blackholeList, f[1:]...)
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		server := &dns.Server{
			BindAddr:    c.String("bind"),
			ForwardAddr: c.String("forward"),
			Blackhole:   blackholeList,
		}
		return server.ListenAndServe()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
