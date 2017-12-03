#include "StdInc.h"
#include <WS2tcpip.h>
#include "RPCAsync.h"

struct local_reply
{
	int id;
	INPRPCMessage* message;
	DWORD start;
};

// rpc state variables
static struct rpc_state_s
{
	// client socket
	SOCKET socket;

	WSAEVENT hSocketEvent;

	// message parsing
	uint32_t totalBytes;
	uint32_t readBytes;
	uint32_t messageType;
	uint32_t messageID;
	uint8_t* messageBuffer;

	// connected flag
	bool initialized = false;
	bool connected;
	bool wasConnected;
	bool reconnecting;
	bool offline = false;

	// message dispatching
	std::vector<rpc_dispatch_handler_s> dispatchHandlers;
	std::vector<NPRPCAsync*> freeNext;
	std::map<int, NPRPCAsync*> asyncHandlers;

	// thread management
	HANDLE hThread;
	HANDLE hShutdownEvent;

	CRITICAL_SECTION recvCritSec;

	// global message id
	LONG sendMessageID;
} g_rpc;

// code
//void Authenticate_Reauthenticate();

static DWORD WINAPI RPC_HandleReconnect(LPVOID param)
{
	while (!g_rpc.connected)
	{
		Log_Print("Reconnecting to RPC...\n");

		Sleep(5000);

		RPC_Shutdown();
		RPC_Init();

		if (DN_Connect(g_np.serverHost, g_np.serverPort))
		{
			//Authenticate_Reauthenticate();
		}
	}

	g_rpc.reconnecting = false;

	return 0;
}

void RPC_Reconnect()
{
	if (g_rpc.offline)
		return;

	if (!g_rpc.wasConnected)
	{
		return;
	}

	if (g_rpc.reconnecting)
	{
		return;
	}

	g_rpc.connected = false;
	g_rpc.reconnecting = true;
	CreateThread(NULL, 0, RPC_HandleReconnect, NULL, NULL, NULL);
}

int RPC_GenerateID()
{
	InterlockedIncrement(&g_rpc.sendMessageID);
	return g_rpc.sendMessageID;
}

void RPC_DispatchMessage(INPRPCMessage* message)
{
	uint32_t type = message->GetType();

	Log_Print("Dispatching RPC message with ID %d and type %d.\n", g_rpc.messageID, type);

	for (std::vector<rpc_dispatch_handler_s>::iterator i = g_rpc.dispatchHandlers.begin(); i != g_rpc.dispatchHandlers.end(); i++)
	{
		if (i->type == type)
		{
			Log_Print("Calling dispatch handler.\n");
			i->callback(message);
		}
	}

	if (g_rpc.messageID > 0 && g_rpc.asyncHandlers[g_rpc.messageID] == NULL)
	{
		Sleep(100);
	}

	if (g_rpc.asyncHandlers[g_rpc.messageID] != NULL)
	{
		Log_Print("Calling async dispatch handler %d and type %d.\n", g_rpc.messageID, type);

		NPRPCAsync* async = g_rpc.asyncHandlers[g_rpc.messageID];
		async->SetResult(message);
		g_rpc.freeNext.push_back(async);
	}

	Log_Print("Finished.\n");
}

static void RPC_HandleMessage()
{
	Log_Print("RPC_HandleMessage: type %d\n", g_rpc.messageType);

	for (int i = 0; i < NUM_RPC_MESSAGE_TYPES; i++)
	{
		if (g_rpcMessageTypes[i].type == g_rpc.messageType)
		{
			INPRPCMessage* message = g_rpcMessageTypes[i].handler();
			message->Deserialize(g_rpc.messageBuffer, g_rpc.totalBytes);
			RPC_DispatchMessage(message);
			message->Free();
		}
	}

	delete[] g_rpc.messageBuffer;
}

static uint32_t facefeed = 0xFACEFEED;
static uint8_t backBuffer[16][2048];
static uint32_t curBackBufferIdx;

bool RPC_ParseMessage(uint8_t* buffer, size_t len)
{
	uint8_t* origin = buffer;
	uint32_t read = len;

	memset(backBuffer[curBackBufferIdx], 0, 2048);
	memcpy(backBuffer[curBackBufferIdx], buffer, len);
	curBackBufferIdx = (curBackBufferIdx + 1) % 16;

	if (*(DWORD*)buffer == 0xFACEFEED)
	{
		send(g_rpc.socket, (const char*)&facefeed, 4, 0);
		return false;
	}

	//Log_Print("[pd] got new buf of size %i\n", len);

	while (read > 0)
	{
		// if we've not read any prior part of this packet before
		if (g_rpc.readBytes == 0)
		{
			rpc_message_header_s* message = (rpc_message_header_s*)origin;

			//Fix nullprt messages
			while (&message == nullptr || &message->signature == nullptr) {
				if (read <= 0)
					return false;
				origin += 4;
				read -= 4;

				message = (rpc_message_header_s*)origin;
			}

			if (message->signature == 0xFACEFEED)
			{
				send(g_rpc.socket, (const char*)&facefeed, 4, 0);
				return false;
			}

			if (message->signature != 0xDEADC0DE)
			{
				Log_Print("Signature (0x%x) doesn't match\n", message->signature);
				return false;
			}

			// set up size/buffer
			g_rpc.totalBytes = message->length;
			g_rpc.readBytes = 0;

			g_rpc.messageBuffer = new uint8_t[message->length];
			g_rpc.messageType = message->type;
			g_rpc.messageID = message->id;

			memset(g_rpc.messageBuffer, 0, message->length);

			// skip the header
			origin = &origin[sizeof(rpc_message_header_s)];
			read -= sizeof(rpc_message_header_s);

			//Log_Print("[pd] new msg: %i totalBytes, msg %i (origin %p)\n", g_rpc.totalBytes, g_rpc.messageID, (origin - buffer));
		}

		//Log_Print("[pd] %i read, %i readBytes, %i totalBytes, msg %i, start %08x, nigiro %p\n", read, g_rpc.readBytes, g_rpc.totalBytes, g_rpc.messageID, *(DWORD*)&origin[0], (origin - buffer));

		int copyLen = min(read, (g_rpc.totalBytes - g_rpc.readBytes));
		memcpy(&g_rpc.messageBuffer[g_rpc.readBytes], origin, copyLen);
		g_rpc.readBytes += copyLen;
		read -= copyLen; //
		origin += copyLen;

		if (g_rpc.readBytes >= g_rpc.totalBytes)
		{
			//Log_Print("[pd] full msg: %i totalBytes, msg %i\n", g_rpc.totalBytes, g_rpc.messageID);

			g_rpc.readBytes = 0;

			// handle the message here
			RPC_HandleMessage();
		}
	}

	return true;
}

// reads a message from the RPC socket
static int RPC_ReadMessage()
{
	if (!g_rpc.initialized)
		return 0;

	EnterCriticalSection(&g_rpc.recvCritSec);

	uint8_t buffer[2048];
	int len = recv(g_rpc.socket, (char*)buffer, sizeof(buffer), 0);

	if (len > 0)
	{
		int retval = RPC_ParseMessage(buffer, len) ? FD_WRITE : 0;

		LeaveCriticalSection(&g_rpc.recvCritSec);

		return retval;
	}

	if (len == SOCKET_ERROR)
	{
		if (WSAGetLastError() != WSAEWOULDBLOCK)
		{
			LeaveCriticalSection(&g_rpc.recvCritSec);
			return FD_CLOSE;
		}
	}

	if (len == 0)
	{
		LeaveCriticalSection(&g_rpc.recvCritSec);
		return FD_CLOSE;
	}

	LeaveCriticalSection(&g_rpc.recvCritSec);
	return 0;
}

#define RPC_EVENT_SHUTDOWN (WAIT_OBJECT_0)
#define RPC_EVENT_SOCKET (WAIT_OBJECT_0 + 1)

static uint32_t RPC_WaitForEvent()
{
	HANDLE waitHandles[] =
	{
		g_rpc.hShutdownEvent, // shutdown event should be first, should other events be signaled at the same time a shutdown occurs
		g_rpc.hSocketEvent
	};

	return WaitForMultipleObjects(2, waitHandles, FALSE, 0);
}

static uint32_t RPC_DetermineEvent()
{
	WSANETWORKEVENTS networkEvents;
	WSAEnumNetworkEvents(g_rpc.socket, g_rpc.hSocketEvent, &networkEvents);

	if (networkEvents.lNetworkEvents & FD_CLOSE)
	{
		return FD_CLOSE;
	}
	else if (networkEvents.lNetworkEvents & FD_READ)
	{
		return FD_READ;
	}

	return 0;
}

static bool RPC_RunThread()
{
	if (g_rpc.offline)
		return true;
	int result = 0;

	do
	{
		result = RPC_ReadMessage();

		if (result == FD_READ)
		{
			RPC_HandleMessage();
		}
	} while (result == FD_WRITE);

	if (result == FD_CLOSE)
	{
		RPC_Reconnect();
	}

	return true;
}

bool RPC_Init()
{
	Log_Print("Initializing RPC\n");

	// startup Winsock
	WSADATA data;
	int result;

	if ((result = WSAStartup(MAKEWORD(2, 2), &data)))
	{
		Log_Print("Couldn't initialize Winsock (0x%x)\n", result);
		return false;
	}

	InitializeCriticalSection(&g_rpc.recvCritSec);

	g_rpc.initialized = true;
	return true;
}

void RPC_Shutdown()
{
	Log_Print("Shutting down RPC\n");

	if (g_rpc.hThread)
	{
		if (g_rpc.hShutdownEvent)
		{
			SetEvent(g_rpc.hShutdownEvent);
			WaitForSingleObject(g_rpc.hThread, INFINITE);
		}

		TerminateThread(g_rpc.hThread, 1);

		CloseHandle(g_rpc.hThread);
		CloseHandle(g_rpc.hShutdownEvent);

		g_rpc.hThread = NULL;
		g_rpc.hShutdownEvent = NULL;
	}

	if (g_rpc.socket)
	{
		closesocket(g_rpc.socket);
		g_rpc.socket = NULL;
	}

	if (g_rpc.hSocketEvent)
	{
		WSACloseEvent(g_rpc.hSocketEvent);
		g_rpc.hSocketEvent = NULL;
	}
	g_rpc.initialized = false;

	//g_rpc.dispatchHandlers.clear();

	WSACleanup();
}

void RPC_RegisterDispatch(uint32_t type, DispatchHandlerCB callback)
{
	rpc_dispatch_handler_s handler;
	memset(&handler, 0, sizeof(handler));
	handler.callback = callback;
	handler.type = type;

	g_rpc.dispatchHandlers.push_back(handler);
}

void RPC_SendMessage(INPRPCMessage* message, int id)
{
	// serialize the message
	size_t len;
	uint8_t* data = message->Serialize(&len, id);

	if (g_rpc.offline) {
		Log_Print("Sending local RPC message with ID %d.\n", id);
		return;
	}

	// send to the socket
	send(g_rpc.socket, (const char*)data, len, 0);

	// free the data pointer
	message->FreePayload();

	// log it
	Log_Print("Sending RPC message with ID %d.\n", id);
}

void RPC_SendMessage(INPRPCMessage* message)
{
	RPC_SendMessage(message, 0);
}

NPAsync<INPRPCMessage>* RPC_SendMessageAsync(INPRPCMessage* message)
{
	int id = RPC_GenerateID();

	NPRPCAsync* async = new NPRPCAsync();
	async->SetAsyncID(id);

	g_rpc.asyncHandlers[id] = async;

	RPC_SendMessage(message, id);

	Log_Print("Sending async RPC message with ID %d.\n", id);

	return async;
}

bool dispatching = false;
void RPC_RunFrame()
{
	// poll the socket
	RPC_RunThread();

	// free to-free items
	/*std::vector<int> removeIDs;

	for (auto& i : g_rpc.asyncHandlers)
	{
	if (!i.second/* || i.second->HandleTimeout()* /)
	{
	if (i.second)
	{
	i.second->Free();
	}

	removeIDs.push_back(i.first);
	}
	}

	for (auto& i : removeIDs)
	{
	g_rpc.asyncHandlers.erase(i);
	}*/

	/*for (auto i = g_rpc.freeNext.begin(); i != g_rpc.freeNext.end(); i++)
	{
	g_rpc.asyncHandlers.erase((*i)->GetAsyncID());

	//(*i)->Free();
	}*/

	g_rpc.freeNext.clear();
}

LIBDN_API bool LIBDN_CALL DN_Connect(const char* serverAddr, uint16_t port) {
	if (g_rpc.connected || g_rpc.offline)
	{
		return true;
	}

	Log_Print("Connecting to the network using background agent on %s:%d.\n", serverAddr, port);

	if (!g_rpc.initialized)
		return false;

	// code to connect to some server
	hostent* hostEntity = gethostbyname(serverAddr);

	if (hostEntity == NULL)
	{
		Log_Print("Could not look up %s with error: %d.\n", serverAddr, WSAGetLastError());
		return false;
	}
	
	Log_Print("Resolved %s:%d.\n", serverAddr, port);
	sockaddr_in server;
	memset(&server, 0, sizeof(server));
	server.sin_family = AF_INET;
	server.sin_addr.s_addr = *(ULONG*)hostEntity->h_addr_list[0];
	server.sin_port = htons(port);

	if (connect(g_rpc.socket, (sockaddr*)&server, sizeof(sockaddr))) {
		u_long nonBlocking = TRUE;
		ioctlsocket(g_rpc.socket, FIONBIO, &nonBlocking);

		BOOL noDelay = TRUE;
		setsockopt(g_rpc.socket, IPPROTO_TCP, TCP_NODELAY, (char*)&noDelay, sizeof(BOOL));

		g_rpc.wasConnected = true;
		g_rpc.connected = true;
		Log_Print("Connected.\n");

		return true;
	}
}