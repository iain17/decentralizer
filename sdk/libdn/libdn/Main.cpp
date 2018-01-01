#include "stdafx.h"
#include "StdInc.h"

namespace libdn {
	LIBDN_API void LIBDN_CALL Init(LogCB logCallback) {
		if (context.initialized) {
			return;
		}
		bool adna = ADNA_Init();
		if (!adna) {
			MessageBoxA(NULL, "Could not start ADNA process.", "libdn", MB_OK);
			exit(0);
			return;
		}
		context.initialized = true;
		context.g_logCB = logCallback;
		Log_Print("Initializing libdn");
	}

	LIBDN_API void LIBDN_CALL Shutdown() {
		return;
	}

	LIBDN_API void LIBDN_CALL RunFrame() {
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