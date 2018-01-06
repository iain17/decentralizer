package app

import (
	"errors"
	"fmt"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/decentralizer/app/sessionstore"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/logger"
	"time"
	"encoding/hex"
	"github.com/iain17/timeout"
	"context"
	"github.com/iain17/kvcache/lttlru"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
)

func (d *Decentralizer) getMatchmakingKey(sessionType uint64) string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_MATCHMAKING_%d", hex.EncodeToString(ih[:]), sessionType)
}

func (d *Decentralizer) initMatchmaking() {
	go d.GetIP()
	d.sessions 					= make(map[uint64]*sessionstore.Store)
	d.sessionIdToSessionType	= make(map[uint64]uint64)
	var err error
	d.searches, err = lttlru.NewTTL(MAX_SESSION_SEARCHES)
	if err != nil {
		logger.Fatal(err)
	}

	d.b.RegisterValidator(DHT_SESSIONS_KEY_TYPE, func(key string, val []byte) error{
		var sessions pb.DNSessions
		err = gogoProto.Unmarshal(val, &sessions)
		if err != nil {
			return err
		}
		return validateDNSessions(&sessions)
	}, true)
}

func validateDNSessions(sessions *pb.DNSessions) error {
	//Check publish time
	now := time.Now().UTC()
	publishedTime := time.Unix(int64(sessions.Published), 0).UTC()
	publishedTimeText := publishedTime.String()
	expireTime := now.Add(-EXPIRE_TIME_SESSION)
	expireTimeText := expireTime.String()
	if publishedTime.Before(expireTime) {
		err := fmt.Errorf("record with publish date %s has expired. It was before %s", publishedTimeText, expireTimeText)
		logger.Warning(err)
		return err
	}
	if publishedTime.After(now) {
		err := errors.New("record with publish date %s was published in the future")
		logger.Warning(err)
		return err
	}
	logger.Infof("successfully validated DNSessions published at: %s", publishedTimeText)
	return nil
}

func (d *Decentralizer) getSessionStorage(sessionType uint64) *sessionstore.Store {
	d.matchmakingMutex.Lock()
	defer func() {
		d.matchmakingMutex.Unlock()
	}()
	if d.sessions[sessionType] == nil {
		var err error
		d.sessions[sessionType], err = sessionstore.New(d.ctx, MAX_SESSIONS, time.Duration(EXPIRE_TIME_SESSION), d.i.Identity)
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
	sessionId, err := sessions.Insert(info)
	if err != nil {
		return 0, err
	}
	d.sessionIdToSessionType[sessionId] = sessionType
	timeout.Do(func(ctx context.Context) {
		err := d.advertise(sessionType)
		if err != nil {
			logger.Errorf("Could not advertise session: %s", err.Error())
		}
	}, 5*time.Second)
	return sessionId, err
}

//Advertise all the session ids we have
func (d *Decentralizer) advertise(sessionType uint64) error {
	//Before we override DHT with our advisement. Let us check others.
	search := d.getSessionSearch(sessionType)
	store := search.fetch()
	localSessions, err := store.FindByPeerId(d.i.Identity.Pretty())
	if err != nil {
		return err
	}
	sessions := &pb.DNSessions{
		Published: uint64(time.Now().UTC().Unix()),
		Results: localSessions,
	}
	err = validateDNSessions(sessions)
	if err != nil {
		return err
	}
	data, err := gogoProto.Marshal(sessions)
	if err != nil {
		return err
	}
	return d.b.PutValue(DHT_SESSIONS_KEY_TYPE, d.getMatchmakingKey(sessionType), data)
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