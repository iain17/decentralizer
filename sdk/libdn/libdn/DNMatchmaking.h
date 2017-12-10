#pragma once
#include <map>

typedef uint64_t DNSID;

class DNSessionInfo {
public:
	DNID dnId;//DecentralizerId
	std::string pId;//PeerId
	
	DNSID sessionId;
	uint64_t type;
	std::string name;
	uint32_t address;
	uint16_t port;
	std::map<std::string, std::string> details;
};

class DNUpsertSessionResult
{
public:
	bool result;
	DNSID sessionId;
};
