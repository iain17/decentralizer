package framed

import (
	"testing"
	"sync"
	"io"
	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()

	// Use io.Pipe to simulate a network connection.
	// A real network application should take care to properly close the
	// underlying connection.
	rp, wp := io.Pipe()
	msg := []byte("A long time ago in a galaxy far, far away...")

	// Start a goroutine to act as the receiver.
	wg.Add(1)
	go func() {
		defer wg.Done()

		data, err := Read(rp)
		assert.NoError(t, err)
		assert.Equal(t, msg, data)
		rp.Close()
	}()

	// Start a goroutine to act as the transmitter.
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := Write(wp, msg)
		assert.NoError(t, err)
	}()
}