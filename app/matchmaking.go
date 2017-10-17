package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/golang/protobuf/proto"
)

func (d *Decentralizer) CreateSession(sessionType uint32, name string, address uint64, port uint32, details map[string]string) error {
	dPeer := pb.GetPeer(d.i.Identity)
	peerInfo, err := proto.Marshal(&pb.DMessage{
		Version: pb.VERSION,
		Msg: &pb.DMessage_UpsertSession{
			UpsertSession: &pb.UpsertSession{
				Info: &pb.SessionInfo{
					Owner: &dPeer,
					Type: sessionType,
					Name: name,
					Address: address,
					Port: port,
					Details: details,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return ipfs.Publish(d.i, key(sessionType), peerInfo)
}

func (d *Decentralizer) GetSession(sessionType int32) {
	//ipfs.Receive(topic, func(peer peer.ID, message string) {
	//	logger.Infof("Received: %s: %s\n", peer.String(), message)
	//})
}

func key(sessionType uint32) string {
	return fmt.Sprintf("MATCHMAKING_%d", sessionType)
}