package api

import (
	"context"
	"github.com/iain17/decentralizer/pb"
	"errors"
)

//
// Address book
//
// Create or update a peer. Takes peer info, returns if it was a success.
func (s *Server) UpsertPeer(context.Context, *pb.RPCUpsertPeerRequest) (*pb.RPCUpsertPeerResponse, error) {
	return nil, errors.New("not implemented")
}

// Get peer ids. takes a key and value to filter the peers by details. If left empty this filter will not apply and all will be fetched.
func (s *Server) GetPeerIds(context.Context, *pb.RPCGetPeerIdsRequest) (*pb.RPCGetPeerIdsResponse, error) {
	return nil, errors.New("not implemented")
}

// Get an individual peer. Takes either a peer id or decentralizer id and returns the peer info.
func (s *Server) GetPeer(context.Context, *pb.RPCGetPeerRequest) (*pb.RPCGetPeerResponse, error) {
	return nil, errors.New("not implemented")
}