package ipfs

import (
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreunix"
	Path "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/path"
	"io/ioutil"
	"context"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"fmt"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/blockservice"
	dag "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/merkledag"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/commands/files"
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
		return nil, err
	}
	r, err := d.api.Unixfs().Cat(d.i.Context(), pth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
