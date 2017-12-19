
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

	if (context.g_logCB)
	{
		context.g_logCB(msgBuffer);
	}
	else {
		printf(msgBuffer);
	}
}