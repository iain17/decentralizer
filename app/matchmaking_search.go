package app

import (
	"context"
	"github.com/iain17/decentralizer/app/sessionstore"
	"sync"
	"github.com/iain17/logger"
	Peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/decentralizer/pb"
	"time"
	"github.com/Akagi201/kvcache/ttlru"
	"github.com/jeffchao/backoff"
	"github.com/pkg/errors"
	"fmt"
)

type search struct {
	running bool
	updating bool
	backoff *backoff.FibonacciBackoff
	mutex sync.Mutex
	fetching sync.Mutex
	d *Decentralizer
	sessionType uint64
	ctx context.Context
	storage *sessionstore.Store
	seen *lru.LruWithTTL
}

func (d *Decentralizer) newSearch(ctx context.Context, sessionType uint64) (*search, error) {
	storage := d.getSessionStorage(sessionType)
	seen, err := lru.NewTTL(MAX_IGNORE)
	if err != nil {
		return nil, err
	}
	f := backoff.Fibonacci()
	f.Interval = 1 * time.Second
	f.Delay = 3 * time.Second
	instance := &search{
		d: d,
		backoff: f,
		sessionType: sessionType,
		ctx: ctx,
		storage: storage,
		seen: seen,
	}
	instance.run()
	d.cron.Every(30).Seconds().Do(instance.update)
	return instance, nil
}

//Looks for new providers. Ran at the start of a search and on a set interval.
func (s *search) run() error {
	if s.running {
		return nil
	}
	s.mutex.Lock()
	s.running = true
	defer func() {
		s.running = false
		s.mutex.Unlock()
	}()
	logger.Infof("Searching for sessions with type %d", s.sessionType)
	providers := s.d.b.Find(s.d.getMatchmakingKey(s.sessionType), MAX_SESSIONS)
	queried := 0
	for provider := range providers {
		//Stop any duplicate queries and peers that are known to not respond to our app.
		id := provider.String()
		if s.seen.Contains(id) {
			continue
		}
		s.seen.Add(id, true)
		if s.d.ignore.Contains(id) {
			continue
		}
		s.d.sessionQueries <- sessionRequest{
			peer: provider,
			sessionType: s.sessionType,
		}
		queried++
	}
	logger.Infof("Queried %d of the %d for sessions of type %d", queried, len(providers), s.sessionType)
	return nil
}

//Fetches updates from existing providers.
//If we again find sessions, we'll also become a provider.
func (s *search) update() error {
	if s.updating {
		return nil
	}
	s.updating = true
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
		s.updating = false
	}()
	logger.Infof("Updating search for sessions with type %d", s.sessionType)
	peers, err := s.d.GetPeersByDetails("sessionProvider", "1")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not update session search: %s", err.Error()))
	}
	for _, peer := range peers {
		provider, err := s.d.decodePeerId(peer.PId)
		if err != nil {
			logger.Warningf("Failed to decode peer id %s: %v", peer.PId, err)
			continue
		}
		s.d.sessionQueries <- sessionRequest{
			peer: provider,
			sessionType: s.sessionType,
		}
	}
	//Become a provider.
	s.d.b.Provide(s.d.getMatchmakingKey(s.sessionType))
	logger.Infof("Finished updating sessions of type %d", s.sessionType)
	return nil
}

func (s *search) add(sessions []*pb.Session, from Peer.ID) error {
	logger.Infof("Received sessions %d from %s", len(sessions), from.Pretty())
	for _, session := range sessions {
		sessionId, err := s.storage.Insert(session)
		if err != nil {
			return err
		}
		s.d.sessionIdToSessionType[sessionId] = s.sessionType
	}
	if len(sessions) > 0 {
		go func() {
			//Add this session provider to our address book. So we can fetch updates from him and quickly get sessions again from him.
			peer, _ := s.d.FindByPeerId(from.Pretty())
			if peer != nil {
				peer.Details["sessionProvider"] = "1"
				logger.Infof("Added %s to our address book as a session provider", peer.PId)
				s.d.peers.Upsert(peer)
			}
		}()
	}
	return nil
}

func (s *search) fetch() *sessionstore.Store {
	s.backoff.Retry(s.update)
	return s.storage
}