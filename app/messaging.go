package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/framed"
	"github.com/golang/protobuf/proto"
	"github.com/giantswarm/retry-go"
	"time"
	"strings"
	"context"
	"fmt"
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

	ctx, cancel := context.WithTimeout(d.ctx, 10 * time.Minute) //TODO: configurable?
	defer cancel()
	var stream inet.Stream
	op := func() (err error) {
		logger.Infof("Trying to open stream to (to: %s:%d)", id.Pretty(), channel)
		d.clearBackOff(id)
		stream, err = d.i.PeerHost.NewStream(ctx, id, SEND_DIRECT_MESSAGE)
		return
	}
	err = retry.Do(op,
		retry.RetryChecker(func(err error) bool {
			//If there is something about dialing. Retry.
			if strings.Contains(err.Error(), "dial") {
				logger.Warningf("Failed to open stream to %s: %s", id.Pretty(), err)
				return true
			}
			return false
		}),
		retry.Timeout(5 * time.Minute),
		retry.Sleep(30 * time.Second))
	if err != nil {
		return err
	}
	defer stream.Close()
	logger.Infof("Opened a dialogue with %s", id.Pretty())

	//Request
	reqData, err := proto.Marshal(&pb.DNDirectMessageRequest{
		Channel: channel,
		Message: message,
	})
	if err != nil {
		err = fmt.Errorf("[%s] Could not marshal request: %s", id.Pretty(), err.Error())
		logger.Warning(err)
		return err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		err = fmt.Errorf("[%s] write failed: %s", id.Pretty(), err.Error())
		logger.Warning(err)
		return err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		err = fmt.Errorf("[%s] read failed: %s", id.Pretty(), err.Error())
		logger.Warning(err)
		return err
	}
	var response pb.DNDirectMessageResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		err = fmt.Errorf("[%s] Could not unmarshal response: %s", id.Pretty(), err.Error())
		logger.Warning(err)
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
