package app

import (
	"github.com/iain17/logger"
	"github.com/iain17/decentralizer/pb"
	"errors"
	"time"
	"io/ioutil"
	"fmt"
	"encoding/hex"
	gogoProto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"github.com/iain17/decentralizer/utils"
	"github.com/iain17/stime"
	"github.com/iain17/decentralizer/vars"
)

func (d *Decentralizer) getPublisherKey() string {
	ih := d.n.InfoHash()
	return fmt.Sprintf("%s_PUBLISHER", hex.EncodeToString(ih[:]))
}

func (d *Decentralizer) initPublisherFiles() {
	d.b.RegisterValidator(vars.DHT_PUBLISHER_KEY_TYPE, func(key string, value []byte) error{
		var record pb.DNPublisherRecord
		err := d.unmarshal(value, &record)
		if err != nil {
			return err
		}
		//Definition should be 0
		if len(record.Definition) != 0 {
			return fmt.Errorf("you're doing it wrong! Definition should not be set on DHT")
		}
		if record.Path == "" {
			return fmt.Errorf("you're doing it wrong! Path should not be empty")
		}
		return d.validateDNPublisherRecord(&record)
	}, func(key string, values [][]byte) (int, error) {
		var currDefinition *pb.PublisherDefinition
		best := 0
		for i, val := range values {
			var record pb.DNPublisherRecord
			err := d.unmarshal(val, &record)
			if err != nil {
				logger.Warning(err)
				continue
			}
			definition, err := d.unmarshalDNPublisherRecord(&record)
			if err != nil {
				logger.Warning(err)
				continue
			}
			if currDefinition == nil || utils.IsNewerRecord(currDefinition.Published, definition.Published) {
				currDefinition = definition
				best = i
			}
		}
		return best, nil
	}, true)

	err := d.restorePublisherDefinition()
	if err != nil {
		logger.Warning(err)
	}
	go d.updatePublisherDefinition()
	d.cron.Every(300).Seconds().Do(func() {
		err := d.updatePublisherDefinition()
		if err != nil {
			logger.Warning(err.Error())
		}
	})
}

func (d *Decentralizer) validateDNPublisherRecord(record *pb.DNPublisherRecord) error {
	definition, err := d.unmarshalDNPublisherRecord(record)
	if d.publisherRecord != nil && d.publisherDefinition.Published != definition.Published && !utils.IsNewerRecord(d.publisherDefinition.Published, definition.Published) {
		return errors.New("definition is older")
	}
	return err
}

func (d *Decentralizer) unmarshalDNPublisherRecord(record *pb.DNPublisherRecord) (*pb.PublisherDefinition, error) {
	err := d.resolveDNPublisherRecord(record)
	if err != nil {
		return nil, err
	}
	var definition pb.PublisherDefinition
	err = d.unmarshal(record.Definition, &definition)
	if err != nil {
		return nil, err
	}
	err = d.n.Verify(record.Definition, record.Signature)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("signature verification '%x':%d failed: %s", record.Signature, len(record.Definition), err.Error()))
	}
	logger.Debugf("Publisher record unmarshaled. Signature '%x':%d validated", record.Signature, len(record.Definition))
	return &definition, err
}

//Makes sure that the record.Definition is filled.
func (d *Decentralizer) resolveDNPublisherRecord(record *pb.DNPublisherRecord) error {
	var err error
	if len(record.Definition) == 0 {
		record.Definition, err = d.filesApi.GetFile(record.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decentralizer) readPublisherRecordFromDisk() ([]byte, error) {
	data, err := configPath.QueryCacheFolder().ReadFile(vars.PUBLISHER_DEFINITION_FILE)
	if err != nil {
		//Check if publisher file is in the same director as us
		data, err = ioutil.ReadFile("./" + vars.PUBLISHER_DEFINITION_FILE)
	}
	if data == nil {
		data, err = Asset("static/publisherDefinition.dat")
	}
	return data, err
}

func (d *Decentralizer) readPublisherRecordFromNetwork() ([]byte, error) {
	d.WaitTilEnoughPeers()
	logger.Debugf("Asking the network for a publisher record")
	data, err := d.b.GetValue(d.ctx, vars.DHT_PUBLISHER_KEY_TYPE, d.getPublisherKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get best publisher record value: %s", err.Error())
	}
	return data, nil
}

func (d *Decentralizer) updatePublisherDefinition() error {
	data, err := d.readPublisherRecordFromNetwork()
	if data == nil || len(data) == 0 {
		err := d.PushPublisherRecord()
		if err != nil {
			logger.Warningf("Could not push publisher record: %s", err.Error())
		}
	}
	if err != nil {
		return fmt.Errorf("could not update publisher file from network: %s", err.Error())
	}
	return d.ReadPublisherDefinition(data)
}

func (d *Decentralizer) ResetPublisherDefinition() {
	d.publisherDefinition = nil
	d.publisherRecord = nil
}

//Restores from disk
func (d *Decentralizer) restorePublisherDefinition() error {
	data, err := d.readPublisherRecordFromDisk()
	if err != nil {
		return fmt.Errorf("could not restore publisher file from disk: %s", err.Error())
	}
	return d.ReadPublisherDefinition(data)
}

func (d *Decentralizer) ReadPublisherDefinition(data []byte) error {
	var record pb.DNPublisherRecord
	err := d.unmarshal(data, &record)
	if err != nil {
		return fmt.Errorf("could not read publisher file: %s", err.Error())
	}
	err = d.loadNewPublisherRecord(&record)
	if err != nil {
		return fmt.Errorf("could not read publisher file: %s", err.Error())
	}
	return nil
}

//Returns path, error
func (d *Decentralizer) savePublisherRecordToIpfs() (string, error) {
	if d.publisherRecord == nil {
		return "", fmt.Errorf("could not save publisher definition because it wasn't initialized")
	}
	var err error
	d.publisherRecord.Path, err = d.filesApi.SaveFile(d.publisherRecord.Definition)
	return d.publisherRecord.Path, err
}

func (d *Decentralizer) savePublisherRecordToDisk() error {
	if d.publisherRecord == nil {
		return fmt.Errorf("could not save publisher definition because it wasn't initialized")
	}
	data, err := gogoProto.Marshal(d.publisherRecord)
	if err != nil {
		return fmt.Errorf("could not marshal publisherRecord: %s", err.Error())
	}
	err = configPath.QueryCacheFolder().WriteFile(vars.PUBLISHER_DEFINITION_FILE, data)
	if err != nil {
		return fmt.Errorf("could not publisherRecord to disk: %s", err.Error())
	}
	return nil
}

//Loads in a new publisherRecord
func (d *Decentralizer) loadNewPublisherRecord(record *pb.DNPublisherRecord) error {
	definition, err := d.unmarshalDNPublisherRecord(record)
	if err != nil {
		return err
	}
	if d.publisherRecord != nil && d.publisherDefinition.Published >= definition.Published {
		return fmt.Errorf("our publisher record(%d) is newer or equal to %d", d.publisherDefinition.Published, definition.Published)
	}
	d.publisherRecord = record
	d.publisherDefinition = definition
	d.savePublisherRecordToDisk()
	d.savePublisherRecordToIpfs()
	d.runPublisherInstructions()
	d.PushPublisherRecord()
	return nil
}

func (d *Decentralizer) signPublisherRecord(definition *pb.PublisherDefinition) (*pb.DNPublisherRecord, error) {
	definition.Published = uint64(stime.Now().Unix())
	data, err := gogoProto.Marshal(definition)
	if err != nil {
		return nil, err
	}
	signature, err := d.n.Sign(data)
	if err != nil {
		return nil, err
	}
	return &pb.DNPublisherRecord{
		Signature: signature,
		Definition: data,
	}, nil
}

func (d *Decentralizer) PublishPublisherRecord(definition *pb.PublisherDefinition) ([]byte, error) {
	update, err := d.signPublisherRecord(definition)
	if err != nil {
		return nil, err
	}
	err = d.loadNewPublisherRecord(update)
	if err != nil {
		return nil, err
	}
	err = d.PushPublisherRecord()
	if err != nil {
		if err.Error() == "failed to find any peer in table" {
			err = nil
		}
	}
	return configPath.QueryCacheFolder().ReadFile(vars.PUBLISHER_DEFINITION_FILE)
}

func (d *Decentralizer) PushPublisherRecord() error {
	if d.limitedConnection {
		return errors.New("can't push on a limited connection")
	}
	if d.publisherRecord == nil {
		return errors.New("no update set")
	}
	if d.publisherRecord.Path == "" {
		return errors.New("no path set")
	}
	//Because we are going to push a update we will remove the binary definition
	d.publisherRecord.Definition = nil

	data, err := gogoProto.Marshal(d.publisherRecord)
	if err != nil {
		return err
	}
	logger.Info("Publishing our publisher version")
	return d.b.PutValue(vars.DHT_PUBLISHER_KEY_TYPE, d.getPublisherKey(), data)
}

//Called when the publisher file has been loaded
func (d *Decentralizer) runPublisherInstructions() {
	//If the publisher file told us to stop. Stop!
	if !d.publisherDefinition.Status {
		panic("Publisher has closed this network!")
		return
	}
	logger.Infof("Publisher instructions loaded: %s", time.Unix(int64(d.publisherDefinition.Published), 0).UTC().Format(time.RFC822))
}

func (d *Decentralizer) PublisherDefinition() *pb.PublisherDefinition {
	return d.publisherDefinition
}