package cli

import (
	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/agent"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) agentMode() cli.Command {
	return cli.Command{
		Name:        "agent",
		Aliases:     []string{"worker", "builder"},
		Usage:       "start agent mode",
		Description: "start agent mode",
		Action: func(ctx *cli.Context) error {
			// Set up a connection to the server.
			conn, err := grpc.Dial("127.0.0.1:9999",
				grpc.WithInsecure())
			if err != nil {
				return errors.E(err, "did not connect")
			}
			defer conn.Close()
			agent := agent.New(1, conn)
			return errors.E(agent.Run(), "error running agent")
		},
	}
}
