#include "StdInc.h"

namespace libdn {
	//Will hang until we are connected and DN is ready.
	LIBDN_API void LIBDN_CALL WaitUntilReady() {
		HealthResult* health;
		health->ready = false;
		while (health == nullptr || !health->ready) {
			health = Health();
			if (health != nullptr && health->ready) {
				break;
			}
			Sleep(100);
		}
	}

	LIBDN_API HealthResult* LIBDN_CALL Health() {
		// Data we are sending to the server.
		pb::RPCHealthRequest request;

		// Container for the data we expect from the server.
		pb::RPCHealthReply reply;

		auto ctx = context.client->getContext();
		grpc::Status status = context.client->stub_->GetHealth(ctx, request, &reply);

		HealthResult* result = new HealthResult();
		if (status.ok()) {
			result->message = reply.message();
			result->ready = reply.ready();
		} else {
			result->message = va("[RPC failed] %s: %s", status.error_code(), status.error_message());
		}
		return result;
	}

}