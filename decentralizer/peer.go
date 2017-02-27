package decentralizer

type Peer struct {
	// The IP address of the peer
	// Required: true
	// Read Only: true
	IP string `json:"ip"`

	// The port of the peer's service.
	// Required: true
	// Read Only: true
	Port int32 `json:"port"`

	// The port internally used by decentralizer to
	RPCPort int32 `json:"rpcPort"`

	// Anything. All the details the peer wants to mention to others.
	Details map[string]string `json:"details,omitempty"`
}

func NewPeer(ip string, RPCPort, port int32, details map[string]string) *Peer {
	return &Peer{
		IP: ip,
		RPCPort: RPCPort,
		Port: port,
		Details: details,
	}
}
