#include "StdInc.h"

namespace libdn {
	//Will hang until we are connected and DN is ready.
	LIBDN_API void LIBDN_CALL WaitUntilReady() {
		HealthResult* health = Health();
		while (!health || !health->ready) {
			health = Health();
			if (health != nullptr && health->ready) {
				break;
			}
			//Did we close improperly?
			if (health->message.find(".lock") != std::string::npos) {
				ADNA_Shutdown();
				std::string lockFile = health->basePath + "\\ipfs\\repo.lock";
				int i = remove(lockFile.c_str());
				if(i == 0){
					Log_Print(fmt::format("Could not delete lock file. Please do so manually: {0}", lockFile.c_str()).c_str());
					Sleep(8000);
				}
				Sleep(1000);
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
		result->ready = reply.ready();
		result->basePath = reply.basepath();
		if (status.ok()) {
			result->message = reply.message();
		} else {
			result->message = fmt::format("[RPC failed: Get health] {0}: {1}", status.error_code(), status.error_message());
			Log_Print(result->message.c_str());
		}
		return result;
	}

}