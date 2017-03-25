package serve

import (
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"github.com/iain17/decentralizer/decentralizer/pb"
	"net"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	logger "github.com/Sirupsen/logrus"
)

type decentralizer struct {

}

func ServeGrpc(addr string) pb.DecentralizerServer {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	instance := &decentralizer{

	}
	s := grpc.NewServer()
	pb.RegisterDecentralizerServer(s, instance)
	reflection.Register(s)

	logger.Infof("Protobuf server listening at %s", addr)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
	return instance
}


// Start the search for peers. Returns messages received from peers.
func (d *decentralizer) Search(req *pb.SearchRequest, stream pb.Decentralizer_SearchServer) error {
	decService := service.GetService(req.GetName())
	if decService != nil {
		return errors.New("Service already exists.")
	}
	service.AddService(req.GetName(), uint32(req.GetPort()))

	for {
		select {
		// If the context is finished, don't bother processing the
		// message.
		case <-stream.Context().Done():
			break
		default:
		}
	}
	return service.StopService(req.GetName())
}

// Get peers found for a specific service.
func (d *decentralizer) GetPeers(ctx context.Context, req *pb.GetPeersRequest) (*pb.GetPeersResponse, error) {
	decService := service.GetService(req.GetName())
	if decService == nil {
		return nil, errors.New("Service does not exist.")
	}
	return &pb.GetPeersResponse{
		Peers: decService.GetPeers(),
	}, nil
}