#pragma once
// ----------------------------------------------------------
// Storage
// ----------------------------------------------------------
namespace libdn {
	// Fetches a publisher file.
	LIBDN_API Promise<std::string>* LIBDN_CALL GetPublisherFile(const char* name);

	// Fetches a peer file
	LIBDN_API Promise<std::string>* LIBDN_CALL GetPeerFile(PeerID pid, const char* name);

	// WRites a peer file
	LIBDN_API Promise<bool>*LIBDN_CALL WritePeerFile(const char * name, std::string data);
}