#pragma once

#ifndef _STDINC
#define _STDINC

#define _CRT_SECURE_NO_WARNINGS

// Windows headers
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

//Protobuf
#include "matchmaking.pb.h"
#include "platform.pb.h"
#include "addressbook.pb.h"
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

const int MAX_SESSIONS = 1024;
// global state
extern struct DN_state_s {
	DNID DuID;
	char serverHost[1024];
	uint16_t serverPort;
	ConnectLogCB g_logCB;

	::google::protobuf::RepeatedField<::google::protobuf::uint64> sessions;
} g_dn;

#endif