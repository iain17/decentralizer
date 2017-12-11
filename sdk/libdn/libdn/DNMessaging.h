#pragma once

// ----------------------------------------------------------
// Direct messaging service
// ----------------------------------------------------------
namespace libdn {
	// sends direct message to another peer.
	LIBDN_API  DNAsync<bool>* LIBDN_CALL DN_SendDirectMessage(PeerID id, const uint8_t* data, uint32_t length);

	// function to register a callback when a direct message has been received
	// arguments: source peer id, data, length
	LIBDN_API void LIBDN_CALL DN_RegisterDirectMessageCallback(void(__cdecl * callback)(PeerID, const uint8_t*, uint32_t));
}