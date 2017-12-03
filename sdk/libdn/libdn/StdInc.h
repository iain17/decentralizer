#pragma once

#ifndef _STDINC
#define _STDINC

#define _CRT_SECURE_NO_WARNINGS

// Windows headers
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

// C/C++ headers
#include <string>
#include <vector>
#include <queue>

// code headers
#include "libdn.h"
#include "Utils.h"
#include "RPC.h"
#include "DNAsyncImpl.h"

// messages
#include "common.pb.h"
#include "MessageDefinition.h"

// global state
extern struct DN_state_s
{
	DNID DuID;
	char serverHost[1024];
	uint16_t serverPort;
	ConnectLogCB g_logCB;
} g_np;

#endif