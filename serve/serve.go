package serve

import (
	desc "github.com/iain17/decentralizer/decentralizer"
)

var service desc.Decentralizer

func Setup() {
	if service != nil {
		return
	}
	var err error
	service, err = desc.New()
	if err != nil {
		panic(err)
	}
}