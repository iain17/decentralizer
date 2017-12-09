package pb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
)

type fakeReadWriter struct{
	data []byte
}

func (f *fakeReadWriter) Write(p []byte) (n int, err error) {
	f.data = append(f.data, p...)
	return len(p), nil
}

func (f *fakeReadWriter) Read(p []byte) (n int, err error) {
	p = f.data
	return len(f.data)-1, nil
}

func TestWrite(t *testing.T) {
	w := &fakeReadWriter{}
	err := Write(w, &RPCMessage{
		Id: 1338,
		Msg: &RPCMessage_HealthReply{
			HealthReply: &HealthReply{
				Ready: true,
				Message: "very nice...",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x8, 0x1, 0x10, 0xba, 0xa, 0x22, 0x10, 0x8, 0x1, 0x12, 0xc, 0x76, 0x65, 0x72, 0x79, 0x20, 0x6e, 0x69, 0x63, 0x65, 0x2e, 0x2e, 0x2e, 255}, w.data)
}

func TestRead(t *testing.T) {
	ibuf := bytes.NewBuffer([]byte{0x8, 0x1, 0x10, 0xba, 0xa, 0x22, 0x10, 0x8, 0x1, 0x12, 0xc, 0x76, 0x65, 0x72, 0x79, 0x20, 0x6e, 0x69, 0x63, 0x65, 0x2e, 0x2e, 0x2e, 255})
	msg, err := Decode(ibuf)
	assert.NoError(t, err)
	assert.Equal(t, msg.Id, int64(1338))
	assert.NotNil(t, msg.GetHealthReply())
	assert.Equal(t, "very nice...", msg.GetHealthReply().Message)
}

