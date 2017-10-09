package discovery

import "github.com/op/go-logging"

type Node struct {
	info map[string]string
	logger        *logging.Logger
}

func (n *Node) String() string {
	return "Bare node."
}