package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/dht-hello/service/restapi/operations"
	"github.com/iain17/dht-hello/dht"
	"github.com/go-openapi/swag"
	"github.com/iain17/dht-hello/service/models"
)

func StopSearch(params operations.StopSearchParams) middleware.Responder {
	err := dht.StopSearch(params.Identifier)
	if err != nil {
		return operations.NewGetPeersDefault(int(err.Code())).WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}
	return operations.NewStopSearchOK()
}