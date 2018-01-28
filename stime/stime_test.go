package stime

import (
	"testing"
	"fmt"
	"time"
)

//This this manually. Set the machine date to something weird or out of date.
func TestNow(t *testing.T) {
	a := Now()
	fmt.Printf("The time is: %s instead of %s", a, time.Now().UTC())
}