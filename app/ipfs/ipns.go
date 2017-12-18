package ipfs

import (
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/path"
)

func FilePublish(n *core.IpfsNode, pth path.Path) error {
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
