package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/dht-hello/service/restapi/operations"
	"github.com/iain17/dht-hello/dht"
	"github.com/iain17/dht-hello/service/models"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/strfmt"
)

func GetPeers(params operations.GetPeersParams) middleware.Responder {
	peers, err := dht.GetClients(params.Identifier)
	if err != nil {
		return operations.NewGetPeersDefault(int(err.Code())).WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}
	results := models.Peers{}
	for _, peer := range peers {
		results = append(results, &models.Peer{
			IP: strfmt.IPv4(peer.IP.String()),
			Port: int32(peer.Port),
		})
	}
	return operations.NewGetPeersOK().WithPayload(results)
}