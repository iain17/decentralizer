package serve

import (
	"http"
	"net"
	"github.com/mustafaakin/gongular"
)

func serveGrpc(lis net.Listener) {
	r := gongular.NewRouter()
	r.GET("/", Index)
	r.GET("/answer", SomePath)

	r.ListenAndServe(":8000")
}

func Index() string {
	return "Welcome to the decentralizer http service."
}

