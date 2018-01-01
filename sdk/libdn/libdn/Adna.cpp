#include "StdInc.h"
#include <direct.h>

const char* adnaExecutable = "adna.exe";
char basePath[FILENAME_MAX];

void ADNA_Shutdown() {
	killProcessByName(adnaExecutable);
}

PROCESS_INFORMATION* NewAdnaInstance() {
	PROCESS_INFORMATION piProcInfo;
	STARTUPINFO siStartInfo;
	bool bSuccess = FALSE;

	// Set up members of the PROCESS_INFORMATION structure. 
	ZeroMemory(&piProcInfo, sizeof(PROCESS_INFORMATION));

	// Set up members of the STARTUPINFO structure. 
	// This structure specifies the STDERR and STDOUT handles for redirection.
	ZeroMemory(&siStartInfo, sizeof(STARTUPINFO));
	siStartInfo.cb = sizeof(STARTUPINFO);
	//siStartInfo.hStdError = g_hChildStd_ERR_Wr;
	//siStartInfo.hStdOutput = g_hChildStd_OUT_Wr;
	siStartInfo.dwFlags |= STARTF_USESTDHANDLES;

	LPSTR params = (LPSTR)"";//--SafeLogging 0
	const char* exec = va("%s\\%s", basePath, adnaExecutable);
	bSuccess = CreateProcess(exec,
		params,
		NULL,          // process security attributes 
		NULL,          // primary thread security attributes 
		TRUE,          // handles are inherited 
		CREATE_NO_WINDOW,             // creation flags 
		NULL,          // use parent's environment 
		(LPSTR)va("%s\\", basePath),          // use parent's current directory 
		&siStartInfo,  // STARTUPINFO pointer 
		&piProcInfo);  // receives PROCESS_INFORMATION
	//CloseHandle(g_hChildStd_ERR_Wr);
	//CloseHandle(g_hChildStd_OUT_Wr);
	// If an error occurs, exit the application. 
	if (!bSuccess) {
		MessageBoxA(NULL, va("Error starting %s.\n", exec), "libdn", MB_OK);
		return nullptr;
	}

	Log_Print("Started adna.\n");
	return &piProcInfo;
}

bool ADNA_Init() {
	if (strlen(basePath) == 0) {
		if (!_getcwd(basePath, sizeof(basePath))) {
			MessageBoxA(NULL, "Could not resolve path.", "libdn", MB_OK);
		}
	}
	if (!IsProcessRunning(adnaExecutable)) {
		PROCESS_INFORMATION* piProcInfo = NewAdnaInstance();
		if (piProcInfo) {
			return true;
		} else {
			return false;
		}
	}
	return true;
}