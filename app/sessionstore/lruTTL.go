package sessionstore

import "github.com/iain17/logger"

func (s *Store) onEvicted(key interface{}, value interface{}) {
	if sSessionId, ok := key.(uint64); ok {
		s.Remove(sSessionId)
		logger.Infof("Deleted session id: %s", sSessionId)
	}
}