#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<bool>*LIBDN_CALL SendDirectMessage(PeerID& pid, std::string& data) {
		return SendDirectMessageLegacy(pid, data.c_str(), data.size());
	}

	LIBDN_API Promise<bool>*LIBDN_CALL SendDirectMessageLegacy(PeerID& pid, const void* data, size_t size) {
		auto result = new Promise<bool>([pid, data, size](Promise<bool>* promise) {
			// Data we are sending to the server.
			pb::RPCDirectMessageRequest request;
			request.set_pid(pid);
			request.set_message(data, size);

			// Container for the data we expect from the server.
			pb::RPCDirectMessageResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->SendDirectMessage(ctx, request, &reply);

			UpsertSessionResult result;
			if (!status.ok()) {
				promise->reject(fmt::format("[Could not send direct message to {0}] {1}: {2}", pid.c_str(), status.error_code(), status.error_message().c_str()));
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