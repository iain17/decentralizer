#include "StdInc.h"

DN_state_s context;

namespace libdn {
	LIBDN_API bool LIBDN_CALL Connect(const char* address) {
		Log_Print("Connecting to adna on %s\n", address);
		context.client = (DecentralizerClient*)&grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
		return true;
	}
}