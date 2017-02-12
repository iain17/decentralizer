package p2p

import (
	corenet "github.com/ipfs/go-ipfs/core/corenet"
	logger "github.com/Sirupsen/logrus"
	"crypto/sha1"
	"github.com/ipfs/go-ipfs/blocks"
	floodsub "github.com/libp2p/go-floodsub"
	"io"
	"golang.org/x/net/context"
)

type p2pApp struct {
	name string
	hash string
	protocol string
	p2p *p2p
}

func newP2pApp(p2p *p2p, name string) *p2pApp {

	hash, err := getHash(name)

	if err != nil {
		logger.Error(err)
		return nil
	}

	app := &p2pApp{
		hash: hash,
		name: name,
		p2p: p2p,
		protocol: "/app/"+hash,
	}
	logger.Infof("Registering app with name '%s'\n", name)
	go app.listen()
	go app.search()
	return app
}

func getHash(identifier string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(identifier))
	if err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}

func (s *p2pApp) listen() error {
	list, err := corenet.Listen(s.p2p.Node, s.protocol)
	if err != nil {
		return err
	}

	for {
		con, err := list.Accept()
		if err != nil {
			logger.Error(err)
			continue
		}
		defer con.Close()

		logger.Info(con, "Hello! This is whyrusleepings awesome ipfs service")
		logger.Infof("Connection from: %s\n", con.Conn().RemotePeer())


		//select {
		//	case <-s.p2p.Ctx.Done():
		//		return s.p2p.Ctx.Err()
		//}
	}
}

func (s *p2pApp) search() {
	sub, err := s.p2p.Node.Floodsub.Subscribe(s.name)//Topic name is set here.
	if err != nil {
		logger.Error(err)
	}

	out := make(chan interface{})

	go func() {
		defer sub.Cancel()
		defer close(out)

		out <- floodsub.Message{}

		for {
			msg, err := sub.Next(s.p2p.Ctx)
			if err == io.EOF || err == context.Canceled {
				return
			} else if err != nil {
				logger.Error(err)
				return
			}

			out <- msg
		}
	}()

	//discover
	go func() {
		blk := blocks.NewBlock([]byte("floodsub:" + sub.Topic()))
		cid, err := s.p2p.Node.Blocks.AddBlock(blk)
		if err != nil {
			logger.Error("pubsub discovery: ", err)
			return
		}
		logger.Debugf("Discovered: %s", cid)
	}()


}