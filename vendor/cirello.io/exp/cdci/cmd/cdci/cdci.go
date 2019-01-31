package main // import "cirello.io/exp/cdci/cmd/cdci"

import (
	"log"

	"cirello.io/exp/cdci/pkg/cli"
)

func main() {
	log.SetPrefix("cdci: ")
	log.SetFlags(0)
	cli.Run()
}
