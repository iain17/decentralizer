package upnp

import "testing"

func TestGetLocalIP(t *testing.T) {
	ips, err := GetLocalIPs()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ips: %s", ips)
}

func TestGetLocalInternetIP(t *testing.T) {
	ip, err := GetLocalIntenetIp()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("local internet ip: %s", ip)
}

func TestPortMap(t *testing.T) {
	u := new(Upnp)
	err := u.AddPortMapping(63010, 63010, "udp")
	if err != nil {
		t.Fatal(err)
	}
}
