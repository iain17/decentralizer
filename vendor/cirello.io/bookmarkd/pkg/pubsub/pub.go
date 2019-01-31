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

package pubsub // import "cirello.io/bookmarkd/pkg/pubsub"

import (
	"math"
	"math/rand"
	"sync"
)

// Broker is responsible to broadcast messages to the websocket pipes.
type Broker struct {
	mu          sync.Mutex
	subscribers map[int64]func(msg interface{})
}

// New starts a new websocket message broadcaster.
func New() *Broker {
	return &Broker{
		subscribers: make(map[int64]func(msg interface{})),
	}
}

// Subscribe registers a websocket listener.
func (b *Broker) Subscribe(f func(msg interface{})) (unsubscribe func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := rand.Int63n(math.MaxInt64)
	b.subscribers[id] = f
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		delete(b.subscribers, id)
	}
}

// Broadcast dispatches the messages to the websocket subscribers.
func (b *Broker) Broadcast(msg interface{}) {
	go func() {
		b.mu.Lock()
		subscribers := make(map[int64]func(msg interface{}))
		for k, v := range b.subscribers {
			subscribers[k] = v
		}
		b.mu.Unlock()

		for _, s := range subscribers {
			s(msg)
		}
	}()
}
