package discovery

import (
	"fmt"
	"github.com/iain17/logger"
	"github.com/iain17/discovery/pb"
	"sync"
)

type Node struct {
	id        string //Unique random string identifier
	info map[string]string
	logger        *logger.Logger
	infoMutex	  sync.Mutex
}

func (n *Node) String() string {
	return fmt.Sprintf("bare node with info: %#v", n.info)
}

func (n *Node) SetInfo(key string, value string) {
	n.infoMutex.Lock()
	n.info[key] = value
	n.infoMutex.Unlock()
}

func (n *Node) GetInfo(key string) string {
	n.infoMutex.Lock()
	defer n.infoMutex.Unlock()
	return n.info[key]
}

func (n *Node) Initialize(info *pb.DPeerInfo) {
	n.infoMutex.Lock()
	n.id = info.Id
	n.info = info.Info
	n.infoMutex.Unlock()
}