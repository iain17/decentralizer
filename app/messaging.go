package app

import (
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-peer"
)

type DirectMessage struct {
	From peer.ID
	Message []byte
}

func (d *Decentralizer) initMessaging() {
	d.i.PeerHost.SetStreamHandler(SEND_DIRECT_MESSAGE, d.directMessageReceived)
}

func (d *Decentralizer) SendMessage(peerId string, message []byte) (error) {
	id, err := libp2pPeer.IDB58Decode(peerId)
	if err != nil {
		return err
	}

	stream, err := d.i.PeerHost.NewStream(d.i.Context(), id, SEND_DIRECT_MESSAGE)
	if err != nil {
		return err
	}

	//Request
	reqData, err := proto.Marshal(&pb.DNDirectMessageRequest{
		Message: message,
	})
	if err != nil {
		return err
	}
	err = Write(stream, reqData)
	if err != nil {
		return err
	}

	//Response
	resData, err := Read(stream)
	if err != nil {
		return err
	}
	var response pb.DNDirectMessageResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decentralizer) directMessageReceived(stream inet.Stream) {
	reqData, err := Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNDirectMessageRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}

	d.directMessage <- &DirectMessage{
		Message: request.Message,
	}

	//Response
	response, err := proto.Marshal(&pb.DNDirectMessageResponse{
		Delivered: true,
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