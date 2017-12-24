// example.cpp : Defines the entry point for the console application.
//

#include "stdafx.h"
#include "libdn.h"
#include <iostream>

const char* NETWORKKEY = "2d2d2d2d2d424547494e205253412050524956415445204b45592d2d2d2d2d0a4d4949456f77494241414b434151454178655a54466a55454a6e6b66385a68355346522b4757786f74663857396f306a6b55613437574b52677a7270555871790a385241766c586f6566314d4d78686d5167523367383151686775312f386e61464a7a704f4e44414e70774d31676e3053323961705a6a696e6758415a705a39450a4672363645326234504c365a756a5136477a745333726d7567494a5a76506c4f6732524f6e44724f6753596c34326c446c5738507342643478415750543762680a7652756534795049613661434c4d7074324f4838694a2f6252686a67486d656c3537354549746a7164697870726769627a5a6c77485a437a74466e304b52755a0a574e70724b7a6156537456625a4f4477523641482f355356424f3749392b61353949786d31506a624b53444c625846766a6d74514c4f56574d6d6459413447740a447555453572713547322f417367654f632b70486b55346148566c467158666a75736e657577494441514142416f494241514338376151496f5668792b6c4f360a713975745a3678797a5149794c584e5973576c784f646b3246314866764a4165443074687842674a5665706a6c332b735a352b4442476c4c4939685354445a480a335570464a753664392f6f7770576d695235474865716d43517a6632354851336e354b375042346367384d64437346722b4977346a79775149616773573055630a6d633251746d5174316835725157594f6375486f655879366d63336249397875572f3739324b5677696e6c636d666371586c356365632b5131616934524d44300a5638415048473442685a774852614865576b7933794e35357a6d704a68506a34335443754f4c384b31452b476f64593276554633526c3068767a466f4453497a0a534e35385a2f55656c6e4d736e5175677a6553486d627955647971675748432b725865733251423159656f6e76373435537a6b773069366c50672f704d5443770a6e684e56552b5768416f4742414d3830534c423972586b795a59687352795048715a7576443150674355586f5852416b6b69686553475631436a542b325753670a763456504235735238432f4765753555444b4655383435623477545a686847467039703349646f5a5058564a3454536e5678763266502b725532515079702f6c0a715565524330596c393463326e764c4b666d7a377843446c3356663039577066384b596c424f726774746f436b4c4154576c4b5939433035416f4742415053420a4755613836447a594c64726a512f535854574a773069304464374f487765456a336e557268364f3330546b566157764e58565a5a5339423332536b42454736510a524c772f452f6a6a676c3048365348615854386a596743596f31537953323262447451534c6d55795369777a565a424164536e4f2b3476365844386c696346700a62736c4f384269617643503951434d6a5141306f58504b615a51685970537a6e3251583336682b54416f47415638526478655232527041435566634c617978330a75326f376f3975534566714b3850754d72577a435862646c79327a6e794b674f642f6b787a34325a6e364d4444314371794f75692f766f4d2f31446b6153656b0a4966573063523266327236676c68304c324e786674697872396b5a36486143365134593873456f457631467a6f6f514461555a376e545041766a4555677971410a564e2f355a5551714c383547573037585134566d614945436759414945372b3341345355686675313047347452565a4d477a67475071675571545a7862704c700a77663967484446774e6c48654f74474c6962576b644745624a71725a5444444a477a685972344e642b5758744e5636424f48554457676544513853554956777a0a43307133457873364c4a5032435073563333325632545a3036354f4b62535934786a2f4f51455a59316750705a542b3362343771674d6b33706c3447687234330a6f55483975774b4267473278577470673769336a322b6c5635476e6172443236486743316f5841624c2f30326467552f2f6c4f72726c39726e786836753369640a4e69715845534d69565253714936544778424a30714969517135767a57683173493239332b42345476457934376275574930656730753841447637533254576a0a5671384e35334863643644463138745131396c66556c32516c744c64384c55626155634b6a6634525a464d63645573316146426b0a2d2d2d2d2d454e44205253412050524956415445204b45592d2d2d2d2d0a";
libdn::Peer* self;

void LogCB(const char* message) {
	printf("[DN] %s\n", message);
}

void MessageCB(libdn::PeerID from, const uint8_t* data, uint32_t size) {
	printf("[DN] Received direct message from %s with a size of %i\n", from.c_str(), size);
}

void createSession() {
	LogCB("createSession");
	//Create session.
	libdn::Session* session = new libdn::Session();
	session->name = "Big tests";
	session->type = 2117;
	session->port = 8080;
	session->details["cool"] = "yes";
	auto request = libdn::UpsertSession(session);
	request->then([](libdn::UpsertSessionResult result) {
		printf("Created session. id = %llx\n", result.sessionId);
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
			self->details["name"] = "iain17";
			auto request = libdn::UpsertPeer(self);
			request->wait();
			self = libdn::GetSelf();

			printf("Self peer id = '%s' decentralized id = '%16llX' and name is %s\n", self->pId.c_str(), self->dnId, self->details["name"].c_str());
		
			//Resolve my decentralized id.
			libdn::PeerID* peerId = libdn::ResolveDecentralizedId(self->dnId);
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
	auto req = libdn::SendDirectMessage(to, data);
	req->wait();
}

int main() {
	printf("DN_Init()\n");
	libdn::Init(LogCB, MessageCB);
	bool status = false;
	while (!status) {
		status = libdn::Connect("10.1.1.34:50010", NETWORKKEY, true);
	}
	getSelf();

	while (true) {
		LogCB("frame");
		libdn::RunFrame();
		matchMakingTest();
		storageTest();
		messagingTest();
		Sleep(100);
	}
}