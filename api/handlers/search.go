package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"fmt"
	"github.com/iain17/dht-hello/service/restapi/operations"
)

func StartSearch(params operations.StartSearchParams) middleware.Responder {
	fmt.Println(*params.Port)
	return middleware.NotImplemented("operation .StartSearch has not yet been implemented")
}