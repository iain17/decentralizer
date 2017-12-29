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
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"fmt"
)

type search struct {
	running bool
	backoff *backoff.ExponentialBackOff
	mutex sync.Mutex
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
	f := backoff.NewExponentialBackOff()
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
	logger.Debug("Trying to lock mutex")
	s.mutex.Lock()
	if s.running {
		logger.Debug("Search run is already running...")
		logger.Debug("Unlocking mutex")
		s.mutex.Unlock()
		return nil
	}
	s.running = true
	logger.Debug("Unlocking mutex")
	s.mutex.Unlock()

	defer func() {
		s.running = false
	}()
	//Keeps looking until we've found at least 1!
	for {
		select {
		case <-s.d.ctx.Done():
			return nil
		default:
			logger.Infof("Searching for sessions with type %d", s.sessionType)
			providers := s.d.b.Find(s.d.getMatchmakingKey(s.sessionType), 512)
			queried := 0
			total := len(providers)
			for provider := range providers {
				info := s.d.i.Peerstore.PeerInfo(provider)//Fetched because bitwise will only save the addrs temp: pstore.TempAddrTTL
				if len(info.Addrs) == 0 {
					logger.Debug("We forgot already how to find this peer")
					continue
				}

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
					search: s,
					peer: info,
				}
				queried++
			}
			logger.Infof("Queried %d of the %d for sessions of type %d", queried, total, s.sessionType)
			sessions, err := s.storage.FindAll()
			if err != nil {
				logger.Warning(err)
			} else {
				if len(sessions) > 0 {
					return nil//Done
				}
			}
		}
	}
	return nil
}

//Fetches updates from existing providers.
//If we again find sessions, we'll also become a provider.
func (s *search) update() error {
	logger.Debug("Trying to lock mutex")
	s.mutex.Lock()
	if s.running {
		logger.Debug("Search run is already running...")
		logger.Debug("Unlocking mutex")
		s.mutex.Unlock()
		return nil
	}
	s.running = true
	logger.Debug("Unlocking mutex")
	s.mutex.Unlock()

	defer func() {
		s.running = false
	}()
	peers, err := s.d.GetPeersByDetails("sessionProvider", "1")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not update session search: %s", err.Error()))
	}
	if len(peers) == 0 {
		return nil
	}
	logger.Infof("Updating search for sessions with type %d", s.sessionType)
	for _, peer := range peers {
		provider, err := s.d.decodePeerId(peer.PId)
		if err != nil {
			logger.Warningf("Failed to decode peer id %s: %v", peer.PId, err)
			continue
		}
		info := s.d.i.Peerstore.PeerInfo(provider)
		if len(info.Addrs) == 0 {
			logger.Debug("We forgot already how to find this peer")
			continue
		}
		s.d.sessionQueries <- sessionRequest{
			search: s,
			peer: info,
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
	duration := s.backoff.NextBackOff()
	if duration != backoff.Stop {
		s.run()
		time.Sleep(duration)
	}
	return s.storage
}