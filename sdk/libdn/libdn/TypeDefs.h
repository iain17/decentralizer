#pragma once

#ifndef CONNECT_TYPEDEFS_
#define CONNECT_TYPEDEFS_

#include "stdint.h"

namespace libdn {
	typedef uint64_t DNID;
	typedef std::string PeerID;

	typedef void(_cdecl * LogCB)(const char* message);
	typedef void(_cdecl * DirectMessageCB)(PeerID&, const uint8_t *, uint32_t);

	#define LIBDN_API extern "C" __declspec(dllexport)

	#define LIBDN_CALL __cdecl
}

#endif