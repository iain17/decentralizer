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

		explicit DecentralizerClient(std::shared_ptr<Channel> channel) : stub_(Decentralizer::NewStub(channel)) {
		}
	};
}