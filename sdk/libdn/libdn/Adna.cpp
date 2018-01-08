#include "StdInc.h"
#include <direct.h>

namespace libdn {
	const char* adnaExecutable = "adna.exe";
	char basePath[FILENAME_MAX];

	void ADNA_Shutdown() {
		killProcessByName(adnaExecutable);
	}

#pragma warning( disable : 4800 ) // stupid warning about bool
#define BUFSIZE 4096
	HANDLE g_hChildStd_OUT_Rd = NULL;
	HANDLE g_hChildStd_OUT_Wr = NULL;
	HANDLE g_hChildStd_ERR_Rd = NULL;
	HANDLE g_hChildStd_ERR_Wr = NULL;

	void ADNA_setupPipe() {
		SECURITY_ATTRIBUTES sa;
		// Set the bInheritHandle flag so pipe handles are inherited. 
		sa.nLength = sizeof(SECURITY_ATTRIBUTES);
		sa.bInheritHandle = TRUE;
		sa.lpSecurityDescriptor = NULL;

		// Create a pipe for the child process's STDERR. 
		if (!CreatePipe(&g_hChildStd_ERR_Rd, &g_hChildStd_ERR_Wr, &sa, 0)) {
			exit(1);
		}
		// Ensure the read handle to the pipe for STDERR is not inherited.
		if (!SetHandleInformation(g_hChildStd_ERR_Rd, HANDLE_FLAG_INHERIT, 0)) {
			exit(1);
		}
		// Create a pipe for the child process's STDOUT. 
		if (!CreatePipe(&g_hChildStd_OUT_Rd, &g_hChildStd_OUT_Wr, &sa, 0)) {
			exit(1);
		}
		// Ensure the read handle to the pipe for STDOUT is not inherited
		if (!SetHandleInformation(g_hChildStd_OUT_Rd, HANDLE_FLAG_INHERIT, 0)) {
			exit(1);
		}
	}

	// Read output from the child process's pipe for STDOUT
	// and write to the parent process's pipe for STDOUT. 
	// Stop when there is no more data. 
	DWORD WINAPI ADNA_read(LPVOID lpParam) {
		//PROCESS_INFORMATION piProcInfo = *(PROCESS_INFORMATION*)lpParam;
		DWORD dwRead;
		CHAR chBuf[BUFSIZE];
		bool bSuccess = FALSE;
		for (;;) {
			bSuccess = ReadFile(g_hChildStd_OUT_Rd, chBuf, BUFSIZE, &dwRead, NULL);
			//bSuccess = ReadFile(g_hChildStd_ERR_Rd, chBuf, BUFSIZE, &dwRead, NULL);
			if (!bSuccess || dwRead == 0) {
				break;
			}
			std::string s(chBuf, dwRead);
			Log_Print("[ADNA]: %s", s.c_str());
		}
		g_hChildStd_OUT_Rd = NULL;
		return 0;
	}

	bool ADNA_reachable() {
		bool reachable = false;
		int tries = 0;
		while (!reachable) {
			reachable = CheckPortTCP((short int)context.port, (char*)context.host);
			if (reachable) {
				break;
			}
			if (tries > 3) {
				break;
			}
			tries++;
			Sleep(500);
		}
		return reachable;
	}

	PROCESS_INFORMATION* NewAdnaInstance() {
		ADNA_Shutdown();
		ADNA_setupPipe();

		PROCESS_INFORMATION piProcInfo;
		STARTUPINFO siStartInfo;
		bool bSuccess = FALSE;

		// Set up members of the PROCESS_INFORMATION structure. 
		ZeroMemory(&piProcInfo, sizeof(PROCESS_INFORMATION));

		// Set up members of the STARTUPINFO structure. 
		// This structure specifies the STDERR and STDOUT handles for redirection.
		ZeroMemory(&siStartInfo, sizeof(STARTUPINFO));
		siStartInfo.cb = sizeof(STARTUPINFO);
		siStartInfo.hStdError = g_hChildStd_ERR_Wr;
		siStartInfo.hStdOutput = g_hChildStd_OUT_Wr;
		siStartInfo.dwFlags |= STARTF_USESTDHANDLES;

		LPSTR params = (LPSTR)va("%s api", adnaExecutable);
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
		CloseHandle(g_hChildStd_ERR_Wr);
		CloseHandle(g_hChildStd_OUT_Wr);
		// If an error occurs, exit the application. 
		if (!bSuccess) {
			//MessageBoxA(NULL, va("Error starting %s.\n", exec), "libdn", MB_OK);
			return nullptr;
		}

		if (&piProcInfo) {
			CreateThread(0, 0, ADNA_read, &piProcInfo, 0, NULL);
		}

		Sleep(1000);
		if (!IsProcessRunning(adnaExecutable)) {
			return nullptr;
		}
		if (!ADNA_reachable()) {
			return nullptr;
		}

		Log_Print("Started adna process.\n");
		return &piProcInfo;
	}

	bool ADNA_Init() {
		if (strlen(basePath) == 0) {
			if (!_getcwd(basePath, sizeof(basePath))) {
				MessageBoxA(NULL, "Could not resolve path.", "libdn", MB_OK);
				return false;
			}
		}
		int tries = 0;
		bool reachable = false;
		//While adna is not reachable.
		while(!reachable) {
			reachable = ADNA_reachable();
			if (reachable) {
				break;
			}
			if (tries > 3) {
				break;
			}
			//If local, spawn one
			if (std::strcmp(context.host, "localhost") == 0 || std::strcmp(context.host, "127.0.0.1") == 0) {
				PROCESS_INFORMATION* piProcInfo = NewAdnaInstance();
				if (!piProcInfo) {
					context.port++;
				}
			}
			Sleep(1000);
			tries++;
		}
		return reachable;
	}
}