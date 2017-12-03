// This file has been automatically generated.

#include "../../libconnect/StdInc.h"
#include "MessageDefinition.h"

HelloRequest* RPCHelloRequest::GetBuffer()
{
	return _buffer.GetPayload();
}

void RPCHelloRequest::Deserialize(const uint8_t* message, size_t length)
{
	_buffer.Deserialize(message, length);
}

uint8_t* RPCHelloRequest::Serialize(size_t* length, uint32_t id)
{
	if (_payload)
	{
		return _payload;
	}

	std::string str = _buffer.Serialize();
	uint32_t type = this->GetType();
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
	header->id = id;
	
	// set return stuff
	*length = buflen;
	_payload = buffer;
	
	return _payload;
}

int RPCHelloRequest::GetType()
{
	return 1000;
}

void RPCHelloRequest::Free()
{
	if (_payload != NULL)
	{
		delete[] _payload;
	}
	
	delete this;
}

void RPCHelloRequest::FreePayload()
{
	if (_payload != NULL)
	{
		delete[] _payload;
		_payload = NULL;
	}
}

RPCHelloRequest* RPCHelloRequest::Create()
{
	return new RPCHelloRequest();
}

HelloReply* RPCHelloReply::GetBuffer()
{
	return _buffer.GetPayload();
}

void RPCHelloReply::Deserialize(const uint8_t* message, size_t length)
{
	_buffer.Deserialize(message, length);
}

uint8_t* RPCHelloReply::Serialize(size_t* length, uint32_t id)
{
	if (_payload)
	{
		return _payload;
	}

	std::string str = _buffer.Serialize();
	uint32_t type = this->GetType();
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
	header->id = id;
	
	// set return stuff
	*length = buflen;
	_payload = buffer;
	
	return _payload;
}

int RPCHelloReply::GetType()
{
	return 1001;
}

void RPCHelloReply::Free()
{
	if (_payload != NULL)
	{
		delete[] _payload;
	}
	
	delete this;
}

void RPCHelloReply::FreePayload()
{
	if (_payload != NULL)
	{
		delete[] _payload;
		_payload = NULL;
	}
}

RPCHelloReply* RPCHelloReply::Create()
{
	return new RPCHelloReply();
}


rpc_message_type_s g_rpcMessageTypes[] = 
{
	{ 1000, (CreateMessageCB)&RPCHelloRequest::Create },
	{ 1001, (CreateMessageCB)&RPCHelloReply::Create },
};