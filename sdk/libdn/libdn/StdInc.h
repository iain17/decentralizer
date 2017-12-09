#pragma once

#ifndef _STDINC
#define _STDINC

#define _CRT_SECURE_NO_WARNINGS

// Windows headers
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

//Protobuf
#include "protocol.pb.h"
using namespace pb;

// C/C++ headers
#include <string>
#include <vector>
#include <queue>
#include <mutex>

// code headers
#include "libdn.h"
#include "Utils.h"
#include "RPC.h"
#include "DNAsyncImpl.h"

// global state
extern struct DN_state_s
{
	DNID DuID;
	char serverHost[1024];
	uint16_t serverPort;
	ConnectLogCB g_logCB;
} g_np;

#endif