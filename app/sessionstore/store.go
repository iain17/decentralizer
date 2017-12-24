package sessionstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/Akagi201/kvcache/ttlru"
	libp2pPeer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"time"
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/app/peerstore"
)

type Store struct {
	self       libp2pPeer.ID
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

func New(size int, expireAt time.Duration, self libp2pPeer.ID) (*Store, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	instance := &Store{
		self:     self,
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

func (s *Store) decodeId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.self, nil
	}
	return libp2pPeer.IDB58Decode(id)
}

func (s *Store) Insert(info *pb.Session) (uint64, error) {
	logger.Infof("Inserting session: %v", info)
	info.SessionId = GetId(info)
	pId, err := s.decodeId(info.PId)
	if err != nil {
		return 0, err
	}
	info.PId, info.DnId = peerstore.PeerToDnId(pId)
	txn := s.db.Txn(true)
	defer txn.Commit()
	err = txn.Insert(TABLE, info)
	if err == nil && info.PId != s.self.Pretty() {
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
	logger.Infof("Find sessions by '%s' = '%s'", key, value)
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