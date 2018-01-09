#pragma once
namespace libdn {
	pb::Session DNSessionToPBSession(libdn::Session dnInfo);
	libdn::Session PBSessionToDNSession(pb::Session pbInfo);
	pb::Peer DNPeerToPBPeer(libdn::Peer dnInfo);
	libdn::Peer PBPeerToDNPeer(pb::Peer pbInfo);
}