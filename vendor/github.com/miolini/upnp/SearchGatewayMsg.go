package upnp

import (
	"log"
	"fmt"
	"net"
	"strings"
)

type Gateway struct {
	GatewayName   string //Gateway Name
	Host          string //Gateway ip and port
	DeviceDescUrl string //Gateway device description path
	Cache         string //cache
	ST            string
	USN           string
	deviceType    string //Urn device "urn: schemas-upnp-org: service: WANIPConnection: 1"
	ControlURL    string //Device port mapping request path
	ServiceType   string //Upnp services provide the type of service
}

type SearchGateway struct {
	searchMessage string
	upnp          *Upnp
}

func (this *SearchGateway) Send() (bool, error) {
	this.buildRequest()
	result, err := this.SendMessage()
	if err != nil {
		return false, err
	}
	if result == "" {
		//Overtime
		this.upnp.Active = false
		return false, nil
	}
	this.resolve(result)

	this.upnp.Gateway.ServiceType = "urn:schemas-upnp-org:service:WANIPConnection:1"
	this.upnp.Active = true
	return true, nil
}
func (this *SearchGateway) SendMessage() (result string, err error) {
	//Send broadcast messages to bring the port, formats such as: "239.255.255.250:1900"
	var conn *net.UDPConn
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic err: %s", r)
		}
	}()
	//go func(conn *net.UDPConn) {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			log.Printf("panic in timeout conn err: %s", err)
	//		}
	//	}()
	//	//Timeout to 3 seconds
	//	time.Sleep(time.Second * 3)
	//	if err := conn.Close(); err != nil {
	//		log.Printf("conn close err: %s", err)
	//	}
	//}(conn)
	remotAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		return "", fmt.Errorf("Multicast address format is incorrect err: %s", err)
	}

	locaAddr, err := net.ResolveUDPAddr("udp", this.upnp.LocalHost+":0")
	if err != nil {
		return "", fmt.Errorf("Local IP address is incorrent err: %s", err)
	}
	conn, err = net.ListenUDP("udp", locaAddr)
	if err != nil {
		return "", fmt.Errorf("Listening udp error err: %s", err)
	}
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			log.Printf("conn close err: %s", err)
		}
	}(conn)

	_, err = conn.WriteToUDP([]byte(this.searchMessage), remotAddr)
	if err != nil {
		return "", fmt.Errorf("Sent to a multicast address err: %s", err)
	}

	buf := make([]byte, 1024)

	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return "", fmt.Errorf("Error message received from a multicast address")
	}

	return string(buf[:n]), nil
}

func (this *SearchGateway) buildRequest() {
	this.searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
}

func (this *SearchGateway) resolve(result string) {
	this.upnp.Gateway = &Gateway{}

	lines := strings.Split(result, "\r\n")
	for _, line := range lines {
		//According to a first colon into two strings
		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			this.upnp.Gateway.ST = nameValues[1]
		case "CACHE-CONTROL":
			this.upnp.Gateway.Cache = nameValues[1]
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			this.upnp.Gateway.Host = urls[0]
			this.upnp.Gateway.DeviceDescUrl = "/" + urls[1]
		case "SERVER":
			this.upnp.Gateway.GatewayName = nameValues[1]
		default:
		}
	}
}
