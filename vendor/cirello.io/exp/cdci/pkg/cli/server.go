package cli

import (
	"net"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"cirello.io/exp/cdci/pkg/server"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v1"
)

func (c *commands) serverMode() cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"dispatcher", "queue"},
		Usage:       "start server mode",
		Description: "start server mode",
		Action: func(ctx *cli.Context) error {
			tasks := make(chan *api.Recipe, 1)
			tasks <- &api.Recipe{
				Id:          1,
				Environment: []string{"WHO=world"},
				Commands:    "echo Hello, $WHO;",
			}
			// close(tasks)
			l, err := net.Listen("tcp", ":9999")
			if err != nil {
				return errors.E(err, "failed to listen")
			}
			s := grpc.NewServer()
			api.RegisterRunnerServer(s, server.New(tasks))
			return errors.E("failed to serve", s.Serve(l))
		},
	}
}
