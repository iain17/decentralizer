package stime

import (
	"testing"
	"fmt"
	"time"
)

//This this manually. Set the machine date to something weird or out of date.
func TestNow(t *testing.T) {
	now := time.Now().UTC()
	ntpNow := Now()
	fmt.Printf("The time is: %s instead of %s", ntpNow, now)
}

func TestSlow(t *testing.T) {
	slow := IsBadNetwork()
	if slow {
		fmt.Println("Slow")
	} else {
		fmt.Println("Not slow.")
	}
}