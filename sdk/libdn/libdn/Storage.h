#pragma once
// ----------------------------------------------------------
// Storage
// ----------------------------------------------------------
namespace libdn {
	enum EGetFileResult {
		GetFileResultOK = 0,
		GetFileResultNotFound = 1,
		GetFileResultNotAllowed = 2,
		GetFileResultServiceError = 3
	};

	enum EWriteFileResult {
		WriteFileResultOK = 0,
		WriteFileResultNotAllowed = 1,
		WriteFileResultServiceError = 2
	};

	class GetPublisherFileResult {
	public:
		// the request result
		EGetFileResult result;

		// the amount of bytes written to the buffer
		uint32_t size;

		// binary data result
		uint8_t* buffer;
	};

	class GetUserFileResult {
	public:
		// the request result
		EGetFileResult result;

		// the amount of bytes written to the buffer
		uint32_t size;

		// binary data result
		uint8_t* buffer;
	};

	class WriteUserFileResult {
	public:
		// the request result
		EWriteFileResult result;
	};

	// obtains a publisher file.
	LIBDN_API Async<GetPublisherFileResult>* LIBDN_CALL GetPublisherFile(const char* fileName, uint8_t* buffer, size_t length);

	// obtains a file from the remote per-user storage
	LIBDN_API Async<GetUserFileResult>* LIBDN_CALL GetUserFile(PeerID id, const char* fileName, uint8_t* buffer, size_t length);

	// uploads a file to the remote per-user storage
	LIBDN_API Async<WriteUserFileResult>* LIBDN_CALL WriteUserFile(PeerID id, const char* fileName, const uint8_t* buffer, size_t length);
}