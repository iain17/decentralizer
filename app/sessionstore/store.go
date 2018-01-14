package sessionstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/kvcache/ttlru"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"time"
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/app/peerstore"
	"context"
	"github.com/iain17/decentralizer/utils"
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

func New(ctx context.Context, size int, expireAt time.Duration, self libp2pPeer.ID) (*Store, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	instance := &Store{
		self:     self,
		db: db,
		expireAt: expireAt,
	}
	instance.sessionIds, err = lru.NewTTLWithEvict(ctx, size, instance.onEvicted)
	return instance, err
}

func (s *Store) onEvicted(key interface{}, value interface{}) {
	//In go routine. Because of deadlock caused when evicted is called, the insert function higher up the call chain has just locked the mutex.
	go func() {
		if isExternalSession, ok := value.(bool); ok {
			if sSessionId, ok := key.(uint64); ok {
				if isExternalSession {
					err := s.Remove(sSessionId)
					if err != nil {
						logger.Warningf("Could not delete session id: %d", sSessionId)
					} else {
						logger.Infof("Deleted session id: %d", sSessionId)
					}
				} else {
					s.sessionIds.Add(sSessionId, false)//Add it back in cuz we can't delete our sessions.
				}
			} else {
				logger.Warningf("Could not delete session id...", sSessionId)
			}
		}
	}()
}

func (s *Store) decodeId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.self, nil
	}
	return libp2pPeer.IDB58Decode(id)
}

func (s *Store) Contains(sessionId uint64) bool {
	return s.sessionIds.Contains(sessionId)
}

func (s *Store) IsExternalSession(sessionId uint64) bool {
	if value, ok := s.sessionIds.Get(sessionId); ok {
		return value.(bool)
	}
	return false
}

func (s *Store) IsNewer(info *pb.Session) bool {
	existingSession, _ := s.FindSessionId(info.SessionId)
	var published uint64
	if existingSession != nil {
		published = existingSession.Published
	}
	return utils.IsNewerRecord(published, info.Published)
}

//Do not call this directly. setSessionIdToType needs to be called as well!
func (s *Store) Insert(info *pb.Session) (uint64, error) {
	if info == nil {
		return 0, errors.New("tried inserting nil value")
	}
	if !s.IsNewer(info) {
		return 0, errors.New("record is older")
	}
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
	if err == nil {
		if info.PId != s.self.Pretty() {
			s.sessionIds.AddWithTTL(info.SessionId, true, s.expireAt)
		} else {
			s.sessionIds.Add(info.SessionId, false)
		}
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

func (s *Store) SessionIds() []interface{} {
	return s.sessionIds.Keys()
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

func (s *Store) Len() int {
	return s.sessionIds.Len()
}

func (s *Store) IsEmpty() bool {
	return s.sessionIds.Len() == 0
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