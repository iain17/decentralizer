// This file has been automatically generated.

#include "../../libdn/StdInc.h"
#include "MessageDefinition.h"

HealthRequest* RPCHealthRequest::GetBuffer()
{
	return _buffer.GetPayload();
}

void RPCHealthRequest::Deserialize(const uint8_t* message, size_t length)
{
	_buffer.Deserialize(message, length);
}

uint8_t* RPCHealthRequest::Serialize(size_t* length, uint32_t id)
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

int RPCHealthRequest::GetType()
{
	return 1000;
}

void RPCHealthRequest::Free()
{
	if (_payload != NULL)
	{
		delete[] _payload;
	}
	
	delete this;
}

void RPCHealthRequest::FreePayload()
{
	if (_payload != NULL)
	{
		delete[] _payload;
		_payload = NULL;
	}
}

RPCHealthRequest* RPCHealthRequest::Create()
{
	return new RPCHealthRequest();
}

HealthReply* RPCHealthReply::GetBuffer()
{
	return _buffer.GetPayload();
}

void RPCHealthReply::Deserialize(const uint8_t* message, size_t length)
{
	_buffer.Deserialize(message, length);
}

uint8_t* RPCHealthReply::Serialize(size_t* length, uint32_t id)
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

int RPCHealthReply::GetType()
{
	return 1001;
}

void RPCHealthReply::Free()
{
	if (_payload != NULL)
	{
		delete[] _payload;
	}
	
	delete this;
}

void RPCHealthReply::FreePayload()
{
	if (_payload != NULL)
	{
		delete[] _payload;
		_payload = NULL;
	}
}

RPCHealthReply* RPCHealthReply::Create()
{
	return new RPCHealthReply();
}


rpc_message_type_s g_rpcMessageTypes[] = 
{
	{ 1000, (CreateMessageCB)&RPCHealthRequest::Create },
	{ 1001, (CreateMessageCB)&RPCHealthReply::Create },
};