package discovery

import (
	"time"
	"sync"
	"context"
	"github.com/op/go-logging"
	"net"
)

type NetTableService struct {
	localNode *LocalNode
	context context.Context
	//waitGroup sync.WaitGroup
	newConn chan *net.UDPAddr

	lock      sync.RWMutex
	blackList map[string]time.Time
	peers     map[string]*RemoteNode

	heartbeatTicker <-chan time.Time

	logger *logging.Logger
}

func (nt *NetTableService) Init(ctx context.Context, ln *LocalNode) error {
	nt.logger = logging.MustGetLogger("NetTable")
	nt.localNode = ln
	nt.context = ctx
	nt.newConn = make(chan *net.UDPAddr, 10)
	nt.blackList = make(map[string]time.Time)
	nt.peers = make(map[string]*RemoteNode)
	nt.Run()
	return nil
}

func (nt *NetTableService) Run() error {
	//Spawn 10 workers.
	for i := 0; i < 10; i++ {
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
	nt.lock.Lock()
	for _, peer := range nt.peers {
		peer.Close()
	}
	nt.lock.Unlock()
}

func (nt *NetTableService) GetNewConnChan() chan<- *net.UDPAddr {
	return nt.newConn
}

func (nt *NetTableService) AddRemoteNode(rn *RemoteNode) {
	addr := rn.conn.RemoteAddr().String()

	nt.lock.Lock()
	nt.peers[addr] = rn
	nt.lock.Unlock()

	nt.logger.Info("Connected to %s", addr)
	go rn.listen(nt.localNode)
}

func (nt *NetTableService) RemoveRemoteNode(addr net.Addr) {
	nt.lock.Lock()
	delete(nt.peers, addr.String())
	nt.lock.Unlock()
}

func (nt *NetTableService) heartbeat() {
	nt.heartbeatTicker = time.Tick(5 * time.Second)
	for {
		select {
		case _, ok := <-nt.heartbeatTicker:
			if !ok {
				break
			}
			nt.lock.Lock()
			for _, peer := range nt.peers {
				if err := peer.sendHeartBeat(); err != nil {
					nt.logger.Error("error on send heartbeat. %v", err)
				}
			}
			nt.lock.Unlock()
		}
	}
}

func (nt *NetTableService) tryConnect(h *net.UDPAddr) error {
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
	nt.lock.Lock()
	nt.blackList[h.String()] = time.Now()
	nt.lock.Unlock()
}

func (nt *NetTableService) isBlackListed(h *net.UDPAddr) bool {
	nt.lock.Lock()
	_, ok := nt.blackList[h.String()]
	nt.lock.Unlock()
	return ok
}