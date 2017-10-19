package sessionstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/app/pb"
	"strconv"
	"github.com/ChrisLundquist/golang-lru"
	"time"
)

type Store struct {
	db *memdb.MemDB
	sessionIds *lru.LruWithTTL
	expireAt time.Duration
}

type Session struct {
	*pb.SessionInfo
	SessionId string
}

const TABLE = "sessions"
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		TABLE: {
			Name: TABLE,
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "SessionId"},
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

func (s *Store) Insert(info *pb.SessionInfo) (uint64, error) {
	sessionId := GetId(info)
	txn := s.db.Txn(true)
	defer txn.Commit()
	err := txn.Insert(TABLE, &Session{
		SessionId: strconv.FormatInt(int64(sessionId), 16),//Because memdb wants a string
		SessionInfo: info,
	})
	if err == nil{
		s.sessionIds.Add(sessionId, true)
	}
	return sessionId, err
}

func (s *Store) Remove(sessionId uint64) error {
	txn := s.db.Txn(true)
	defer txn.Commit()
	_, err := txn.DeleteAll(TABLE, "SessionId", strconv.FormatInt(int64(sessionId), 16))//Because memdb wants a string
	s.sessionIds.Remove(sessionId)
	return err
}

func (s *Store) FindByDetails(index string, key, value string) (result []*Session, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, index, key, value)
	for {
		if session, ok := p.Next().(*Session); ok {
			result = append(result, session)
		}
	}
	return
}

func (s *Store) FindByPeerId(peerId string) (result []*pb.SessionInfo, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "peerId", peerId)
	for {
		if session, ok := p.Next().(*Session); ok {
			result = append(result, session.SessionInfo)
		}
	}
	return
}