package app

import (
	"errors"
	"gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/core/coreapi"
	"io/ioutil"
)

func (d * Decentralizer) GetPublisherFile(name string) ([]byte, error) {
	if d.publisherUpdate == nil {
		return nil, errors.New("Publisher definition not defined")
	}
	var result []byte
	//First check the files
	if d.publisherUpdate.Definition.Files[name] != nil {
		result = d.publisherUpdate.Definition.Files[name]
	}
	//Try links
	if result == nil {

		//Fetch from IPFS
		if d.publisherUpdate.Definition.Links[name] != "" {
			path := d.publisherUpdate.Definition.Links[name]
			pth := coreapi.ResolvedPath(path, nil, nil)
			api := coreapi.NewCoreAPI(d.i)
			r, err := api.Unixfs().Cat(d.i.Context(), pth)
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