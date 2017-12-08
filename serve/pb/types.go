// This file has been automatically generated.
package pb

import "reflect"

type MessageType int32
var types map[reflect.Type]MessageType

const (
	RPCHealthRequest	MessageType = 1000
	RPCHealthReply	MessageType = 1001
)

func init() {
	types = make(map[reflect.Type]MessageType)
	registerType((*RPCMessage_HealthRequest)(nil), RPCHealthRequest)
	registerType((*RPCMessage_HealthReply)(nil), RPCHealthReply)
}

func registerType(x interface{}, msgType MessageType) {
	types[reflect.TypeOf(x)] = msgType
}