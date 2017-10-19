package sessionstore

import "github.com/ethersphere/go-ethereum/logger"

func (s *Store) onEvicted(key interface{}, value interface{}) {
	if sSessionId, ok := key.(uint64); ok {
		s.Remove(sSessionId)
		logger.Info("Deleted session id: %s", sSessionId)
	}
}