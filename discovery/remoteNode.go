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
	transfer, err := proto.Marshal(&pb.Hearbeat{})
	if err != nil {
		return err
	}
	return rn.write(transfer)
}

func (rn *RemoteNode) Send(message string) error {
	transfer, err := proto.Marshal(&pb.Transfer{
		Data: message,
	})
	if err != nil {
		return err
	}
	return rn.write(transfer)
}

func (rn *RemoteNode) write(data []byte) error {
	rn.logger.Debug("sending message...")
	_, err := rn.conn.Write(data)
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

		switch packet.GetMsg().(type) {
		case *pb.Message_Heartbeat :
			rn.logger.Debug("heard beat received")
			rn.lastHeartbeat = time.Now()
			break
		}
	}
}
