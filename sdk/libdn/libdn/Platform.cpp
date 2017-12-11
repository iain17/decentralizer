#include "StdInc.h"

namespace libdn {

	//Will hang until we are connected and DN is ready.
	LIBDN_API void LIBDN_CALL DN_WaitUntilReady() {
		DNHealthResult* health;
		health->ready = false;
		while (!g_dn.connected || health == nullptr || !health->ready) {
			health = DN_Health();
			if (health != nullptr && health->ready) {
				break;
			}
			Sleep(100);
		}
	}

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
}