package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/decentralizer/service/restapi/operations"
	"github.com/go-openapi/swag"
	"github.com/iain17/decentralizer/service/models"
)

func StopSearch(params operations.StopSearchParams) middleware.Responder {
	err := dService.StopService(params.AppName)
	if err != nil {
		return operations.NewGetPeersDefault(int(500)).WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}
	return operations.NewStopSearchOK()
}