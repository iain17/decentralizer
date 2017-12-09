package sessionstore

import "github.com/iain17/decentralizer/pb"

func GetId(info *pb.SessionInfo) uint64 {
	return uint64(info.Type)+uint64(info.Port)+info.DId
}