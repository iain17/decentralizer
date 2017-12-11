// example.cpp : Defines the entry point for the console application.
//

#include "stdafx.h"
#include "libdn.h"
#include <iostream>

void LogCB(const char* message) {
	printf("[NP] %s", message);
}
using namespace libdn;

bool sessionCreated = false;
void matchMakingTest() {
	if (!sessionCreated) {
		sessionCreated = true;
		//Create session.
		DNSessionInfo info;
		info.name = "Big tests";
		info.type = 1337;
		info.port = 8080;
		info.details["cool"] = "yes";
		auto request = DN_UpsertSession(&info);
		request->SetCallback([](DNAsync<DNUpsertSessionResult>* async) {
			DNUpsertSessionResult* answer = async->GetResult();
			printf("session id: %d\n", answer->sessionId);
		}, NULL);
	}

	//Get all sessions with type 1337 and all session that are cool.
	auto request = DN_GetNumSessions(1337, "cool", "yes");
	if (request->Wait(7500) != nullptr) {
		int* num = request->GetResult();
		printf("Received %i session ids\n", *num);
		//For each session id we have received back. Fetch it.
		for(int i = 0; i < *num; i++) {
			//Fetch that session
			auto session = DN_GetSessionByIndex(i);
			if (session != nullptr) {
				char address[4];

				address[0] = session->address & 0xFF;
				address[1] = (session->address >> 8) & 0xFF;
				address[2] = (session->address >> 16) & 0xFF;
				address[3] = (session->address >> 24) & 0xFF;
				//Print something cool
				printf("Session %s on address %i.%i.%i.%i:%i\n", session->name.c_str(), address[3], address[2], address[1], address[0], session->port);
			}

		}
		printf("-----\n");
		Sleep(3000);
	}
}

int main() {
	printf("DN_Init()\n");
	DN_Init(LogCB);
	bool status = false;
	while (!status) {
		status = DN_Connect("10.1.1.34", 3036);
	}

	while (true) {
		DN_RunFrame();
		matchMakingTest();
	}
}