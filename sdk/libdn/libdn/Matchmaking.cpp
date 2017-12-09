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

LIBDN_API DNAsync<DNUpsertSessionResult>* LIBDN_CALL DN_UpsertSession(DNSessionInfo * info)
{
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

LIBDN_API DNAsync<bool>*LIBDN_CALL DN_DeleteSession(DNSID sid)
{
	return NULL;
}

LIBDN_API DNAsync<bool>* LIBDN_CALL DN_RefreshSessions(uint32_t type)
{
	return NULL;
}

LIBDN_API DNAsync<DNSID[]>* LIBDN_CALL DN_GetNumSessions(uint32_t type, std::map<std::string, std::string> details)
{
	return NULL;
}

LIBDN_API void LIBDN_CALL DN_GetSessionData(DNSID sessionId, DNSessionInfo* out)
{
	return;
}
