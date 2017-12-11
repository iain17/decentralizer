#pragma once

#include "TypeDefs.h"
#include <vector>
#include "Async.h"
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
	LIBDN_API bool LIBDN_CALL Init(ConnectLogCB callback);

	// cleans up and shuts down the network platform
	LIBDN_API bool LIBDN_CALL Shutdown();

	// connects to a DN server
	LIBDN_API bool LIBDN_CALL Connect(const char* server, uint16_t port);

	// Should be called in the games loop
	LIBDN_API bool LIBDN_CALL RunFrame();
}