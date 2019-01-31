// Copyright 2018 github.com/ucirello/blackholedns
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

// Package dns implements a simple DNS forwarder that divert undesired requests
// to blackhole.
package dns // import "cirello.io/blackholedns/dns"

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

// Server takes DNS requests, and those matching special forward rules are
// handled specially. Otherwise, requests are plainly forwarded.
type Server struct {
	// BindAddr defines the TCP and UDP ports this server will listen to.
	BindAddr string `json:"bind_addr,omitempty"`

	// ForwwardAddr is the target DNS service to which DNS requests are
	// forwarded to if they are not diverted to the blackhole.
	ForwardAddr string `json:"forward_addr,omitempty"`

	// Blackhole is a list of DNS entries to which query calls will be
	// diverted to an invalid address (0.0.0.0)
	Blackhole []string `json:"blackhole,omitempty"`

	mu           sync.RWMutex
	knownDomains map[string]struct{}
}

func (s *Server) forward(m *dns.Msg) *dns.Msg {
	aM, err := dns.Exchange(m, s.ForwardAddr)
	if err != nil {
		log.Println(err)
		return nil
	}
	return aM
}

func (s *Server) handler(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name[:len(r.Question[0].Name)-1]

	s.mu.RLock()
	_, isBlackholed := s.knownDomains[domain]
	s.mu.RUnlock()
	if isBlackholed {
		s.sendToBlackhole(w, r, domain)
		return
	}

	for _, target := range s.Blackhole {
		if !strings.HasSuffix(domain, target) {
			continue
		}

		s.mu.Lock()
		s.knownDomains[domain] = struct{}{}
		s.mu.Unlock()
		s.sendToBlackhole(w, r, domain)
		return
	}

	fwd := s.forward(r)
	if fwd == nil {
		dns.HandleFailed(w, r)
		return
	}
	w.WriteMsg(fwd)
}

func (s *Server) sendToBlackhole(w dns.ResponseWriter, r *dns.Msg, domain string) {
	a, err := dns.NewRR(fmt.Sprintf("%s IN A 0.0.0.0", domain))
	if err != nil {
		log.Println(err)
		dns.HandleFailed(w, r)
		return
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Answer = []dns.RR{a}
	w.WriteMsg(m)
}

// ListenAndServe DNS requests.
func (s *Server) ListenAndServe() error {
	s.knownDomains = make(map[string]struct{})

	var wg sync.WaitGroup
	handler := dns.HandlerFunc(s.handler)
	errCh := make(chan error)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := dns.ListenAndServe(s.BindAddr, "udp", handler)
		if err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := dns.ListenAndServe(s.BindAddr, "tcp", handler)
		if err != nil {
			select {
			case errCh <- err:
			default:
			}

		}
	}()
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
