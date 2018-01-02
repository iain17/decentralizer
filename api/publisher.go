package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"github.com/pkg/errors"
)

//Publish a new publisher update. (Only if you have the private key!)
func (s *Server) PublishPublisherUpdate(ctx context.Context, req *pb.RPCPublishPublisherUpdateRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.app.PublishPublisherUpdate(req.Definition)
}

// Get the full publisher definition
func (s *Server) GetPublisherDefinition(context.Context, *pb.GetPublisherDefinitionRequest) (*pb.PublisherDefinition, error) {
	definition := s.app.PublisherDefinition()
	if definition == nil {
		return nil, errors.New("No publisher definition set.")
	}
	return definition, nil
}

// Get a publisher file.
func (s *Server) GetPublisherFile(ctx context.Context, req *pb.RPCGetPublisherFileRequest) (*pb.RPCGetPublisherFileResponse, error) {
	file, err := s.app.GetPublisherFile(req.Name)
	if err != nil {
		logger.Warning(err)
	}
	return &pb.RPCGetPublisherFileResponse{
		File: file,
	}, err
}