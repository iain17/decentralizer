package ipfs

import (
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/path"
	"bytes"
	"errors"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
	"github.com/iain17/logger"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreunix"
	Path "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/path"
	"io/ioutil"
	"context"
	libp2pPeer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	coreiface "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi/interface"
)
//Simplifies all the interactions with IPFS.
const CONCURRENT_PUBLISH = 2
type FilesAPI struct {
	ctx					   context.Context
	newPathToPublish       chan path.Path
	i 				   	   *core.IpfsNode
	api					   coreiface.CoreAPI
}

func NewFilesAPI(ctx context.Context, node *core.IpfsNode, api coreiface.CoreAPI) (*FilesAPI, error) {
	instance := &FilesAPI{
		ctx: ctx,
		i: node,
		api: api,
		newPathToPublish: make(chan path.Path, CONCURRENT_PUBLISH*2),
	}
	logger.Debugf("Running %d user file publish workers", CONCURRENT_PUBLISH)
	for i := 0; i < CONCURRENT_PUBLISH; i++ {
		go instance.processPublication()
	}
	return instance, nil
}

func (d *FilesAPI) processPublication() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case path, ok := <-d.newPathToPublish:
			if !ok {
				return
			}
			err := FilePublish(d.i, path)
			if err != nil {
				logger.Warningf("Failed to publish %s: %s", path, err)
			}
		}
	}
}

func (d *FilesAPI) SavePeerFile(name string, data []byte) (string, error) {
	logger.Infof("Saving peer file %s", name)
	location, path, err := coreunix.AddWrapped(d.i, bytes.NewBuffer(data), name)
	if err != nil {
		return "", err
	}
	ph := Path.FromCid(path.Cid())
	if err != nil {
		return "", err
	}
	d.newPathToPublish <- ph
	return "/ipfs/"+location, nil
}

func (d *FilesAPI) GetPeerFiles(owner libp2pPeer.ID) ([]*iface.Link, error) {
	logger.Infof("Get peer files of peer id %s", owner.Pretty())
	rawPath := "/ipns/" + owner.Pretty()
	pth := coreapi.ResolvedPath(rawPath, nil, nil)
	return d.api.Unixfs().Ls(d.i.Context(), pth)
}

func (d *FilesAPI) GetPeerFile(owner libp2pPeer.ID, name string) ([]byte, error) {
	logger.Infof("Get peer file %s of peer id %s", name, owner.String())
	files, err := d.GetPeerFiles(owner)
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
