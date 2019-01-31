package utils

import (
	"time"
	"fmt"
	"github.com/iain17/logger"
	"errors"
	"github.com/iain17/stime"
)

//Returns true if the second passed record is newer than the first one.
//If its the same it will return true as well
func IsNewerRecord(current uint64, new uint64) bool {
	if new == 0 {
		return false
	}
	if current == 0 && new != 0 {
		return true
	}
	now := stime.Now()
	publishedTime := time.Unix(int64(new), 0).UTC()
	publishedTimeText := publishedTime.String()
	expireTime := time.Unix(int64(current), 0).UTC()
	expireTimeText := expireTime.String()
	if !publishedTime.After(expireTime) {
		err := fmt.Errorf("record with publish date %s is not newer than %s", publishedTimeText, expireTimeText)
		logger.Debug(err)
		return false
	}
	if publishedTime.After(now) {
		err := errors.New("new peer with publish date %s was published in the future")
		logger.Debug(err)
		return false
	}
	logger.Debugf("record with publish date %s IS newer than %s", publishedTimeText, expireTimeText)
	return true
}
