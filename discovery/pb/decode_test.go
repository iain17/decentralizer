package pb

import (
	"github.com/iain17/decentralizer/discovery/env"
	"github.com/stretchr/testify/assert"
	"github.com/golang/protobuf/proto"
	"net"
	"testing"
)

func TestDecode(t *testing.T) {

	l, err := net.Listen("tcp", ":1235")
	if err != nil {
		t.Fatal(err)
	}

	running := true
	go func() {
		conn, err := net.Dial("tcp", ":1235")
		assert.NoError(t, err)

		defer conn.Close()

		heartbeat, err := proto.Marshal(&Message{
			Version: env.VERSION,
			Msg: &Message_Heartbeat{
				Heartbeat: &Hearbeat{
					Message: "This now works fine!",
				},
			},
		})
		assert.NoError(t, err)
		err = Write(conn, heartbeat)
		assert.NoError(t, err)
		println("sent. done")
		for running {
		}
	}()

	defer l.Close()
	for {
		conn, err := l.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		println("Accepted")

		res, err := Decode(conn)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, res.GetHeartbeat())
		assert.Equal(t, res.GetHeartbeat().Message, "This now works fine!")
		break
	}
	running = false
}

func TestWrite(t *testing.T) {
	heartbeat, err := proto.Marshal(&Message{
		Version: 123,
		Msg: &Message_Heartbeat{
			Heartbeat: &Hearbeat{
				Message: "This actually works",
			},
		},
	})
	assert.NoError(t, err)
	var res Message
	err = proto.Unmarshal(heartbeat, &res)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.GetHeartbeat())
	assert.Equal(t, res.GetHeartbeat().Message, "This actually works")
}
