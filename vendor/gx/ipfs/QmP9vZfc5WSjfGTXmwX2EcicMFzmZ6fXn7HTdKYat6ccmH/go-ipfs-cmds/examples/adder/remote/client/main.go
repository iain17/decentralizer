package main

import (
	"os"

	"gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds/examples/adder"

	cmds "gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds"
	cli "gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds/cli"
	http "gx/ipfs/QmP9vZfc5WSjfGTXmwX2EcicMFzmZ6fXn7HTdKYat6ccmH/go-ipfs-cmds/http"
	cmdkit "gx/ipfs/QmQp2a2Hhb7F6eK2A5hN8f9aJy4mtkEikL9Zj4cgB7d1dD/go-ipfs-cmdkit"
)

func main() {
	// parse the command path, arguments and options from the command line
	req, cmd, _, err := cli.Parse(os.Args[1:], os.Stdin, adder.RootCmd)
	if err != nil {
		panic(err)
	}

	// create http rpc client
	client := http.NewClient(":6798")

	// send request to server
	res, err := client.Send(req)
	if err != nil {
		panic(err)
	}

	req.SetOption("encoding", cmds.Text)

	// create an emitter
	re, retCh := cli.NewResponseEmitter(os.Stdout, os.Stderr, cmd.Encoders["Text"], req)

	if pr, ok := cmd.PostRun[cmds.CLI]; ok {
		re = pr(req, re)
	}

	wait := make(chan struct{})
	// copy received result into cli emitter
	go func() {
		err = cmds.Copy(re, res)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal|cmdkit.ErrFatal)
		}
		close(wait)
	}()

	// wait until command has returned and exit
	ret := <-retCh
	<-wait
	os.Exit(ret)
}
