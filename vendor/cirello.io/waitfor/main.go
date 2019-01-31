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


// Command waitfor waits for a network target to be available (times out in 1
// minute) and run the given command.
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var timeout = flag.Duration("timeout", 1*time.Minute, "time before it gives up")

func init() {
	flag.Parse()
}

func main() {
	log.SetPrefix("waitfor:")
	log.SetFlags(0)

	if len(flag.Args()) < 2 {
		log.Fatal("not enough parameters")
	}

	target := flag.Arg(0)
	if !dial(*timeout, target) {
		log.Fatal("timeout")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-c
		log.Println("shutting down")
		cancel()
	}()

	cmd := strings.Join(flag.Args()[1:], " ")
	cmdExec := exec.CommandContext(ctx, "sh", "-c", cmd)
	cmdExec.Stdin = os.Stdin
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdout = os.Stdout
	cmdExec.Run()
}

func dial(timeout time.Duration, target string) bool {
	to := time.After(timeout)
	for {
		select {
		case <-to:
			return false
		case <-time.After(250 * time.Millisecond):
			c, err := net.Dial("tcp", target)
			if err == nil {
				c.Close()
				return true
			}
		}
	}
}
