package app

import (
	"github.com/ipfs/go-ipfs/core/coreapi"
	//"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	"bytes"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/ipfs/go-ipfs/core/coreapi/interface"
	"github.com/ipfs/go-ipfs/core/coreunix"
	Path "github.com/ipfs/go-ipfs/path"
	libp2pPeer "gx/ipfs/QmWNY7dV54ZDYmTA1ykVdwNCqC11mpU4zSUp6XDpLTH9eG/go-libp2p-peer"
	"io/ioutil"
	"errors"
	"github.com/iain17/logger"
)

func (d *Decentralizer) SavePeerFile(name string, data []byte) (string, error) {
	logger.Infof("Saving peer file %s", name)
	location, path, err := coreunix.AddWrapped(d.i, bytes.NewBuffer(data), name)
	if err != nil {
		return "", err
	}
	ph := Path.FromCid(path.Cid())
	if err != nil {
		return "", err
	}
	logger.Infof("Path %s", ph)
	err = ipfs.FilePublish(d.i, ph)
	if err != nil {
		return "", err
	}
	return location, nil
}

func (d *Decentralizer) GetPeerFiles(peerId string) ([]*iface.Link, error) {
	logger.Infof("Get peer files of peer id %s", peerId)
	id, err := libp2pPeer.IDB58Decode(peerId)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(d.i)
	rawPath := "/ipns/" + id.Pretty()
	pth := coreapi.ResolvedPath(rawPath, nil, nil)
	return api.Unixfs().Ls(d.i.Context(), pth)
}

func (d *Decentralizer) GetPeerFile(peerId string, name string) ([]byte, error) {
	logger.Infof("Get peer file %s of per id %s", name, peerId)
	files, err := d.GetPeerFiles(peerId)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name == name {
			return d.GetFile(file.Cid.String())
		}
	}
	return nil, errors.New("could not find peer file")
}

//Path could be "/ipfs/QmQy2Dw4Wk7rdJKjThjYXzfFJNaRKRHhHP5gHHXroJMYxk"
func (d *Decentralizer) GetFile(path string) ([]byte, error) {
	logger.Infof("Get file", path)
	api := coreapi.NewCoreAPI(d.i)

	pth := coreapi.ResolvedPath(path, nil, nil)
	_, err := api.ResolvePath(d.i.Context(), pth)
	if err != nil {
		return nil, err
	}
	r, err := api.Unixfs().Cat(d.i.Context(), pth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
