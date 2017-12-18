package repo

import (
	"errors"
	"io"

	filestore "gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/filestore"
	keystore "gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/keystore"
	config "gx/ipfs/QmTxUjSZnG7WmebrX2U7furEPNSy33pLgA53PtpJYJSZSn/go-ipfs/repo/config"

	ma "gx/ipfs/QmW8s4zTsUoX1Q6CeYxVKPyqSKbF7H1YDUyTostBtZ8DaG/go-multiaddr"
	ds "gx/ipfs/QmdHG8MAuARdGHxx4rPQASLcvhz24fzjSQq7AJRAQEorq5/go-datastore"
)

var (
	ErrApiNotRunning = errors.New("api not running")
)

type Repo interface {
	Config() (*config.Config, error)
	SetConfig(*config.Config) error

	SetConfigKey(key string, value interface{}) error
	GetConfigKey(key string) (interface{}, error)

	Datastore() Datastore
	GetStorageUsage() (uint64, error)

	Keystore() keystore.Keystore

	FileManager() *filestore.FileManager

	// SetAPIAddr sets the API address in the repo.
	SetAPIAddr(addr ma.Multiaddr) error

	SwarmKey() ([]byte, error)

	io.Closer
}

// Datastore is the interface required from a datastore to be
// acceptable to FSRepo.
type Datastore interface {
	ds.Batching // should be threadsafe, just be careful
	io.Closer
}
