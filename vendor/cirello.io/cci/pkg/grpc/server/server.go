// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package server implements the server-side GRPC interface of the coordinator
// and workers.
package server // import "cirello.io/cci/pkg/grpc/server"

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cirello.io/cci/pkg/coordinator"
	"cirello.io/cci/pkg/grpc/api"
	"cirello.io/cci/pkg/models"
	"cirello.io/errors"
)

// Server interprets worker calls to coordinator.
type Server struct {
	coord         *coordinator.Coordinator
	configuration models.Configuration

	lockRefresh *sync.Cond
	locks       map[string]*lock // map of repoNames+workerID to lock state
}

// New instantiates a new server.
func New(coord *coordinator.Coordinator, configuration models.Configuration) *Server {
	s := &Server{
		coord:         coord,
		configuration: configuration,
		lockRefresh:   sync.NewCond(&sync.Mutex{}),
		locks:         make(map[string]*lock),
	}
	for repoName, recipe := range configuration {
		for i := 0; i < int(recipe.Concurrency); i++ {
			lockName := fmt.Sprintf("%v-%v", repoName, i)
			s.locks[lockName] = &lock{}
		}
	}
	go s.expireLocks()
	return s
}

// Configuration allows for server-side initiated worker configuration.
func (s *Server) Configuration(ctx context.Context, _ *api.ConfigurationRequest) (*api.ConfigurationResponse, error) {
	resp := &api.ConfigurationResponse{
		Configuration: make(map[string]*api.Recipe),
	}
	for k, v := range s.configuration {
		resp.Configuration[k] = &v
	}
	return resp, nil
}

// Run coordinates both the delivery of build job to workers and their actual
// liveness.
func (s *Server) Run(srv api.Runner_RunServer) error {
	// TODO: convert errors to GRPC's.
	// Run has two phases: a) handshake in which the worker declares which
	// repo it is observing; b) a continuous keep-alive to prove that the
	// worker is still alive.
	// The handshake takes the declared repository and tries to grab an
	// execution slot for the work. Even there are more workers than slot,
	// only n-slotted workers will actually have work to do.
	req, err := srv.Recv()
	if err != nil {
		return errors.E(err, "error receiving call from client")
	}
	buildReq := req.GetBuild()
	if buildReq == nil {
		return errors.E("client did not send the handshake message")
	}

	repoName := buildReq.GetRepoFullName()

	lockIndex, lockSeq, err := s.waitForLock(repoName)
	if err != nil {
		return errors.E(err, "cannot grab lock")
	}

	ctx, cancel := context.WithCancel(srv.Context())
	defer cancel()

	go func() {
		defer cancel()
		t := time.Tick(1 * time.Second)
		for {
			select {
			case <-srv.Context().Done():
				return
			case <-t:
				if !s.isLockOwner(lockIndex, lockSeq) {
					cancel()
				}
			default:
				req, err := srv.Recv()
				if err != nil {
					err := errors.E(err, "error receiving call from client")
					log.Println(err)
					return
				}
				switch req.GetCommand().(type) {
				case *api.JobRequest_KeepAlive:
					if err := s.refreshLock(lockIndex, lockSeq); err != nil {
						err := errors.E(err, "lost lock")
						log.Println(err)
						return
					}
				case *api.JobRequest_MarkInProgress:
					build := req.GetMarkInProgress()
					err := s.markInProgress(build, lockIndex, lockSeq)
					if err != nil {
						err := errors.Wrapf(err, "cannot mark build %d as in progress", build.ID)
						log.Println(err)
						return
					}
				case *api.JobRequest_MarkComplete:
					build := req.GetMarkComplete()
					err := s.markComplete(build, lockIndex, lockSeq)
					if err != nil {
						err := errors.Wrapf(err, "cannot mark build %d as completed", build.ID)
						log.Println(err)
						return
					}
				}
			}
		}
	}()

	for {
		log.Println("GRPC server dispatching for", repoName)
		select {
		case <-ctx.Done():
			return nil
		case build := <-s.coord.Next(repoName):
			if err := s.dispatchBuild(srv, repoName, lockIndex, lockSeq, build); err != nil {
				cancel()
				log.Println("cannot dispatch build:", err)
				return err
			}
		}
	}
}

func (s *Server) refreshLock(lockName string, seq int) error {
	s.lockRefresh.L.Lock()
	defer s.lockRefresh.L.Unlock()
	defer s.lockRefresh.Broadcast()
	l, ok := s.locks[lockName]
	if !ok {
		return errors.Errorf("cannot find lock for %s", lockName)
	}
	if l.seq != seq {
		return errors.E("not current lock owner")
	}
	l.refresh(seq)
	return nil
}

func (s *Server) markInProgress(build *api.Build, lockIndex string, lockSeq int) error {
	if err := s.refreshLock(lockIndex, lockSeq); err != nil {
		return errors.E(err)
	}
	err := s.coord.MarkInProgress(&models.Build{
		Build: build,
	})
	return errors.E(err)
}

func (s *Server) markComplete(build *api.Build, lockIndex string, lockSeq int) error {
	if err := s.refreshLock(lockIndex, lockSeq); err != nil {
		return errors.E(err)
	}
	err := s.coord.MarkComplete(&models.Build{
		Build: build,
	})
	return errors.E(err)
}

func (s *Server) dispatchBuild(srv api.Runner_RunServer, repoName, lockIndex string, lockSeq int, build *models.Build) error {
	if err := s.refreshLock(lockIndex, lockSeq); err != nil {
		return errors.E(err, "lost lock before build dispatch")
	}
	jobResp := &api.JobResponse{
		Build:  build.Build,
		Recipe: build.Recipe,
	}
	err := srv.Send(jobResp)
	return errors.E(err, "cannot dispatch build to client")
}

func (s *Server) waitForLock(repoName string) (lockName string, seq int, err error) {
	s.lockRefresh.L.Lock()
	defer s.lockRefresh.L.Unlock()
	// TODO: implement using lock server
	cfg, ok := s.configuration[repoName]
	if !ok {
		return "", -1, errors.Errorf("invalid repository: %s", repoName)
	}
	for {
		for i := 0; i < int(cfg.Concurrency); i++ {
			lockName := fmt.Sprintf("%v-%v", repoName, i)
			l, ok := s.locks[lockName]
			if !ok {
				log.Println("lockName not found", lockName)
				return "", -1, errors.Errorf("cannot find lock for %s", lockName)
			}
			seq, err := l.tryLock()
			if err != nil {
				log.Println("tryLock:", err)
				continue
			}
			return lockName, seq, nil
		}
		s.lockRefresh.Wait()
	}

}

func (s *Server) expireLocks() {
	const ttl = 1 * time.Minute
	t := time.Tick(time.Second)
	for range t {
		s.lockRefresh.L.Lock()
		for lockName, l := range s.locks {
			if l.expire(ttl) {
				log.Println(lockName, "expired")
			}
		}
		s.lockRefresh.Broadcast()
		s.lockRefresh.L.Unlock()
	}
}

func (s *Server) isLockOwner(lockName string, seq int) bool {
	s.lockRefresh.L.Lock()
	defer s.lockRefresh.L.Unlock()
	defer s.lockRefresh.Broadcast()
	l, ok := s.locks[lockName]
	if !ok {
		return false
	}
	return l.seq == seq
}
