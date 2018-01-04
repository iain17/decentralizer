package app

import (
	"context"
	"github.com/iain17/decentralizer/app/sessionstore"
	"sync"
	"github.com/iain17/logger"
	Peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	net "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	"github.com/iain17/decentralizer/pb"
	"time"
	"github.com/iain17/kvcache/lttlru"
	"github.com/pkg/errors"
	"fmt"
	"github.com/iain17/timeout"
)

type search struct {
	running bool
	mutex sync.Mutex
	d *Decentralizer
	sessionType uint64
	ctx context.Context
	storage *sessionstore.Store
	seen *lttlru.LruWithTTL
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
	instance.fetch()
	d.cron.Every(30).Seconds().Do(instance.update)
	return instance, nil
}

func (s *search) run() error {
	logger.Info("Trying to lock mutex")
	s.mutex.Lock()
	if s.running {
		logger.Debug("Search run is already running...")
		s.mutex.Unlock()
		return nil
	}
	s.running = true
	logger.Info("Unlocking mutex")
	s.mutex.Unlock()

	defer func() {
		logger.Info("Running to false")
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
			total := 0
			for provider := range providers {
				info := s.d.i.Peerstore.PeerInfo(provider) //Fetched because bitwise will only save the addrs temp: pstore.TempAddrTTL
				total++

				connectedNess := s.d.i.PeerHost.Network().Connectedness(provider)
				if connectedNess == net.CannotConnect {
					logger.Debug("Known to be terrible to connect to")
					continue
				}

				//Check if we've got a way to connect
				if connectedNess == net.NotConnected {
					if len(info.Addrs) == 0 {
						logger.Debugf("We've forgotten already how to find this peer: %s", provider.Pretty())
						continue
					}
				}

				//Stop any duplicate queries and peers that are known to not respond to our app.
				id := provider.String()
				if s.seen.Contains(id) {
					continue
				}
				s.seen.AddWithTTL(id, true, 60 * time.Second)
				if s.d.ignore.Contains(id) {
					continue
				}
				s.d.sessionQueries <- sessionRequest{
					search: s,
					peer: info,
					connected: connectedNess == net.Connected,
				}
				queried++
			}
			logger.Infof("Queried %d of the %d for sessions of type %d", queried, total, s.sessionType)
			if queried == 0 {
				logger.Debug("Nothing left to query.")
				return nil
			}
			logger.Infof("1")
			if !s.storage.IsEmpty() {
				logger.Debug("We found at least one. quitting.")
				return nil
			} else {
				logger.Debug("None found")
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

//Fetches updates from existing providers.
//If we again find sessions, we'll also become a provider.
func (s *search) update() error {
	logger.Info("Trying to lock mutex")
	s.mutex.Lock()
	if s.running {
		logger.Debug("Search run is already running...")
		s.mutex.Unlock()
		return nil
	}
	s.running = true
	logger.Info("Unlocking mutex")
	s.mutex.Unlock()

	defer func() {
		logger.Info("Running to false")
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
	logger.Debug("Fetching")
	go s.run()
	timeout.Do(func(ctx context.Context) {
		for s.storage.IsEmpty() {
			select {
			case <- ctx.Done():
				return
			default:
				logger.Debug("Still empty")
				time.Sleep(1 * time.Second)
			}
		}
	}, 5 * time.Second)
	logger.Debug("Returning")
	return s.storage
}