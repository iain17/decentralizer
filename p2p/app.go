package p2p

import (
	corenet "github.com/ipfs/go-ipfs/core/corenet"
	"time"
	logger "github.com/Sirupsen/logrus"
	"crypto/sha1"
	"fmt"
)

const interval = 30

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
		protocol: "/app/whyrusleeping",
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
		t := time.Now().UTC()
		key := getKey(t, s.hash)

		go s.broadcastUs(key)

		logger.Debug("Find others. %s", key)
		res, err := s.p2p.Node.Routing.GetValues(s.p2p.Ctx, key, 10)
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

		t = time.Now().UTC()
		secLeft := interval - (t.Second() % interval)
		fmt.Printf("%d", secLeft)
		time.Sleep(time.Duration(secLeft) * time.Second)
	}
}

func (s *p2pApp) broadcastUs(key string) {
	logger.Debugf("broadcastUs")
	err := s.p2p.Node.Routing.PutValue(s.p2p.Ctx, key, []byte("test"))
	if err != nil {
		logger.Error(err)
	}
}

func getKey(t time.Time, hash string) string {
	//TODO: Optimize. Sprintf is slow and unnecessary here.
	return fmt.Sprintf("%s-%d-%02d-%02dT%02d:%02d:%d", hash, t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second() / interval)
}