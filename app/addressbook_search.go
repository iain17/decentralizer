package app

import (
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"strconv"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"github.com/iain17/decentralizer/utils"
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
	self, _ := d.peers.FindByPeerId("self")
	if self.DnId == decentralizedId {
		return self, nil
	}
	logger.Infof("Querying network for peer %d", decentralizedId)
	values, err := d.b.GetValues(d.i.Context(), DHT_PEER_KEY_TYPE, getDecentralizedIdKey(decentralizedId), 1)
	if err != nil {
		logger.Warningf("Could not find peer with id %d: %s", err.Error(), decentralizedId)
		return nil, err
	}
	logger.Infof("Found %d possible values for peer %d", len(values), decentralizedId)
	result, _ := d.peers.FindByDecentralizedId(decentralizedId)
	updated := false
	for _, value := range values {
		var record pb.DNPeerRecord
		err = gogoProto.Unmarshal(value.Val, &record)
		if err != nil {
			logger.Warning(err)
			continue
		}
		if result == nil || utils.IsNewerRecord(result.Published, record.Peer.Published) {
			updated = true
			result = record.Peer
		}
	}
	if result == nil {
		return nil, errors.New("could not find peer in the network")
	}
	if updated {
		err = d.peers.Upsert(result)
	}
	return result, err
}