#pragma once

// ----------------------------------------------------------
// Addressbook service
// ----------------------------------------------------------
namespace libdn {
	class Peer {
	public:
		PeerID pId;
		DNID dnId;
		std::map<std::string, std::string> details;
	};

	LIBDN_API Async<int>* LIBDN_CALL UpsertContact(Peer* peer);

	// gets the number of peers using search in details
	// the key is the key of details and the value is the value of the details.
	LIBDN_API Async<int>* LIBDN_CALL GetNumContacts(const char* key, const char* value);

	// gets a single peer by index
	LIBDN_API Peer* LIBDN_CALL GetPeerByIndex(int index);

	// gets a single session's info by either peer id or decentralized id
	LIBDN_API Async<Peer>* LIBDN_CALL GetPeerById(DNID dId, PeerID pId);
}