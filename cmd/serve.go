package cmd

import (
	logger "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/iain17/decentralizer/serve"
	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start serving the api",
	Long: `Start serving the api`,
	Run: run,
}

var (
	// Addr is the server listen address.
	addr string
	typeService string
)

func run(cmd *cobra.Command, args []string) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.DebugLevel)

	//Components
	go serve.Serve(addr)

	logger.Info("Server is up and running.")

	select {}
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&addr,"listen", "l", ":8080", "The network interface and port to listen on. Default :8080")
	serveCmd.PersistentFlags().StringVarP(&typeService,"type", "t", "proto", "The choice of how you'd like to interact with the service. Proto or HTTP")
}
