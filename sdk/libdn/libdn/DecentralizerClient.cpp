#include "StdInc.h"

DN_state_s context;

namespace libdn {
	LIBDN_API void LIBDN_CALL WaitUntilReady();
	LIBDN_API bool LIBDN_CALL Connect(const char* host, int port, const char* networkKey, bool isPrivateKey, bool limited) {
		std::lock_guard<std::mutex> lock(context.mutex);
		if (!context.initialized) {
			return false;
		}
		context.host = host;
		context.port = port;
		bool adna = ADNA_Ensure_Process(false);
		if (!adna) {
			MessageBoxA(NULL, "Could not connect.", "error", MB_OK);
			exit(0);
		}
		const char* address = va("%s:%i", context.host, context.port);
		Log_Print("Connecting to adna on %s", address);
		auto channel = grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
		auto p2 = std::chrono::system_clock::now() + std::chrono::seconds(30);
		bool result = channel->WaitForConnected(p2);
		if (result) {
			context.client = new DecentralizerClient(channel, networkKey, isPrivateKey, limited);
			Log_Print("Connected. Waiting until ready");
			WaitUntilReady();
			Log_Print("Ready");
		}
		return result;
	}
}