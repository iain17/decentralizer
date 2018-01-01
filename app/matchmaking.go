package app

import (
	"errors"
	"fmt"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/logger"
	"gx/ipfs/QmT6n4mspWYEya864BhCUJEgyxiRfmiSY9ruQwTUNpRKaM/protobuf/proto"
	inet "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	libp2pPeerStore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	"time"
	"encoding/hex"
	"github.com/iain17/timeout"
	"context"
	"github.com/iain17/framed"
	"github.com/iain17/kvcache/lttlru"
	"io"
	"github.com/iain17/decentralizer/app/ipfs"
	"strings"
	"github.com/giantswarm/retry-go"
)

type sessionRequest struct{
	peer   libp2pPeerStore.PeerInfo
	search *search
	connected bool
}

func (d *Decentralizer) getMatchmakingKey(sessionType uint64) string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_MATCHMAKING_%d", hex.EncodeToString(ih[:]), sessionType)
}

func (d *Decentralizer) initMatchmaking() {
	go d.GetIP()
	d.sessionQueries 			= make(chan sessionRequest, CONCURRENT_SESSION_REQUEST)
	d.sessions 					= make(map[uint64]*sessionstore.Store)
	d.sessionIdToSessionType	= make(map[uint64]uint64)
	var err error
	d.searches, err = lttlru.NewTTL(MAX_SESSION_SEARCHES)
	if err != nil {
		logger.Fatal(err)
	}
	d.i.PeerHost.SetStreamHandler(GET_SESSION_REQ, d.getSessionResponse)

	//Spawn some workers
	logger.Debugf("Running %d session request workers", CONCURRENT_SESSION_REQUEST)
	for i := 0; i < CONCURRENT_SESSION_REQUEST; i++ {
		go d.processSessionRequest()
	}
}

func (d *Decentralizer) processSessionRequest() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case req, ok := <-d.sessionQueries:
			logger.Infof("Querying %s for sessions of type: %d", req.peer.ID.Pretty(), req.search.sessionType)
			if !ok {
				return
			}
			addrs := req.peer.Addrs
			if !req.connected {
				//Add the addrs temp while we try and connect.
				addrs = ipfs.FilterNonReachableAddrs(req.peer.Addrs, true, false,false)
				if len(addrs) == 0 {
					logger.Debugf("No reachable addrs found for sessions provider: %s", req.peer.ID.Pretty())
					continue
				}
			}
			d.i.Peerstore.AddAddrs(req.peer.ID, addrs, 10 * time.Second)
			sessions, err := d.getSessionsRequest(req.peer.ID, req.search.sessionType)
			if err != nil {
				if err.Error() == "protocol not supported" || err == io.EOF {
					d.ignore.Add(req.peer.ID.Pretty(), true)
					continue
				}
				//give em another go.
				req.search.seen.Remove(req.peer.ID.Pretty())
				if err.Error() == "i/o deadline reached" {
					continue
				}
				logger.Debugf("Failed to get sessions from %s: %v", req.peer.ID.Pretty(), err)
			} else {
				logger.Debug("Received sessions. Adding them!")
				err = req.search.add(sessions, req.peer.ID)
				if err != nil {
					logger.Warning(err)
				}
			}
		}
	}
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
	d.matchmakingMutex.Lock()
	defer func() {
		d.matchmakingMutex.Unlock()
	}()
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(MAX_SESSIONS, time.Duration((EXPIRE_TIME_SESSION*1.5)*time.Second), d.i.Identity)
		if err != nil {
			logger.Warningf("Could not get session storage: %v", err)
			return nil
		}
	}
	return d.sessions[sessionType]
}

func (d *Decentralizer) getSessionSearch(sessionType uint64) (result *search) {
	d.searchMutex.Lock()
	defer func() {
		d.searchMutex.Unlock()
	}()
	if !d.searches.Contains(sessionType) {
		var err error
		result, err = d.newSearch(d.i.Context(), sessionType)
		d.searches.Add(sessionType, result)
		if err != nil {
			logger.Warningf("Could not start session search: %v", err)
			return nil
		}
	} else {
		value, ok := d.searches.Get(sessionType)
		if !ok {
			return nil
		}
		result = value.(*search)
	}
	return result
}

func (d *Decentralizer) UpsertSession(sessionType uint64, name string, port uint32, details map[string]string) (uint64, error) {
	sessions := d.getSessionStorage(sessionType)
	pId, dId := peerstore.PeerToDnId(d.i.Identity)
	info := &pb.Session{
		DnId:    dId,
		PId:     pId,
		Type:    sessionType,
		Name:    name,
		Address: uint32(utils.Inet_aton(d.GetIP())),
		Port:    port,
		Details: details,
	}
	//Also look for other sessions. So we can become a provider for more than just ourselves.
	go d.getSessionSearch(sessionType)
	sessionId, err := sessions.Insert(info)
	if err != nil {
		return 0, err
	}
	d.sessionIdToSessionType[sessionId] = sessionType
	timeout.Do(func(ctx context.Context) {
		d.b.Provide(d.getMatchmakingKey(sessionType))
	}, 3*time.Second)
	return sessionId, err
}

func (d *Decentralizer) setSessionIdToType(sessionId uint64, sessionType uint64) {
	d.matchmakingMutex.Lock()
	defer d.matchmakingMutex.Unlock()
	d.sessionIdToSessionType[sessionId] = sessionType
}

func (d *Decentralizer) DeleteSession(sessionId uint64) error {
	if d.sessionIdToSessionType[sessionId] == 0 {
		return errors.New("no such session exists")
	}
	sessionType := d.sessionIdToSessionType[sessionId]
	sessions := d.getSessionStorage(sessionType)
	return sessions.Remove(sessionId)
}

func (d *Decentralizer) GetSession(sessionId uint64) (*pb.Session, error) {
	if d.sessionIdToSessionType[sessionId] == 0 {
		return nil, errors.New("no such session exists")
	}
	sessionType := d.sessionIdToSessionType[sessionId]
	sessions := d.getSessionStorage(sessionType)
	return sessions.FindSessionId(sessionId)
}

func (d *Decentralizer) GetSessions(sessionType uint64) ([]*pb.Session, error) {
	search := d.getSessionSearch(sessionType)
	if search != nil {
		storage := search.fetch()
		return storage.FindAll()
	}
	return nil, errors.New("could not get session search")
}

func (d *Decentralizer) GetSessionsByDetails(sessionType uint64, key, value string) ([]*pb.Session, error) {
	search := d.getSessionSearch(sessionType)
	if search != nil {
		storage := search.fetch()
		return storage.FindByDetails(key, value)
	}
	return nil, errors.New("could not get session search")
}

func (d *Decentralizer) GetSessionsByPeer(peerId string) ([]*pb.Session, error) {
	var result []*pb.Session
	for _, key := range d.searches.Keys() {
		value, ok := d.searches.Get(key)
		if ! ok {
			continue
		}
		storage := value.(*search).fetch()
		peers, err := storage.FindByPeerId(peerId)
		if err != nil {
			logger.Warning(err)
			continue
		}
		result = append(result, peers...)
	}
	return result, nil
}

//Receive a request to give our sessions for a certain type
func (d *Decentralizer) getSessionResponse(stream inet.Stream) {
	reqData, err := framed.Read(stream)
	defer stream.Close()
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNSessionRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("%s requested our sessions of type %d...", stream.Conn().RemotePeer().Pretty(), request.Type)
	sessionsStorage := d.getSessionStorage(request.Type)
	sessions, err := sessionsStorage.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Sending %d sessions back of type %d", len(sessions), request.Type)
	//Response
	response, err := proto.Marshal(&pb.DNSessionResponse{
		Results: sessions,
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

// Get in contact with a peer and ask it for session of a certain type
func (d *Decentralizer) getSessionsRequest(peer peer.ID, sessionType uint64) ([]*pb.Session, error) {
	var stream inet.Stream
	op := func() (err error) {
		d.clearBackOff(peer)
		stream, err = d.i.PeerHost.NewStream(d.i.Context(), peer, GET_SESSION_REQ)
		return
	}
	err := retry.Do(op,
		retry.RetryChecker(func(err error) bool {
			//If there is something about dialing. Retry.
			if strings.Contains(err.Error(), "dial") {
				logger.Warning(err)
				return true
			}
			return false
		}),
		retry.MaxTries(30),
		retry.Timeout(10 * time.Minute),
		retry.Sleep(2 * time.Second))
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	logger.Debugf("Requesting %s for any sessions", peer.Pretty())
	//Request
	reqData, err := proto.Marshal(&pb.DNSessionRequest{
		Type: sessionType,
	})
	if err != nil {
		return nil, err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		return nil, err
	}
	var response pb.DNSessionResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}
	return response.Results, nil
}
