// Copyright 2017 github.com/ucirello
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

// Command fswatch listens for file modifications and triggers a commmand.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

var (
	ignore = flag.String("ignore", "", "comma-separated list of directories to ignore")
)

func init() {
	flag.Parse()
}

func main() {
	log.SetPrefix("watch:")
	log.SetFlags(0)
	flag.Args()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-c
		log.Println("shutting down")
		cancel()
	}()

	if len(flag.Args()) == 0 {
		log.Fatalln("no file specified")
	}

	patterns, ignores, cmd := trimSpaces(strings.Split(flag.Arg(0), ",")), trimSpaces(strings.Split(*ignore, ",")), strings.Join(flag.Args()[1:], " ")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("cannot stat current working directory:", err)
	}
	updates, err := monitorWorkDir(ctx, wd, patterns, ignores)
	if err != nil {
		log.Fatalln("cannot listen for file changes:", err)
	}

	for {
		c, cancel := context.WithCancel(ctx)
		go func() {
			cmdExec := exec.CommandContext(c, "sh", "-c", cmd)
			cmdExec.Stdin = os.Stdin
			cmdExec.Stderr = os.Stderr
			cmdExec.Stdout = os.Stdout
			cmdExec.Run()
		}()
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-updates:
			cancel()
		}
	}
}

func trimSpaces(strs []string) []string {
	var ret []string
	for _, str := range strs {
		ret = append(ret, strings.TrimSpace(str))
	}
	return ret
}
