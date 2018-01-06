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
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
)

type search struct {
	running bool
	mutex sync.Mutex
	d *Decentralizer
	sessionType uint64
	ctx context.Context
	storage *sessionstore.Store
	seen *lttlru.LruWithTTL//Tells us if we've already pased the result of this peer.
	publication *lttlru.LruWithTTL//Tell us if the record we are about to insert is from a newer source.
}

func (d *Decentralizer) newSearch(ctx context.Context, sessionType uint64) (*search, error) {
	storage := d.getSessionStorage(sessionType)
	seen, err := lttlru.NewTTL(MAX_IGNORE)
	if err != nil {
		return nil, err
	}
	publication, err := lttlru.NewTTL(MAX_SESSIONS)
	if err != nil {
		return nil, err
	}
	instance := &search{
		d: d,
		sessionType: sessionType,
		ctx: ctx,
		storage: storage,
		seen: seen,
		publication: publication,
	}
	instance.fetch()
	//d.cron.Every(30).Seconds().Do(instance.fetch)
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
	logger.Infof("Searching for sessions with type %d", s.sessionType)
	values, err := s.d.b.GetValues(DHT_SESSIONS_KEY_TYPE, s.d.getMatchmakingKey(s.sessionType), 512)
	if err != nil {
		return err
	}
	queried := 0
	total := 0
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
		s.seen.AddWithTTL(id, true, 1 * time.Minute)

		var response pb.DNSessions
		err = gogoProto.Unmarshal(value.Val, &response)
		if err != nil {
			logger.Warning(err)
			continue
		}
		logger.Infof("Got %d sessions from %s", len(response.Results), id)
		publishedTime := time.Unix(int64(response.Published), 0)
		for _, session := range response.Results {
			//if session.PId != value.From.Pretty() {
			//	logger.Warningf("We can't accept sessions that aren't yours %s", id)
			//	s.d.ignore.AddWithTTL(id, true, 5 * time.Minute)
			//	break
			//}
			if rawLastInserted, ok := s.publication.Get(session.SessionId); ok {
				LastInserted := time.Unix(rawLastInserted.(int64), 0)
				if publishedTime.Before(LastInserted) || publishedTime.Equal(LastInserted) {
					continue
				}
			}
			s.publication.Add(session.SessionId, time.Now().UTC().Unix())
			_, err := s.storage.Insert(session)
			if err != nil {
				logger.Warning(err)
			}
		}
		queried++
	}
	logger.Infof("Queried %d of the %d for sessions of type %d", queried, total, s.sessionType)
	return nil
}

func (s *search) refresh() {
	err := s.run()
	if err != nil {
		logger.Warning(err)
	}
}

func (s *search) fetch() *sessionstore.Store {
	timeout.Do(func(ctx context.Context) {
		tries := 0
		s.refresh()
		for s.storage.IsEmpty() {
			select {
			case <- ctx.Done():
				return
			default:
				if tries > 5 {
					break
				}
				s.refresh()
				time.Sleep(1 * time.Second)
				tries++
			}
		}
	}, 5 * time.Second)
	return s.storage
}