package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iain17/dht-hello/service/restapi/operations"
	"github.com/iain17/dht-hello/dht"
	"github.com/iain17/dht-hello/service/models"
	"github.com/go-openapi/swag"
	logger "github.com/Sirupsen/logrus"
	"github.com/go-openapi/strfmt"
)

func GetPeers(params operations.GetPeersParams) middleware.Responder {
	peers, err := dht.GetPeers(params.Identifier)
	//If it doesn't exist create it.
	if err != nil && err.Code() == 404 {
		logger.Warn(err)
		dht.Search(params.Identifier, 0, true)
		peers, err = dht.GetPeers(params.Identifier)
	}
	if err != nil {
		return operations.NewGetPeersDefault(int(err.Code())).WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}
	var results models.Peers
	for _, peer := range peers {
		results = append(results, &models.Peer{
			IP: strfmt.IPv4(peer.IP.String()),
			Port: int32(peer.Port),
		})
	}
	return operations.NewGetPeersOK().WithPayload(results)
}