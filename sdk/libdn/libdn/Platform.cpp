#include "StdInc.h"

static void GetHealthCB(DNAsync<RPCMessage>* async) {
	NPAsyncImpl<DNHealthResult>* asyncResult = (NPAsyncImpl<DNHealthResult>*)async->GetUserData();
	RPCMessage* message = async->GetResult();
	const HealthReply& reply = message->healthreply();

	DNHealthResult* result = new DNHealthResult();
	result->message = reply.message();
	result->ready = reply.ready();
	asyncResult->SetResult(result);
}

LIBDN_API DNHealthResult* LIBDN_CALL DN_Health() {
	//build request.
	HealthRequest* request = new HealthRequest();
	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_healthrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsyncCache("health", msg);

	//Set callback.
	NPAsyncImpl<DNHealthResult>* result = new NPAsyncImpl<DNHealthResult>();
	async->SetCallback(GetHealthCB, result);
	async->Wait();

	//clear up
	//request->Clear();

	return result->GetResult();
}