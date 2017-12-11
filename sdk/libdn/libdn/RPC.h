#pragma once

// dispatch handler callback
typedef void(*DispatchHandlerCB)(RPCMessage*);

struct rpc_dispatch_handler_s
{
	uint32_t type;
	uint32_t id;
	DispatchHandlerCB callback;
};

// initialize the RPC system
bool RPC_Init();

// shut down the RPC system
void RPC_Shutdown();

// register a dispatch handler
void RPC_RegisterDispatch(uint32_t type, DispatchHandlerCB callback);

// send a message
bool RPC_SendMessage(RPCMessage* message);
Async<RPCMessage>* RPC_SendMessageAsync(RPCMessage* message);
Async<RPCMessage>* RPC_SendMessageAsyncCache(std::string key, RPCMessage* message);

// increment and return a new sequence ID
uint64_t RPC_GenerateID();

// initialize storage service RPC components
void Storage_Init();

// initialize messaging service RPC components
void Messaging_Init();

// run RPC frame
void RPC_RunFrame();