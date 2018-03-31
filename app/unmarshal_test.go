package app

import (
	"testing"
	"hash/crc32"
	"github.com/hashicorp/golang-lru"
	"github.com/iain17/decentralizer/pb"
	"gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
	"github.com/iain17/stime"
)

var benchData = buildPb()

func buildPb() []byte {
	var data, err = proto.Marshal(&pb.Peer{
		Published: uint64(stime.Now().Unix()),
		PId: "test123",
		DnId: 1231231,
		Details: map[string]string{
			"this": "is",
			"very": "cool",
		},
	})
	if err != nil {
		panic(err)
	}
	return data
}

func BenchmarkDecentralizer_unmarshal1(b *testing.B) {
	unmarshalCache, err := lru.New(MAX_UNMARSHAL_CACHE)
	if err != nil {
		panic(err)
	}
	instance := &Decentralizer{
		unmarshalCache:			unmarshalCache,
		crcTable:				crc32.NewIEEE(),
	}

	for n := 0; n < b.N; n++ {
		var peer pb.Peer
		err = instance.unmarshal(benchData, &peer)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkDecentralizer_unmarshal2(b *testing.B) {
	_, err := lru.New(MAX_UNMARSHAL_CACHE)
	if err != nil {
		panic(err)
	}
	instance := &Decentralizer{
		//unmarshalCache:			unmarshalCache,
		//crcTable:				crc32.NewIEEE(),
	}

	for n := 0; n < b.N; n++ {
		var peer pb.Peer
		err = instance.unmarshal(benchData, &peer)
		if err != nil {
			panic(err)
		}
	}
}