// example.cpp : Defines the entry point for the console application.
//

#include "stdafx.h"
#include "libdn.h"
#include <iostream>

const char* NETWORKKEY = "2d2d2d2d2d424547494e20525341205055424c4943204b45592d2d2d2d2d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b4341514541736364546c7a386669314338504a38436e386c380a493245726534494e79424264663679706d694e64794554554d34304f4338462b44355376775166594d51514d614161412b65527a715a7279354a785439596a660a7642467a3873313134597a796a6b743853463666434534686778555156564e30326e416e396c4d525164683567433859513641724e7a43774949316a7a4a6d670a45506154563757552f4b6e4d6753476b583941754c5431692b57552b5174695158687964676745684835506f7855514c4b58325766434769476b6871317a41370a776c5559564f78615763553259596a714f657456683063367646476b53476b4263794d6e725652445771705954533451582b38736837444e796f326b463632330a49464f682f31364f784354684f4c2f674366666c455a767862464e612b6f50754f446463456b36326658486b71484947613267353178324f665377752f474c640a4c514944415141420a2d2d2d2d2d454e4420525341205055424c4943204b45592d2d2d2d2d0a";
libdn::Peer* self;

void LogCB(const char* message) {
	printf("[DN] %s\n", message);
}

void MessageCB(libdn::PeerID& from, const uint8_t* data, uint32_t size) {
	printf("[DN] Received direct message from %s: %s with a size of %i\n", from.c_str(), data, size);
}

void createSession() {
	LogCB("createSession");
	//Create session.
	libdn::Session session;
	session.name = "Big tests";
	session.type = 2117;
	session.port = 8080;
	session.details["cool"] = "yes";
	auto request = libdn::UpsertSession(session);
	request->then([](libdn::DNSID sessionId) {
		printf("Created session. id = %llx\n", sessionId);
	});
	request->wait();
}

void findSessions() {
	LogCB("findSessions");
	//Get all sessions with type 1337 and all session that are cool.
	auto request = libdn::GetNumSessions(2117, "cool", "yes");
	if (request->wait()) {
		int num = request->get();
		printf("Received %i session ids\n", num);
		//For each session id we have received back. Fetch it.
		for (int i = 0; i < num; i++) {
			//Fetch that session
			auto session = libdn::GetSessionByIndex(i);
			if (session != nullptr) {
				unsigned char address[4];

				address[0] = session->address & 0xFF;
				address[1] = (session->address >> 8) & 0xFF;
				address[2] = (session->address >> 16) & 0xFF;
				address[3] = (session->address >> 24) & 0xFF;
				//Print something cool
				printf("Session %s on address %i.%i.%i.%i:%i\n", session->name.c_str(), address[3], address[2], address[1], address[0], session->port);
			}

		}
	}
	printf("-----\n");
	Sleep(3000);
}

void matchMakingTest() {
	LogCB("matchMakingTest");
	createSession();
	findSessions();
}

void writePeerFile(std::string expected) {
	LogCB("writePeerFile");
	auto request = libdn::WritePeerFile("test.txt", expected);
	if (request->wait()) {
		LogCB("Storage write succeeded!");
	} else {
		LogCB("Storage write failed!");
	}
}

void getPeerFile(std::string expected) {
	LogCB("getPeerFile");
	std::string pid = "self";
	auto request = libdn::GetPeerFile(pid, "test.txt");
	if (request->wait()) {
		if (request->get().compare(expected)) {
			LogCB("Storage get failed!");
		} else {
			LogCB("Storage get succeeded!");
		}
	} else {
		LogCB("Storage get failed!");
	}
}

void getSelf() {
	LogCB("getSelf");
	std::string pid = "self";
	while (self == nullptr) {
		self = libdn::GetSelf();
		if (self == nullptr) {
			LogCB("Could not get self?!?");
		} else {

			//Update my name
			self->details["name"] = "example";
			auto request = libdn::UpsertPeer(*self);
			request->wait();
			self = libdn::GetSelf();

			printf("Self peer id = '%s' decentralized id = '%16llX' and name is %s\n", self->pId.c_str(), self->dnId, self->details["name"].c_str());
		
			//Resolve my decentralized id.
			auto peerId = libdn::ResolveDecentralizedId(self->dnId);
			if (!peerId->empty()) {
				printf("Resolving works %s\n", peerId->c_str());
			} else {
				printf("Resolving is fucked\n");
			}
		}
	}
}

void storageTest() {
	LogCB("storageTest");
	std::string expected = "hey there";
	writePeerFile(expected);
	getPeerFile(expected);
}

void messagingTest() {
	std::string to = "self";
	std::string data = "Cool";
	auto req = libdn::SendDirectMessage(1337, to, data);
	req->wait();
}

int main() {
	printf("DN_Init()\n");
	libdn::Init(LogCB);
	bool status = false;
	while (!status) {
		status = libdn::Connect("localhost", 50010, NETWORKKEY, false, false);
	}
	printf("Starting example app\n");
	libdn::RegisterDirectMessageCallback(1337, MessageCB);

	getSelf();

	while (true) {
		LogCB("frame");
		libdn::RunFrame();
		messagingTest();
		matchMakingTest();
		storageTest();
		Sleep(100);
	}
}