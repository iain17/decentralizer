package peerstore

import (
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/ChrisLundquist/golang-lru"
	"time"
	"errors"
	"github.com/iain17/logger"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

type Store struct {
	self libp2pPeer.ID
	db *memdb.MemDB
	sessionIds *lru.LruWithTTL
	expireAt time.Duration
}

const TABLE = "peers"
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		TABLE: {
			Name: TABLE,
			Indexes: map[string]*memdb.IndexSchema{
				"id": {
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "PId"},
				},
				"DnId":{
					Name:    "DnId",
					Unique:  true,
					Indexer: &memdb.UintFieldIndex{Field: "DnId"},
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
		self: self,
		db: db,
		expireAt: expireAt,
	}
	instance.sessionIds, err = lru.NewTTLWithEvict(size, instance.onEvicted)
	return instance, err
}

func (s *Store) decodeId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.self, nil
	}
	return libp2pPeer.IDB58Decode(id)
}

func (s *Store) onEvicted(key interface{}, value interface{}) {
	//In go routine. Because of deadlock caused when evicted is called, the insert function higher up the call chain has just locked the mutex.
	go func() {
		if peerId, ok := key.(string); ok {
			err := s.Remove(peerId)
			if err != nil {
				logger.Warningf("Could not delete peer id: %s", peerId)
			} else {
				logger.Infof("Deleted peer id: %s", peerId)
			}
		} else {
			logger.Warningf("Could not delete peer id...", peerId)
		}
	}()
}

func (s *Store) Upsert(info *pb.Peer) error {
	peerId, err := s.decodeId(info.PId)
	if err != nil {
		return err
	}
	info.PId, info.DnId = PeerToDnId(peerId)
	txn := s.db.Txn(true)
	defer txn.Commit()
	err = txn.Insert(TABLE, info)
	if err == nil && info.PId != s.self.Pretty() {
		s.sessionIds.AddWithTTL(info.PId, true, s.expireAt)
	}
	return err
}

func (s *Store) Remove(peerId string) error {
	id, err := s.decodeId(peerId)
	if err != nil {
		return err
	}
	txn := s.db.Txn(true)
	defer txn.Commit()
	_, err = txn.DeleteAll(TABLE, "id", id.Pretty())
	s.sessionIds.Remove(id.Pretty())
	return err
}

func (s *Store) FindAll() (result []*pb.Peer, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "id")
	if err != nil {
		return nil, err
	}
	for {
		if peer, ok := p.Next().(*pb.Peer); ok {
			result = append(result, peer)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByDetails(key, value string) (result []*pb.Peer, err error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "details", key, value)
	if err != nil {
		return nil, err
	}
	for {
		if peer, ok := p.Next().(*pb.Peer); ok {
			result = append(result, peer)
		} else {
			break
		}
	}
	return
}

func (s *Store) FindByDecentralizedId(decentralizedId uint64) (*pb.Peer, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "DnId", decentralizedId)
	if err != nil {
		return nil, err
	}
	record := p.Next()
	if record == nil {
		return nil, errors.New("could not find peer")
	}
	if peer, ok := record.(*pb.Peer); ok {
		return peer, nil
	}
	return nil, err
}

func (s *Store) FindByPeerId(peerId string) (*pb.Peer, error) {
	id, err := s.decodeId(peerId)
	if err != nil {
		return nil, err
	}
	txn := s.db.Txn(false)
	defer txn.Abort()
	p, err := txn.Get(TABLE, "id", id.Pretty())
	if err != nil {
		return nil, err
	}
	record := p.Next()
	if record == nil {
		return nil, errors.New("could not find peer")
	}
	if peer, ok := record.(*pb.Peer); ok {
		return peer, nil
	}
	return nil, err
}