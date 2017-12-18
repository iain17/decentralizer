package main

import (
	nethttp "net/http"

	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"

	"gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds/examples/adder"

	cmds "gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds"
	http "gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds/http"
)

func main() {
	h := http.NewHandler(
		cmds.Context{
			ConstructNode: func() (*core.IpfsNode, error) {
				return &core.IpfsNode{}, nil
			},
		},
		adder.RootCmd,
		http.NewServerConfig())

	// create http rpc server
	err := nethttp.ListenAndServe(":6798", h)
	if err != nil {
		panic(err)
	}
}
