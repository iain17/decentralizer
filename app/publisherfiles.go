package app

import (
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	inet "gx/ipfs/QmU4vCDZTPLDqSDKguWbHCiUe46mZUtmM2g2suBZ9NE8ko/go-libp2p-net"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"github.com/golang/protobuf/proto"
	"errors"
	"time"
)

func (d *Decentralizer) initPublisherFiles() {
	d.i.PeerHost.SetStreamHandler(GET_PUBLISHER_UPDATE, d.getPublisherUpdateResponse)
	d.downloadPublisherDefinition()
	go d.updatePublisherDefinition()
	d.cron.AddFunc("* 30 * * * *", d.updatePublisherDefinition)
}

func (d *Decentralizer) downloadPublisherDefinition() {
	data, err := configPath.QueryCacheFolder().ReadFile(PUBLISHER_DEFINITION_FILE)
	if err != nil {
		//logger.Warningf("Could not restore publisher file: %v", err)
		return
	}
	var definition pb.PublisherUpdate
	err = proto.Unmarshal(data, &definition)
	if err != nil {
		logger.Warningf("Could not restore publisher file: %v", err)
		return
	}
	err = d.loadNewPublisherUpdate(&definition)
	if err != nil {
		logger.Warningf("Could not restore publisher file: %v", err)
		return
	}
}

func (d *Decentralizer) savePublisherUpdate() {
	if d.publisherUpdate.Created == 0 {
		logger.Warning("could not save publisher definition because it wasn't initialized")
		return
	}
	data, err := proto.Marshal(d.publisherUpdate)
	if err != nil {
		logger.Error(err)
		return
	}
	err = configPath.QueryCacheFolder().WriteFile(PUBLISHER_DEFINITION_FILE, data)
	if err != nil {
		logger.Error(err)
	}
	d.runPublisherInstructions()
}

func (d *Decentralizer) verifyPublisherUpdate(update *pb.PublisherUpdate) error {
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
	return nil
}

func (d *Decentralizer) loadNewPublisherUpdate(update *pb.PublisherUpdate) error {
	err := d.verifyPublisherUpdate(update)
	if err != nil {
		return err
	}
	d.publisherUpdate = update
	d.savePublisherUpdate()
	return nil
}

func (d *Decentralizer) signDefinition(definition *pb.PublisherDefinition) (*pb.PublisherUpdate, error) {
	data, err := proto.Marshal(definition)
	if err != nil {
		return nil, err
	}
	signature, err := d.n.Sign(data)
	if err != nil {
		return nil, err
	}
	return &pb.PublisherUpdate{
		Created: time.Now().Unix(),
		Signature: signature,
		Definition: definition,
	}, nil
}

func (d *Decentralizer) PublishPublisherUpdate(definition *pb.PublisherDefinition) error {
	update, err := d.signDefinition(definition)
	if err != nil {
		return err
	}
	err = d.loadNewPublisherUpdate(update)
	if err != nil {
		return err
	}
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
		lenPeers := len(d.i.PeerHost.Network().Peers())
		if lenPeers >= 3 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	checked := 0
	for _, peer := range d.i.PeerHost.Network().Peers() {
		if checked >= 3 {
			break
		}
		definition, err := d.getPublisherUpdateRequest(peer)
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

func (d *Decentralizer) getPublisherUpdateRequest(peer peer.ID) (*pb.PublisherUpdate, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_PUBLISHER_UPDATE)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(300 * time.Millisecond))
	defer stream.Close()
	logger.Infof("Requesting %s for their publisher file.", peer.Pretty())

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


func (d *Decentralizer) getPublisherUpdateResponse(stream inet.Stream) {
	if d.publisherUpdate == nil {
		logger.Info("Someone asked for our publisher update. But we ourselves don't have it yet.")
		stream.Conn().Close()
		return
	}
	logger.Info("Responding with our publisher update.")

	reqData, err := Read(stream)
	if err != nil {
		logger.Error(err)
		return
	}
	var request pb.DNPublisherUpdateRequest
	err = proto.Unmarshal(reqData, &request)
	if err != nil {
		logger.Error(err)
		return
	}

	//Response
	response, err := proto.Marshal(&pb.DNPublisherUpdateResponse{
		Update: d.publisherUpdate,
	})
	if err != nil {
		logger.Error(err)
		return
	}
	err = Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}