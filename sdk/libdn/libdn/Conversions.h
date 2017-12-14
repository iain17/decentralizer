#pragma once
namespace libdn {
	pb::Session* DNSessionToPBSession(libdn::Session * dnInfo);
	libdn::Session* PBSessionToDNSession(pb::Session * pbInfo);
}