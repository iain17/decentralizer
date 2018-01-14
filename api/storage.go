package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"time"
	"github.com/iain17/timeout"
)

//
// Storage
//
// Write a user file. Takes a file name and the data it should save.
func (s *Server) WritePeerFile(ctx context.Context, request *pb.RPCWritePeerFileRequest) (*pb.RPCWritePeerFileResponse, error) {
	var err error
	timeout.Do(func(ctx context.Context) {
		_, err = s.App.SavePeerFile(request.Name, request.File)
		if err != nil {
			logger.Warning(err)
		}
	}, 5 * time.Second)
	return &pb.RPCWritePeerFileResponse{
		Success: err != nil,
	}, err
}

// Get a user file. Takes a file name, returns the file.
func (s *Server) GetPeerFile(ctx context.Context, request *pb.RPCGetPeerFileRequest) (*pb.RPCGetPeerFileResponse, error) {
	time_start := time.Now()
	file, err := s.App.GetPeerFile(request.PId, request.Name)
	if err != nil {
		logger.Warning(err)
	}
	logger.Infof("Responded get peer file request in: %s", time.Since(time_start).String())
	return &pb.RPCGetPeerFileResponse{
		File: file,
	}, err
}