package server // import "cirello.io/exp/cdci/pkg/server"
import (
	"sync"
	"time"

	"cirello.io/errors"
	"cirello.io/exp/cdci/pkg/api"
	"github.com/davecgh/go-spew/spew"
)

// Server dispatches new builds and collects the logs.
type Server struct {
	agents sync.Map // map[int64]time.Time - map of agentID to last ping

	tasks chan *api.Recipe
}

// New prepares a new server.
func New(tasks chan *api.Recipe) *Server {
	return &Server{
		tasks: tasks,
	}
}

// Run is the pipe where new tasks are dispatched to agents.
func (s *Server) Run(srv api.Runner_RunServer) error {
	spew.Dump("client connected to me")
	ctx := srv.Context()
	var agentErr error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				r, err := srv.Recv()
				if err != nil {
					agentErr = errors.E(err,
						"cannot read response from the agent")
					return
				}
				switch v := r.GetResponse().(type) {
				case *api.RunResponse_Result:
					s.agents.Store(v.Result.GetAgentID(), time.Now())
					spew.Dump("got result back (and updated pong)", v)
				case *api.RunResponse_Pong:
					s.agents.Store(v.Pong.GetAgentID(), time.Now())
					spew.Dump("got pong back", v.Pong)
				}
			}
		}
	}()

mainloop:
	for {
		select {
		case <-ctx.Done():
			spew.Dump(ctx.Err())
			break mainloop
		case <-time.After(2 * time.Second):
			err := srv.Send(&api.RunRequest{
				Action: &api.RunRequest_Ping{
					Ping: &api.Ping{},
				},
			})
			spew.Dump("pinging...", err)
		case recipe, ok := <-s.tasks:
			if !ok {
				break mainloop
			}
			err := srv.Send(&api.RunRequest{
				Action: &api.RunRequest_Recipe{
					Recipe: recipe,
				},
			})
			spew.Dump("worked dispatched...", err)
		}
	}
	if agentErr != nil {
		return agentErr
	}
	return nil
}
