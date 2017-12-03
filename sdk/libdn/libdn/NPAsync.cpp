#include "StdInc.h"

std::queue<DNAsyncCallback*> _asyncCallbacks;

void Async_RunCallbacks()
{
	RPC_RunFrame();

	while (!_asyncCallbacks.empty())
	{
		DNAsyncCallback* callback = _asyncCallbacks.front();

		if (callback->RunCallback())
		{
			callback->Free();
		}

		_asyncCallbacks.pop();
	}
}