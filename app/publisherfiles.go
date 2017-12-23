package app

import (
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"github.com/golang/protobuf/proto"
	"errors"
	"time"
)

func (d *Decentralizer) initPublisherFiles() {
	d.downloadPublisherDefinition()
	go d.updatePublisherDefinition()
	d.cron.AddFunc("* 30 * * * *", d.updatePublisherDefinition)
}

func (d *Decentralizer) downloadPublisherDefinition() {
	data, err := configPath.QueryCacheFolder().ReadFile(PUBLISHER_DEFINITION_FILE)
	if err != nil {
		logger.Warningf("Could not restore address book: %v", err)
		return
	}
	var definition pb.PublisherUpdate
	err = proto.Unmarshal(data, &definition)
	if err != nil {
		logger.Warningf("Could not restore address book: %v", err)
		return
	}
	d.loadNewPublisherUpdate(&definition)
}

func (d *Decentralizer) savePublisherUpdate() {
	if d.publisherUpdate.Created == 0 {
		logger.Warning("could not save publisher definition because it wasn't initialized")
		return
	}
	data, err := proto.Marshal(d.publisherUpdate)
	if err == nil {
		logger.Error(err)
		return
	}
	err = configPath.QueryCacheFolder().WriteFile(PUBLISHER_DEFINITION_FILE, data)
	if err == nil {
		logger.Error(err)
	}
}

func (d *Decentralizer) loadNewPublisherUpdate(update *pb.PublisherUpdate) error {
	if d.publisherUpdate != nil && d.publisherUpdate.Created >= update.Created {
		return errors.New("definition is older")
	}

	data, err := proto.Marshal(update.Definition)
	if err != nil {
		return err
	}
	err = d.n.Verify(data, update.Signature)
	if err != nil {
		return err
	}

	d.publisherUpdate = update
	d.savePublisherUpdate()
	return nil
}

//Anything here we wanna do.
//Called when the publisher file has been loaded
func (d *Decentralizer) runPublisherInstructions() {
	//If the publisher file told us to stop. Stop!
	if !d.publisherUpdate.Definition.Status {
		panic("Publisher has closed this network!")
		return
	}
	logger.Infof("Publisher instructions loaded: %s", time.Unix(d.publisherUpdate.Created, 0).Format(time.RFC822))
}

//Will go through each connected peer and try and connect. Find out what publisher update they are running.
//If we've got 3 responses we'll stop trying. And take that this is the latest
func (d *Decentralizer) updatePublisherDefinition() {
	//Wait until we've got enough connections.
	for {
		if len(d.i.PeerHost.Network().Peers()) < 3 {
			time.Sleep(1 * time.Second)
		}
	}
	checked := 0
	for _, peer := range d.i.PeerHost.Network().Peers() {
		if checked >= 3 {
			break
		}
		definition, err := d.getPublisherUpdate(peer)
		if err != nil {
			logger.Debugf("Could not get publisher update: %v", err)
			continue
		}
		checked++
		err = d.loadNewPublisherUpdate(definition)
		//We've updated.
		if err == nil {
			break
		}
	}
}

func (d *Decentralizer) getPublisherUpdate(peer peer.ID) (*pb.PublisherUpdate, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_PUBLISHER_UPDATE)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(300 * time.Millisecond))
	defer stream.Close()

	//Request
	reqData, err := proto.Marshal(&pb.DNPublisherUpdateRequest{})
	if err != nil {
		return nil, err
	}
	err = Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := Read(stream)
	if err != nil {
		return nil, err
	}
	var response pb.DNPublisherUpdateResponse
	err = proto.Unmarshal(resData, &response)
	if err != nil {
		return nil, err
	}
	return response.Update, nil
}