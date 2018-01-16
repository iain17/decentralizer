#include "StdInc.h"
namespace libdn {

	void fixImproperShutdown(const char* basePath) {
		Log_Print("Improper shutdown detected. Autofixing...");
		ADNA_Shutdown();
		Sleep(5000);
		bool adna = ADNA_Ensure_Process(true);
		if (!adna) {
			exit(0);
		}
		HealthResult health;
		auto request = Health(false);
		if (request->wait()) {
			health = request->get();
		}
		if (health.message.length() == 0) {
			health.message = "Could not fetch health";
		}
		if (!health.ready) {
			MessageBoxA(NULL, health.message.c_str(), "could not start", MB_OK);
			exit(0);
		}
	}

	//Will hang until we are connected and DN is ready.
	LIBDN_API void LIBDN_CALL WaitUntilReady() {
		std::lock_guard<std::mutex> lock(context.healthMutex);
		
		HealthResult health;
		while (!health.ready) {
			auto request = Health(false);
			if (request->wait()) {
				health = request->get();
				if (health.ready) {
					break;
				}
				//Did we close improperly?
				if (health.message.find(".lock") != std::string::npos) {
					fixImproperShutdown(health.basePath.c_str());
					continue;
				}
				//Something about binding.
				if (health.message.find("bind") != std::string::npos) {
					fixImproperShutdown(health.basePath.c_str());
					continue;
				}
				Sleep(100);
			}
		}
	}

	LIBDN_API std::shared_ptr<Promise<HealthResult>> LIBDN_CALL Health(bool waitforminconnections) {
		auto result = std::make_shared<Promise<HealthResult>>([waitforminconnections](auto promise) {
			auto result = HealthResult();
			if (context.client == nullptr) {
				result.message = "Not connected";
				return result;
			}

			bool adna = ADNA_Ensure_Process(false);
			if (!adna) {
				result.message = "Failed to connect";
				return result;
			}

			// Data we are sending to the server.
			pb::RPCHealthRequest request;
			request.set_waitforminconnections(waitforminconnections);

			// Container for the data we expect from the server.
			pb::RPCHealthReply reply;

			auto ctx = context.client->getContext();
			grpc::Status status = context.client->stub_->GetHealth(ctx, request, &reply);

			result.ready = reply.ready();
			result.basePath = reply.basepath();
			if (status.ok()) {
				result.message = reply.message();
			} else {
				result.message = fmt::format("[RPC failed: Get health] {0}: {1}", status.error_code(), status.error_message());
				Log_Print(result.message.c_str());
			}
			return result;
		});
		return result;
	}

}