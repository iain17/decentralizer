package pb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net"
	"github.com/iain17/decentralizer/discovery/env"
	"github.com/golang/protobuf/proto"
)

func TestDecode(t *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", ":1235")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		heartbeat, err := proto.Marshal(&Message{
			Version: env.VERSION,
			Msg: &Message_Heartbeat{
				Heartbeat: &Hearbeat{
					Message: "This actually works",
				},
			},
		})
		assert.NoError(t, err)
		err = Write(conn, heartbeat)
		println("sent. done")
	}()

	l, err := net.Listen("tcp", ":1235")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		res, err := Decode(conn)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, res.GetHeartbeat())
		assert.Equal(t, res.GetHeartbeat().Message, "This actually works")
		break
	}
}

func TestWrite(t *testing.T) {
	heartbeat, err := proto.Marshal(&Message{
		Version: 123,
		Msg: &Message_Heartbeat{
			Heartbeat: &Hearbeat{},
		},
	})
	assert.NoError(t, err)
	var result Message
	err = proto.Unmarshal(heartbeat, &result)
	assert.NoError(t, err)
}