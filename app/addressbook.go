package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
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

func (d *Decentralizer) GetPeersByDetails(key, value string) ([]*pb.Peer, error) {
	return d.peers.FindByDetails(key, value)
}

func (d *Decentralizer) GetPeers() ([]*pb.Peer, error) {
	return d.peers.FindAll()
}

func (d *Decentralizer) FindByPeerId(peerId string) (p *pb.Peer, err error) {
	p, err = d.peers.FindByPeerId(peerId)
	if err != nil {
		var id peer.ID
		id, err = libp2pPeer.IDB58Decode(peerId)
		if err != nil {
			return nil, err
		}
		p, err = d.getPeerRequest(id)
		d.peers.Upsert(p)
	}
	return p, err
}

func (d *Decentralizer) getPeerRequest(peer peer.ID) (*pb.Peer, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_PEER_REQ)
	if err != nil {
		return nil, err
	}
	//Request
	reqData, err := proto.Marshal(&pb.DNPeerRequest{})
	if err != nil {
		return nil, err
	}
	err = Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := Read(stream)
	if err != nil {
		return nil, err
	}
	var response pb.DNPeerResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}
	return response.Peer, nil
}

func (d *Decentralizer) FindByDecentralizedId(decentralizedId uint64) (*pb.Peer, error) {
	return d.peers.FindByDecentralizedId(decentralizedId)
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