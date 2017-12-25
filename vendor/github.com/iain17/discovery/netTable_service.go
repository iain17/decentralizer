package discovery

import (
	"context"
	"github.com/iain17/discovery/pb"
	"github.com/iain17/logger"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/golang-lru"
	"errors"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
)

type NetTableService struct {
	localNode *LocalNode
	context   context.Context
	newConn   chan *net.UDPAddr

	blackList *lru.Cache
	seen      *lru.Cache
	peers     *lru.Cache

	heartbeatTicker <-chan time.Time

	logger *logger.Logger
}

func (nt *NetTableService) Init(ctx context.Context, ln *LocalNode) error {
	nt.logger = logger.New("NetTable")
	nt.localNode = ln
	nt.context = ctx
	nt.newConn = make(chan *net.UDPAddr, CONCCURENT_NEW_CONNECTION*2)
	var err error
	nt.blackList, err = lru.New(1000)
	if err != nil {
		return err
	}
	nt.seen, err = lru.New(1000)
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
		go nt.processNewConnection()
	}
	//Send a heart beat to the peers we are connected to
	go nt.heartbeat()
	return nil
}

func (nt *NetTableService) processNewConnection() {
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
	key := addr.String()
	if nt.seen.Contains(key) {
		return
	}
	nt.seen.Add(key, true)
	nt.logger.Debugf("new potential peer %q discovered", addr)
	nt.newConn <- addr
}

func (nt *NetTableService) AddRemoteNode(rn *RemoteNode) {
	addr := rn.conn.RemoteAddr().String()
	nt.peers.Add(addr, rn)
	nt.logger.Infof("Connected to %s", addr)
	go rn.listen(nt.localNode)
	err := nt.Save()
	if err != nil {
		nt.logger.Warningf("Error saving peers: %s", err)
	}
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
	nt.heartbeatTicker = time.Tick(HEARTBEAT_DELAY * time.Second)
	for {
		select {
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
	return nt.blackList.Contains(h) || (h.IP.String() == nt.localNode.ip && h.Port == nt.localNode.port)
}
