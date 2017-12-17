#include "StdInc.h"
#include "Conversions.h"

namespace libdn {
	pb::Session* DNSessionToPBSession(libdn::Session * dnInfo) {
		pb::Session* result = new pb::Session();
		result->set_dnid(dnInfo->dnId);
		result->set_pid(dnInfo->pId);

		result->set_sessionid(dnInfo->sessionId);
		result->set_type(dnInfo->type);
		result->set_name(dnInfo->name);
		result->set_address(dnInfo->address);
		result->set_port(dnInfo->port);
		auto pbDetails = result->mutable_details();

		for (auto const &ent1 : dnInfo->details) {
			pbDetails->insert(google::protobuf::MapPair<std::string, std::string>(ent1.first, ent1.second));
		}
		return result;
	}

	libdn::Session* PBSessionToDNSession(pb::Session * pbInfo) {
		libdn::Session* result = new libdn::Session();
		result->dnId = pbInfo->dnid();
		result->pId = pbInfo->pid();

		result->sessionId = pbInfo->sessionid();
		result->type = pbInfo->type();
		result->name = pbInfo->name();
		result->address = pbInfo->address();
		result->port = pbInfo->port();
		for (auto const &ent1 : pbInfo->details()) {
			result->details[ent1.first] = ent1.second;
		}
		return result;
	}

	pb::Peer* DNPeerToPBPeer(libdn::Peer * dnInfo) {
		pb::Peer* result = new pb::Peer();
		result->set_dnid(dnInfo->dnId);
		result->set_pid(dnInfo->pId);
		auto pbDetails = result->mutable_details();
		for (auto const &ent1 : dnInfo->details) {
			pbDetails->insert(google::protobuf::MapPair<std::string, std::string>(ent1.first, ent1.second));
		}
		return result;
	}

	libdn::Peer* PBPeerToDNPeer(pb::Peer * pbInfo) {
		libdn::Peer* result = new libdn::Peer();
		result->dnId = pbInfo->dnid();
		result->pId = pbInfo->pid();
		for (auto const &ent1 : pbInfo->details()) {
			result->details[ent1.first] = ent1.second;
		}
		return result;
	}

}