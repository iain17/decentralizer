package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
)

func (d *Decentralizer) initAddressbook() {
	d.i.PeerHost.SetStreamHandler(GET_PEER_REQ, d.getPeerResponse)
}

func (d *Decentralizer) UpsertPeer(pId string, details map[string]string) error {
	err := d.peers.Upsert(&pb.Peer{
		PId:     pId,
		Details: details,
	})
	if err != nil {
		return err
	}
	return err
}

func (d *Decentralizer) getPeerResponse(stream inet.Stream) {
	reqData, err := Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNPeerRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}
	peer, err := d.peers.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		logger.Error(err)
		return
	}

	//Response
	response, err := proto.Marshal(&pb.DNPeerResponse{
		Peer: peer,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	err = Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}