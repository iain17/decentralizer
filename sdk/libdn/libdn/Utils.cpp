
#include "StdInc.h"
#include <stdio.h>
#include <intrin.h>

namespace libdn {
	#define VA_BUFFER_COUNT		4
	#define VA_BUFFER_SIZE		32768

	static char g_vaBuffer[VA_BUFFER_COUNT][VA_BUFFER_SIZE];
	static int g_vaNextBufferIndex = 0;

	const char *va(const char *fmt, ...) {
		va_list ap;
		char *dest = &g_vaBuffer[g_vaNextBufferIndex][0];
		g_vaNextBufferIndex = (g_vaNextBufferIndex + 1) % VA_BUFFER_COUNT;
		va_start(ap, fmt);
		int res = _vsnprintf(dest, VA_BUFFER_SIZE, fmt, ap);
		dest[VA_BUFFER_SIZE - 1] = '\0';
		va_end(ap);

		if (res < 0 || res >= VA_BUFFER_SIZE) {
			Log_Print("Attempted to overrun string in call to va() - return address 0x%x", _ReturnAddress());
		}

		return dest;
	}

	void Log_Print(const char* message, ...) {
		static char msgBuffer[32768];
		va_list ap;
		va_start(ap, message);
		vsnprintf(msgBuffer, sizeof(msgBuffer), message, ap);
		va_end(ap);

		OutputDebugStringA(msgBuffer);

		if (context.g_logCB) {
			context.g_logCB(msgBuffer);
		} else {
			printf(msgBuffer);
		}
	}

	void killProcessByName(const char *filename) {
		HANDLE hSnapShot = CreateToolhelp32Snapshot(TH32CS_SNAPALL, NULL);
		PROCESSENTRY32 pEntry;
		pEntry.dwSize = sizeof(pEntry);
		BOOL hRes = Process32First(hSnapShot, &pEntry);
		while (hRes) {
			if (strcmp(pEntry.szExeFile, filename) == 0) {
				HANDLE hProcess = OpenProcess(PROCESS_TERMINATE, 0,
					(DWORD)pEntry.th32ProcessID);
				if (hProcess != NULL) {
					TerminateProcess(hProcess, 9);
					CloseHandle(hProcess);
				}
			}
			hRes = Process32Next(hSnapShot, &pEntry);
		}
		CloseHandle(hSnapShot);
	}

	bool IsProcessRunning(const char *filename) {
		HANDLE hSnapShot = CreateToolhelp32Snapshot(TH32CS_SNAPALL, NULL);
		PROCESSENTRY32 pEntry;
		pEntry.dwSize = sizeof(pEntry);
		BOOL hRes = Process32First(hSnapShot, &pEntry);
		while (hRes) {
			if (strcmp(pEntry.szExeFile, filename) == 0) {
				HANDLE hProcess = OpenProcess(PROCESS_TERMINATE, 0,
					(DWORD)pEntry.th32ProcessID);
				if (hProcess != NULL) {
					return true;
				}
			}
			hRes = Process32Next(hSnapShot, &pEntry);
		}
		CloseHandle(hSnapShot);
		return false;
	}

	#include <ws2def.h>
	#include <WinSock2.h>

	BOOL CheckPortTCP(short int port, char* hostname) {
		struct sockaddr_in client;
		int sock;

		hostent * record = gethostbyname(hostname);
		if (record == NULL) {
			Log_Print("Could not look up '%s'. Error code: %d.\n", hostname, WSAGetLastError());
			return false;
		}

		client.sin_family = AF_INET;
		client.sin_port = htons(port);
		client.sin_addr.s_addr = *(ULONG*)record->h_addr_list[0];

		int iTimeout = 1600;
		sock = (int)socket(AF_INET, SOCK_STREAM, 0);
		setsockopt(sock,
			SOL_SOCKET,
			SO_RCVTIMEO,
			/*
			reinterpret_cast<char*>(&tv),
			sizeof(timeval) );
			*/
			(const char *)&iTimeout,
			sizeof(iTimeout));
		return (connect(sock, (struct sockaddr *) &client, sizeof(client)) == 0);
	}

	void setupUtils() {
		WSADATA wsadata;
		if (WSAStartup(MAKEWORD(2, 2), &wsadata) != 0) {
			MessageBox(NULL, TEXT("WSAStartup failed!"), TEXT("Error"), MB_OK);
			exit(0);
		}
	}

	#ifdef _WIN32
	#include <io.h> 
	#define access    _access_s
	#else
	#include <unistd.h>
	#endif

	bool FileExists(const std::string &Filename) {
		return access(Filename.c_str(), 0) == 0;
	}
}