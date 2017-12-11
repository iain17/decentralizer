#pragma once
#include "StdInc.h"
#include "Async.h"

class RPCAsync : public Async<RPCMessage>
{
private:
	void(__cdecl* _callback)(Async<RPCMessage>*);
	RPCMessage* _result;
	void* _userData;

	void(__cdecl* _timeoutCallback)(Async<RPCMessage>*);
	unsigned int _timeout;
	DWORD _start;
	bool _completed;

public:
	RPCAsync()
	{
		_callback = NULL;
		_result = NULL;
		_userData = NULL;
		_timeout = -1;
		_timeoutCallback = nullptr;
		_start = 0;
		_completed = false;
	}

	// implementations for base Async
	virtual RPCMessage* Wait()
	{
		return Wait(-1);
	}

	virtual RPCMessage* Wait(unsigned int timeout)
	{
		DWORD start = GetTickCount();

		while (!HasCompleted())
		{
			RPC_RunFrame();

			Sleep(1);

			DWORD elapsed = GetTickCount() - start;

			if (elapsed > timeout)
			{
				return nullptr;
			}
		}

		return GetResult();
	}

	virtual bool HasCompleted()
	{
		return _completed;
	}

	virtual bool HasResult()
	{
		return (_result != NULL);
	}

	virtual RPCMessage* GetResult()
	{
		return _result;
	}

	virtual void SetCallback(void(__cdecl* callback)(Async<RPCMessage>*), void* userData)
	{
		_callback = callback;
		_userData = userData;
		_start = GetTickCount();
	}

	virtual void SetTimeoutCallback(void(__cdecl* callback)(Async<RPCMessage>*), unsigned int timeout)
	{
		_timeout = timeout;
		_start = GetTickCount();
		_timeoutCallback = callback;
	}

	virtual void* GetUserData()
	{
		return _userData;
	}

	virtual void Free()
	{
		delete this;
	}

	// additional definitions
	// set the result
private:

	uint64_t _asyncID;

public:
	uint64_t GetAsyncID()
	{
		return _asyncID;
	}

	void SetAsyncID(uint64_t value)
	{
		_asyncID = value;
	}

	void SetResult(RPCMessage* result)
	{
		_completed = true;
		_result = result;

		this->RunCallback();
	}

	// run the callback function (if completed)
	bool HandleTimeout()
	{
		if (!HasCompleted())
		{
			DWORD elapsed = GetTickCount() - _start;

			if (elapsed > _timeout)
			{
				if (_timeoutCallback)
				{
					_timeoutCallback(this);
				}

				return true;
			}
		}

		return false;
	}

	void RunCallback()
	{
		if (HasResult())
		{
			if (_callback)
			{
				_callback(this);
			}
		}
	}
};