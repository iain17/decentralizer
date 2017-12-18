package tcp

import (
	"testing"

	tpt "gx/ipfs/QmQGRkVSA5vTXt9WpJ6nBGnrvq9zRNsfjoNPq8uQrhnBoq/go-libp2p-transport"
	utils "gx/ipfs/QmQGRkVSA5vTXt9WpJ6nBGnrvq9zRNsfjoNPq8uQrhnBoq/go-libp2p-transport/test"
	ma "gx/ipfs/QmW8s4zTsUoX1Q6CeYxVKPyqSKbF7H1YDUyTostBtZ8DaG/go-multiaddr"
)

func TestTcpTransport(t *testing.T) {
	ta := NewTCPTransport()
	tb := NewTCPTransport()

	zero := "/ip4/127.0.0.1/tcp/0"
	utils.SubtestTransport(t, ta, tb, zero)
}

func TestTcpTransportCantListenUtp(t *testing.T) {
	utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/0/utp")
	if err != nil {
		t.Fatal(err)
	}

	tpt := NewTCPTransport()
	_, err = tpt.Listen(utpa)
	if err == nil {
		t.Fatal("shouldnt be able to listen on utp addr with tcp transport")
	}
}

func TestCorrectIPVersionMatching(t *testing.T) {
	ta := NewTCPTransport()

	addr4, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		t.Fatal(err)
	}
	addr6, err := ma.NewMultiaddr("/ip6/::1/tcp/0")
	if err != nil {
		t.Fatal(err)
	}

	d4, err := ta.Dialer(addr4, tpt.ReuseportOpt(true))
	if err != nil {
		t.Fatal(err)
	}

	d6, err := ta.Dialer(addr6, tpt.ReuseportOpt(true))
	if err != nil {
		t.Fatal(err)
	}

	if d4.Matches(addr6) {
		t.Fatal("tcp4 dialer should not match ipv6 address")
	}

	if d6.Matches(addr4) {
		t.Fatal("tcp4 dialer should not match ipv6 address")
	}
}
