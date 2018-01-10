package app

import (
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"strconv"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"time"
	"fmt"
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

//Returns true if the second passed peer is newer than the first one.
func (d *Decentralizer) isNewerPeer(existingPeer *pb.Peer, newPeer *pb.Peer) bool {
	if newPeer == nil {
		return false
	}
	if existingPeer == nil && newPeer != nil {
		return true
	}
	now := time.Now().UTC()
	publishedTime := time.Unix(int64(newPeer.Published), 0).UTC()
	publishedTimeText := publishedTime.String()
	expireTime := time.Unix(int64(existingPeer.Published), 0).UTC()
	expireTimeText := expireTime.String()
	if publishedTime.Before(expireTime) {
		err := fmt.Errorf("new peer with publish date %s has expired. It was before %s", publishedTimeText, expireTimeText)
		logger.Warning(err)
		return false
	}
	if publishedTime.After(now) {
		err := errors.New("new peer with publish date %s was published in the future")
		logger.Warning(err)
		return false
	}
	logger.Infof("successfully accepted new peer updated at: %s", publishedTimeText)
	return true
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
		if d.isNewerPeer(result, record.Peer) {
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