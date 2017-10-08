package discovery

import (
	"time"
	"context"
	"github.com/op/go-logging"
	"net"
	"github.com/pkg/errors"
	"github.com/hashicorp/golang-lru"
)
const CONCCURENT_NEW_CONNECTION = 10

type NetTableService struct {
	localNode *LocalNode
	context context.Context
	newConn chan *net.UDPAddr
	newPeer chan *RemoteNode

	blackList *lru.Cache
	peers     *lru.Cache

	heartbeatTicker <-chan time.Time

	logger *logging.Logger
}

func (nt *NetTableService) Init(ctx context.Context, ln *LocalNode) error {
	nt.logger = logging.MustGetLogger("NetTable")
	nt.localNode = ln
	nt.context = ctx
	nt.newConn = make(chan *net.UDPAddr, CONCCURENT_NEW_CONNECTION)
	var err error
	nt.blackList, err = lru.New(1000)
	if err != nil {
		return err
	}
	nt.peers, err = lru.NewWithEvict(nt.localNode.discovery.max, evicted)
	if err != nil {
		return err
	}
	nt.newPeer = make(chan *RemoteNode)
	nt.Run()
	return nil
}

func evicted(_ interface{}, value interface{}) {
	if node, ok := value.(*RemoteNode); ok {
		node.Close()
	}
}

func (nt *NetTableService) Run() error {
	//Spawn some workers
	for i := 0; i < CONCCURENT_NEW_CONNECTION; i++ {
		go nt.processDHTIn()
	}
	//Send a heart beat to the peers we are connected to
	go nt.heartbeat()
	return nil
}

func (nt *NetTableService) processDHTIn() {
	defer nt.Stop()
	for {
		select {
		case <-nt.context.Done():
			return
		case host, ok := <-nt.newConn:
			if !ok {
				return
			}

			if nt.isBlackListed(host) {
				continue
			}

			if err := nt.tryConnect(host); err != nil {
				nt.logger.Debugf("unable connect %s err: %s", host, err)
			}
		}
	}
}

func (nt *NetTableService) Stop() {
	for _, key := range nt.peers.Keys() {
		value, res := nt.peers.Get(key)
		if !res {
			continue
		}
		value.(*RemoteNode).Close()
	}
}

func (nt *NetTableService) GetNewConnChan() chan<- *net.UDPAddr {
	return nt.newConn
}

func (nt *NetTableService) AddRemoteNode(rn *RemoteNode) {
	addr := rn.conn.RemoteAddr().String()
	nt.peers.Add(addr, rn)
	nt.logger.Info("Connected to %s", addr)
	go rn.listen(nt.localNode)
	nt.newPeer <- rn
}

func (nt *NetTableService) RemoveRemoteNode(addr net.Addr) {
	nt.peers.Remove(addr.String())
}

func (nt *NetTableService) GetPeers() []*RemoteNode {
	var result []*RemoteNode
	for _, key := range nt.peers.Keys() {
		value, res := nt.peers.Get(key)
		if !res {
			continue
		}
		result = append(result, value.(*RemoteNode))
	}
	return result
}

func (nt *NetTableService) isEnoughPeers() bool {
	return nt.peers.Len() >= nt.localNode.discovery.max
}

func (nt *NetTableService) heartbeat() {
	nt.heartbeatTicker = time.Tick(5 * time.Second)
	for {
		select {
		case _, ok := <-nt.heartbeatTicker:
			if !ok {
				break
			}
			i := 0
			for _, key := range nt.peers.Keys() {
				value, res := nt.peers.Get(key)
				if !res {
					continue
				}
				peer := value.(*RemoteNode)
				if time.Since(peer.lastHeartbeat).Seconds() >= 10 {
					nt.logger.Debugf("Closing peer connection. Haven't received a heartbeat for far too long")
					peer.Close()
				}
				i++
				if err := peer.sendHeartBeat(); err != nil {
					nt.logger.Error("error on send heartbeat. %v", err)
				}
			}
		}
	}
}

func (nt *NetTableService) tryConnect(h *net.UDPAddr) error {
	if nt.isEnoughPeers() {
		return errors.New("Will not connected. Reached max.")
	}
	rn, err := connect(h, nt.localNode)
	if err != nil {
		nt.addToBlackList(h)
		return err
	}
	nt.logger.Debug("adding remote node...")
	nt.AddRemoteNode(rn)
	return nil
}

//The black list is just a list of nodes we've already tried and or are connected to.
//TODO: Fix that we don't connect to ourselves.
func (nt *NetTableService) addToBlackList(h *net.UDPAddr) {
	nt.blackList.Add(h.String(), 0)
}

func (nt *NetTableService) isBlackListed(h *net.UDPAddr) bool {
	return nt.blackList.Contains(h)
}