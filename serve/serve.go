package serve

import (
	"github.com/iain17/decentralizer/app"
	"github.com/iain17/decentralizer/pb"
	"reflect"
)

type Handler func(msg *pb.RPCMessage)(*pb.RPCMessage, error)

type Serve struct {
	app *app.Decentralizer
	handlers map[reflect.Type]Handler
}

func New(app *app.Decentralizer) *Serve {
	i := &Serve {
		app: app,
		handlers: make(map[reflect.Type]Handler),
	}
	//Platform
	i.registerHandler((*pb.RPCMessage_HealthRequest)(nil), i.handleHealthRequest)

	//Sessions
	i.registerHandler((*pb.RPCMessage_UpsertSessionRequest)(nil), i.handleUpsertSessionRequest)
	i.registerHandler((*pb.RPCMessage_DeleteSessionRequest)(nil), i.handleDeleteSessionRequest)
	i.registerHandler((*pb.RPCMessage_SessionIdsRequest)(nil), i.handleSessionIdsRequest)
	i.registerHandler((*pb.RPCMessage_GetSessionRequest)(nil), i.handleGetSessionRequest)
	return i
}

func (s *Serve) registerHandler(x interface{}, handler Handler) {
	s.handlers[reflect.TypeOf(x)] = handler
}