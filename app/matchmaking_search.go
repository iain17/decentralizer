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
	"gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
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
	instance.fetch()
	d.cron.Every(30).Seconds().Do(instance.fetch)
	return instance, nil
}

func (s *search) run(ctx context.Context) error {
	logger.Info("Trying to lock mutex")
	s.mutex.Lock()

	defer s.mutex.Unlock()

	logger.Infof("Searching for sessions with type %d", s.sessionType)
	values, err := s.d.b.GetValues(ctx, DHT_SESSIONS_KEY_TYPE, s.d.getMatchmakingKey(s.sessionType), 1024)
	if err != nil {
		return err
	}
	queried := 0
	total := 0
	s.seen.Purge()
	for _, value := range values {
		total++
		id := value.From.Pretty()
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
			s.d.setSessionIdToType(session.SessionId, session.Type)
			s.storage.Insert(session)
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
	return s.storage, err
}