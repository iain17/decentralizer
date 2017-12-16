// example.cpp : Defines the entry point for the console application.
//

#include "stdafx.h"
#include "libdn.h"
#include <iostream>

void LogCB(const char* message) {
	printf("[DN] %s\n", message);
}

using namespace libdn;

bool sessionCreated = false;
void matchMakingTest() {
	if (!sessionCreated) {
		//Create session.
		Session* session = new Session();
		session->name = "Big tests";
		session->type = 1337;
		session->port = 8080;
		session->details["cool"] = "yes";
		auto request = UpsertSession(session);
		//request->fail(LogRPC);
		request->then([](UpsertSessionResult result) {
			printf("session id: %d\n", result.sessionId);
			sessionCreated = true;
		});
		request->fail(LogCB);
		request->wait();
	}
	//Get all sessions with type 1337 and all session that are cool.
	auto request = GetNumSessions(1337, "cool", "yes");
	request->fail(LogCB);
	if (!request->wait()) {
		return;
	}
	int num = request->get();
	printf("Received %i session ids\n", num);
	//For each session id we have received back. Fetch it.
	for(int i = 0; i < num; i++) {
		//Fetch that session
		auto session = GetSessionByIndex(i);
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

int main() {
	printf("DN_Init()\n");
	libdn::Init(LogCB);
	bool status = false;
	while (!status) {
		status = libdn::Connect("10.1.1.34:50051");//10.1.1.34
	}

	while (true) {
		libdn::RunFrame();
		matchMakingTest();
		Sleep(100);
	}
}