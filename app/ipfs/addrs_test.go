package ipfs

import (
	"testing"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	"gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	"net"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
)

func TestFilterNonReachableAddrs(t *testing.T) {
	srv := &http.Server{Addr: ":1337"}
	go func() {
		srv.ListenAndServe()
	}()
	srv1 := &http.Server{Addr: ":1338"}
	go func() {
		srv1.ListenAndServe()
	}()
	defer func() {
		srv.Close()
		srv1.Close()
	}()

	time.Sleep(1 * time.Second)

	addr1, err := manet.FromNetAddr(&net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 1337,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr1)
	addr2, err := manet.FromNetAddr(&net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 1338,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr2)
	addr3, err := manet.FromNetAddr(&net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr3)
	addr4, err := ma.NewMultiaddr("/p2p-circuit/ipfs/QmVs1rnP17aDgbqHjqkfh3iQtMXPqpSU43gdJQ398pmEee")
	assert.NoError(t, err)
	assert.NotNil(t, addr4)

	result := FilterNonReachableAddrs([]ma.Multiaddr{addr1, addr3, nil, addr2, addr3, addr4},true, false, true)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, addr1.String(), result[0].String())
	assert.Equal(t, addr2.String(), result[1].String())
	assert.Equal(t, addr4.String(), result[2].String())

	result = FilterNonReachableAddrs([]ma.Multiaddr{addr3}, true, false, true)
	assert.Equal(t, 0, len(result))

	result = FilterNonReachableAddrs([]ma.Multiaddr{}, true,false, true)
	assert.Equal(t, 0, len(result))

	result = FilterNonReachableAddrs([]ma.Multiaddr{addr1}, true,false, true)
	assert.Equal(t, 1, len(result))
}

func TestFilterNonReachableAddrs2(t *testing.T) {
	a, err := net.LookupIP("google.com")
	assert.NoError(t, err)
	addr1, err := manet.FromNetAddr(&net.TCPAddr{
		IP: a[0],
		Port: 80,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr1)
	addr2, err := manet.FromNetAddr(&net.TCPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 80,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr1)
	addr3, err := manet.FromNetAddr(&net.TCPAddr{
		IP: a[1],
		Port: 1337,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addr1)
	result := FilterNonReachableAddrs([]ma.Multiaddr{addr1, addr3, nil, addr2, addr3},false, false, true)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, addr1.String(), result[0].String())
}