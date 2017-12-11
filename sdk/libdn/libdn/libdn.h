#pragma once

#include "DNTypeDefs.h"
#include <vector>
#include "DNAsync.h"
#include "DNPlatform.h"
#include "DNMatchMaking.h"
#include "DNMessaging.h"
#include "DNAddressbook.h"
#include "DNStorage.h"

namespace libdn {
	// ----------------------------------------------------------
	// Initialization/shutdown functions
	// ----------------------------------------------------------

	// starts up the network platform functions
	LIBDN_API bool LIBDN_CALL DN_Init(ConnectLogCB callback);

	// cleans up and shuts down the network platform
	LIBDN_API bool LIBDN_CALL DN_Shutdown();

	// connects to a DN server
	LIBDN_API bool LIBDN_CALL DN_Connect(const char* server, uint16_t port);

	// Should be called in the games loop
	LIBDN_API bool LIBDN_CALL DN_RunFrame();
}