#pragma once
#include <string>

const uint64_t API_VERSION = 1;

class DNHealthResult {
public:
	bool ready;
	std::string message;
};