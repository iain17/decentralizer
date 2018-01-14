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
)

func (d *Decentralizer) initAddressbook() {
	var err error
	d.peers, err = peerstore.New(d.ctx, MAX_CONTACTS, time.Duration((EXPIRE_TIME_CONTACT*1.5)*time.Second), d.i.Identity)
	if err != nil {
		logger.Fatal(err)
	}
	d.downloadPeers()
	d.saveSelf()
	d.cron.Every(30).Seconds().Do(d.uploadPeers)
	d.cron.Every(5).Minutes().Do(func() {
		if !d.IsEnoughPeers() {
			return
		}
		d.advertisePeerRecord()
	})

	d.b.RegisterValidator(DHT_PEER_KEY_TYPE, func(rawKey string, val []byte) error {
		var record pb.DNPeerRecord
		err = gogoProto.Unmarshal(val, &record)
		if err != nil {
			return fmt.Errorf("record invalid. could not unmarshal: %s", err.Error())
		}
		//Check key
		key, err := d.b.DecodeKey(rawKey)
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
	}, true)

	d.b.RegisterSelector(DHT_PEER_KEY_TYPE, func(key string, values [][]byte) (int, error) {
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
	})
}

func (d *Decentralizer) downloadPeers() {
	data, err := Base.ReadFile(ADDRESS_BOOK_FILE)
	if err != nil {
		//logger.Warningf("Could not restore address book: %v", err)
		return
	}
	var addressbook pb.DNAddressbook
	err = gogoProto.Unmarshal(data, &addressbook)
	if err != nil {
		logger.Warningf("Could not restore address book: %v", err)
		return
	}
	for _, peer := range addressbook.Peers {
		err = d.peers.Upsert(peer)
		if err != nil {
			logger.Warningf("Error saving peer: %s", peer.PId)
			continue
		}
	}
	logger.Info("Restored address book")
}

func (d *Decentralizer) advertisePeerRecord() error {
	d.WaitTilEnoughPeers()
	peer, err := d.FindByPeerId("self")
	if err != nil {
		logger.Warningf("Could not provide self: %v", err)
		return err
	}
	data, err := gogoProto.Marshal(&pb.DNPeerRecord{
		Peer: peer,
	})
	if err != nil {
		return err
	}
	err = d.b.PutValue(DHT_PEER_KEY_TYPE, getDecentralizedIdKey(peer.DnId), data)
	if err != nil {
		logger.Warning(err)
	} else {
		logger.Info("Successfully provided self")
	}
	return err
}

func (d *Decentralizer) uploadPeers() {
	if !d.addressBookChanged {
		return
	}
	peers, err := d.peers.FindAll()
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	data, err := gogoProto.Marshal(&pb.DNAddressbook{
		Peers: peers,
	})
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	err = Base.WriteFile(ADDRESS_BOOK_FILE, data)
	if err != nil {
		logger.Warningf("Could not save address book: %v", err)
		return
	}
	d.addressBookChanged = false
	logger.Info("Saved address book")
}

//Save our self at least in the address book.
func (d *Decentralizer) saveSelf() error {
	self, err := d.peers.FindByPeerId("self")
	var details map[string]string
	if err != nil {
		details = map[string]string{}
	} else {
		details = self.Details
	}
	if details["name"] == "" {
		details["name"] = randomdata.SillyName()
	}
	info, err := ipinfo.GetIpInfo()
	if err != nil {
		logger.Warningf("Could not find ip info for our session: %s", err)
	}
	if info != nil {
		details["country"] = info.CountryCode
		details["ip"] = info.Ip
	}

	//Add self
	go func() {
		err = d.UpsertPeer("self", details)
		if err != nil {
			logger.Warningf("Could no save self: %s", err.Error())
		}
	}()
	d.uploadPeers()
	return nil
}

func (d *Decentralizer) UpsertPeer(pId string, details map[string]string) error {
	err := d.peers.Upsert(&pb.Peer{
		Published: uint64(time.Now().UTC().Unix()),
		PId:     pId,
		Details: details,
	})
	d.addressBookChanged = true
	if pId == "self" {
		err = d.advertisePeerRecord()
	}
	return err
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