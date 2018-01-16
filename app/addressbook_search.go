package app

import (
	"errors"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"strconv"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/decentralizer/app/peerstore"
	"time"
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
	self, _ := d.peers.FindByPeerId("self")
	if self.DnId == decentralizedId {
		return self, nil
	}

	logger.Infof("Querying network for peer %d", decentralizedId)
	data, err := d.b.GetValue(d.i.Context(), DHT_PEER_KEY_TYPE, getDecentralizedIdKey(decentralizedId))
	if err != nil {
		logger.Warningf("Could not find peer with id %d: %s", err.Error(), decentralizedId)
		return nil, err
	}
	result, _ := d.peers.FindByDecentralizedId(decentralizedId)
	updated := false
	var record pb.DNPeerRecord
	err = d.unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	if result == nil || utils.IsNewerRecord(result.Published, record.Peer.Published) {
		updated = true
		result = record.Peer
	}
	if result == nil {
		return nil, errors.New("could not find peer in the network")
	}
	if updated {
		err = d.peers.Upsert(result)
	}
	if result == nil {
		result, err = d.getPeerFromNetworkBackup(decentralizedId)
	}
	return result, err
}

//Using providers we try and figure it out.
func (d *Decentralizer) getPeerFromNetworkBackup(decentralizedId uint64) (*pb.Peer, error) {
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
			logger.Infof("Resolved %d == %s", id)
			return &pb.Peer{
				Published: uint64(time.Now().UTC().Unix()),
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