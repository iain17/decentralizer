#include "StdInc.h"

namespace libdn {
	LIBDN_API std::shared_ptr<Promise<DNSID>> LIBDN_CALL UpsertSession(libdn::Session& session) {
		auto result = std::make_shared<Promise<DNSID>>([session](auto promise) {
			// Data we are sending to the server.
			pb::RPCUpsertSessionRequest request;
			auto s = DNSessionToPBSession(session);
			request.set_allocated_session(&s);

			// Container for the data we expect from the server.
			pb::RPCUpsertSessionResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->UpsertSession(ctx, request, &reply);
			request.release_session();

			DNSID result = 0;
			if (status.ok()) {
				result = reply.sessionid();
			} else {
				promise->reject(fmt::format("[Could not upsert session] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}
			return result;
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL DeleteSession(DNSID sid) {
		auto result = std::make_shared<Promise<bool>>([sid](auto promise) {
			//build request.
			pb::RPCDeleteSessionRequest request;
			request.set_sessionid(sid);

			// Container for the data we expect from the server.
			pb::RPCDeleteSessionResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->DeleteSession(ctx, request, &reply);

			if (status.ok()) {
				return true;
			} else {
				promise->reject(fmt::format("[Could not delete session] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}

			return false;
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise<int>> LIBDN_CALL GetNumSessions(uint32_t type, const char* key, const char* value) {
		auto result = std::make_shared<Promise<int>>([type, key, value](auto promise) {
			//build request.
			pb::RPCGetSessionIdsByDetailsRequest request;
			request.set_type(type);
			request.set_key(key);
			request.set_value(value);

			// Container for the data we expect from the server.
			pb::RPCGetSessionIdsResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetSessionIdsByDetails(ctx, request, &reply);

			if (status.ok()) {
				context.sessions = reply.sessionids();
				int size = context.sessions.size();
				return size;
			} else {
				promise->reject(fmt::format("[Could not get session ids] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}

			return 0;
		});
		return result;
	}

	LIBDN_API std::shared_ptr<Promise<libdn::Session>> LIBDN_CALL GetSessionBySessionId(DNSID sessionId) {
		auto result = std::make_shared<Promise<libdn::Session>>([sessionId](auto promise) {
			//build request.
			pb::RPCGetSessionRequest request;
			request.set_sessionid(sessionId);

			// Container for the data we expect from the server.
			pb::RPCGetSessionResponse reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetSession(ctx, request, &reply);

			if (!status.ok()) {
				promise->reject(fmt::format("[Could not get session] {0}: {1}", status.error_code(), status.error_message().c_str()));
			}
			return PBSessionToDNSession(reply.session());
		});
		return result;
	}


	LIBDN_API std::shared_ptr<Session> LIBDN_CALL GetSessionByIndex(int index) {
		if (index > context.sessions.size() - 1) {
			return NULL;
		}
		auto req = GetSessionBySessionId(context.sessions.Get(index));
		req->fail([](std::string reason) {
			Log_Print(reason.c_str());
		});
		if (req->wait()) {
			return std::make_shared<Session>(req->get());
		}
		return NULL;
	}
}