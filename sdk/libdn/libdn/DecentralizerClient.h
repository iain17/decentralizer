#pragma once
#include <grpc/support/log.h>

namespace libdn {
	using grpc::Channel;
	using grpc::ClientAsyncResponseReader;
	using grpc::ClientContext;
	using grpc::CompletionQueue;
	using grpc::Status;
	using pb::Decentralizer;

	class DecentralizerClient {
	public:
		// Out of the passed in Channel comes the stub, stored here, our view of the
		// server's exposed services.
		std::unique_ptr<Decentralizer::Stub> stub_;

		// The producer-consumer queue we use to communicate asynchronously with the
		// gRPC runtime.
		CompletionQueue cq_;

		const char* networkKey;
		bool isPrivateKey;

		grpc::ClientContext* getContext() {
			grpc::ClientContext* ctx = new grpc::ClientContext();
			ctx->AddMetadata("cver", "0.1.0");
			ctx->AddMetadata("netkey", networkKey);
			ctx->AddMetadata("privkey", isPrivateKey ? "1" : "0");
			return ctx;
		}

		explicit DecentralizerClient(std::shared_ptr<Channel> channel, const char* networkKey, bool isPrivateKey) : stub_(Decentralizer::NewStub(channel)) {
			this->networkKey = networkKey;
			this->isPrivateKey = isPrivateKey;
		}
	};
}