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
	httpListen string
	protoListen string
)

func run(cmd *cobra.Command, args []string) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.DebugLevel)

	//Components
	serve.Setup()
	go serve.ServeHttp(httpListen)
	go serve.ServeGrpc(protoListen)

	select {}
}

func init() {
	serveCmd.Flags().StringVarP(&httpListen,"httpListen", "", "127.0.0.1:8080", "The network interface and port that the http server will listen on. Default 127.0.0.1:8080")
	serveCmd.Flags().StringVarP(&protoListen,"ProtoListen", "", "127.0.0.1:8081", "The network interface and port that the protobuf server will listen on. Default 127.0.0.1:8081")
	RootCmd.AddCommand(serveCmd)
}
