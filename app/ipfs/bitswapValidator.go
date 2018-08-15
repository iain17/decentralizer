package ipfs

import (
	"errors"
	"github.com/hashicorp/golang-lru"
	"hash"
	"hash/crc32"
)

type validateFunc_t func(key string, value []byte) error
type selectFunc_t func(key string, value [][]byte) (int, error)

type decentralizedValidator struct{
	dhtCache *lru.Cache //Cache our result to certain DHT values.
	crcTable hash.Hash32

	validateFunc validateFunc_t
	selectFunc selectFunc_t
}

func newDecentralizedValidator(validateFunc validateFunc_t, selectFunc selectFunc_t) (*decentralizedValidator, error) {
	dhtCache, err := lru.New(4096)
	if err != nil {
		return nil, err
	}
	return &decentralizedValidator{
		dhtCache: dhtCache,
		crcTable: crc32.NewIEEE(),
		validateFunc: validateFunc,
		selectFunc: selectFunc,
	}, nil
}

func (v *decentralizedValidator) getValidatorKey(keyType string, data []byte) uint32 {
	v.crcTable.Reset()
	v.crcTable.Write(data)
	//return fmt.Sprintf("%s/%d", keyType, b.crcTable.Sum32())
	return v.crcTable.Sum32()
}


func (v *decentralizedValidator) Validate(key string, value []byte) error {
	cacheKey := v.getValidatorKey(key, value)
	if cacheVal, ok := v.dhtCache.Get(cacheKey); ok {
		if val, ok := cacheVal.(bool); ok {
			if val {
				return nil
			} else {
				return errors.New("cache validator return error previously")
			}
		}
	}
	result := v.validateFunc(key, value)
	v.dhtCache.Add(cacheKey, result == nil)
	return result
}

func (v *decentralizedValidator) Select(key string, value [][]byte) (int, error) {
	return v.selectFunc(key, value)
}