package app

import (
	"context"
	"github.com/iain17/decentralizer/app/sessionstore"
	"sync"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"time"
	"github.com/iain17/timeout"
	"github.com/iain17/kvcache/lttlru"
	"gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	pstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	"fmt"
)

type search struct {
	mutex sync.Mutex
	d *Decentralizer
	sessionType uint64
	ctx context.Context
	storage *sessionstore.Store
	seen *lttlru.LruWithTTL//Tells us if we've already pased the result of this peer.
}

func (d *Decentralizer) newSearch(ctx context.Context, sessionType uint64) (*search, error) {
	storage := d.getSessionStorage(sessionType)
	seen, err := lttlru.NewTTL(MAX_IGNORE)
	if err != nil {
		return nil, err
	}
	instance := &search{
		d: d,
		sessionType: sessionType,
		ctx: ctx,
		storage: storage,
		seen: seen,
	}
	go instance.connectToProviders()
	d.cron.Every(30).Seconds().Do(instance.search)
	d.cron.Every(30).Seconds().Do(storage.Save)
	return instance, nil
}

//Takes as long as it needs
func (s *search) search() error {
	ctx, cancel := context.WithTimeout(s.ctx, 30 * time.Second)
	defer cancel()
	s.connectToProviders()
	return s.run(ctx)
}

func (s *search) run(ctx context.Context) error {
	logger.Info("Trying to lock mutex")
	s.mutex.Lock()

	defer s.mutex.Unlock()

	logger.Infof("Searching for sessions with type %d", s.sessionType)
	values, err := s.d.b.GetValues(ctx, DHT_SESSIONS_KEY_TYPE, s.d.getMatchmakingKey(s.sessionType), 3)
	if err != nil {
		return fmt.Errorf("could not find session with type %d: %s", s.sessionType, err.Error())
	}
	queried := 0
	total := 0
	s.seen.Purge()
	self := s.d.i.Identity.Pretty()
	for _, value := range values {
		total++
		id := value.From.Pretty()
		if id == self {
			continue
		}
		if s.d.ignore.Contains(id) {
			logger.Debugf("%s is on the ignore list", id)
			continue
		}
		if s.seen.Contains(id) {
			logger.Debugf("%s is on the seen list", id)
			continue
		}
		s.seen.Add(id, true)

		var record pb.DNSessionsRecord
		err = s.d.unmarshal(value.Val, &record)
		if err != nil {
			logger.Warning(err)
			continue
		}
		logger.Infof("Got %d sessions from %s", len(record.Results), id)
		for _, session := range record.Results {
			if session.PId == self {
				continue
			}
			s.d.setSessionIdToType(session.SessionId, session.Type)
			_, err := s.storage.Insert(session)
			if err != nil {
				logger.Warning(err)
			}
		}
		if s.d.i.PeerHost.Network().Connectedness(value.From) == net.Connected {
			s.d.sessionQueries <- sessionRequest{
				search: s,
				id:     value.From,
			}
			queried++
		}
	}
	logger.Infof("Queried %d of the %d for sessions of type %d", queried, total, s.sessionType)
	return nil
}

func (s *search) connectToProviders() {
	values := s.d.b.Find(s.d.getMatchmakingKey(s.sessionType), 1)
	seen := make(map[string]bool)
	for value := range values {
		id := value.ID.Pretty()
		if seen[id] {
			continue
		}
		seen[id] = true
		go s.d.i.PeerHost.Connect(s.d.i.Context(), pstore.PeerInfo{
			ID:    value.ID,
			Addrs: value.Addrs,
		})
		if s.d.i.PeerHost.Network().Connectedness(value.ID) == net.Connected {
			s.d.sessionQueries <- sessionRequest{
				search: s,
				id:     value.ID,
			}
		}
	}
}

func (s *search) refresh(ctx context.Context) error {
	err := s.run(ctx)
	if err != nil {
		logger.Warning(err)
	}
	return err
}

func (s *search) fetch() (*sessionstore.Store, error) {
	searchCtx, cancel := context.WithTimeout(s.d.i.Context(), 5 * time.Second)
	var err error
	timeout.Do(func(ctx context.Context) {
		tries := 0
		err = s.refresh(searchCtx)
		for s.storage.IsEmpty() {
			select {
			case <- ctx.Done():
				return
			default:
				if tries > 5 {
					break
				}
				err = s.refresh(searchCtx)
				time.Sleep(1 * time.Second)
				tries++
			}
		}
	}, 5 * time.Second)
	cancel()

	//Keep it to yourself eh. If we have results. Show em!
	if err != nil && s.storage.Len() > 0 {
		logger.Warning(err)
		err = nil
	}
	s.storage.Save()
	return s.storage, err
}