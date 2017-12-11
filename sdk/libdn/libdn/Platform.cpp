#include "StdInc.h"

namespace libdn {

	//Will hang until we are connected and DN is ready.
	LIBDN_API void LIBDN_CALL WaitUntilReady() {
		HealthResult* health;
		health->ready = false;
		while (!g_dn.connected || health == nullptr || !health->ready) {
			health = Health();
			if (health != nullptr && health->ready) {
				break;
			}
			Sleep(100);
		}
	}

	LIBDN_API HealthResult* LIBDN_CALL Health() {
		//build request.
		RPCHealthRequest* request = new RPCHealthRequest();
		pb::RPCMessage* msg = new pb::RPCMessage();
		msg->set_allocated_healthrequest(request);

		//Send request.
		Async<RPCMessage>* async = RPC_SendMessageAsyncCache("health", msg);

		//Set callback.
		AsyncImpl<HealthResult>* result = new AsyncImpl<HealthResult>();
		async->SetCallback([](Async<RPCMessage>* async) {
			AsyncImpl<HealthResult>* asyncResult = (AsyncImpl<HealthResult>*)async->GetUserData();
			RPCMessage* message = async->GetResult();
			const RPCHealthReply& reply = message->healthreply();

			HealthResult* result = new HealthResult();
			result->message = reply.message();
			result->ready = reply.ready();
			asyncResult->SetResult(result);
		}, result);
		async->Wait();

		return result->GetResult();
	}
}