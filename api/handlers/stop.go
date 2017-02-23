package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/dht-hello/service/restapi/operations"
)

func StopSearch(params operations.StopSearchParams) middleware.Responder {
	return middleware.NotImplemented("sorry")
	//err := dht.StopSearch(params.Identifier)
	//if err != nil {
	//	code := err.Code()
	//	return operations.NewGetPeersDefault(int(code)).WithPayload(&models.Error{
	//		Message: swag.String(err.Error()),
	//		Code: &code,
	//	})
	//}
	//return operations.NewStopSearchOK()
}