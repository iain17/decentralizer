package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/decentralizer/service/restapi/operations"
	"github.com/go-openapi/swag"
	"github.com/iain17/decentralizer/service/models"
)

func StartSearch(params operations.StartSearchParams) middleware.Responder {
	err := dService.AddService(params.AppName, uint32(*params.Port))
	if err != nil {
		return operations.NewStartSearchDefault(500).WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}
	return operations.NewStartSearchOK()
}