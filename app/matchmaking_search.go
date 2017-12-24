package app

import (
	"context"
	"github.com/iain17/decentralizer/app/sessionstore"
	"sync"
	"github.com/iain17/logger"
	Peer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/timeout"
	"time"
)

type search struct {
	running bool
	updating bool
	mutex sync.Mutex
	d *Decentralizer
	sessionType uint64
	ctx context.Context
	storage *sessionstore.Store
	seen map[string]bool
}

func (d *Decentralizer) newSearch(ctx context.Context, sessionType uint64) (*search, error) {
	storage := d.getSessionStorage(sessionType)
	instance := &search{
		d: d,
		sessionType: sessionType,
		ctx: ctx,
		storage: storage,
		seen: make(map[string]bool),
	}
	timeout.Do(instance.run, 10*time.Second)//Initial search will last 15 seconds max. 10 here + 5 fetch.
	d.cron.AddFunc("30 * * * * *", func() {
		//Just wrap another function around it. because we can't check if the mutex is blocked and if the run takes more than 30 seconds, it'll slowly build up the cron go routines.
		if instance.running {
			return
		}
		instance.running = true
		timeout.Do(instance.run, EXPIRE_TIME_SESSION*time.Second)
		instance.running = false
	})
	return instance, nil
}

//Looks for new providers. Ran at the start of a search and on a set interval.
func (s *search) run(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	logger.Infof("Searching for sessions with type %d", s.sessionType)
	var wg sync.WaitGroup
	providers := s.d.b.Find(s.d.getMatchmakingKey(s.sessionType), MAX_SESSIONS)
	for provider := range providers {
		//Stop any duplicate queries and peers that are known to not respond to our app.
		id := provider.String()
		if s.d.ignore[id] {
			continue
		}
		if s.seen[id] {
			continue
		}
		s.seen[id] = true

		wg.Add(1)
		go func(id Peer.ID) {
			defer wg.Done()
			_, err := s.request(id)
			if err != nil {
				logger.Infof("Failed to save sessions received from %s: %v", id.Pretty(), err)
			}
		}(provider)
	}
	wg.Wait()
	logger.Infof("Search complete for sessions with type %d", s.sessionType)
}

//Fetches updates from existing providers.
//If we again find sessions, we'll also become a provider.
func (s *search) update(ctx context.Context) {
	if s.updating {
		return
	}
	s.updating = true
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
		s.updating = false
	}()
	logger.Infof("Updating search for sessions with type %d", s.sessionType)
	var wg sync.WaitGroup
	peers, err := s.d.GetPeersByDetails("sessionProvider", "1")
	if err != nil {
		logger.Warningf("Could not update session search: %v", err)
		return
	}
	total := 0
	for _, peer := range peers {
		provider, err := s.d.decodePeerId(peer.PId)
		if err != nil {
			logger.Warningf("Failed to decode peer id %s: %v", peer.PId, err)
			continue
		}
		wg.Add(1)
		go func(id Peer.ID) {
			defer wg.Done()
			num, err := s.request(id)
			if err != nil {
				if err.Error() == "i/o deadline reached" {
					return
				}
				if err.Error() == "protocol not supported" {
					s.d.ignore[id.Pretty()] = true
					return
				}
				logger.Debugf("Failed to get sessions from %s: %v", id.Pretty(), err)
			}
			total += num
		}(provider)
	}
	wg.Wait()
	//Become a provider.
	if total > 0 {
		s.d.b.Provide(s.d.getMatchmakingKey(s.sessionType))
	}
	logger.Infof("Updated %d sessions of type %d", total, s.sessionType)
}

func (s *search) request(id Peer.ID) (int, error) {
	sessions, err := s.d.getSessionsRequest(id, s.sessionType)
	if err != nil {
		return 0, err
	}
	return len(sessions), s.add(sessions, id)
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
				logger.Infof("Added %s to our addressbook as a session provider", peer.PId)
				s.d.peers.Upsert(peer)
			}
		}()
	}
	return nil
}

func (s *search) fetch() *sessionstore.Store {
	timeout.Do(s.update, 5*time.Second)
	return s.storage
}