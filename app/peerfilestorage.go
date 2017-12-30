package app

import (
	"github.com/iain17/decentralizer/app/ipfs"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"github.com/iain17/logger"
	"github.com/spf13/afero"
	"time"
	"fmt"
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"github.com/shibukawa/configdir"
)

func (d *Decentralizer) initStorage() {
	var err error
	d.filesApi, err = ipfs.NewFilesAPI(d.ctx, d.i, d.api)
	if err != nil {
		logger.Fatalf("Could not start filesapi: %s", err.Error())
	}
	paths := configPath.QueryFolders(configdir.Global)
	if len(paths) == 0 {
		logger.Fatal("Could not resolve config path")
	}
	base := afero.NewBasePathFs(afero.NewOsFs(), paths[0].Path+"/peer-data")
	layer := afero.NewMemMapFs()
	d.ufs = afero.NewCacheOnReadFs(base, layer, 100 * time.Second)


	go d.restorePeerFiles()
}

func (d *Decentralizer) restorePeerFiles() {
	d.WaitTilEnoughPeers()
	reveries, _ := Asset("reveries.flac")
	d.SavePeerFile("reveries.flac", reveries)
	d.GetPeerFile("self", "reveries.flac")
}

func (d *Decentralizer) getPeerFilePath(owner libp2pPeer.ID, name string) string {
	basePath := "/"+owner.Pretty()
	_, err := d.ufs.Stat(basePath)
	if os.IsNotExist(err) {
		d.ufs.MkdirAll(basePath, 0777)
	}
	return fmt.Sprintf("%s/%s", basePath, name)
}

//Save our peer file
func (d *Decentralizer) SavePeerFile(name string, data []byte) (string, error) {
	id := d.i.Identity
	path := d.getPeerFilePath(id, name)
	err := d.writeFile(path, data)
	if err != nil {
		return "", err
	}
	return d.filesApi.SavePeerFile(name, data)
}

func (d *Decentralizer) writeFile(path string, data []byte) error {
	f, err := d.ufs.Create(path)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if n != len(data) {
		return errors.New("partial write")
	}
	return err
}

func (d *Decentralizer) getFile(path string) ([]byte, error) {
	f, err := d.ufs.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func (d *Decentralizer) getIPFSFile(path string) ([]byte, error) {
	return d.filesApi.GetFile(path)
}

//Get a particular peer file from someone.
func (d *Decentralizer) GetPeerFile(peerId string, name string) ([]byte, error) {
	var result []byte
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return nil, err
	}
	path := d.getPeerFilePath(id, name)
	info, err := d.ufs.Stat(path)
	if err != nil || info.ModTime().After(time.Now().Add(FILE_EXPIRE)) {
		//Time to get a fresh copy
		result, err = d.filesApi.GetPeerFile(id, name)
		if err == nil && result != nil {
			err = d.writeFile(path, result)
		}
	}
	//No result yet?
	if result == nil {
		result, err = d.getFile(path)
	}
	return result, err
}