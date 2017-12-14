#include "StdInc.h"

namespace libdn {
	LIBDN_API Promise<UpsertSessionResult>* LIBDN_CALL UpsertSession(libdn::Session * session) {
		pb::Session* test = DNSessionToPBSession(session);
		auto result = new Promise<UpsertSessionResult>(
			[test](Promise<UpsertSessionResult>* promise) {
				// Data we are sending to the server.
				pb::RPCUpsertSessionRequest request;
				request.set_allocated_session(test);

				// Container for the data we expect from the server.
				pb::RPCUpsertSessionResponse reply;

				// Context for the client. It could be used to convey extra information to
				// the server and/or tweak certain RPC behaviors.
				grpc::ClientContext ctx;
				grpc::Status status = context.client->stub_->UpsertSession(&ctx, request, &reply);

				UpsertSessionResult result;
				if (status.ok()) {
					result.sessionId = reply.sessionid();
				} else {
					//result->message = va("[RPC failed] %s: %s", status.error_code(), status.error_message());
				}
				return result;
			}
		);
		return result;
	}
}