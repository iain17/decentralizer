package pb

import (
	"github.com/stretchr/testify/assert"
	"github.com/golang/protobuf/proto"
	"testing"
)

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
