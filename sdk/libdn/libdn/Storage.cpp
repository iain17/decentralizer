#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<std::string>* LIBDN_CALL GetPublisherFile(const char* name) {
		auto result = new Promise<std::string>([name](Promise<std::string>* promise) {
			// Data we are sending to the server.
			pb::RPCGetPublisherFileRequest request;
			request.set_name(name);

			// Container for the data we expect from the server.
			pb::RPCGetPublisherFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPublisherFile(ctx, request, &reply);

			if (status.ok()) {
				return reply.file();
			} else {
				promise->reject(va("[Could not get publisher file] %i: %s", status.error_code(), status.error_message().c_str()));
			}
		});
		return result;
	}

	LIBDN_API Promise<std::string>*LIBDN_CALL GetPeerFile(PeerID pid, const char * name) {
		auto result = new Promise<std::string>([pid, name](Promise<std::string>* promise) {
			// Data we are sending to the server.
			pb::RPCGetPeerFileRequest request;
			request.set_name(name);
			request.set_pid(pid);

			// Container for the data we expect from the server.
			pb::RPCGetPeerFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPeerFile(ctx, request, &reply);

			if (status.ok()) {
				return reply.file();
			} else {
				promise->reject(va("[Could not get peer file %s] %i: %s", name, status.error_code(), status.error_message().c_str()));
			}
		});
		return result;
	}

	LIBDN_API Promise<bool>*LIBDN_CALL WritePeerFile(const char * name, std::string data) {
		auto result = new Promise<bool>([name, data](Promise<bool>* promise) {
			// Data we are sending to the server.
			pb::RPCWritePeerFileRequest request;
			request.set_name(name);
			request.set_file(data);

			// Container for the data we expect from the server.
			pb::RPCWritePeerFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->WritePeerFile(ctx, request, &reply);

			if (status.ok()) {
				return reply.success();
			} else {
				promise->reject(va("[Could not write peer file %s] %i: %s", name, status.error_code(), status.error_message().c_str()));
			}
			return false;
		});
		return result;
	}
}