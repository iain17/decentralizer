package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
)

func (s *Server) PublishPublisherUpdate(ctx context.Context, req *pb.RPCPublishPublisherUpdateRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.app.PublishPublisherUpdate(req.Definition)
}