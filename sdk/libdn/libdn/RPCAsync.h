#pragma once
#include "StdInc.h"
#include "DNAsync.h"

class NPRPCAsync : public NPAsync<RPCMessage>
{
private:
	void(__cdecl* _callback)(NPAsync<RPCMessage>*);
	RPCMessage* _result;
	void* _userData;

	void(__cdecl* _timeoutCallback)(NPAsync<RPCMessage>*);
	unsigned int _timeout;
	DWORD _start;

public:
	NPRPCAsync()
	{
		_callback = NULL;
		_result = NULL;
		_userData = NULL;
		_timeout = -1;
		_timeoutCallback = nullptr;
		_start = 0;
	}

	// implementations for base NPAsync
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
		return (_result != NULL);
	}

	virtual RPCMessage* GetResult()
	{
		return _result;
	}

	virtual void SetCallback(void(__cdecl* callback)(NPAsync<RPCMessage>*), void* userData)
	{
		_callback = callback;
		_userData = userData;
		_start = GetTickCount();
	}

	virtual void SetTimeoutCallback(void(__cdecl* callback)(NPAsync<RPCMessage>*), unsigned int timeout)
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

	int _asyncID;

public:
	int GetAsyncID()
	{
		return _asyncID;
	}

	void SetAsyncID(int value)
	{
		_asyncID = value;
	}

	void SetResult(RPCMessage* result)
	{
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
		if (HasCompleted())
		{
			if (_callback)
			{
				_callback(this);
			}
		}
	}
};