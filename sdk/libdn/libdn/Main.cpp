#include "stdafx.h"
#include "StdInc.h"

namespace libdn {
	LIBDN_API void LIBDN_CALL Init(LogCB logCallback, DirectMessageCB directMessageCallback) {
		if (context.initialized) {
			return;
		}
		context.initialized = true;
		context.g_logCB = logCallback;
		context.g_dmCB = directMessageCallback;
		Log_Print("Initializing libdn");
	}

	LIBDN_API void LIBDN_CALL Shutdown() {
		return;
	}

	void ListenToDirectMessages();
	LIBDN_API void LIBDN_CALL RunFrame() {
		if (!context.DMListening) {
			ListenToDirectMessages();
		}
		return;
	}
}

BOOL APIENTRY DllMain(HMODULE hModule,
	DWORD  ul_reason_for_call,
	LPVOID lpReserved
) {
	switch (ul_reason_for_call) {
	case DLL_PROCESS_ATTACH:
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}