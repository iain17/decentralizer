package ipfs

import (
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/path"
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
