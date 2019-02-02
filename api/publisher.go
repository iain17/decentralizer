package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"github.com/pkg/errors"
)

// Load a publisher definition. Will not work if its older!
func (s *Server) ReadPublisherDefinition(ctx context.Context, req *pb.LoadPublisherDefinitionRequest) (*pb.Empty, error) {
	if len(req.Definition) == 0 {
		return &pb.Empty{}, s.App.PushPublisherRecord()
	}
	s.App.ResetPublisherDefinition()
	err := s.App.ReadPublisherDefinition(req.Definition)
	return &pb.Empty{}, err
}

//Publish a new publisher update. (Only if you have the private key!)
func (s *Server) PublishPublisherUpdate(ctx context.Context, req *pb.RPCPublishPublisherUpdateRequest) (*pb.RPCPublishPublisherUpdateResponse, error) {
	definition, err := s.App.PublishPublisherRecord(req.Definition)
	return &pb.RPCPublishPublisherUpdateResponse{
		Definition: definition,
	}, err
}

// Get the full publisher definition
func (s *Server) GetPublisherDefinition(context.Context, *pb.GetPublisherDefinitionRequest) (*pb.PublisherDefinition, error) {
	definition := s.App.PublisherDefinition()
	if definition == nil {
		return nil, errors.New("No publisher definition set.")
	}
	return definition, nil
}

// Get a publisher file.
func (s *Server) GetPublisherFile(ctx context.Context, req *pb.RPCGetPublisherFileRequest) (*pb.RPCGetPublisherFileResponse, error) {
	file, err := s.App.GetPublisherFile(req.Name)
	if err != nil {
		logger.Warning(err)
	}
	return &pb.RPCGetPublisherFileResponse{
		File: file,
	}, err
}