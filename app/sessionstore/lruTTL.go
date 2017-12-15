package sessionstore

import "github.com/iain17/logger"

func (s *Store) onEvicted(key interface{}, value interface{}) {
	//In go routine. Because of deadlock caused when evicted is called, the insert function higher up the call chain has just locked the mutex.
	go func() {
		if sSessionId, ok := key.(uint64); ok {
			err := s.Remove(sSessionId)
			if err != nil {
				logger.Warningf("Could not delete session id: %s", sSessionId)
			} else {
				logger.Infof("Deleted session id: %s", sSessionId)
			}
		} else {
			logger.Warningf("Could not delete session id...", sSessionId)
		}
	}()
}