package main

import (
	"github.com/iain17/discovery"
	"github.com/iain17/discovery/network"
	"github.com/iain17/logger"
	"time"
	"os"
	"context"
	"fmt"
	signal "os/signal"
)

const testNetwork = "2d2d2d2d2d424547494e20525341205055424c4943204b45592d2d2d2d2d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b4341514541736364546c7a386669314338504a38436e386c380a493245726534494e79424264663679706d694e64794554554d34304f4338462b44355376775166594d51514d614161412b65527a715a7279354a785439596a660a7642467a3873313134597a796a6b743853463666434534686778555156564e30326e416e396c4d525164683567433859513641724e7a43774949316a7a4a6d670a45506154563757552f4b6e4d6753476b583941754c5431692b57552b5174695158687964676745684835506f7855514c4b58325766434769476b6871317a41370a776c5559564f78615763553259596a714f657456683063367646476b53476b4263794d6e725652445771705954533451582b38736837444e796f326b463632330a49464f682f31364f784354684f4c2f674366666c455a767862464e612b6f50754f446463456b36326658486b71484947613267353178324f665377752f474c640a4c514944415141420a2d2d2d2d2d454e4420525341205055424c4943204b45592d2d2d2d2d0a"

func init() {
	logger.AddOutput(logger.Stdout{
		MinLevel: logger.DEBUG, //logger.DEBUG,
		Colored:  true,
	})
}

func main(){
	//newNetwork()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())
	go func(){
		<-c
		fmt.Println("halting example...")
		cancel()
	}()

	n, err := network.Unmarshal(testNetwork)
	if err != nil {
		panic(err)
	}
	host, _ := os.Hostname()
	d, err := discovery.New(ctx, n, 10, nil,false, map[string]string {
		"host": host,
		"example": "yes",
	})
	d.SetNetworkMessage("hey cool:"+time.Now().UTC().String())
	if err != nil {
		panic(err)
	}

	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			d.Stop()
			return
		case <-ticker:
			d.LocalNode.SetInfo("Updated", time.Now().Format(time.RFC822))
			logger.Info("peers:")
			for _, peer := range d.WaitForPeers(1, 0*time.Second) {
				logger.Infof("%s", peer)
			}
			logger.Infof("%v", d.GetNetworkMessages())
		}
	}
}

func newNetwork() {
	n, err := network.New()
	if err != nil {
		panic(err)
	}
	logger.Info("New network:")
	logger.Infof("Private key: %s", n.MarshalFromPrivateKey())
	logger.Infof("Public key: %s", n.Marshal())

}