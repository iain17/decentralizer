package main

import (
	"testing"

	"gx/ipfs/QmNUKMfTHQQpEwE8bUdv5qmKC3ymdW7zw82LFS8D6MQXmu/go-ipfs/commands"
)

func TestIsCientErr(t *testing.T) {
	t.Log("Catch both pointers and values")
	if !isClientError(commands.Error{Code: commands.ErrClient}) {
		t.Errorf("misidentified value")
	}
	if !isClientError(&commands.Error{Code: commands.ErrClient}) {
		t.Errorf("misidentified pointer")
	}
}
