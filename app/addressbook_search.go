package app

import (
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"strconv"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/decentralizer/app/peerstore"
	"github.com/iain17/stime"
	"fmt"
	"github.com/iain17/decentralizer/vars"
)

func getDecentralizedIdKey(decentralizedId uint64) string {
	return strconv.FormatUint(decentralizedId, 10)
}

func reverseDecentralizedIdKey(key string) (uint64, error) {
	value, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

//Fetches a new peer from the network. If there isn't a new record it'll return nil
func (d *Decentralizer) getPeerFromNetwork(decentralizedId uint64) (*pb.Peer, error) {
	//Never query self from network
	if d.peers.Self.DnId == decentralizedId {
		return d.peers.Self, nil
	}
	var err error
	networkRecord, err := d.getPeerFromDHT(decentralizedId)
	if err != nil {
		logger.Warning(err)
	}
	if networkRecord == nil {
		networkRecord, err = d.getPeerFromNetworkBackup(decentralizedId)
	}
	if networkRecord == nil {
		return nil, errors.New("could not find peer in the network")
	}
	//Check if existing db record is older. if so update it.
	dbRecord, _ := d.peers.FindByDecentralizedId(decentralizedId)
	if dbRecord == nil {
		d.peers.Insert(networkRecord)
	} else {
		if utils.IsNewerRecord(dbRecord.Published, networkRecord.Published) {
			dbRecord.Details = networkRecord.Details
		}
	}
	return networkRecord, err
}

func (d *Decentralizer) getPeerFromDHT(decentralizedId uint64) (*pb.Peer, error) {
	logger.Infof("Querying DHT network for peer %d", decentralizedId)
	data, err := d.b.GetValue(d.i.Context(), vars.DHT_PEER_KEY_TYPE, getDecentralizedIdKey(decentralizedId))
	if err != nil {
		return nil, fmt.Errorf("could not find peer with id %d: %s", decentralizedId, err.Error())
	}
	var record pb.DNPeerRecord
	err = d.unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return record.Peer, err
}

//Using providers we try and figure it out.
func (d *Decentralizer) getPeerFromNetworkBackup(decentralizedId uint64) (*pb.Peer, error) {
	logger.Infof("Querying backup network for peer %d", decentralizedId)
	values := d.b.Find(getDecentralizedIdKey(decentralizedId), 1024)
	seen := make(map[string]bool)
	for value := range values {
		id := value.ID.Pretty()
		if seen[id] {
			continue
		}
		seen[id] = true
		_, possibleId := peerstore.PeerToDnId(value.ID)
		if possibleId == decentralizedId {
			logger.Infof("Resolved using backup %d == %s", decentralizedId, id)
			return &pb.Peer{
				Published: uint64(stime.Now().Unix()),
				PId: id,
				DnId: decentralizedId,
				Details: map[string]string{
					"backup": "true",
				},
			}, nil
		}
	}
	return nil, errors.New("could not find peer in network. Even with trying the backup method")
}