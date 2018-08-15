package app

import (
	"errors"
	"io/ioutil"
	//"github.com/ipfs/go-ipfs/core/coreunix"
	//ipfsfiles "gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit/files"
	//"os"
	//"github.com/iain17/old/ipLookup/src/utils/logger"
	//"fmt"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core/coreapi/interface"
)

func (d *Decentralizer) GetPublisherFile(name string) ([]byte, error) {
	if d.publisherRecord == nil {
		return nil, errors.New("Publisher definition not defined")
	}
	var result []byte
	//First check the files
	if d.publisherDefinition.Files[name] != nil {
		result = d.publisherDefinition.Files[name]
	}
	//Try links
	if result == nil {

		//Fetch from IPFS
		if d.publisherDefinition.Links[name] != "" {
			path := d.publisherDefinition.Links[name]
			pth, err := iface.ParsePath(path)
			if err != nil {
				return nil, err
			}
			r, err := d.api.Unixfs().Cat(d.i.Context(), pth)
			if err != nil {
				return nil, err
			}
			result, err = ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
		}
	}

	if result == nil {
		return nil, errors.New("could not find specified publisher file")
	}

	return result, nil
}

//func (d *Decentralizer) AddPublisherFiles(path string) error {
//	fileAdder, err := coreunix.NewAdder(d.i.Context(), d.i.Pinning, d.i.Blockstore, d.i.DAG)
//	if err != nil {
//		return err
//	}
//	fileAdder.NoCopy = true
//	fileAdder.RawLeaves = true
//	fileAdder.Pin = true
//	//fileAdder.Progress = true
//	//fileAdder.Out = out
//	fileAdder.Wrap = true
//
//	info, err := os.Lstat(path)
//	if err != nil {
//		return err
//	}
//	file, err := ipfsfiles.NewSerialFile(path, path, false, info)
//	if err != nil {
//		return err
//	}
//	err = fileAdder.AddFile(file)
//	if err != nil {
//		logger.Error(err)
//	}
//
//	nd, err := fileAdder.Finalize()
//	if err != nil {
//		logger.Error(err)
//	}
//
//	asd, _ := fileAdder.RootNode()
//	asd.Links()
//
//	pth := coreapi.ResolvedPath("/ipfs/"+nd.Cid().String(), nil, nil)
//	wtf, _ := d.api.Unixfs().Ls(d.i.Context(), pth)
//	for _, o := range wtf {
//		fmt.Printf("(iain) %s = %s", o.Name, o.Cid)
//	}
//
//	return nil
//}
