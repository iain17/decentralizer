#pragma once
#include <map>

/// ---------------------------------------------------------
 // Matchmaking service
 // ----------------------------------------------------------
namespace libdn {
	typedef uint64_t DNSID;
	class Session {
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
		DNSID sessionId;
	};

	// creates/updates a session
	LIBDN_API Promise<UpsertSessionResult>* LIBDN_CALL UpsertSession(libdn::Session * session);

	// deletes a session
	LIBDN_API Promise<bool>* LIBDN_CALL DeleteSession(DNSID sessionId);

	// gets the number of sessions
	LIBDN_API Promise<int>* LIBDN_CALL GetNumSessions(uint32_t type, const char* key, const char* value);

	// gets a single session's info by index
	LIBDN_API Session* LIBDN_CALL GetSessionByIndex(int index);

	// gets a single session's info by sessionId
	LIBDN_API Promise<Session*>* LIBDN_CALL GetSessionBySessionId(DNSID sessionId);
}