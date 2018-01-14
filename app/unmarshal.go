package app

import (
	"github.com/gogo/protobuf/proto"
	"reflect"
	"errors"
)

//Cache un-marshalling. This does cost a lot more memory
func (b *Decentralizer) getUnmarshalKey(buf []byte) uint32 {
	b.crcTable.Reset()
	b.crcTable.Write(buf)
	return b.crcTable.Sum32()
}

func (d *Decentralizer) unmarshal(buf []byte, pb proto.Message) error {
	//Experimental.
	var cacheKey uint32
	if d.unmarshalCache != nil {
		cacheKey = d.getUnmarshalKey(buf)
		if cacheVal, ok := d.unmarshalCache.Get(cacheKey); ok {
			if val, ok := cacheVal.(proto.Message); ok {
				v := reflect.ValueOf(pb).Elem()
				v.Set(reflect.ValueOf(val).Elem())
			}
			if val, ok := cacheVal.(string); ok {
				return errors.New(val)
			}
			return nil
		}
	}
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		if d.unmarshalCache != nil {
			d.unmarshalCache.Add(cacheKey, err.Error())
		}
		return err
	}
	if d.unmarshalCache != nil {
		d.unmarshalCache.Add(cacheKey, pb)
	}
	return err
}