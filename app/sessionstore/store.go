package sessionstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/ChrisLundquist/golang-lru"
	"time"
	"errors"
	"github.com/iain17/logger"
)

type Store struct {
	db *memdb.MemDB
	sessionIds *lru.LruWithTTL
	expireAt time.Duration
}

const TABLE = "sessions"
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		//pb.SessionInfo
		TABLE: {
			Name: TABLE,
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.UintFieldIndex{Field: "SessionId"},
				},
				"peerId":{
					Name:    "peerId",
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "PId"},
				},
				"details":{
					Name:    "details",
					Unique:  false,
					Indexer: &memdb.StringMapFieldIndex{Field: "Details"},
				},
			},
		},
	},
}

func New(size int, expireAt time.Duration) (*Store, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	instance := &Store{
		db: db,
		expireAt: expireAt,
	}
	instance.sessionIds, err = lru.NewTTLWithEvict(size, instance.onEvicted)
	return instance, err
}

func (s *Store) onEvicted(key interface{}, value interface{}) {
	//In go routine. Because of deadlock caused when evicted is called, the insert function higher up the call chain has just locked the mutex.
	go func() {
		if sSessionId, ok := key.(uint64); ok {
			err := s.Remove(sSessionId)
			if err != nil {
				logger.Warningf("Could not delete session id: %d", sSessionId)
			} else {
				logger.Infof("Deleted session id: %d", sSessionId)
			}
		} else {
			logger.Warningf("Could not delete session id...", sSessionId)
		}
	}()
}

func (s *Store) Insert(info *pb.Session) (uint64, error) {
	logger.Infof("Inserting session: %v", info)
	info.SessionId = GetId(info)
	txn := s.db.Txn(true)
	defer txn.Commit()
	err := txn.Insert(TABLE, info)
	if err == nil{
		s.sessionIds.AddWithTTL(info.SessionId, true, s.expireAt)
	}
	return info.SessionId, err
}

func (s *Store) Remove(sessionId uint64) error {
	txn := s.db.Txn(true)
	defer txn.Commit()
	_, err := txn.DeleteAll(TABLE, "id", sessionId)
	s.sessionIds.Remove(sessionId)
	return err
}

func (s *Store) FindAll() (result []*pb.Session, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "details")
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.Session); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByDetails(key, value string) (result []*pb.Session, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "details", key, value)
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.Session); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByPeerId(peerId string) (result []*pb.Session, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "peerId", peerId)
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.Session); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindSessionId(sessionId uint64) (*pb.Session, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "id", sessionId)
	if err != nil {
		return nil, err
	}
	record := p.Next()
	if record == nil {
		return nil, errors.New("Could not find session.")
	}
	if session, ok := record.(*pb.Session); ok {
		return session, nil
	}
	return nil, err
}