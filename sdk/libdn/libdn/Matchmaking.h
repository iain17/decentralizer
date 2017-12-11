#pragma once
#include <map>

/// ---------------------------------------------------------
 // Matchmaking service
 // ----------------------------------------------------------
namespace libdn {
	typedef uint64_t DNSID;
	class SessionInfo {
	public:
		PeerID pId;
		DNID dnId;

		DNSID sessionId;
		uint64_t type;
		std::string name;
		uint32_t address;
		uint16_t port;
		std::map<std::string, std::string> details;
	};

	class UpsertSessionResult {
	public:
		bool result;
		DNSID sessionId;
	};

	// creates/updates a session
	LIBDN_API Async<UpsertSessionResult>* LIBDN_CALL UpsertSession(SessionInfo* data);

	// deletes a session
	LIBDN_API Async<bool>* LIBDN_CALL DeleteSession(DNSID sessionId);

	// gets the number of sessions
	LIBDN_API Async<int>* LIBDN_CALL GetNumSessions(uint32_t type, const char* key, const char* value);

	// gets a single session's info by index
	LIBDN_API SessionInfo* LIBDN_CALL GetSessionByIndex(int index);

	// gets a single session's info by sessionId
	LIBDN_API Async<SessionInfo>* LIBDN_CALL GetSessionBySessionId(DNSID sessionId);
}