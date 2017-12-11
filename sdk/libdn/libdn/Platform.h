#pragma once
#include <string>

// ----------------------------------------------------------
// Platform global functions
// ----------------------------------------------------------
namespace libdn {
	const uint64_t API_VERSION = 1;
	class HealthResult {
	public:
		bool ready;
		std::string message;
	};

	// Fetch the health of the DN server.
	LIBDN_API void LIBDN_CALL WaitUntilReady();

	// Fetch the health of the DN server.
	LIBDN_API HealthResult* LIBDN_CALL Health();
}