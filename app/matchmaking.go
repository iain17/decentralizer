package app

import (
	"fmt"
	"github.com/iain17/decentralizer/app/pb"
	"github.com/iain17/decentralizer/utils"
	"time"
	"github.com/golang/protobuf/proto"
	"github.com/iain17/decentralizer/app/ipfs"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/app/sessionstore"
	"errors"
)

func (d *Decentralizer) UpsertSession(sessionType uint32, name string, port uint32, details map[string]string) (uint64, error) {
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(1000, time.Duration((EXPIRE_TIME_SESSION * 1.5) * time.Second))
		if err != nil {
			return 0, err
		}
	}
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
	return d.sessions[sessionType].Insert(info)
}

func (d *Decentralizer) DeleteSession(sessionType uint32, sessionId uint64) error {
	if d.sessions[sessionType] == nil {
		return errors.New("no such sessionType exists")
	}
	return d.sessions[sessionType].Remove(sessionId)
}

//Every x amount of seconds advertise the sessions we are hosting.
func (d *Decentralizer) Advertise() {
	for {
		time.Sleep(EXPIRE_TIME_SESSION * time.Second)

		for _, sessionStore := range d.sessions {

			localSessions, err := sessionStore.FindByPeerId(d.i.Identity.Pretty())
			if err != nil {
				logger.Warning(err)
				continue
			}
			logger.Info("Advertising %d of sessions", len(localSessions))
			for _, sessionInfo := range localSessions {
				msg, err := proto.Marshal(&pb.DMessage{
					Version: pb.VERSION,
					Msg: &pb.DMessage_UpsertSession{
						UpsertSession: &pb.UpsertSession{
							Info: sessionInfo,
						},
					},
				})
				if err != nil {
					logger.Warning(err)
					continue
				}
				err = ipfs.Publish(d.i, getPubSubKey(sessionInfo.Type), msg)
			}

		}
	}
}

//TODO: Validate the session first.
func (d *Decentralizer) GetSessions(sessionType uint32, key, value string) error {
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(1000, time.Duration((EXPIRE_TIME_SESSION * 1.5) * time.Second))
		if err != nil {
			return err
		}
	}
	//Start listening to these kind of sessions.
	if d.subscriptions[sessionType] == nil {
		ipfs.Subscribe(d.i, getPubSubKey(sessionType), d.receivedSession)
	}
	localSessions, err := d.sessions[sessionType].FindByDetails(key, value)
}

func (d *Decentralizer) receivedSession(peer peer.ID, message []byte) {
	logger.Infof("Received: %s: %s\n", peer.String(), message)
}

func getPubSubKey(sessionType uint32) string {
	return fmt.Sprintf("MATCHMAKING_%d", sessionType)
}