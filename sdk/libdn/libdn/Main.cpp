#include "stdafx.h"
#include "StdInc.h"

namespace libdn {
	DWORD WINAPI RefreshSelf(LPVOID lpParam) {
		while (true) {
			WaitUntilReady();
			context.selfMutex.lock();
			refreshSelf();
			context.selfMutex.unlock();
			Sleep(15 * 1000);
		}
	}

	LIBDN_API void LIBDN_CALL Init(LogCB logCallback) {
		if (context.initialized) {
			return;
		}
		setupUtils();
		context.initialized = true;
		context.g_logCB = logCallback;
		CreateThread(0, 0, RefreshSelf, NULL, 0, NULL);
		Log_Print("Initializing libdn");
	}

	LIBDN_API void LIBDN_CALL Shutdown() {
		return;
	}

	LIBDN_API void LIBDN_CALL RunFrame() {
		return;
	}
}
/*
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
*/