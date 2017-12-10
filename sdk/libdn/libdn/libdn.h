#pragma once

#include "DNTypeDefs.h"
#include <vector>
#include "DNAsync.h"
#include "DNPlatform.h"
#include "DNMatchMaking.h"

// ----------------------------------------------------------
// Initialization/shutdown functions
// ----------------------------------------------------------

// starts up the network platform functions
LIBDN_API bool LIBDN_CALL DN_Init(ConnectLogCB callback);

// cleans up and shuts down the network platform
LIBDN_API bool LIBDN_CALL DN_Shutdown();

// connects to a DN server
LIBDN_API bool LIBDN_CALL DN_Connect(const char* server, uint16_t port);

// ----------------------------------------------------------
// Callback handling
// ----------------------------------------------------------

// handles and dispatches callbacks
// must be called every frame since it now handles sockets
LIBDN_API bool LIBDN_CALL DN_RunFrame();

// ----------------------------------------------------------
// Platform global functions
// ----------------------------------------------------------

// Fetch the health of the DN server.
LIBDN_API void LIBDN_CALL DN_WaitUntilReady();

// Fetch the health of the DN server.
LIBDN_API DNHealthResult* LIBDN_CALL DN_Health();

/// ---------------------------------------------------------
// Matchmaking service
// ----------------------------------------------------------
// creates/updates a session
LIBDN_API DNAsync<DNUpsertSessionResult>* LIBDN_CALL DN_UpsertSession(DNSessionInfo* data);

// deletes a session
LIBDN_API DNAsync<bool>* LIBDN_CALL DN_DeleteSession(DNSID sessionId);

// gets the number of sessions
LIBDN_API DNAsync<std::vector<DNSID>>* LIBDN_CALL DN_GetSessionIds(uint32_t type, const char* key, const char* value);

// gets a single session's info
LIBDN_API DNAsync<DNSessionInfo>* LIBDN_CALL DN_GetSession(DNSID sessionId);

/*
// ----------------------------------------------------------
// Storage service
// ----------------------------------------------------------

// obtains a file from the remote global per-title storage
LIBDN_API DNAsync<NPGetPublisherFileResult>* LIBDN_CALL DN_GetPublisherFile(const char* fileName, uint8_t* buffer, size_t bufferLength);

// obtains a file from the remote per-user storage
LIBDN_API DNAsync<NPGetUserFileResult>* LIBDN_CALL DN_GetUserFile(const char* fileName, NPID npID, uint8_t* buffer, size_t bufferLength);

// uploads a file to the remote per-user storage
LIBDN_API DNAsync<NPWriteUserFileResult>* LIBDN_CALL DN_WriteUserFile(const char* fileName, NPID npID, const uint8_t* buffer, size_t bufferLength);

// sends a random string to the NP server
LIBDN_API void LIBDN_CALL DN_SendRandomString(const char* str);

// similar to above, except allowing a result
LIBDN_API DNAsync<NPSendRandomStringResult>* LIBDN_CALL DN_SendRandomStringExt(const char* str, char* outStr, size_t outLength);

// server push messages from external services
LIBDN_API void LIBDN_CALL DN_RegisterRandomStringCallback(void(__cdecl * callback)(const char*));

// obtains important/relatevant and what else would be public data the np server has on a specfic npid
LIBDN_API DNAsync<NPGetConnectionDetailsResult>* LIBDN_CALL DN_GetConnectionDetails(NPID npID);

LIBDN_API void LIBDN_CALL DN_RegisterScreenShotCallback(void(__cdecl * callback)(NPID));

LIBDN_API void LIBDN_CALL DN_SendScreenShot(uint64_t id, const char* data);

// ----------------------------------------------------------
// Friends service
// ----------------------------------------------------------

// returns if the friends API is available
LIBDN_API bool LIBDN_CALL DN_FriendsConnected();

// gets the number of registered friends
LIBDN_API int LIBDN_CALL DN_GetNumFriends();

// obtains profile data for a list of NPIDs
LIBDN_API DNAsync<NPFriendResult>* LIBDN_CALL DN_GetFriends();

//Used by steam api (t5, iw5m)
// gets a specific friend's NPID - index is from 0...[DN_GetNumFriends() - 1]
LIBDN_API NPID LIBDN_CALL DN_GetFriend(NPID index);

// gets the name for a friend
// will only work with NPIDs known to the client - currently only if they're friends
LIBDN_API const char* LIBDN_CALL DN_GetFriendName(NPID npID);

// gets the presence state for a friend
LIBDN_API EPresenceState LIBDN_CALL DN_GetFriendPresence(NPID npID);

// sets a presence key/value pair
// a value of NULL removes the key, if existent
LIBDN_API void LIBDN_CALL DN_SetRichPresence(const char* key, const char* value);

// sets the presence body
LIBDN_API void LIBDN_CALL DN_SetRichPresenceBody(const char* body);

// uploads the rich presence data
LIBDN_API void LIBDN_CALL DN_StoreRichPresence();

// gets a rich presence value for a friend
// will only work with friends, not other known NPIDs
LIBDN_API const char* LIBDN_CALL DN_GetFriendRichPresence(NPID npID, const char* key);

// obtains extended profile data for a list of NPIDs
// to write this data, write a user file named 'profile[profileType]' (max 64kb currently, will be truncated on fetch)
//Difference between DN_GetFriendRichPresence is that it will ask the np server for a value outside our normal NPFriend struct.
LIBDN_API DNAsync<NPGetExtProfileDataResult>* LIBDN_CALL DN_GetExtProfileData(uint32_t numIDs, const char* profileType, const NPID* npIDs, NPExtProfileData* outData);

// gets an avatar for any client
LIBDN_API DNAsync<NPGetUserAvatarResult>* LIBDN_CALL DN_GetUserAvatar(int id, uint8_t* buffer, size_t bufferLength);

// ----------------------------------------------------------
// Messaging service
// ----------------------------------------------------------

// sends arbitrary data to a client
LIBDN_API void LIBDN_CALL DN_SendMessage(NPID npid, const uint8_t* data, uint32_t length);

// function to register a callback when a message is received
// arguments: source NPID, data, length
LIBDN_API void LIBDN_CALL DN_RegisterMessageCallback(void(__cdecl * callback)(NPID, const uint8_t*, uint32_t));
*/