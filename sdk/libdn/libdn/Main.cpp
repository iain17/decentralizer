#include "stdafx.h"
#include "StdInc.h"

namespace libdn {
	LIBDN_API void LIBDN_CALL Init(ConnectLogCB callback) {
		if (context.initialized) {
			return;
		}
		context.initialized = true;
		context.g_logCB = callback;
		Log_Print("Initializing libdn");

		//Test
		libdn::Promise<int>* test = new libdn::Promise<int>([](libdn::Promise<int>* promise) {
			promise->reject("Not cool");
			return 1;
		});
		test->wait();

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