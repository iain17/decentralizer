package discovery

import (
	"context"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/logger"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/golang-lru"
	ttlru "github.com/iain17/kvcache/lttlru"
	"errors"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
	"sync"
	"fmt"
)

type NetTableService struct {
	localNode *LocalNode
	context   context.Context
	newConn   chan *net.UDPAddr

	blackList *ttlru.LruWithTTL
	seen      *ttlru.LruWithTTL
	peers     *lru.Cache
	mutex 	  sync.Mutex

	heartbeatTicker <-chan time.Time

	logger *logger.Logger
}

func (d *NetTableService) String() string {
	return "NetTable"
}

func (nt *NetTableService) init(ctx context.Context) error {
	defer func() {
		if nt.localNode.wg != nil {
			nt.localNode.wg.Done()
		}
	}()
	nt.logger = logger.New(nt.String())
	nt.context = ctx
	nt.newConn = make(chan *net.UDPAddr, BACKLOG_NEW_CONNECTION)
	var err error
	nt.blackList, err = ttlru.NewTTL(1000)
	if err != nil {
		return err
	}
	nt.seen, err = ttlru.NewTTL(1024)
	if err != nil {
		return err
	}
	nt.peers, err = lru.NewWithEvict(nt.localNode.discovery.max, evicted)
	if err != nil {
		return err
	}
	err = nt.Restore()
	if err != nil {
		nt.logger.Warningf("Could not restore previous connected peers: %s", err)
	}
	return nil
}

func evicted(_ interface{}, value interface{}) {
	if node, ok := value.(*RemoteNode); ok {
		node.Close()
	}
}

func (nt *NetTableService) Serve(ctx context.Context) {
	nt.localNode.waitTilCoreReady()
	defer nt.Stop()

	if err := nt.init(ctx); err != nil {
		nt.localNode.lastError = err
		panic(err)
	}
	nt.localNode.waitTilReady()
	//Spawn some workers
	for i := 0; i < CONCCURENT_NEW_CONNECTION; i++ {
		go nt.processNewConnection()
	}
	//Send a heart beat to the peers we are connected to
	nt.heartbeat()
}

func (nt *NetTableService) processNewConnection() {
	for {
		select {
		case <-nt.context.Done():
			return
		case host, ok := <-nt.newConn:
			if !ok {
				return
			}

			key := host.String()
			if _, ok := nt.seen.Get(key); ok {
				continue
			}
			if nt.isBlackListed(host) {
				continue
			}
			nt.seen.AddWithTTL(key, true, 15 * time.Minute)
			nt.logger.Debugf("new potential peer %q discovered", host)

			if err := nt.tryConnect(host); err != nil {
				nt.logger.Debugf("unable connect %s err: %s", host, err)
			}
		}
	}
}

func (nt *NetTableService) Stop() {
	for _, peer := range nt.GetPeers() {
		peer.Close()
	}
}

func (nt *NetTableService) Save() error {
	var peers []*pb.DPeer
	for _, peer := range nt.GetPeers() {
		ipPort := strings.Split(peer.conn.RemoteAddr().String(), ":")
		if len(ipPort) != 2 {
			continue
		}
		port, err := strconv.Atoi(ipPort[1])
		if err != nil {
			continue
		}
		peers = append(peers, &pb.DPeer{
			Ip:   ipPort[0],
			Port: int32(port),
		})
	}
	data, err := proto.Marshal(&pb.Peers{
		Peers: peers,
	})
	if err != nil {
		return err
	}
	nt.logger.Infof("Saved %d peers to net table file", len(peers))
	return configPath.QueryCacheFolder().WriteFile(NET_TABLE_FILE, data)
}

func (nt *NetTableService) Restore() error {
	file, err := configPath.QueryCacheFolder().Open(NET_TABLE_FILE)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	peers := &pb.Peers{}
	err = proto.Unmarshal(data, peers)
	if err != nil {
		return err
	}
	nt.logger.Infof("Restored %d peers from net table file", len(peers.Peers))
	for _, peer := range peers.Peers {
		nt.Discovered(&net.UDPAddr{
			IP:   net.ParseIP(peer.Ip),
			Port: int(peer.Port),
		})
	}
	return nil
}

func (nt *NetTableService) Discovered(addr *net.UDPAddr) {
	if addr.IP.String() == nt.localNode.ip && addr.Port == nt.localNode.port {
		return
	}
	nt.newConn <- addr
}

func (nt *NetTableService) AddRemoteNode(rn *RemoteNode) {
	nt.mutex.Lock()
	defer nt.mutex.Unlock()
	if rn.id == "" || nt.peers.Contains(rn.id) {
		nt.logger.Warningf("Already connected to %s", rn.id)
		rn.conn.Close()
		return
	}

	nt.peers.Add(rn.id, rn)

	addr := rn.conn.RemoteAddr().String()
	nt.logger.Infof("Connected to %s: %s", addr, rn.id)
	go rn.listen(nt.localNode)
	err := nt.Save()
	if err != nil {
		nt.logger.Warningf("Error saving peers: %s", err)
	}
}

func (nt *NetTableService) isConnected(id string) bool {
	return nt.peers.Contains(id)
}

func (nt *NetTableService) RemoveRemoteNode(rn *RemoteNode) {
	nt.peers.Remove(rn.id)
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
	nt.heartbeatTicker = time.Tick(HEARTBEAT_DELAY * time.Second)
	for {
		select {
		case <-nt.context.Done():
			return
		case _, ok := <-nt.heartbeatTicker:
			if !ok {
				break
			}
			i := 0
			for _, peer := range nt.GetPeers() {
				if time.Since(peer.lastHeartbeat).Seconds() >= HEARTBEAT_DELAY*2 {
					nt.logger.Debugf("Closing peer connection. Haven't received a heartbeat for far too long")
					peer.Close()
					continue
				}
				i++
				if err := peer.sendHeartBeat(); err != nil {
					peer.Close()
					nt.logger.Errorf("error on send heartbeat. %v", err)
				}
			}
		}
	}
}

func (nt *NetTableService) tryConnect(h *net.UDPAddr) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tryConnect panic: %s", r)
		}
	}()
	if nt.isEnoughPeers() {
		return errors.New("will not connected. Reached max")
	}
	var rn *RemoteNode
	rn, err = connect(h, nt.localNode)
	if err != nil {
		nt.addToBlackList(h)
		return err
	}
	nt.logger.Debug("adding remote node...")
	nt.AddRemoteNode(rn)
	return err
}

//The black list is just a list of nodes we've already tried and or are connected to.
func (nt *NetTableService) addToBlackList(h *net.UDPAddr) {
	nt.blackList.AddWithTTL(h.String(), 0, 10 * time.Minute)
}

func (nt *NetTableService) isBlackListed(h *net.UDPAddr) bool {
	return nt.blackList.Contains(h) || (h.IP.String() == nt.localNode.ip && h.Port == nt.localNode.port)
}
