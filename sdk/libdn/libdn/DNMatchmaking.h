#pragma once
#include <map>

/// ---------------------------------------------------------
 // Matchmaking service
 // ----------------------------------------------------------
namespace libdn {
	typedef uint64_t DNSID;
	class DNSessionInfo {
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

	class DNUpsertSessionResult {
	public:
		bool result;
		DNSID sessionId;
	};

	// creates/updates a session
	LIBDN_API DNAsync<DNUpsertSessionResult>* LIBDN_CALL DN_UpsertSession(DNSessionInfo* data);

	// deletes a session
	LIBDN_API DNAsync<bool>* LIBDN_CALL DN_DeleteSession(DNSID sessionId);

	// gets the number of sessions
	LIBDN_API DNAsync<int>* LIBDN_CALL DN_GetNumSessions(uint32_t type, const char* key, const char* value);

	// gets a single session's info by index
	LIBDN_API DNSessionInfo* LIBDN_CALL DN_GetSessionByIndex(int index);

	// gets a single session's info by sessionId
	LIBDN_API DNAsync<DNSessionInfo>* LIBDN_CALL DN_GetSessionBySessionId(DNSID sessionId);
}