// +build linux darwin freebsd netbsd openbsd
// +build nofuse

package node

import (
	"errors"

	core "gx/ipfs/QmUvjLCSYy7t4msRzrxfsfj99wboPhTUy7WktCv2LxS7BT/go-ipfs/core"
)

func Mount(node *core.IpfsNode, fsdir, nsdir string) error {
	return errors.New("not compiled in")
}
