package peerstore

import (
	"errors"
	"github.com/iain17/kvcache/ttlru"
	"github.com/hashicorp/go-memdb"
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"time"
	"context"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"sync"
)

type Store struct {
	selfId 		libp2pPeer.ID
	Self		*pb.Peer
	SelfMutex	sync.RWMutex

	db         *memdb.MemDB
	sessionIds *lru.LruWithTTL
	expireAt   time.Duration
	Changed    bool
	path       string
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
				"DnId": {
					Name:    "DnId",
					Unique:  true,
					Indexer: &memdb.UintFieldIndex{Field: "DnId"},
				},
				"details": {
					Name:    "details",
					Unique:  false,
					Indexer: &memdb.StringMapFieldIndex{Field: "Details"},
				},
			},
		},
	},
}

func New(ctx context.Context, size int, expireAt time.Duration, selfId libp2pPeer.ID, path string) (*Store, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	instance := &Store{
		selfId:   selfId,
		db:       db,
		expireAt: expireAt,
		path: path,
	}
	instance.sessionIds, err = lru.NewTTLWithEvict(ctx, size, instance.onEvicted)
	instance.restore()
	return instance, err
}

func (s *Store) restore() {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		logger.Warningf("Could not restore peer store: %s", err.Error())
		return
	}
	var addressbook pb.DNAddressbook
	err = proto.Unmarshal(data, &addressbook)
	if err != nil {
		logger.Warningf("Could not restore peer store: %v", err)
		return
	}
	for _, peer := range addressbook.Peers {
		err = s.Insert(peer)
		if err != nil {
			logger.Warningf("Error saving peer %s: %s", peer.PId, err.Error())
			continue
		}
	}
	logger.Info("Restored peer store")
}

func (s *Store) Save() {
	if !s.Changed {
		return
	}
	s.Changed = false
	if s.path == "" {
		return
	}
	peers, err := s.FindAll()
	if err != nil {
		logger.Warningf("Could not save peer store: %v", err)
		return
	}
	data, err := proto.Marshal(&pb.DNAddressbook{
		Peers: peers,
	})
	if err != nil {
		logger.Warningf("Could not save peer store: %v", err)
		return
	}
	err = ioutil.WriteFile(s.path, data, 0777)
	if err != nil {
		logger.Warningf("Could not save peer store: %v", err)
		return
	}
	logger.Info("Saved peer store")
}

func (s *Store) decodeId(id string) (libp2pPeer.ID, error) {
	if id == "self" {
		return s.selfId, nil
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
/**
Inserts a record. Takes a pointer peer. This pointer value is directly saved, meaning changes to this object will also change the value in the db.
 */
func (s *Store) Insert(info *pb.Peer) error {
	if info == nil {
		return errors.New("peer info not defined")
	}
	peerId, err := s.decodeId(info.PId)
	if err != nil {
		return err
	}
	info.PId, info.DnId = PeerToDnId(peerId)
	//TODO: Remove later
	existingPeer, _ := s.FindByPeerId(info.PId)
	if existingPeer != nil {
		panic("Trying to insert duplicate item...")
	}
	txn := s.db.Txn(true)
	defer txn.Commit()
	err = txn.Insert(TABLE, info)
	if err == nil && info.PId != s.selfId.Pretty() {
		s.sessionIds.AddWithTTL(info.PId, true, s.expireAt)
	}
	s.Changed = true
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
	if s.Self != nil && peerId == "self" {
		return s.Self, nil
	}
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
