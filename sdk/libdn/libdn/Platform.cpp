#include "StdInc.h"

LIBDN_API DNHealthResult* LIBDN_CALL DN_Health() {
	//build request.
	HealthRequest* request = new HealthRequest();
	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_healthrequest(request);
	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);
	
	DNHealthResult* result = new DNHealthResult();
	async->Wait(30 * 1000);
	if (async->HasCompleted()) {
		RPCMessage* message = async->GetResult();
		const HealthReply& reply = message->healthreply();
		result->message = reply.message();
		result->ready = reply.ready();
	}
	async->Free();
	return result;
}

/*
LIBDN_API DNAsync<DNHealthResult>* LIBDN_CALL DN_Health() {
	//build request.
	HealthRequest* request = new HealthRequest();
	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_healthrequest(request);
	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsync(msg);

	//Set callback.
	NPAsyncImpl<DNHealthResult>* result = new NPAsyncImpl<DNHealthResult>();
	async->SetTimeoutCallback([](DNAsync<RPCMessage>* async) {
		NPAsyncImpl<DNHealthResult>* asyncResult = (NPAsyncImpl<DNHealthResult>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		const HealthReply& reply = message->healthreply();

		DNHealthResult* result = new DNHealthResult();
		result->message = reply.message();
		result->ready = reply.ready();
		asyncResult->SetResult(result);
	}, 30 * 1000);

	//clear up
	return result;
}
*/