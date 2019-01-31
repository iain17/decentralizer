package cli // import "cirello.io/exp/cdci/pkg/cli"

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/urfave/cli.v1"
)

type commands struct {
	db *sqlx.DB
}

func (c *commands) bootstrap(ctx *cli.Context) error {
	// db, err := sqlx.Open("sqlite3", "cdci.db")
	// if err != nil {
	// 	return nil, errors.E(err)
	// }
	// TODO: models.Task
	// TODO: task.Bootstrap()
	return nil
}

// Run executes the application in CLI mode
func Run() {
	app := cli.NewApp()
	app.Name = "cdci"
	app.Usage = "dirty continuous integration service"
	app.Version = "0.0.1"
	cmds := &commands{}
	app.Before = cmds.bootstrap
	app.Commands = []cli.Command{
		cmds.agentMode(),
		cmds.serverMode(),
	}
	sort.Slice(app.Commands, func(i, j int) bool {
		return strings.Compare(app.Commands[i].Name, app.Commands[j].Name) < 0
	})
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
