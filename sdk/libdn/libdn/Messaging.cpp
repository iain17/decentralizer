#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<bool>*LIBDN_CALL SendDirectMessage(PeerID pid, const uint8_t * data, uint32_t length) {
		auto result = new Promise<bool>([pid, data, length](Promise<bool>* promise) {
			// Data we are sending to the server.
			pb::RPCDirectMessageRequest request;
			request.set_pid(pid);

			// Container for the data we expect from the server.
			pb::RPCDirectMessageResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->SendDirectMessage(ctx, request, &reply);

			UpsertSessionResult result;
			if (!status.ok()) {
				promise->reject(va("[Could not send direct message to %s] %i: %s", pid.c_str(), status.error_code(), status.error_message().c_str()));
				return false;
			}
			return true;
		});
		return result;
	}

	//TODO: Implement this with streams.
	LIBDN_API void LIBDN_CALL RegisterDirectMessageCallback(void(*callback)(PeerID, const uint8_t *, uint32_t)) {
		return;
	}
}