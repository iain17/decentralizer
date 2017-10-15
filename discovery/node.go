package discovery

import (
	"github.com/op/go-logging"
	"fmt"
)

type Node struct {
	info map[string]string
	logger        *logging.Logger
}

func (n *Node) String() string {
	return fmt.Sprintf("bare node with info: %#v", n.info)
}

func (n *Node) SetInfo(key string, value string) {
	n.info[key] = value
}

func (n *Node) GetInfo(key string) string {
	return n.info[key]
}