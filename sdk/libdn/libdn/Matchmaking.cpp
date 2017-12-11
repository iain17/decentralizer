#include "StdInc.h"

namespace libdn {
	pb::SessionInfo* DNSessionToPBSession(libdn::SessionInfo * dnInfo) {
		pb::SessionInfo* result = new pb::SessionInfo();
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

			//->insert(ent1.first.c_str(), ent1.second.c_str());
		}
		return result;
	}

	libdn::SessionInfo* PBSessionToDNSession(pb::SessionInfo * pbInfo) {
		libdn::SessionInfo* result = new libdn::SessionInfo();
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

	LIBDN_API Async<UpsertSessionResult>* LIBDN_CALL UpsertSession(SessionInfo * info) {
		//build request.
		pb::SessionInfo* sessInfo = DNSessionToPBSession(info);
		RPCUpsertSessionRequest* request = new RPCUpsertSessionRequest();
		request->set_allocated_info(sessInfo);

		pb::RPCMessage* msg = new pb::RPCMessage();
		msg->set_allocated_upsertsessionrequest(request);

		//Send request.
		Async<RPCMessage>* async = RPC_SendMessageAsync(msg);

		//Set callback.
		AsyncImpl<UpsertSessionResult>* result = new AsyncImpl<UpsertSessionResult>();
		async->SetCallback([](Async<RPCMessage>* async) {
			AsyncImpl<UpsertSessionResult>* asyncResult = (AsyncImpl<UpsertSessionResult>*)async->GetUserData();
			RPCMessage* message = async->GetResult();
			const RPCUpsertSessionResponse& reply = message->upsertsessionresponse();

			UpsertSessionResult* result = new UpsertSessionResult();
			result->result = reply.result();
			result->sessionId = (DNSID)reply.sessionid();
			asyncResult->SetResult(result);
		}, result);

		return result;
	}

	LIBDN_API Async<bool>*LIBDN_CALL DeleteSession(DNSID sid) {
		//build request.
		RPCDeleteSessionRequest* request = new RPCDeleteSessionRequest();
		request->set_sessionid(sid);

		pb::RPCMessage* msg = new pb::RPCMessage();
		msg->set_allocated_deletesessionrequest(request);

		//Send request.
		Async<RPCMessage>* async = RPC_SendMessageAsync(msg);

		//Set callback.
		AsyncImpl<bool>* result = new AsyncImpl<bool>();
		async->SetCallback([](Async<RPCMessage>* async) {
			AsyncImpl<bool>* asyncResult = (AsyncImpl<bool>*)async->GetUserData();
			RPCMessage* message = async->GetResult();
			auto reply = message->deletesessionresponse();
			bool result = reply.result();
			asyncResult->SetResult(&result);
		}, result);

		return result;
	}

	LIBDN_API Async<int>* LIBDN_CALL GetNumSessions(uint32_t type, const char* key, const char* value) {
		//build request.
		RPCSessionIdsRequest* request = new RPCSessionIdsRequest();
		request->set_type(type);
		request->set_key(key);
		request->set_value(value);

		pb::RPCMessage* msg = new pb::RPCMessage();
		msg->set_allocated_sessionidsrequest(request);

		//Send request.
		Async<RPCMessage>* async = RPC_SendMessageAsync(msg);

		//Set callback.
		auto result = new AsyncImpl<int>();
		async->SetCallback([](Async<RPCMessage>* async) {
			auto asyncResult = (AsyncImpl<int>*)async->GetUserData();
			RPCMessage* message = async->GetResult();
			auto reply = message->sessionidsresponse();
			g_dn.sessions = reply.sessionids();
			int size = g_dn.sessions.size() - 1;
			asyncResult->SetResult(&size);
		}, result);

		return result;
	}

	LIBDN_API Async<SessionInfo>* LIBDN_CALL GetSessionBySessionId(DNSID sessionId) {
		//build request.
		RPCGetSessionRequest* request = new RPCGetSessionRequest();
		request->set_sessionid(sessionId);

		pb::RPCMessage* msg = new pb::RPCMessage();
		msg->set_allocated_getsessionrequest(request);

		//Send request.
		char cache[255];
		sprintf(cache, "matchmaking_%lld", sessionId);
		Async<RPCMessage>* async = RPC_SendMessageAsyncCache(cache, msg);

		//Set callback.
		auto result = new AsyncImpl<SessionInfo>();
		async->SetCallback([](Async<RPCMessage>* async) {
			auto asyncResult = (AsyncImpl<SessionInfo>*)async->GetUserData();
			RPCMessage* message = async->GetResult();
			auto reply = message->getsessionresponse();
			pb::SessionInfo info = reply.result();
			auto session = PBSessionToDNSession(&info);
			asyncResult->SetResult(session);
		}, result);

		return result;
	}


	LIBDN_API SessionInfo* LIBDN_CALL GetSessionByIndex(int index) {
		if (index > MAX_SESSIONS || index > g_dn.sessions.size() - 1) {
			return NULL;
		}
		auto sessionRequest = GetSessionBySessionId(g_dn.sessions.Get(index));
		if (sessionRequest->Wait(7500) != nullptr) {
			return sessionRequest->GetResult();
		}
		return NULL;
	}
}