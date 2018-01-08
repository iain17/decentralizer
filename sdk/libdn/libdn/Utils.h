namespace libdn {
	void Log_Print(const char* message, ...);
	void killProcessByName(const char *filename);
	bool IsProcessRunning(const char *filename);
	BOOL CheckPortTCP(short int dwPort, char*ipAddressStr);
	const char *va(const char *fmt, ...);
}

#include <stdarg.h>  // For va_start, etc.
#include <memory>    // For std::unique_ptr
#include <Tlhelp32.h>
