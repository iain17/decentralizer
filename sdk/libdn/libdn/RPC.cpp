#include "StdInc.h"
#include <WS2tcpip.h>
#include "RPCAsync.h"

static struct rpc_state_s
{
	// client socket
	SOCKET socket;
	WSAEVENT hSocketEvent;

	// connected flag
	bool connected;
	bool wasConnected;
	bool reconnecting;

	// message dispatching
	std::vector<rpc_dispatch_handler_s> dispatchHandlers;
	std::vector<NPRPCAsync*> freeNext;
	std::map<uint64_t, NPRPCAsync*> asyncHandlers;

	// thread management
	HANDLE hThread;
	CRITICAL_SECTION recvCritSec;
	
	// global message id
	uint64_t sendMessageID;
} g_rpc;

//Will hang until we are connected and DN is ready.
LIBDN_API void LIBDN_CALL DN_WaitUntilReady()
{
	DNHealthResult* health;
	health->ready = false;
	while (!g_rpc.connected || health == nullptr || !health->ready) {
		health = DN_Health();
		if (health != nullptr && health->ready) {
			break;
		}
		Sleep(100);
	}
}

static DWORD WINAPI RPC_HandleReconnect(LPVOID param)
{
	while (!g_rpc.connected)
	{
		Log_Print("Reconnecting to RPC...\n");

		Sleep(5000);

		RPC_Shutdown();
		RPC_Init();
		if (DN_Connect(g_dn.serverHost, g_dn.serverPort))
		{
			Log_Print("Connected to RPC.\n");
			g_rpc.sendMessageID = 0;
			//Authenticate_Reauthenticate();
		}
	}

	g_rpc.reconnecting = false;

	return 0;
}

void RPC_Reconnect()
{
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

uint64_t RPC_GenerateID()
{
	InterlockedIncrement(&g_rpc.sendMessageID);
	return g_rpc.sendMessageID;
}

static void RPC_DispatchMessage(RPCMessage* message)
{
	Log_Print("Dispatching RPC message with ID %d.\n", message->id());

	for (std::vector<rpc_dispatch_handler_s>::iterator i = g_rpc.dispatchHandlers.begin(); i != g_rpc.dispatchHandlers.end(); i++)
	{
		if (i->type == message->msg_case() && message->id() == 0)
		{
			Log_Print("Dispatching RPC message to dispatch handler.\n");
			i->callback(message);
		}
	}

	if (message->id() > 0 && g_rpc.asyncHandlers[message->id()] == NULL)
	{
		Sleep(100);
	}

	if (g_rpc.asyncHandlers[message->id()] != NULL)
	{
		Log_Print("Dispatching RPC message to async handler.\n");

		NPRPCAsync* async = g_rpc.asyncHandlers[message->id()];
		async->SetResult(message);
		g_rpc.freeNext.push_back(async);
	}
}

BYTE DELIMITER[3] = { '\n', '\r', '\n' };
const uint32_t MAXMESSAGESIZE = 32768;

char buffer[2048];
static char backBuffer[MAXMESSAGESIZE];
static int backBufferRead = 0;
static int read = 0;
static int len = 0;

void resetBackBuffer() {
	memset(backBuffer, 0, sizeof(backBuffer));
	backBufferRead = 0;
}

//Reads all of the bytes until delimiter
int readBytes(SOCKET s, char* result, char* delimiter) {
	while (true) {
		while (read < len) {
			const char byte = buffer[read];
			//Log_Print("reading %d of %d: %d.\n", read, len, byte);
			if (read + 2 < len && buffer[read] == delimiter[0] && buffer[read + 1] == delimiter[1] && buffer[read + 2] == delimiter[2]) {
				//Log_Print("delimiter found: %d.\n", delimiter);
				int len = backBufferRead;
				memcpy(result, backBuffer, len);
				resetBackBuffer();
				read += 3;
				return len;
			}
			//Log_Print("copying %d.\n", byte);
			backBuffer[backBufferRead] = byte;
			read++;
			backBufferRead++;
			//Overflow.
			if (backBufferRead > MAXMESSAGESIZE) {
				backBufferRead = 0;
			}
		}
		//Log_Print("Receiving...\n", read, len);
		len = recv(s, buffer, sizeof(buffer), 0);
		read = 0;
		if (len <= 0) {
			return len;
		}
	}
	return len;
}

bool RPC_ParseMessage(char* buffer, size_t len)
{
	RPCMessage message;
 	if (!message.ParseFromArray(buffer, len)) {
		Log_Print("Failed to parse RPCMessage.\n");
		return false;
	}
	Log_Print("Received reply to RPCMessage with id: %d \n", message.id());
	//Check version
	if (message.version() != API_VERSION) {
		Log_Print("Version mismatch v0x%x != v0x%x \n", message.version(), API_VERSION);
		return false;
	}

	RPC_DispatchMessage(&message);
}

// reads a message from the RPC socket
static int RPC_ReadMessage()
{
	EnterCriticalSection(&g_rpc.recvCritSec);
	char buffer[MAXMESSAGESIZE];
	int len = readBytes(g_rpc.socket, buffer, (char*)DELIMITER);

	if (len > 0)
	{
		int retval = RPC_ParseMessage(buffer, len) ? FD_WRITE : 0;
		LeaveCriticalSection(&g_rpc.recvCritSec);
		return retval;
	}
	LeaveCriticalSection(&g_rpc.recvCritSec);

	if (len == SOCKET_ERROR)
	{
		if (WSAGetLastError() != WSAEWOULDBLOCK)
		{
			return FD_CLOSE;
		}
	}

	if (len == 0)
	{
		return FD_CLOSE;
	}
	return 0;
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
	int result = 0;
	do
	{
		result = RPC_ReadMessage();

		if (result == FD_READ)
		{
			//RPC_HandleMessage();
		}
	} while (result == FD_WRITE);

	if (result == FD_CLOSE)
	{
		RPC_Reconnect();
	}

	return true;
}

void resetBackBuffer();
bool RPC_Init()
{
	resetBackBuffer();
	Log_Print("Initializing RPC\n");

	// startup Winsock
	WSADATA data;
	int result;

	if ((result = WSAStartup(MAKEWORD(2, 2), &data)))
	{
		Log_Print("Couldn't initialize Winsock (0x%x)\n", result);
		return false;
	}
	// create RPC socket
	g_rpc.socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

	if (g_rpc.socket == INVALID_SOCKET)
	{
		Log_Print("Couldn't create RPC socket\n");
		RPC_Shutdown();
		return false;
	}
	InitializeCriticalSection(&g_rpc.recvCritSec);

	return true;
}

void RPC_Shutdown()
{
	Log_Print("Shutting down RPC\n");

	if (g_rpc.hThread)
	{
		TerminateThread(g_rpc.hThread, 1);
		CloseHandle(g_rpc.hThread);
		g_rpc.hThread = NULL;
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

	for (auto const &ent1 : g_rpc.asyncHandlers) {
		NPRPCAsync* async = ent1.second;
		async->SetResult(NULL);
	}

	g_rpc.dispatchHandlers.clear();

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

//todo: memory leak?
bool RPC_SendMessage(RPCMessage message, int id) {
	EnterCriticalSection(&g_rpc.recvCritSec);
	message.set_id(id);
	message.set_version(API_VERSION);
	try {
		char buffer[MAXMESSAGESIZE];
		int size = message.ByteSizeLong();
		bool res = message.SerializeToArray((void *)buffer, size);
		if (!res) {
			Log_Print("Failed to properly serialize protobuf message.\n");
			LeaveCriticalSection(&g_rpc.recvCritSec);
			return false;
		}

		// send to the socket
		send(g_rpc.socket, (const char*)buffer, size, 0);

		//Send delimiter
		send(g_rpc.socket, (const char*)DELIMITER, 3, 0);

		// cleanup
		//message->Clear();
		//delete buffer;
	}
	catch (...) {
		return false;
	}

	// log it
	Log_Print("Sending RPC message with ID %d.\n", id);
	LeaveCriticalSection(&g_rpc.recvCritSec);
	return true;
}

bool RPC_SendMessage(RPCMessage* message) {
	return RPC_SendMessage(*message, 0);
}

DNAsync<RPCMessage>* RPC_SendMessageAsync(RPCMessage* message) {
	if (message->msg_case() != RPCMessage::MsgCase::kHealthRequest) {
		DN_WaitUntilReady();
	}
	else {
		while (!g_rpc.connected) {
			Sleep(100);
		}
	}

	uint64_t id = RPC_GenerateID();

	NPRPCAsync* async = new NPRPCAsync();
	async->SetAsyncID(id);

	g_rpc.asyncHandlers[id] = async;

	bool res = RPC_SendMessage(*message, id);
	if (!res) {
		async->SetResult(NULL);
	}

	Log_Print("Sending async RPC message with ID %d.\n", id);

	return async;
}

DNAsync<RPCMessage>* RPC_SendMessageAsyncCache(std::string key, RPCMessage* message) {
	/*
	if (g_rpc.cache->contains(key)) {
		Log_Print("Calling back Async RPC call from cache using key %s.\n", key.c_str());
		NPRPCAsync* async = new NPRPCAsync();
		async->SetResult(g_rpc.cache->lookup(key));
		return async;
	}
	*/
	return RPC_SendMessageAsync(message);
}

void RPC_RunFrame()
{
	// poll the socket
	RPC_RunThread();

	g_rpc.freeNext.clear();
}

LIBDN_API bool LIBDN_CALL DN_Connect(const char* serverAddr, uint16_t port) {
	if (g_rpc.connected)
	{
		return true;
	}

	// store server name and port
	strncpy(g_dn.serverHost, serverAddr, sizeof(g_dn.serverHost));
	g_dn.serverPort = port;

	// code to connect to some server
	hostent* hostEntity = gethostbyname(serverAddr);

	if (hostEntity == NULL)
	{
		Log_Print("Could not look up %s: %d.\n", serverAddr, WSAGetLastError());
		return false;
	}

	sockaddr_in server;
	memset(&server, 0, sizeof(server));
	server.sin_family = AF_INET;
	server.sin_addr.s_addr = *(ULONG*)hostEntity->h_addr_list[0];
	server.sin_port = htons(port);

	if (connect(g_rpc.socket, (sockaddr*)&server, sizeof(sockaddr)))
	{
		Log_Print("Connecting failed: %d\n", WSAGetLastError());
		return false;
	}

	u_long nonBlocking = TRUE;
	ioctlsocket(g_rpc.socket, FIONBIO, &nonBlocking);

	BOOL noDelay = TRUE;
	setsockopt(g_rpc.socket, IPPROTO_TCP, TCP_NODELAY, (char*)&noDelay, sizeof(BOOL));

	g_rpc.wasConnected = true;
	g_rpc.connected = true;

	return true;
}