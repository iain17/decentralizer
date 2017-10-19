package app

import (
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/decentralizer/utils"
	"time"
	"github.com/golang/protobuf/proto"
	"github.com/ethersphere/go-ethereum/logger"
	"github.com/iain17/decentralizer/app/ipfs"
)

func (d *Decentralizer) UpsertSession(sessionType uint32, name string, port uint32, details map[string]string) (uint64, error) {
	pId, dId := pb.GetPeer(d.i.Identity)
	info := &pb.SessionInfo{
		DId: dId,
		PId: pId,
		Type: sessionType,
		Name: name,
		Address: utils.Inet_aton(d.d.GetIP()),
		Port: port,
		Details: details,
	}
	return d.sessions.Insert(info)
}

func (d *Decentralizer) DeleteSession(sessionId uint64) error {
	return d.sessions.Remove(sessionId)
}

//Every 30 seconds advertise the sessions we've got going on.
func (d *Decentralizer) Advertise() {
	for {
		time.Sleep(EXPIRE_TIME_SESSION * time.Second)

		localSessions, err := d.sessions.FindByPeerId(d.i.Identity.Pretty())
		if err != nil {
			logger.Warn(err)
			continue
		}
		logger.Info("Advertising %d of sessions", len(localSessions))
		for _, sessInfo := range localSessions {
			msg, err := proto.Marshal(&pb.DMessage{
				Version: pb.VERSION,
				Msg: &pb.DMessage_UpsertSession{
					UpsertSession: &pb.UpsertSession{
						Info: sessInfo,
					},
				},
			})
			if err != nil {
				logger.Warn(err)
				continue
			}
			err = ipfs.Publish(d.i, getPubSubKey(sessInfo.Type), msg)
		}
	}
}

//TODO: Validate the session first.
func (d *Decentralizer) GetSessions(sessionType uint32, details map[string]string) {
	//ipfs.Receive(key(sessionType), func(peer peer.ID, message string) {
	//	logger.Infof("Received: %s: %s\n", peer.String(), message)
	//})
}

func getPubSubKey(sessionType uint32) string {
	return fmt.Sprintf("MATCHMAKING_%d", sessionType)
}