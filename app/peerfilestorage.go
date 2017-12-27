package app

import (
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	//"gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	"bytes"
	"errors"
	"github.com/iain17/decentralizer/app/ipfs"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreunix"
	Path "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/path"
	"io/ioutil"
)
//See: https://gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/blob/master/core/commands/files/files.go
//that is also what is called from the http api.
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
	return "/ipfs/"+location, nil
}

func (d *Decentralizer) GetPeerFiles(peerId string) ([]*iface.Link, error) {
	logger.Infof("Get peer files of peer id %s", peerId)
	id, err := d.decodePeerId(peerId)
	if err != nil {
		return nil, err
	}
	api := coreapi.NewCoreAPI(d.i)
	rawPath := "/ipns/" + id.Pretty()
	pth := coreapi.ResolvedPath(rawPath, nil, nil)
	return api.Unixfs().Ls(d.i.Context(), pth)
}

func (d *Decentralizer) GetPeerFile(peerId string, name string) ([]byte, error) {
	logger.Infof("Get peer file %s of peer id %s", name, peerId)
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
	logger.Infof("Get file: %s", path)
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
