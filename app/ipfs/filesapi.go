package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/iain17/logger"
	libp2pPeer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit/files"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/blockservice"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreunix"
	dag "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/merkledag"
	Path "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	ipld "gx/ipfs/QmZtNq8dArGfnpCZfx2pUNY7UcjGhVp5qqwQ4hH6mpTMRQ/go-ipld-format"
	"io/ioutil"
)

//Simplifies all the interactions with IPFS.
type FilesAPI struct {
	ctx context.Context
	i   *core.IpfsNode
	api iface.CoreAPI
}

func NewFilesAPI(ctx context.Context, node *core.IpfsNode, api iface.CoreAPI) (*FilesAPI, error) {
	instance := &FilesAPI{
		ctx: ctx,
		i:   node,
		api: api,
	}
	return instance, nil
}

func (d *FilesAPI) SaveFile(data []byte) (string, error) {
	path, err := coreunix.Add(d.i, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	return "/ipfs/" + path, err
}

func (d *FilesAPI) PublishPeerFiles(files []files.File) (string, error) {
	bserv := blockservice.New(d.i.Blockstore, d.i.Exchange)
	dserv := dag.NewDAGService(bserv)
	fileAdder, err := coreunix.NewAdder(d.i.Context(), d.i.Pinning, d.i.Blockstore, dserv)
	fileAdder.Wrap = true
	for _, file := range files {
		logger.Infof("Saving peer file %s", file.FileName())
		fileAdder.AddFile(file)
	}
	// copy intermediary nodes from editor to our actual dagservice
	_, err = fileAdder.Finalize()
	if err != nil {
		return "", err
	}
	err = fileAdder.PinRoot()
	if err != nil {
		return "", err
	}
	root, err := fileAdder.RootNode()
	ph := Path.FromCid(root.Cid())
	if err != nil {
		return "", err
	}
	err = PublishPath(d.i, ph)
	return "/ipns/" + d.i.Identity.Pretty(), err
}

func (d *FilesAPI) GetPeerFiles(owner libp2pPeer.ID) ([]*ipld.Link, error) {
	logger.Infof("Get peer files of peer id %s", owner.Pretty())
	rawPath := "/ipns/" + owner.Pretty()
	pth, err := iface.ParsePath(rawPath)
	if err != nil {
		return nil, err
	}
	//d.api.ResolveNode(d.i.Context(), pth)
	//return coreapi.NewCoreAPI(d.i).Unixfs().Ls(d.i.Context(), pth)
	return d.api.Unixfs().Ls(d.i.Context(), pth)
}

func (d *FilesAPI) GetPeerFile(owner libp2pPeer.ID, name string) ([]byte, error) {
	logger.Infof("Get peer file %s of peer id %s", name, owner.Pretty())
	peerFiles, err := d.GetPeerFiles(owner)
	if err != nil {
		return nil, err
	}
	var fileMap []string
	for _, file := range peerFiles {
		logger.Infof("Checking %s == %s", file.Name, name)
		if file.Name == name {
			return d.GetFile(file.Cid.String())
		}
		fileMap = append(fileMap, file.Name)
	}
	err = fmt.Errorf("could not find '%s' from %s in a list of: %v", name, owner.Pretty(), fileMap)
	logger.Warning(err)
	return nil, err
}

//Path could be "/ipfs/QmQy2Dw4Wk7rdJKjThjYXzfFJNaRKRHhHP5gHHXroJMYxk"
func (d *FilesAPI) GetFile(path string) ([]byte, error) {
	logger.Infof("Get file: %s", path)

	pth, err := iface.ParsePath(path)
	_, err = d.api.ResolvePath(d.i.Context(), pth)
	if err != nil {
		return nil, fmt.Errorf("could not get file %s. Could not resolve path: %s", path, err.Error())
	}
	r, err := d.api.Unixfs().Cat(d.i.Context(), pth)
	if err != nil {
		return nil, fmt.Errorf("could not get file %s: %s", path, err.Error())
	}
	return ioutil.ReadAll(r)
}
