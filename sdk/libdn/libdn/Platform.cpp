#include "StdInc.h"


LIBDN_API DNHealthResult* LIBDN_CALL DN_Health() {
	//build request.
	RPCHealthRequest* request = new RPCHealthRequest();
	pb::RPCMessage* msg = new pb::RPCMessage();
	msg->set_allocated_healthrequest(request);

	//Send request.
	DNAsync<RPCMessage>* async = RPC_SendMessageAsyncCache("health", msg);

	//Set callback.
	NPAsyncImpl<DNHealthResult>* result = new NPAsyncImpl<DNHealthResult>();
	async->SetCallback([](DNAsync<RPCMessage>* async) {
		NPAsyncImpl<DNHealthResult>* asyncResult = (NPAsyncImpl<DNHealthResult>*)async->GetUserData();
		RPCMessage* message = async->GetResult();
		const RPCHealthReply& reply = message->healthreply();

		DNHealthResult* result = new DNHealthResult();
		result->message = reply.message();
		result->ready = reply.ready();
		asyncResult->SetResult(result);
	}, result);
	async->Wait();

	return result->GetResult();
}