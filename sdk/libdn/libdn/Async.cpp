#include "StdInc.h"

namespace libdn {
	std::queue<AsyncCallback*> _asyncCallbacks;
	void Async_RunCallbacks() {
		RPC_RunFrame();

		while (!_asyncCallbacks.empty()) {
			AsyncCallback* callback = _asyncCallbacks.front();

			if (callback->RunCallback()) {
				callback->Free();
			}

			_asyncCallbacks.pop();
		}
	}
}