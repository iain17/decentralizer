package p2p

import (
	corenet "github.com/ipfs/go-ipfs/core/corenet"
	"time"
	logger "github.com/Sirupsen/logrus"
	"crypto/sha1"
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

func (s *p2pApp) search() error {
	for {
		s.broadcastUs()

		logger.Debug("Find others.")
		res, err := s.p2p.Node.Routing.GetValues(s.p2p.Ctx, s.hash, 2)
		if err != nil {
			logger.Error(err)
		}

		for _, a := range res {
			logger.Debugf("From %s says %s", a.From.Pretty(), string(a.Val))

			con, err := corenet.Dial(s.p2p.Node, a.From, s.protocol)
			if err != nil {
				logger.Warn(err)
			}
			logger.Debug(con)

		}

		time.Sleep(1 * time.Second)
	}
}

func (s *p2pApp) broadcastUs() {
	logger.Debug("broadcastUs")
	t1 := time.Now()
	err := s.p2p.Node.Routing.PutValue(s.p2p.Ctx, s.hash, []byte(t1.Format(time.UnixDate)))
	if err != nil {
		logger.Error(err)
	}
}