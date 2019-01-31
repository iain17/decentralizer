package ipfs

import (
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
)

func PublishPath(n *core.IpfsNode, pth path.Path) error {
	// verify the path exists
	_, err := core.Resolve(n.Context(), n.Namesys, n.Resolver, pth)
	if err != nil {
		return err
	}
	k, err := n.GetKey("self")
	if err != nil {
		return err
	}
	return n.Namesys.Publish(n.Context(), k, pth)
}
