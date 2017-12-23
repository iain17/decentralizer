package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
)

//
// Storage
//
// Write a user file. Takes a file name and the data it should save.
func (s *Server) WritePeerFile(ctx context.Context, request *pb.RPCWritePeerFileRequest) (*pb.RPCWritePeerFileResponse, error) {
	_, err := s.app.SavePeerFile(request.Name, request.File)
	return &pb.RPCWritePeerFileResponse{
		Success: err == nil,
	}, err
}

// Get a user file. Takes a file name, returns the file.
func (s *Server) GetPeerFile(ctx context.Context, request *pb.RPCGetPeerFileRequest) (*pb.RPCGetPeerFileResponse, error) {
	file, err := s.app.GetPeerFile(request.PId, request.Name)
	return &pb.RPCGetPeerFileResponse{
		File: file,
	}, err
}

// Get a publisher file.
func (s *Server) GetPublisherFile(ctx context.Context, req *pb.RPCGetPublisherFileRequest) (*pb.RPCGetPublisherFileResponse, error) {
	file, err := s.app.GetPublisherFile(req.Name)
	return &pb.RPCGetPublisherFileResponse{
		File: file,
	}, err
}