package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"errors"
)

//
// Storage
//
// Write a user file. Takes a file name and the data it should save.
func (s *Server) WriteUserFile(context.Context, *pb.RPCWriteUserFileRequest) (*pb.RPCWriteUserFileResponse, error) {
	return nil, errors.New("Unimplemented")
}

// Get a user file. Takes a file name, returns the file.
func (s *Server)  GetUserFile(context.Context, *pb.RPCGetUserFileRequest) (*pb.RPCGetUserFileResponse, error) {
	return nil, errors.New("Unimplemented")
}

// Get a publisher file.
func (s *Server) GetPublisherFile(context.Context, *pb.RPCGetPublisherFileRequest) (*pb.RPCGetPublisherFileResponse, error) {
	return nil, errors.New("Unimplemented")
}