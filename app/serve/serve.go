package serve

import (
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"github.com/iain17/dht-hello/p2p"
)

func Run(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ipfsRoot, err := cmd.Flags().GetString("ipfsRoot")

	if err != nil {
		panic(err)
	}

	p2p, err := p2p.New(ctx, ipfsRoot)
	if err != nil {
		panic(err)
	}

	p2p.NewApp("test")
	select{}
	//
	//router := gin.Default()
	//v1 := router.Group("/v1/:app")
	//{
	//	v1.GET("/peers", getPeersEndpoint)
	//	v1.POST("/", )
	//	v1.POST("/read", readEndpoint)
	//}
	//
	//router.Run()
}