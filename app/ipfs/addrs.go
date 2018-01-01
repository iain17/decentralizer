package ipfs

import (
	"time"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	"gx/ipfs/QmX3U3YXCQ6UYBxq2LVWF8dARS1hPUTEYLrSx654Qyxyw6/go-multiaddr-net"
	"net"
	"errors"
	"github.com/hashicorp/golang-lru"
	"fmt"
	"github.com/iain17/logger"
)

var reachableCache = setupReachableCache()

func setupReachableCache() *lru.Cache {
	reachable, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	return reachable
}

func isReachable(host string) bool {
	conn, err := net.DialTimeout("tcp", host, 300*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func isPrivate(IP net.IP) error {
	var err error
	private := false
	if IP == nil {
		err = errors.New("invalid ip")
	} else {
		if !IP.IsGlobalUnicast() {
			return errors.New("multicast or loopback")
		}
		_, privateIPV61BitBlock, _ := net.ParseCIDR("2aa1::1/32")
		_, privateIPV62BitBlock, _ := net.ParseCIDR("fc00::/7")
		_, privateNineBitBlock, _ := net.ParseCIDR("9.0.0.0/8")
		_, privateLocalBitBlock, _ := net.ParseCIDR("127.0.0.0/8")
		_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
		_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
		_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
		private = privateIPV61BitBlock.Contains(IP) || privateIPV62BitBlock.Contains(IP) || privateNineBitBlock.Contains(IP) || privateLocalBitBlock.Contains(IP) || private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
	}
	if private {
		err = errors.New("private ip")
	}
	return err
}

func resolveIPFromNet(addr net.Addr) net.IP {
	switch addr.(type) {
	//case *net.IPAddr:
	//	return addr.(*net.IPAddr).IP
	case *net.TCPAddr:
		return addr.(*net.TCPAddr).IP
	case *net.UDPAddr:
		return addr.(*net.UDPAddr).IP
	}
	return nil
}

func checkAddr(addr ma.Multiaddr, acceptPrivate bool, onlyPlain bool, debug bool) bool {
	protocols := addr.Protocols()
	if len(protocols) == 0 {
		return false
	}
	isPlain := true
	isTcp := false
	for _, protocol := range protocols {
		if protocol.Code == ma.P_IPFS {
			isPlain = false
		}
		if protocol.Code == ma.P_TCP {
			isTcp = true
		}
	}
	if onlyPlain && !isPlain {
		return false
	}
	logger.Debugf("%s is isPlain(%t) isTcp(%t)", addr.String(), isPlain, isTcp)
	if isPlain {
		netAddr, err := manet.ToNetAddr(addr)
		if err != nil {
			return false
		}
		netAddrText := netAddr.String()

		ip := resolveIPFromNet(netAddr)
		if !acceptPrivate {
			err = isPrivate(ip)
			if err != nil {
				logger.Debugf("%s is private", netAddrText)
				return false
			}
		}
		if isTcp {
			if !isReachable(netAddrText) {
				logger.Debugf("%s is NOT reachable", netAddrText)
				return false
			}
		}
	}
	logger.Debugf("%s is reachable!", addr.String())
	return true
}

//Little cache layer
func IsAddrReachable(addr ma.Multiaddr, acceptPrivate bool, onlyPlain bool, debug bool) bool {
	if addr == nil {
		return false
	}
	key := fmt.Sprintf("%s/%t/%t", addr.String(), acceptPrivate, onlyPlain)
	if reachableCache.Contains(key) {
		value, ok := reachableCache.Get(key)
		if ok {
			logger.Debugf("returning cache: %s = %t", key, value.(bool))
			return value.(bool)
		}
	}
	value := checkAddr(addr, acceptPrivate, onlyPlain, debug)
	reachableCache.Add(key, value)
	return value
}

//Only plain means. Only tcp and udp connections.
func FilterNonReachableAddrs(addrs []ma.Multiaddr, acceptPrivate bool, onlyPlain bool, debug bool) []ma.Multiaddr {
	i := 0
	//if addrs == nil {
	//	return addrs
	//}
	var result []ma.Multiaddr
	for _, addr := range addrs {
		if !IsAddrReachable(addr, acceptPrivate, onlyPlain, debug) {
			//addrs[i] = addrs[len(addrs)-1] // Replace it with the last one.
			//addrs = addrs[:len(addrs)-1]   // Chop off the last one.
			continue
		}
		result = append(result, addr)
		i++
	}
	return result
}
