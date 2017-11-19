package app

import (
	"gx/ipfs/QmSn9Td7xgxm9EV7iEjTckpUWmWApggzPxu7eFGWkkpwin/go-block-format"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/keystore"
	"gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	"gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"github.com/ipfs/go-ipfs/path"
)

func (d *Decentralizer) SaveUserFile(name string, data []byte) error {
	block := blocks.NewBlock(data)
	err := d.i.Filestore.Put(block)
	if err != nil {
		return err
	}
	k, err := keylookup(d.i, "self")
	pth, err := path.ParsePath(block.Multihash().B58String())
	if err != nil {
		return err
	}
	// verify the path exists
	_, err = core.Resolve(d.i.Context(), d.i.Namesys, d.i.Resolver, pth)
	if err != nil {
		return err
	}
	return d.i.Namesys.Publish(d.i.Context(), k, pth)
}

//func (d *Decentralizer) GetUserFile(name string) (uint64, error) {
//	d.i.Filestore.Put(blocks.NewBlock(data))
//}

func keylookup(n *core.IpfsNode, k string) (crypto.PrivKey, error) {
	res, err := n.GetKey(k)
	if res != nil {
		return res, nil
	}

	if err != nil && err != keystore.ErrNoSuchKey {
		return nil, err
	}

	keys, err := n.Repo.Keystore().List()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		privKey, err := n.Repo.Keystore().Get(key)
		if err != nil {
			return nil, err
		}

		pubKey := privKey.GetPublic()

		pid, err := peer.IDFromPublicKey(pubKey)
		if err != nil {
			return nil, err
		}

		if pid.Pretty() == k {
			return privKey, nil
		}
	}
	return nil, fmt.Errorf("no key by the given name or PeerID was found")
}