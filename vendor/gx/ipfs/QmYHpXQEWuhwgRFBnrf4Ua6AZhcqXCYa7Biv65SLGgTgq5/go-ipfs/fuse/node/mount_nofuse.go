// +build linux darwin freebsd netbsd openbsd
// +build nofuse

package node

import (
	"errors"

	core "gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core"
)

func Mount(node *core.IpfsNode, fsdir, nsdir string) error {
	return errors.New("not compiled in")
}
