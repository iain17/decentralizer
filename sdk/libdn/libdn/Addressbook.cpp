#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<bool>*LIBDN_CALL UpsertPeer(libdn::Peer * peer) {
		auto result = new Promise<bool>([peer](Promise<bool>* promise) {
			// Data we are sending to the server.
			pb::RPCUpsertPeerRequest request;
			request.set_allocated_peer(DNPeerToPBPeer(peer));

			// Container for the data we expect from the server.
			pb::RPCUpsertPeerResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->UpsertPeer(ctx, request, &reply);

			UpsertSessionResult result;
			if (!status.ok()) {
				promise->reject(fmt::format("[Could not upsert peer] {0}: {1}", status.error_code(), status.error_message().c_str()));
				return false;
			}
			return true;
		});
		return result;
	}

	LIBDN_API Promise<int>*LIBDN_CALL GetNumPeers(const char * key, const char * value) {
		auto result = new Promise<int>([key, value](Promise<int>* promise) {
			//build request.
			pb::RPCGetPeerIdsRequest request;
			request.set_key(key);
			request.set_value(value);

			// Container for the data we expect from the server.
			pb::RPCGetPeerIdsResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPeerIds(ctx, request, &reply);

			if (status.ok()) {
				context.peers = reply.peerids();
				int size = context.peers.size();
				return size;
			} else {
				promise->reject(fmt::format("[Could not get peer ids] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}

			return 0;
		});
		return result;
	}

	//If pId is set to "self" it will automatically resolve to the local peer id.
	//If both are given. Perference is given to pId.
	LIBDN_API Promise<Peer*>*LIBDN_CALL GetPeerById(DNID dId, PeerID& pId) {
		auto result = new Promise<Peer*>([dId, pId](Promise<Peer*>* promise) {
			//build request.
			pb::RPCGetPeerRequest request;
			request.set_dnid(dId);
			request.set_pid(pId);

			// Container for the data we expect from the server.
			pb::RPCGetPeerResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPeer(ctx, request, &reply);

			if (!status.ok()) {
				promise->reject(fmt::format("[Could not get peer] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}
			auto peer = reply.peer();
			return PBPeerToDNPeer(&peer);
		});
		return result;
	}

	LIBDN_API Peer* LIBDN_CALL GetSelf() {
		std::string pid = "self";
		auto request = GetPeerById(0, pid);
		return request->get();
	}

	LIBDN_API Peer *LIBDN_CALL GetPeerByIndex(int index) {
		if (index > context.peers.size() - 1) {
			return NULL;
		}
		std::string pid = context.peers.Get(index);
		auto req = GetPeerById(0, pid);
		req->fail([](std::string reason) {
			Log_Print(reason.c_str());
		});
		if (req->wait()) {
			return req->get();
		}
		return NULL;
	}

	LIBDN_API PeerID* LIBDN_CALL ResolveDecentralizedId(DNID dId) {
		std::string empty = "";
		auto request = GetPeerById(dId, empty);
		request->fail([](const char* reason) {
			Log_Print("Could not resolve DecentralizedId: %s", reason);
		});
		if (request->wait()) {
			PeerID* peerId = (PeerID*) new std::string(request->get()->pId);
			return peerId;
		}
		return new std::string("");
	}
}