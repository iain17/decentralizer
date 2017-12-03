package decentralizer

import (
	"github.com/iain17/decentralizer/decentralizer/pb"
	logger "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"crypto/sha1"
	"github.com/smallnest/rpcx"
	kcp "github.com/xtaci/kcp-go"
	kcpLog "github.com/smallnest/rpcx/log"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"net"
	"github.com/Sirupsen/logrus"
	"time"
	"io/ioutil"
)

type rpcServer struct {
	d *decentralizer
}

//TODO: Later make these a configurable vars
const cryptKey = "jock"
const cryptSalt = "flora"
var bc kcp.BlockCrypt

func init() {
	//Disable annyoing kcp logging.
	kcpLogger := logrus.New()
	kcpLogger.Out = ioutil.Discard
	kcpLog.SetLogger(kcpLogger)

	var err error
	pass := pbkdf2.Key([]byte(cryptKey), []byte(cryptSalt), 4096, 32, sha1.New)
	bc, err = kcp.NewAESBlockCrypt(pass)
	if err != nil {
		panic(err)
	}
}

/*
- The rpc server is a TCP rpcx-kcp server that is used to exchange messages between nodes.
 */
func (d *decentralizer) listenRpcServer() error {
	conn, _, err := getUdpConn()
	if err != nil {
		return err
	}
	conn.Close()

	port := conn.LocalAddr().(*net.UDPAddr).Port
	sPort := strconv.Itoa(port)
	d.rpcPort = uint16(port)

	//block for encryption
	ln, err := kcp.ListenWithOptions(":"+sPort, bc, 10, 3)
	//ln, err := kcp.ListenWithOptions(":"+sPort, nil, 10, 3)
	if err != nil {
		return err
	}

	server := rpcx.NewServer()
	rpc := &rpcServer{
		d: d,
	}
	server.RegisterName("Decentralizer", rpc)

	go func() {
		server.ServeListener(ln)
	}()
	logger.Infof("RPC server listening at %d", port)
	return nil
}

func (s *rpcServer) GetService(args *pb.GetServiceRequest, reply *pb.GetServiceResponse) error {
	service := s.d.services[args.Hash]
	if service == nil {
		return errors.New("No such service registered under that hash")
	}
	if(args.Me != nil) {
		service.PeerDiscovered(args.Me)
	}
	reply.Result = service.self.Peer
	reply.Peers = service.GetPeers()
	return nil
}

func (s *service) getServiceRequest(addr string) (*pb.GetServiceResponse, error) {
	ss := &rpcx.DirectClientSelector{
		Network: "kcp",
		Address: addr,
		DialTimeout: 300 * time.Millisecond,
	}
	client := rpcx.NewClient(ss)
	//Make it fail fast! The majority of these calls will go nowhere anyways. No point to keep waiting
	client.FailMode = rpcx.Failtry
	client.Retries = 0
	client.Timeout = 300 * time.Millisecond
	client.ReadTimeout = 300 * time.Millisecond
	client.WriteTimeout = 300 * time.Millisecond

	//Encryption
	client.Block = bc
	defer client.Close()

	args := &pb.GetServiceRequest{
		Me: s.self.Peer,
		Hash: s.hash,
	}
	var reply pb.GetServiceResponse
	err := client.Call("Decentralizer.GetService", args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
