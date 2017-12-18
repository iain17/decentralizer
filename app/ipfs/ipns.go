package ipfs

import (
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/path"
	"github.com/iain17/timeout"
	"time"
	"context"
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

	completed := false
	timeout.Do(func(ctx context.Context) {
		err = n.Namesys.Publish(n.Context(), k, pth)
		completed = true
	}, 5*time.Second)
	//if !completed {
	//	err = errors.New("could not publish file in under 15 seconds. Check if you are connected to enough peers")
	//}
	return err
}
