package app

import (
	"errors"
	"gx/ipfs/QmYHpXQEWuhwgRFBnrf4Ua6AZhcqXCYa7Biv65SLGgTgq5/go-ipfs/core/coreapi"
	"io/ioutil"
)

func (d * Decentralizer) GetPublisherFile(name string) ([]byte, error) {
	if d.publisherUpdate == nil {
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
			pth := coreapi.ResolvedPath(path, nil, nil)
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