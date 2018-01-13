package app

import (
	"errors"
	"fmt"
)

func (d *Decentralizer) Health() (bool, error) {
	if d.publisherRecord == nil {
		return false, errors.New(fmt.Sprintf("Not ready yet. Waiting for publisher file..."))
	}
	if !d.publisherDefinition.Status {
		return false, errors.New("closed")
	}
	return true, nil
}