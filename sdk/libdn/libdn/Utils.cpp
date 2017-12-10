
#include "StdInc.h"
#include <stdio.h>
#include <intrin.h>  

void Log_Print(const char* message, ...)
{
	static char msgBuffer[32768];
	va_list ap;
	va_start(ap, message);
	vsnprintf(msgBuffer, sizeof(msgBuffer), message, ap);
	va_end(ap);

	OutputDebugStringA(msgBuffer);

	if (g_dn.g_logCB)
	{
		g_dn.g_logCB(msgBuffer);
	}
	else {
		printf(msgBuffer);
	}
}


#define VA_BUFFER_COUNT		4
#define VA_BUFFER_SIZE		32768
static char g_vaBuffer[VA_BUFFER_COUNT][VA_BUFFER_SIZE];
static int g_vaNextBufferIndex = 0;

const char *va(const char *fmt, ...)
{
	va_list ap;
	char *dest = &g_vaBuffer[g_vaNextBufferIndex][0];
	g_vaNextBufferIndex = (g_vaNextBufferIndex + 1) % VA_BUFFER_COUNT;
	va_start(ap, fmt);
	int res = _vsnprintf(dest, VA_BUFFER_SIZE, fmt, ap);
	dest[VA_BUFFER_SIZE - 1] = '\0';
	va_end(ap);

	if (res < 0 || res >= VA_BUFFER_SIZE)
	{
		Log_Print("Attempted to overrun string in call to va() - return address 0x%x", _ReturnAddress());
	}

	return dest;
}