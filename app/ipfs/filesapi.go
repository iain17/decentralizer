package ipfs

import (
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core/coreapi"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
	"github.com/iain17/logger"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core/coreunix"
	Path "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/path"
	"io/ioutil"
	"context"
	libp2pPeer "gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	"fmt"
	 "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/blockservice"
	dag "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/merkledag"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit/files"
	"bytes"
)

//Simplifies all the interactions with IPFS.
type FilesAPI struct {
	ctx					   context.Context
	i 				   	   *core.IpfsNode
	api					   iface.CoreAPI
}

func NewFilesAPI(ctx context.Context, node *core.IpfsNode, api iface.CoreAPI) (*FilesAPI, error) {
	instance := &FilesAPI{
		ctx: ctx,
		i: node,
		api: api,
	}
	return instance, nil
}

func (d *FilesAPI) SaveFile(data []byte) (string, error) {
	path, err := coreunix.Add(d.i, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	return "/ipfs/"+path, err
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
	return "/ipns/"+d.i.Identity.Pretty(), err
}

func (d *FilesAPI) GetPeerFiles(owner libp2pPeer.ID) ([]*iface.Link, error) {
	logger.Infof("Get peer files of peer id %s", owner.Pretty())
	rawPath := "/ipns/" + owner.Pretty()
	pth := coreapi.ResolvedPath(rawPath, nil, nil)
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

	pth := coreapi.ResolvedPath(path, nil, nil)
	_, err := d.api.ResolvePath(d.i.Context(), pth)
	if err != nil {
		return nil, fmt.Errorf("could not get file %s. Could not resolve path: %s", path, err.Error())
	}
	r, err := d.api.Unixfs().Cat(d.i.Context(), pth)
	if err != nil {
		return nil, fmt.Errorf("could not get file %s: %s", path, err.Error())
	}
	return ioutil.ReadAll(r)
}
