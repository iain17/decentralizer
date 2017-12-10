#include "StdInc.h"


pb::SessionInfo* DNSessionToPBSession(DNSessionInfo * dnInfo) {
	pb::SessionInfo* result = new pb::SessionInfo();
	result->set_dnid(dnInfo->dnId);
	result->set_pid(dnInfo->pId);

	result->set_sessionid(dnInfo->sessionId);
	result->set_type(dnInfo->type);
	result->set_name(dnInfo->name);
	result->set_address(dnInfo->address);
	result->set_port(dnInfo->port);
	::google::protobuf::Map< ::std::string, ::std::string > pbDetails = result->details();
	for (auto const &ent1 : dnInfo->details) {
		pbDetails[ent1.first] = ent1.second;
	}
	return result;
}

DNSessionInfo* PBSessionToDNSession(SessionInfo * pbInfo) {
	DNSessionInfo* result = new DNSessionInfo();
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

LIBDN_API DNAsync<DNUpsertSessionResult>* LIBDN_CALL DN_UpsertSession(DNSessionInfo * info) {
	//build request.
	pb::SessionInfo* sessInfo = DNSessionToPBSession(info);
	RPCUpsertSessionRequest* request = new RPCUpsertSessionRequest();
	request->set_allocated_info(sessInfo);

	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_upsertsessionrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);

	//Set callback.
	NPAsyncImpl<DNUpsertSessionResult>* result = new NPAsyncImpl<DNUpsertSessionResult>();
	async->SetCallback([](DNAsync<RPCMessage>* async) {
		NPAsyncImpl<DNUpsertSessionResult>* asyncResult = (NPAsyncImpl<DNUpsertSessionResult>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		const RPCUpsertSessionResponse& reply = message->upsertsessionresponse();

		DNUpsertSessionResult* result = new DNUpsertSessionResult();
		result->result = reply.result();
		result->sessionId = (DNSID)reply.sessionid();
		asyncResult->SetResult(result);
	}, result);

	return result;
}

LIBDN_API DNAsync<bool>*LIBDN_CALL DN_DeleteSession(DNSID sid) {
	//build request.
	RPCDeleteSessionRequest* request = new RPCDeleteSessionRequest();
	request->set_sessionid(sid);

	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_deletesessionrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);

	//Set callback.
	NPAsyncImpl<bool>* result = new NPAsyncImpl<bool>();
	async->SetCallback([](DNAsync<RPCMessage>* async) {
		NPAsyncImpl<bool>* asyncResult = (NPAsyncImpl<bool>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		auto reply = message->deletesessionresponse();
		bool result = reply.result();
		asyncResult->SetResult(&result);
	}, result);

	return result;
}

LIBDN_API DNAsync<std::vector<DNSID>>* LIBDN_CALL DN_GetSessionIds(uint32_t type, const char* key, const char* value)
{
	//build request.
	RPCSessionIdsRequest* request = new RPCSessionIdsRequest();
	request->set_type(type);
	request->set_key(key);
	request->set_value(value);

	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_sessionidsrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);

	//Set callback.
	auto result = new NPAsyncImpl<std::vector<DNSID>>();
	async->SetCallback([](DNAsync<RPCMessage>* async) {
		auto asyncResult = (NPAsyncImpl<std::vector<DNSID>>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		auto reply = message->sessionidsresponse();
		std::vector<DNSID>* result = new std::vector<DNSID>();
		for (::google::protobuf::uint64 sessionId : reply.sessionids()) {
			result->push_back((DNSID)sessionId);
		}
		asyncResult->SetResult(result);
	}, result);

	return result;
}

LIBDN_API DNAsync<DNSessionInfo>* LIBDN_CALL DN_GetSession(DNSID sessionId) {
	//build request.
	RPCGetSessionRequest* request = new RPCGetSessionRequest();
	request->set_sessionid(sessionId);

	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_getsessionrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);

	//Set callback.
	auto result = new NPAsyncImpl<DNSessionInfo>();
	async->SetCallback([](DNAsync<RPCMessage>* async) {
		auto asyncResult = (NPAsyncImpl<DNSessionInfo>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		auto reply = message->getsessionresponse();
		pb::SessionInfo info = reply.result();
		auto session = PBSessionToDNSession(&info);
		asyncResult->SetResult(session);
	}, result);

	return result;
}