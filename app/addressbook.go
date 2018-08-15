package app

import (
	"github.com/iain17/decentralizer/pb"
	"github.com/iain17/logger"
	"github.com/Pallinder/go-randomdata"
	"time"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/ipinfo"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"fmt"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/stime"
)

func (d *Decentralizer) initAddressbook() {
	var err error
	d.peers, err = peerstore.New(d.ctx, MAX_CONTACTS, time.Duration((EXPIRE_TIME_CONTACT*1.5)*time.Second), d.i.Identity, Base.Path+"/"+ADDRESS_BOOK_FILE)
	if err != nil {
		logger.Fatal(err)
	}
	d.saveSelf()
	d.cron.Every(30).Seconds().Do(d.peers.Save)
	d.cron.Every(5).Minutes().Do(func() {
		if !d.IsEnoughPeers() {
			return
		}
		d.AdvertisePeerRecord()
	})
	go d.AdvertisePeerRecord()

	d.b.RegisterValidator(DHT_PEER_KEY_TYPE, func(key string, value []byte) error {
		var record pb.DNPeerRecord
		err = d.unmarshal(value, &record)
		if err != nil {
			return fmt.Errorf("peer record invalid. could not unmarshal: %s", err.Error())
		}
		//Check key
		key, err := d.b.DecodeKey(key)
		if err != nil {
			return err
		}
		expectedDecentralizedId, err := reverseDecentralizedIdKey(key)
		if err != nil {
			return fmt.Errorf("failed to reverse '%s' to decentralized id: %s", key, err.Error())
		}
		if expectedDecentralizedId != record.Peer.DnId {
			return fmt.Errorf("reversing decentralized key id failed. Expected %d, received %d", expectedDecentralizedId, record.Peer.DnId)
		}
		return nil
	}, func(key string, values [][]byte) (int, error) {
		var currPeer pb.Peer
		best := 0
		for i, val := range values {
			var record pb.DNPeerRecord
			err = d.unmarshal(val, &record)
			if err != nil {
				logger.Warning(err)
				continue
			}
			if utils.IsNewerRecord(currPeer.Published, record.Peer.Published) {
				currPeer = *record.Peer
				best = i
			}
		}
		return best, nil
	}, true)

	self, _ := d.peers.FindByPeerId("self")
	if self != nil {
		logger.Infof("Initialized: PeerID '%s', decentralized id '%d': %v", self.PId, self.DnId, self.Details)
	}
}

func (d *Decentralizer) AdvertisePeerRecord() error {
	d.WaitTilEnoughPeers()
	d.peers.SelfMutex.RLock()
	defer d.peers.SelfMutex.RUnlock()

	data, err := gogoProto.Marshal(&pb.DNPeerRecord{
		Peer: d.peers.Self,
	})
	if err != nil {
		return err
	}
	err = d.b.PutValue(DHT_PEER_KEY_TYPE, getDecentralizedIdKey(d.peers.Self.DnId), data)
	if err != nil {
		logger.Warningf("Could not provide self: %s", err)
	} else {
		logger.Info("Successfully provided self")
	}
	//Back-up for if it fails.
	err = d.b.Provide(getDecentralizedIdKey(d.peers.Self.DnId))
	if err != nil {
		logger.Warning(err)
	} else {
		logger.Info("Successfully provided self")
	}
	return err
}

//Save our self at least in the address book.
func (d *Decentralizer) saveSelf() error {
	d.peers.SelfMutex.Lock()
	defer d.peers.SelfMutex.Unlock()

	var err error
	d.peers.Self, err = d.peers.FindByPeerId("self")
	if err != nil {
		d.peers.Self, err = d.InsertPeer("self", map[string]string{
			"name": randomdata.SillyName(),
		})
		if err != nil {
			return fmt.Errorf("could no save self: %s", err.Error())
		}
	}
	info, err := ipinfo.GetIpInfo()
	if err != nil {
		logger.Warningf("Could not find ip info for our session: %s", err)
	}
	if info != nil {
		d.peers.Self.Details["country"] = info.CountryCode
		d.peers.Self.Details["ip"] = info.Ip
	}
	d.peers.Changed = true
	d.peers.Save()
	return nil
}

func (d *Decentralizer) UpsertPeer(pId string, details map[string]string) (peer *pb.Peer, err error) {
	peer, err = d.peers.FindByPeerId(pId)
	isSelf := peer == d.peers.Self
	if peer != nil {
		peer.Details = details
		d.peers.Changed = true
		d.peers.Save()
	} else {
		peer, err = d.InsertPeer(pId, details)
	}
	if isSelf {
		d.AdvertisePeerRecord()
	}
	return peer, err
}

func (d *Decentralizer) InsertPeer(pId string, details map[string]string) (*pb.Peer, error) {
	peer := &pb.Peer{
		Published: uint64(stime.Now().Unix()),
		PId:     pId,
		Details: details,
	}
	return peer, d.peers.Insert(peer)
}

func (d *Decentralizer) GetPeersByDetails(key, value string) ([]*pb.Peer, error) {
	return d.peers.FindByDetails(key, value)
}

func (d *Decentralizer) GetPeers() ([]*pb.Peer, error) {
	return d.peers.FindAll()
}

func (d *Decentralizer) FindByPeerId(peerId string) (p *pb.Peer, err error) {
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return nil, err
	}
	var decentralizedId uint64
	peerId, decentralizedId = peerstore.PeerToDnId(id)
	return d.FindByDecentralizedId(decentralizedId)
}

func (d *Decentralizer) FindByDecentralizedId(decentralizedId uint64) (*pb.Peer, error) {
	//Try and fetch from network
	peer, err := d.getPeerFromNetwork(decentralizedId)
	if err != nil {
		logger.Warningf("Could not fetch peer from network: %s", err.Error())
		peer, err = d.peers.FindByDecentralizedId(decentralizedId)
	}
	return peer, err
}