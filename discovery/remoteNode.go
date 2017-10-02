package discovery

import (
	"fmt"
	"time"
	"net"
	"github.com/op/go-logging"
	"github.com/iain17/decentralizer/discovery/pb"
	"github.com/golang/protobuf/proto"
	"io"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	lastHeartbeat time.Time
}

func NewRemoteNode(conn net.Conn) *RemoteNode {
	return &RemoteNode{
		Node: Node{
			logger:        logging.MustGetLogger(fmt.Sprintf("RemoteNode(%s)", conn.RemoteAddr().String())),
		},
		conn:          conn,
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) sendHeartBeat() error {
	rn.logger.Debug("Sending heartbeat")
	return rn.write(pb.HeartBeatMessage, []byte{'L', 'O', 'V', 'E'})
}

func (rn *RemoteNode) Send(message string) error {
	transfer, err := proto.Marshal(&pb.Transfer{
		Data: message,
	})
	if err != nil {
		return err
	}
	return rn.write(pb.TransferMessage, transfer)
}

func (rn *RemoteNode) write(messageType pb.MessageType, data []byte) error {
	rn.logger.Debug("sending message...")
	packet := pb.NewPacket(messageType, data)
	err := packet.Write(rn.conn)
	if err != nil {
		return err
	}
	rn.logger.Debug("message sent")
	return nil
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Debug("listener stopped...")
	defer func() {
		ln.netTableService.RemoveRemoteNode(rn.conn.RemoteAddr())
	}()

	rn.logger.Debug("listening...")
	for {
		packet, err := pb.Decode(rn.conn)
		if err != nil {
			rn.logger.Error("decode error, %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		rn.logger.Debug("received, %+v", packet)

		switch packet.Body.Type {
		case pb.HeartBeatMessage:
			rn.logger.Debug("heard beat received: %s", packet.Body.Data)
			rn.lastHeartbeat = time.Now()
		}
	}
}
