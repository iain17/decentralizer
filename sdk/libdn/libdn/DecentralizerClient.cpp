#include "StdInc.h"

DN_state_s context;

namespace libdn {
	LIBDN_API bool LIBDN_CALL Connect(const char* address) {
		Log_Print("Connecting to adna on %s", address);
		auto channel = grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
		auto p2 = std::chrono::system_clock::now() + std::chrono::seconds(30);
		bool result = channel->WaitForConnected(p2);
		//grpc::TimePoint<std::chrono::system_clock::time_point>(30);
		if (result) {
			context.client = new DecentralizerClient(channel);
			Log_Print("Connected");
		}
		return result;
	}
}