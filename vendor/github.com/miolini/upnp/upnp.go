package upnp

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

/*
 * Obtain Gateway
 */

//All the ports management
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

//Adding a port mapping record
//Only the mapping management
func (this *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		this.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := this.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		this.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	this.mappingPorts[protocol] = [][]int{one, two}
}

//Delete a mapping record
//Only the mapping management
func (this *MappingPortStruct) delMapping(remotePort int, protocol string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := this.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			//Mapping to be deleted
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	this.mappingPorts = tmp.mappingPorts
}
func (this *MappingPortStruct) GetAllMapping() map[string][][]int {
	return this.mappingPorts
}

type Upnp struct {
	Active             bool              //This protocol is available upnp
	LocalHost          string            //The machine ip address
	GatewayInsideIP    string            //LAN gateway ip
	GatewayOutsideIP   string            //Gateway public network ip
	OutsideMappingPort map[string]int    //Mapping external port
	InsideMappingPort  map[string]int    //Mapping the local port
	Gateway            *Gateway          //Gateway Information
	CtrlUrl            string            //Control request url
	MappingPort        MappingPortStruct //Mapping has been added {"TCP":[1990],"UDP":[1991]}
}

// Get the ip address of the local network
// Get the LAN gateway ip
func (this *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp module being given", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if this.LocalHost == "" {
		this.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
			// mappingPorts: map[string][][]int{},
		}
		this.LocalHost, err = GetLocalIntenetIp()
		if err != nil {
			return err
		}
	}
	searchGateway := SearchGateway{upnp: this}
	ok, err := searchGateway.Send()
	if err != nil {
		return err
	} else if ok {
		return nil
	}
	return errors.New("No gateway device")
}

func (this *Upnp) deviceStatus() {

}

//See Device Description, get control request url
func (this *Upnp) deviceDesc() (err error) {
	if this.GatewayInsideIP == "" {
		if err := this.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: this}
	if _, err := device.Send(); err != nil {
		return err
	}
	this.Active = true
	// log.Println("Gain control request url:", this.CtrlUrl)
	return
}

//View the ip address
func (this *Upnp) ExternalIPAddr() (err error) {
	if this.CtrlUrl == "" {
		if err := this.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()
	return nil
	// log.Println("Obtain public network ip address：", this.GatewayOutsideIP)
}

//Adding a port mapping
func (this *Upnp) AddPortMapping(localPort, remotePort int, protocol string) (err error) {
	defer func(err *error) {
		if errTemp := recover(); errTemp != nil {
			*err = fmt.Errorf("panic err: %s", err)
		}
	}(&err)
	if this.GatewayOutsideIP == "" {
		if err := this.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: this}
	if issuccess := addPort.Send(localPort, remotePort, protocol); issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol)
		return nil
	} else {
		this.Active = false
		return errors.New("Adding a port mapping failed")
	}
}

func (this *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: this}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		this.MappingPort.delMapping(remotePort, protocol)
		log.Println("To delete a port mapping: remote:", remotePort)
	}
	return issuccess
}

//Recycling port
func (this *Upnp) Reclaim() {
	mappings := this.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			this.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			this.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

func (this *Upnp) GetAllMapping() map[string][][]int {
	return this.MappingPort.GetAllMapping()
}
