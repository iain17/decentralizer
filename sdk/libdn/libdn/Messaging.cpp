#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<bool>*LIBDN_CALL SendDirectMessage(uint32_t channel, PeerID& pid, std::string& data) {
		return SendDirectMessageLegacy(channel, pid, data.c_str(), data.size());
	}

	LIBDN_API Promise<bool>*LIBDN_CALL SendDirectMessageLegacy(uint32_t channel, PeerID& pid, const void* data, size_t size) {
		auto result = new Promise<bool>([channel, pid, data, size](Promise<bool>* promise) {
			// Data we are sending to the server.
			pb::RPCDirectMessage request;
			request.set_channel(channel);
			request.set_pid(pid);
			request.set_message(data, size);

			// Container for the data we expect from the server.
			pb::empty reply;

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

	LIBDN_API Promise<bool>* LIBDN_CALL RegisterDirectMessageCallback(uint32_t channel, DirectMessageCB callback) {
		return new Promise<bool>([channel, callback](Promise<bool>* promise) {
			pb::RPCReceiveDirectMessageRequest request;
			request.set_channel(channel);
			auto ctx = context.client->getContext();
			std::unique_ptr< ::grpc::ClientReader< ::pb::RPCDirectMessage>> reader = context.client->stub_->ReceiveDirectMessage(ctx, request);
			pb::RPCDirectMessage message;
			while (reader->Read(&message)) {
				std::string data = message.message();
				callback(message.pid(), (uint8_t*)data.c_str(), data.size());
			}
			return true;
		});
	}


}