package sessionstore

import "github.com/iain17/decentralizer/pb"

func GetId(info *pb.Session) uint64 {
	return uint64(info.Port)+uint64(info.Address)
}