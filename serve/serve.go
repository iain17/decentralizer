package serve

import (
	"github.com/iain17/decentralizer/app"
	"github.com/iain17/decentralizer/serve/pb"
)

type Handler func(msg *pb.RPCMessage)(*pb.RPCMessage, error)

type Serve struct {
	app *app.Decentralizer
	handlers map[pb.MessageType]Handler
}

func New(app *app.Decentralizer) *Serve {
	i := &Serve {
		app: app,
		handlers: make(map[pb.MessageType]Handler),
	}
	i.registerHandler(pb.RPCHealthRequest, i.handleHealthRequest)
	return i
}

func (s *Serve) registerHandler(messageType pb.MessageType, handler Handler) {
	s.handlers[messageType] = handler
}