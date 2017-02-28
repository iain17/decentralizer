package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/decentralizer/service/restapi/operations"
	"github.com/go-openapi/strfmt"
	"github.com/iain17/decentralizer/service/models"
	"github.com/go-openapi/swag"
)

func GetPeers(params operations.GetPeersParams) middleware.Responder {
	service := dService.GetService(params.AppName)
	if service != nil {
		return operations.NewGetPeersDefault(404).WithPayload(&models.Error{
			Message: swag.String("Service does not exist."),
		})
	}
	peers := service.GetPeers()
	results := models.Peers{}
	for _, peer := range peers {
		var details []*models.Detail
		for key, value := range peer.Details {
			details = append(details, &models.Detail{
				Name: key,
				Value: value,
			})
		}
		results = append(results, &models.Peer{
			Details: details,
			IP: strfmt.IPv4(peer.Ip),
			Port: int32(peer.Port),
		})
	}
	return operations.NewGetPeersOK().WithPayload(results)
}