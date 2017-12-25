package app

import (
	"gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	inet "gx/ipfs/QmU4vCDZTPLDqSDKguWbHCiUe46mZUtmM2g2suBZ9NE8ko/go-libp2p-net"
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"github.com/golang/protobuf/proto"
	"errors"
	"time"
	"io/ioutil"
	"strings"
	"io"
	"github.com/iain17/framed"
)

func (d *Decentralizer) initPublisherFiles() {
	d.i.PeerHost.SetStreamHandler(GET_PUBLISHER_UPDATE, d.getPublisherUpdateResponse)
	d.downloadPublisherDefinition()
	go d.updatePublisherDefinition()
	d.cron.AddFunc("* 30 * * * *", d.updatePublisherDefinition)
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
	return nil
}

//Anything here we wanna do.
//Called when the publisher file has been loaded
func (d *Decentralizer) runPublisherInstructions() {
	//If the publisher file told us to stop. Stop!
	if !d.publisherDefinition.Status {
		panic("Publisher has closed this network!")
		return
	}
	logger.Infof("Publisher instructions loaded: %s", time.Unix(d.publisherDefinition.Created, 0).Format(time.RFC822))
}

//Will go through each connected peer and try and connect. Find out what publisher update they are running.
//If we've got 3 responses we'll stop trying. And take that this is the latest
func (d *Decentralizer) updatePublisherDefinition() {
	if d.searchingForPublisherUpdate {
		return
	}
	d.searchingForPublisherUpdate = true
	defer func() {
		d.searchingForPublisherUpdate = false
	}()

	select {
	case <-d.ctx.Done():
		break
	default:
		peers := d.i.PeerHost.Network().Peers()
		if len(peers) == 0 {
			logger.Info("Can't look for publisher definition. Not enough connected peers.")
			time.Sleep(10 * time.Second)
		} else {
			logger.Info("Looking for a updated publisher definition.")
			checked := 0
			for _, peer := range peers {
				if d.publisherDefinition != nil && checked >= MIN_CONNECTED_PEERS {
					break
				}
				id := peer.Pretty()
				if d.ignore[id] {
					continue
				}
				definition, err := d.getPublisherUpdateRequest(peer)
				if err != nil {
					if err.Error() == "i/o deadline reached" {
						continue
					}
					if err.Error() == "protocol not supported" {
						d.ignore[id] = true
						continue
					}
					if strings.Contains(err.Error(), "dial attempt failed") {
						d.ignore[id] = true
						continue
					}
					if err == io.EOF {
						continue
					}
					logger.Warningf("Could not get publisher update: %v", err)
					continue
				}
				checked++
				err = d.loadNewPublisherUpdate(definition)
				if err != nil {
					if err.Error() == "definition is older" {
						continue
					}
					logger.Warningf("Could not load new publisher update: %v", err)
				} else {
					break //updated.
				}
			}
			if d.publisherUpdate == nil || d.publisherDefinition == nil {
				logger.Warning("Could not find publisher definition. Retrying....")
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}
	}
	if d.publisherDefinition == nil {
		logger.Info("ok wtf.")
	}
}

func (d *Decentralizer) getPublisherUpdateRequest(peer peer.ID) (*pb.PublisherUpdate, error) {
	stream, err := d.i.PeerHost.NewStream(d.i.Context(), peer, GET_PUBLISHER_UPDATE)
	if err != nil {
		return nil, err
	}
	stream.SetDeadline(time.Now().Add(1 * time.Second))
	defer stream.Close()
	logger.Infof("Requesting %s for their publisher file.", peer.Pretty())

	//Request
	reqData, err := proto.Marshal(&pb.DNPublisherUpdateRequest{})
	if err != nil {
		return nil, err
	}
	err = framed.Write(stream, reqData)
	if err != nil {
		return nil, err
	}

	//Response
	resData, err := framed.Read(stream)
	if err != nil {
		return nil, err
	}
	if resData[0] == byte('N') && resData[1] == byte('O') && resData[2] == byte('P') {
		return nil, errors.New("peer doesn't have a publisher def")
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
		framed.Write(stream, []byte("NOP"))
		return
	}
	logger.Info("Responding with our publisher update.")

	reqData, err := framed.Read(stream)
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
	err = framed.Write(stream, response)
	if err != nil {
		logger.Error(err)
		return
	}
}