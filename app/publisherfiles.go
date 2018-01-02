package app

import (
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"github.com/golang/protobuf/proto"
	"errors"
	"time"
	"io/ioutil"
	"github.com/iain17/decentralizer/app/ipfs"
	"fmt"
	"encoding/hex"
	"github.com/jeffchao/backoff"
	"strings"
)

func (d *Decentralizer) getPublisherTopic() string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_PUBLISHER", hex.EncodeToString(ih[:]), ih)
}

func (d *Decentralizer) initPublisherFiles() {
	f := backoff.Fibonacci()
	f.Interval = 1 * time.Second
	_, err := ipfs.Subscribe(d.i, d.getPublisherTopic(), func(peer peer.ID, data []byte) {
		call := func() error {
			return d.receivedUpdate(peer, data)
		}
		err := f.Retry(call)
		if err != nil {
			logger.Warning(err)
		}
	})
	if err != nil {
		logger.Fatal(err)
	}
	d.downloadPublisherDefinition()
	d.cron.Every(10).Minutes().Do(func() {
		if d.i == nil {
			return
		}
		lenPeers := len(d.i.PeerHost.Network().Peers())
		if lenPeers <= MIN_CONNECTED_PEERS {
			return
		}
		d.PushPublisherUpdate()
	})
}

func (d *Decentralizer) readPublisherUpdateFromDisk() ([]byte, error) {
	data, err := configPath.QueryCacheFolder().ReadFile(PUBLISHER_DEFINITION_FILE)
	if err != nil {
		//Check if publisher file is in the same director as us
		data, err = ioutil.ReadFile("./" + PUBLISHER_DEFINITION_FILE)
	}
	return data, err
}

func (d *Decentralizer) downloadPublisherDefinition() {
	data, err := d.readPublisherUpdateFromDisk()
	if err != nil {
		logger.Warningf("Could not read publisher file: %v", err)
		return
	}
	var update pb.PublisherUpdate
	err = proto.Unmarshal(data, &update)
	if err != nil {
		logger.Warningf("Could not restore publisher file: %v", err)
		return
	}
	err = d.loadNewPublisherUpdate(&update)
	if err != nil {
		logger.Warningf("Could not restore publisher file: %v", err)
		return
	}
}

func (d *Decentralizer) savePublisherUpdate() {
	if d.publisherUpdate == nil {
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

func (d *Decentralizer) unmarshalPublisherDefinition(update *pb.PublisherUpdate) (*pb.PublisherDefinition, error) {
	var definition pb.PublisherDefinition
	err := proto.Unmarshal(update.Definition, &definition)
	if err != nil {
		return nil, err
	}
	if d.publisherUpdate != nil && d.publisherDefinition.Created >= definition.Created {
		return nil, errors.New("definition is older")
	}
	err = d.n.Verify(update.Definition, update.Signature)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("signature verification failed: %s", err.Error()))
	}
	return &definition, err
}

func (d *Decentralizer) loadNewPublisherUpdate(update *pb.PublisherUpdate) error {
	definition, err := d.unmarshalPublisherDefinition(update)
	if err != nil {
		return err
	}
	d.publisherUpdate = update
	d.publisherDefinition = definition
	d.savePublisherUpdate()
	return nil
}

func (d *Decentralizer) signDefinition(definition *pb.PublisherDefinition) (*pb.PublisherUpdate, error) {
	definition.Created = time.Now().Unix()
	data, err := proto.Marshal(definition)
	if err != nil {
		return nil, err
	}
	signature, err := d.n.Sign(data)
	if err != nil {
		return nil, err
	}
	return &pb.PublisherUpdate{
		Signature: signature,
		Definition: data,
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
	err = d.PushPublisherUpdate()
	if err == nil {
		d.cron.Every(1).Second().Do(d.PushPublisherUpdate)
	}
	return err
}

func (d *Decentralizer) PushPublisherUpdate() error {
	if d.publisherUpdate == nil {
		return errors.New("no update set")
	}
	data, err := proto.Marshal(d.publisherUpdate)
	if err != nil {
		return err
	}
	logger.Info("Publishing publisher update")
	return ipfs.Publish(d.i, d.getPublisherTopic(), data)
}

//Called when the publisher file has been loaded
func (d *Decentralizer) runPublisherInstructions() {
	//If the publisher file told us to stop. Stop!
	if !d.publisherDefinition.Status {
		panic("Publisher has closed this network!")
		return
	}
	logger.Infof("Publisher instructions loaded: %s", time.Unix(d.publisherDefinition.Created, 0).Format(time.RFC822))
}

func (d *Decentralizer) receivedUpdate(peer peer.ID, data []byte) error {
	pId := peer.Pretty()
	if d.ignore.Contains(pId) {
		return errors.New("on ignore list")
	}
	var update pb.PublisherUpdate
	err := proto.Unmarshal(data, &update)
	if err != nil {
		d.ignore.Add(pId, true)
		return err
	}
	err = d.loadNewPublisherUpdate(&update)
	if err != nil {
		if strings.Contains(err.Error(), "signature verification failed") {
			d.ignore.Add(pId, true)
			return err
		}
	} else {
		return d.PushPublisherUpdate()//means it updated. we'll republish.
	}
	return nil
}

func (d *Decentralizer) PublisherDefinition() *pb.PublisherDefinition {
	return d.publisherDefinition
}