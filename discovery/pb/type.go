package pb

type MessageType byte

const (
	HearBeatMessage MessageType = iota
	PeerInfoMessage
	TransferMessage
)