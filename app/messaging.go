package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
	inet "gx/ipfs/QmU4vCDZTPLDqSDKguWbHCiUe46mZUtmM2g2suBZ9NE8ko/go-libp2p-net"
	libp2pPeer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"github.com/iain17/framed"
)

type DirectMessage struct {
	From    libp2pPeer.ID
	Message []byte
}

func (d *Decentralizer) initMessaging() {
	d.i.PeerHost.SetStreamHandler(SEND_DIRECT_MESSAGE, d.directMessageReceived)
}

func (d *Decentralizer) GetMessagingChan(channel uint32) chan *pb.RPCDirectMessage {
	if d.directMessageChannels[channel] == nil {
		d.directMessageChannels[channel] = make(chan *pb.RPCDirectMessage, 10)
	}
	return d.directMessageChannels[channel]
}

func (d *Decentralizer) SendMessage(channel uint32, peerId string, message []byte) error {
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return err
	}
	messageChannel := d.GetMessagingChan(channel)
	logger.Infof("Sending direct message (to: %s:%d)", id.Pretty(), channel)

	if id.Pretty() == d.i.Identity.Pretty() {
		messageChannel <- &pb.RPCDirectMessage{
			Channel: channel,
			PId: id.Pretty(),
			Message: message,
		}
		return nil
	}

	stream, err := d.i.PeerHost.NewStream(d.i.Context(), id, SEND_DIRECT_MESSAGE)
	if err != nil {
		return err
	}
	defer stream.Close()

	//Request
	reqData, err := proto.Marshal(&pb.DNDirectMessageRequest{
		Channel: channel,
		Message: message,
	})
	if err != nil {
		return err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		return err
	}

	//Response
	resData, err := framed.Read(stream)
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
	reqData, err := framed.Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	from := stream.Conn().RemotePeer()
	logger.Infof("Received direct message from %s", from.Pretty())
	var request pb.DNDirectMessageRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}

	messageChannel := d.GetMessagingChan(request.Channel)
	messageChannel <- &pb.RPCDirectMessage{
		Channel: request.Channel,
		PId: from.Pretty(),
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
	err = framed.Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}
