#pragma once
#include <string>

// ----------------------------------------------------------
// Platform global functions
// ----------------------------------------------------------
namespace libdn {
	extern const char * VERSION;
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