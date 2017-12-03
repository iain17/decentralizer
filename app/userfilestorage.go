package app

import (
	//"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/iain17/logger"
	//"strings"
)

func (d *Decentralizer) SaveUserFile() {
	d.b.Provide("cool")
	logger.Infof("OUR id: %s", d.i.Identity.Pretty())
	go d.search()
}

func (d *Decentralizer) search() {
	peerIds := d.b.Find("cool", 0)
	for id := range peerIds {
		logger.Infof("Provider: %s", id.Pretty())
	}
	logger.Info("done")
}

func test() {
	//api := coreapi.NewCoreAPI(d.i)
	//path, err := api.Unixfs().Add(d.i.Context(), strings.NewReader("fuck off"))
	//if err != nil {
	//	panic(err)
	//}
}