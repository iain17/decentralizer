#pragma once

#include "TypeDefs.h"
#include <vector>
#include "Promise.h"
#include <map>

namespace libdn {
	// ----------------------------------------------------------
	// Initialization/shutdown functions
	// ----------------------------------------------------------

	// starts up the network platform functions
	LIBDN_API void LIBDN_CALL Init(LogCB logCallback);

	// cleans up and shuts down the network platform
	LIBDN_API void LIBDN_CALL Shutdown();

	// connects to a DN server
	LIBDN_API bool LIBDN_CALL Connect(const char* host, int port, const char* networkKey, bool isPrivateKey, bool limited);

	// Should be called in the games loop
	LIBDN_API void LIBDN_CALL RunFrame();

	// ----------------------------------------------------------
	// Platform global functions
	// ----------------------------------------------------------
	class HealthResult {
	public:
		bool ready;
		std::string message;
		std::string basePath;
	};

	// Fetch the health of the DN server.
	LIBDN_API void LIBDN_CALL WaitUntilReady();

	// Fetch the health of the DN server.
	LIBDN_API std::shared_ptr<HealthResult> LIBDN_CALL Health();

	// ----------------------------------------------------------
	// Storage
	// ----------------------------------------------------------

	// Fetches a publisher file.
	LIBDN_API std::shared_ptr<Promise<std::string>> LIBDN_CALL GetPublisherFile(const char* name);

	// Fetches a peer file
	LIBDN_API std::shared_ptr<Promise<std::string>>LIBDN_CALL GetPeerFile(PeerID& pid, const char * name);

	// Writes a peer file
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL WritePeerFileLegacy(const char * name, const void* data, size_t size);
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL WritePeerFile(const char * name, std::string& data);

	// ----------------------------------------------------------
	// Addressbook service
	// ----------------------------------------------------------
	class Peer {
	public:
		PeerID pId;
		DNID dnId;
		std::map<std::string, std::string> details;
	};
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL UpsertPeer(Peer& peer);

	// gets the number of peers using search in details
	// the key is the key of details and the value is the value of the details.
	LIBDN_API std::shared_ptr<Promise<int>> LIBDN_CALL GetNumPeers(const char* key, const char* value);

	// gets a single peer by index
	LIBDN_API std::shared_ptr<Peer> LIBDN_CALL GetPeerByIndex(int index);

	// gets a single session's info by either peer id or decentralized id
	LIBDN_API std::shared_ptr<Promise<Peer>> LIBDN_CALL GetPeerById(DNID dId, PeerID& pId);

	//Get yourself.
	LIBDN_API Peer* LIBDN_CALL GetSelf();

	//Refresh ourself
	LIBDN_API bool LIBDN_CALL refreshSelf();

	//Resolve a decentralized id to a peer id
	LIBDN_API PeerID LIBDN_CALL ResolveDecentralizedId(DNID dId);


	// ----------------------------------------------------------
	// Direct messaging service
	// ----------------------------------------------------------
	// sends direct message to another peer.
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL SendDirectMessage(uint32_t channel, PeerID& pid, std::string& data);
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL SendDirectMessageLegacy(uint32_t channel, PeerID& pid, const void* data, size_t size);
	// function to register a callback when a direct message has been received
	// Default channel is 0. For dedicated servers the channel will be the session id.
	LIBDN_API void LIBDN_CALL RegisterDirectMessageCallback(uint32_t channel, DirectMessageCB callback);

	/// ---------------------------------------------------------
	// Matchmaking service
	// ----------------------------------------------------------
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

	// creates/updates a session
	LIBDN_API std::shared_ptr<Promise<DNSID>> LIBDN_CALL UpsertSession(libdn::Session& session);

	// deletes a session
	LIBDN_API std::shared_ptr<Promise<bool>> LIBDN_CALL DeleteSession(DNSID sessionId);

	// gets the number of sessions
	LIBDN_API std::shared_ptr<Promise<int>> LIBDN_CALL GetNumSessions(uint32_t type, const char* key, const char* value);

	// gets a single session's info by index
	LIBDN_API std::shared_ptr<Session> LIBDN_CALL GetSessionByIndex(int index);

	// gets a single session's info by sessionId
	LIBDN_API std::shared_ptr<Promise<libdn::Session>> LIBDN_CALL GetSessionBySessionId(DNSID sessionId);
}