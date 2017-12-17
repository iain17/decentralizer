#pragma once

#include "TypeDefs.h"
#include <vector>
#include "Promise.h"
#include "Platform.h"
#include "MatchMaking.h"
#include "Messaging.h"
#include "Addressbook.h"
#include "Storage.h"

namespace libdn {
	// ----------------------------------------------------------
	// Initialization/shutdown functions
	// ----------------------------------------------------------

	// starts up the network platform functions
	LIBDN_API void LIBDN_CALL Init(ConnectLogCB callback);

	// cleans up and shuts down the network platform
	LIBDN_API void LIBDN_CALL Shutdown();

	// connects to a DN server
	LIBDN_API bool LIBDN_CALL Connect(const char* address, const char* networkKey, bool isPrivateKey);

	// Should be called in the games loop
	LIBDN_API void LIBDN_CALL RunFrame();
}