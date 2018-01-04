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


	go d.republishPeerFiles()
}

func (d *Decentralizer) republishPeerFiles() {
	d.WaitTilEnoughPeers()
	files, err := d.GetPeerFiles("self")
	if err != nil {
		logger.Warning(err)
	}
	for name, _ := range files {
		data, err := d.GetPeerFile("self", name)
		if err != nil {
			logger.Warning(err)
			continue
		}
		d.SavePeerFile(name, data)
	}
}

func (d *Decentralizer) getPeerPath(owner libp2pPeer.ID) string {
	basePath := "/"+owner.Pretty()
	_, err := d.ufs.Stat(basePath)
	if os.IsNotExist(err) {
		d.ufs.MkdirAll(basePath, 0777)
	}
	return basePath
}

func (d *Decentralizer) getPeerFilePath(owner libp2pPeer.ID, name string) string {
	return fmt.Sprintf("%s/%s", d.getPeerPath(owner), name)
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
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return nil, err
	}
	var result []byte
	refresh := false
	path := d.getPeerFilePath(id, name)
	info, err := d.ufs.Stat(path)
	if info != nil && info.ModTime().After(time.Now().Add(FILE_EXPIRE)) {
		refresh = true
	}
	if id.Pretty() != d.i.Identity.Pretty() {
		refresh = true
	}
	if err == nil && !refresh {
		result, err = d.getFile(path)
		if err != nil {
			logger.Warning(err)
		}
	}
	if result == nil || refresh {
		//Time to get a fresh copy
		var fresh []byte
		fresh, err = d.filesApi.GetPeerFile(id, name)
		if err == nil && fresh != nil {
			result = fresh
			err = d.writeFile(path, result)
			if err != nil {
				logger.Warning(err)
			}
		}
	}
	return result, err
}

func (d *Decentralizer) GetPeerFiles(peerId string) (map[string]uint64, error) {
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return nil, err
	}
	//fetch locally
	path := d.getPeerPath(id)
	result := map[string]uint64{}
	err = afero.Walk(d.ufs, path, func(path string, info os.FileInfo, err error) error{
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		result[info.Name()] = (uint64)(info.Size())
		return nil
	})
	if err != nil {
		return nil, err
	}
	//fetch from peer
	links, err := d.filesApi.GetPeerFiles(id)
	if err != nil {
		logger.Warningf("Could not fetch fresh peer files")
	} else {
		for _, link := range links {
			result[link.Name] = link.Size
		}
	}
	return result, nil
}