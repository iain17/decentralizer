#include "StdInc.h"

DN_state_s context;

namespace libdn {
	LIBDN_API void LIBDN_CALL WaitUntilReady();
	LIBDN_API bool LIBDN_CALL Connect(const char* address, const char* networkKey, bool isPrivateKey, bool limited) {
		if (!context.initialized) {
			return false;
		}
		Log_Print("Connecting to adna on %s", address);
		auto channel = grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
		auto p2 = std::chrono::system_clock::now() + std::chrono::seconds(30);
		bool result = channel->WaitForConnected(p2);
		if (result) {
			context.client = new DecentralizerClient(channel, networkKey, isPrivateKey, limited);
			Log_Print("Connected. Waiting until ready");
			Sleep(1000);
			WaitUntilReady();
			Log_Print("Ready");
		}
		return result;
	}
}