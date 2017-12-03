#pragma once

#include "DNTypeDefs.h"

#include "DNAsync.h"
/*
#include "NPAuthenticate.h"
#include "NPEncryption.h"
#include "NPStorage.h"
#include "NPFriends.h"
#include "NPServers.h"
*/

// ----------------------------------------------------------
// Initialization/shutdown functions
// ----------------------------------------------------------

// starts up the network platform functions
LIBDN_API bool LIBDN_CALL DN_Init(ConnectLogCB callback);

// cleans up and shuts down the network platform
LIBDN_API bool LIBDN_CALL DN_Shutdown();

// connects to a NP server
LIBDN_API bool LIBDN_CALL DN_Connect(const char* server, uint16_t port);

/*
// ----------------------------------------------------------
// Dictionary type creation
// ----------------------------------------------------------

// creates a string/string dictionary
// use NPDictionary constructor instead!
LIBDN_API NPDictionaryInternal* LIBDN_CALL DN_CreateDictionary();

// ----------------

// ----------------------------------------------------------
// Callback handling
// ----------------------------------------------------------

// handles and dispatches callbacks
// must be called every frame since it now handles sockets
LIBDN_API bool LIBDN_CALL DN_RunFrame();

// ----------------------------------------------------------
// Authentication service
// ----------------------------------------------------------

// authenticates using an external auth token
LIBDN_API NPAsync<NPAuthenticateResult>* LIBDN_CALL DN_AuthenticateWithToken(const char* authToken);

// authenticates using a username/password
LIBDN_API NPAsync<NPAuthenticateResult>* LIBDN_CALL DN_AuthenticateWithDetails(const char* username, const char* password);

// authenticates using a license key
LIBDN_API NPAsync<NPAuthenticateResult>* LIBDN_CALL DN_AuthenticateWithLicenseKey(const char* licenseKey);

// Returns the license key the client is using
LIBDN_API char* LIBDN_CALL DN_ReturnLicenseKey();

// registers a game server license key
LIBDN_API NPAsync<NPRegisterServerResult>* LIBDN_CALL DN_RegisterServer(const char* configPath);

// validates a user ticket
LIBDN_API NPAsync<NPValidateUserTicketResult>* LIBDN_CALL DN_ValidateUserTicket(const void* ticket, size_t ticketSize, uint32_t clientIP, NPID clientID, const char* name);

// obtains a user ticket for server authentication
LIBDN_API bool LIBDN_CALL DN_GetUserTicket(void* buffer, size_t bufferSize, NPID targetServer);

// gets the NPID for the current client. returns false (and does not change the output buffer) if not yet authenticated
LIBDN_API bool LIBDN_CALL DN_GetNPID(NPID* pID);

// gets the user group for the current client, returns 0 if not authenticated
LIBDN_API int LIBDN_CALL DN_GetUserGroup();

// function to register a callback to kick a client by NPID
LIBDN_API void LIBDN_CALL DN_RegisterKickCallback(void(__cdecl * callback)(NPID, const char*));

// function to register a callback to kick a client by NPID
LIBDN_API void LIBDN_CALL DN_RegisterRemoteConsoleCallback(void(__cdecl * callback)(const char*));

// function to register a callback when external authentication status changes
LIBDN_API void LIBDN_CALL DN_RegisterEACallback(void(__cdecl * callback)(EExternalAuthState));

// loads the game DLL for the specified version number
LIBDN_API void* LIBDN_CALL DN_LoadGameModule(int version);

// ----------------------------------------------------------
// Storage service
// ----------------------------------------------------------

// obtains a file from the remote global per-title storage
LIBDN_API NPAsync<NPGetPublisherFileResult>* LIBDN_CALL DN_GetPublisherFile(const char* fileName, uint8_t* buffer, size_t bufferLength);

// obtains a file from the remote per-user storage
LIBDN_API NPAsync<NPGetUserFileResult>* LIBDN_CALL DN_GetUserFile(const char* fileName, NPID npID, uint8_t* buffer, size_t bufferLength);

// uploads a file to the remote per-user storage
LIBDN_API NPAsync<NPWriteUserFileResult>* LIBDN_CALL DN_WriteUserFile(const char* fileName, NPID npID, const uint8_t* buffer, size_t bufferLength);

// sends a random string to the NP server
LIBDN_API void LIBDN_CALL DN_SendRandomString(const char* str);

// similar to above, except allowing a result
LIBDN_API NPAsync<NPSendRandomStringResult>* LIBDN_CALL DN_SendRandomStringExt(const char* str, char* outStr, size_t outLength);

// server push messages from external services
LIBDN_API void LIBDN_CALL DN_RegisterRandomStringCallback(void(__cdecl * callback)(const char*));

// obtains important/relatevant and what else would be public data the np server has on a specfic npid
LIBDN_API NPAsync<NPGetConnectionDetailsResult>* LIBDN_CALL DN_GetConnectionDetails(NPID npID);

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
LIBDN_API NPAsync<NPFriendResult>* LIBDN_CALL DN_GetFriends();

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
LIBDN_API NPAsync<NPGetExtProfileDataResult>* LIBDN_CALL DN_GetExtProfileData(uint32_t numIDs, const char* profileType, const NPID* npIDs, NPExtProfileData* outData);

// gets an avatar for any client
LIBDN_API NPAsync<NPGetUserAvatarResult>* LIBDN_CALL DN_GetUserAvatar(int id, uint8_t* buffer, size_t bufferLength);

/// ---------------------------------------------------------
// Server list service
// ----------------------------------------------------------

// creates a remote session
LIBDN_API NPAsync<NPCreateSessionResult>* LIBDN_CALL DN_CreateSession(NPSessionInfo* data);

// updates a session
LIBDN_API NPAsync<EServersResult>* LIBDN_CALL DN_UpdateSession(NPSID sid, NPSessionInfo* data);

// deletes a session
LIBDN_API NPAsync<EServersResult>* LIBDN_CALL DN_DeleteSession(NPSID sid);

// refreshes the session list - tags are separated by single spaces
LIBDN_API NPAsync<bool>* LIBDN_CALL DN_RefreshSessions(NPDictionary& infos);

// gets the number of sessions
LIBDN_API int32_t LIBDN_CALL DN_GetNumSessions();

// gets a single session's info
LIBDN_API void LIBDN_CALL DN_GetSessionData(int32_t index, NPSessionInfo* out);

// ----------------------------------------------------------
// Messaging service
// ----------------------------------------------------------

// sends arbitrary data to a client
LIBDN_API void LIBDN_CALL DN_SendMessage(NPID npid, const uint8_t* data, uint32_t length);

// function to register a callback when a message is received
// arguments: source NPID, data, length
LIBDN_API void LIBDN_CALL DN_RegisterMessageCallback(void(__cdecl * callback)(NPID, const uint8_t*, uint32_t));

// ----------------------------------------------------------
// Offline service
// ----------------------------------------------------------

// Enables offline mode
LIBDN_API void LIBDN_CALL DN_GoOffline();

// ----------------------------------------------------------
// API service. For for example C#
// ----------------------------------------------------------

// Enables offline mode
LIBDN_API NPAuthenticateResult* LIBDN_CALL Proxy_DN_AuthenticateWithDetails(const char * username, const char* password);
*/