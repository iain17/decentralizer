package sessionstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/ChrisLundquist/golang-lru"
	"time"
	"errors"
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

func (s *Store) Insert(info *pb.SessionInfo) (uint64, error) {
	info.SessionId = GetId(info)
	txn := s.db.Txn(true)
	defer txn.Commit()
	err := txn.Insert(TABLE, info)
	if err == nil{
		s.sessionIds.Add(info.SessionId, true)
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

func (s *Store) FindAll() (result []*pb.SessionInfo, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "details")
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.SessionInfo); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByDetails(key, value string) (result []*pb.SessionInfo, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "details", key, value)
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.SessionInfo); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByPeerId(peerId string) (result []*pb.SessionInfo, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "peerId", peerId)
	if err != nil {
		return nil, err
	}
	for {
		if session, ok := p.Next().(*pb.SessionInfo); ok {
			result = append(result, session)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindSessionId(sessionId uint64) (*pb.SessionInfo, error) {
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
	if session, ok := record.(*pb.SessionInfo); ok {
		return session, nil
	}
	return nil, err
}