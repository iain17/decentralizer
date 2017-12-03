package decentralizer

import (
	"github.com/iain17/decentralizer/decentralizer/pb"
	"strconv"
	"time"
)

type Peer struct {
	*pb.Peer
	seen time.Time
}

func NewPeer(ip string, RPCPort, port uint32, details map[string]string) *Peer {
	return &Peer{
		Peer: &pb.Peer{
			Ip: ip,
			RpcPort: RPCPort,
			Port: port,
			Details: details,
		},
		seen: time.Now(),
	}
}

func (s *Peer) GetKey() string {
	return s.Ip + ":" + strconv.Itoa(int(s.RpcPort))
}

func (s *Peer) GetAddress() string {
	return s.Ip + ":" + strconv.Itoa(int(s.RpcPort))
}