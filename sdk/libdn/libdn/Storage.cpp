#include "StdInc.h"

namespace libdn {
	LIBDN_API std::shared_ptr<Promise<std::string>> LIBDN_CALL GetPublisherFile(const char* name) {
		auto result = std::make_shared<Promise<std::string>>([name](auto promise) {
			// Data we are sending to the server.
			pb::RPCGetPublisherFileRequest request;
			request.set_name(name);

			// Container for the data we expect from the server.
			pb::RPCGetPublisherFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPublisherFile(ctx, request, &reply);

			if (!status.ok()) {
				promise->reject(fmt::format("[Could not get publisher file] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}
			return reply.file();
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise< std::string >> LIBDN_CALL GetPeerFile(PeerID& pid, const char * name) {
		auto result = std::make_shared<Promise< std::string>>([pid, name](auto promise) {
			// Data we are sending to the server.
			pb::RPCGetPeerFileRequest request;
			request.set_name(name);
			request.set_pid(pid);
			
			// Container for the data we expect from the server.
			pb::RPCGetPeerFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetPeerFile(ctx, request, &reply);
			if (!status.ok()) {
				promise->reject(fmt::format("[Could not get peer file {0}] {1}: {2}", name, status.error_code(), status.error_message().c_str()));
			}
			return reply.file();
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL WritePeerFileLegacy(const char * name, const void* data, size_t size) {
		auto result = std::make_shared<Promise<bool>>([name, data, size](auto promise) {
			// Data we are sending to the server.
			pb::RPCWritePeerFileRequest request;
			request.set_name(name);
			request.set_file(data, size);

				// Container for the data we expect from the server.
				pb::RPCWritePeerFileResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->WritePeerFile(ctx, request, &reply);

			if (status.ok()) {
				return reply.success();
			} else {
				promise->reject(fmt::format("[Could not write peer file {0}] {1}: {2}", name, status.error_code(), status.error_message().c_str()));
			}
			return false;
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL WritePeerFile(const char * name, std::string& data) {
		return WritePeerFileLegacy(name, data.c_str(), data.size());
	}
}