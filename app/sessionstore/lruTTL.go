package sessionstore

import "github.com/iain17/logger"

func (s *Store) onEvicted(key interface{}, value interface{}) {
	println("damn")
	if sSessionId, ok := key.(uint64); ok {
		println("k")
		err := s.Remove(sSessionId)
		if err != nil {
			logger.Warningf("Could not delete session id: %s", sSessionId)
		} else {
			logger.Infof("Deleted session id: %s", sSessionId)
		}
	}
}