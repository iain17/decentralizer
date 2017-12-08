#include "StdInc.h"
#include <WS2tcpip.h>
#include "RPCAsync.h"

static struct rpc_state_s
{
	// client socket
	SOCKET socket;
	WSAEVENT hSocketEvent;

	// message parsing
	uint32_t messageType;
	uint32_t messageID;

	// connected flag
	bool connected;
	bool wasConnected;
	bool reconnecting;

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

static void RPC_DispatchMessage(INPRPCMessage* message)
{
	uint32_t type = message->GetType();

	Log_Print("Dispatching RPC message with ID %d and type %d.\n", g_rpc.messageID, type);

	for (std::vector<rpc_dispatch_handler_s>::iterator i = g_rpc.dispatchHandlers.begin(); i != g_rpc.dispatchHandlers.end(); i++)
	{
		if (i->type == type && g_rpc.messageID == 0)
		{
			Log_Print("Dispatching RPC message to dispatch handler.\n");
			i->callback(message);
		}
	}

	if (g_rpc.messageID > 0 && g_rpc.asyncHandlers[g_rpc.messageID] == NULL)
	{
		Sleep(100);
	}

	if (g_rpc.asyncHandlers[g_rpc.messageID] != NULL)
	{
		Log_Print("Dispatching RPC message to async handler.\n");

		NPRPCAsync* async = g_rpc.asyncHandlers[g_rpc.messageID];
		async->SetResult(message);
		g_rpc.freeNext.push_back(async);
	}
}

static void RPC_HandleMessage(RPCMessage message)
{
	Log_Print("RPC_HandleMessage: type %d\n", message.msg_case());
	/*
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
	*/
}

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
int readBytes(SOCKET s, char* result, char delimiter) {
	while (true) {
		while (read < len) {
			const char byte = buffer[read];
			//Log_Print("reading %d of %d: %d.\n", read, len, byte);
			if (byte == delimiter) {
				//Log_Print("delimiter found: %d.\n", delimiter);
				int len = backBufferRead;
				memcpy(result, backBuffer, len);
				resetBackBuffer();
				read++;
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
	Log_Print("Parsed RPCMessage with id: %d \n", message.id());
	//Check version
	if (message.version() != VERSION) {
		Log_Print("Version mismatch v0x%x != v0x%x \n", message.version(), VERSION);
		return false;
	}

	g_rpc.messageType = message.msg_case();
	g_rpc.messageID = message.id();
	RPC_HandleMessage(message);
}

// reads a message from the RPC socket
static int RPC_ReadMessage()
{
	EnterCriticalSection(&g_rpc.recvCritSec);
	char buffer[MAXMESSAGESIZE];
	int len = readBytes(g_rpc.socket, buffer, DELIMITER);
	LeaveCriticalSection(&g_rpc.recvCritSec);

	if (len > 0)
	{
		int retval = RPC_ParseMessage(buffer, len) ? FD_WRITE : 0;
		return retval;
	}

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

/*static void RPC_SendMessage(int type, IRPCMessage* message)
{
std::string str = message->Serialize();
uint32_t datalen = str.length();
uint32_t buflen = sizeof(rpc_message_header_s) + datalen;

// allocate a response buffer and copy data to it
uint8_t* buffer = new uint8_t[buflen];
const char* data = str.c_str();
memcpy(&buffer[sizeof(rpc_message_header_s)], data, datalen);

// set the response header data
rpc_message_header_s* header = (rpc_message_header_s*)buffer;
header->signature = 0xDEADC0DE;
header->length = datalen;
header->type = type;

// send to the socket
send(g_rpc.socket, (const char*)buffer, buflen, 0);

// free the buffer
delete[] buffer;
}*/

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

/*
static void RPC_HandleHello(INPRPCMessage* message)
{
	RPCHelloMessage* hello = (RPCHelloMessage*)message;
	Log_Print("got %d %d %s %s\n", hello->GetBuffer()->number(), hello->GetBuffer()->number2(), hello->GetBuffer()->name().c_str(), hello->GetBuffer()->stuff().c_str());
}

static void RPC_HandleClose(INPRPCMessage* message)
{
	RPCCloseAppMessage* close = (RPCCloseAppMessage*)message;

	ExitProcess(0x76767676);
}
*/
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

	// create RPC socket event
	//g_rpc.hSocketEvent = WSACreateEvent();
	g_rpc.hShutdownEvent = CreateEvent(NULL, FALSE, FALSE, NULL);

	InitializeCriticalSection(&g_rpc.recvCritSec);

	// blah
	//RPC_RegisterDispatch(RPCHelloMessage::Type, RPC_HandleHello);
	//RPC_RegisterDispatch(RPCCloseAppMessage::Type, RPC_HandleClose);

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

void RPC_RunFrame()
{
	// poll the socket
	RPC_RunThread();

	g_rpc.freeNext.clear();
}

LIBDN_API bool LIBDN_CALL DN_Connect(const char* serverAddr, uint16_t port)
{
	if (g_rpc.connected)
	{
		return true;
	}

	// store server name and port
	strncpy(g_np.serverHost, serverAddr, sizeof(g_np.serverHost));
	g_np.serverPort = port;

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