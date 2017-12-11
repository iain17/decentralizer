#pragma once

// ----------------------------------------------------------
// Addressbook service
// ----------------------------------------------------------
namespace libdn {
	class DNPeer {
	public:
		PeerID pId;
		DNID dnId;
		std::map<std::string, std::string> details;
	};

	LIBDN_API DNAsync<int>* LIBDN_CALL DN_UpsertContact(libdn::DNPeer* peer);

	// gets the number of peers using search in details
	// the key is the key of details and the value is the value of the details.
	LIBDN_API DNAsync<int>* LIBDN_CALL DN_GetNumContacts(const char* key, const char* value);

	// gets a single peer by index
	LIBDN_API libdn::DNPeer* LIBDN_CALL DN_GetPeerByIndex(int index);

	// gets a single session's info by either peer id or decentralized id
	LIBDN_API DNAsync<DNPeer>* LIBDN_CALL DN_GetPeerById(DNID dId, PeerID pId);
}