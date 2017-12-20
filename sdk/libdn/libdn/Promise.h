#pragma once
#include <future>

namespace libdn {

	template <class T>
	class Promise {
		typedef std::function<T(Promise*)> Job;
		typedef std::function<void(T)> resultCallback;
		typedef std::function<void(const char*)> failCallback;
		typedef std::function<void()> finallyCallback;
	private:
		std::future<T> _promise;
		std::vector<resultCallback> _resultCBs;
		std::vector<failCallback> _failCBs;
		std::vector<finallyCallback> _finallyCBs;
		int state;

		void callFinallyCBs() {
			for (auto const& cb : this->_finallyCBs) {
				cb();
			}
			//free();
		}

		bool isbad(void* p) {
			MEMORY_BASIC_INFORMATION mbi = { 0 };
			if (::VirtualQuery(p, &mbi, sizeof(mbi))) {
				DWORD mask = (PAGE_READONLY | PAGE_READWRITE | PAGE_WRITECOPY | PAGE_EXECUTE_READ | PAGE_EXECUTE_READWRITE | PAGE_EXECUTE_WRITECOPY);
				bool b = !(mbi.Protect & mask);
				// check the page is not a guard page
				if (mbi.Protect & (PAGE_GUARD | PAGE_NOACCESS)) b = true;

				return b;
			}
			return true;
		}

	public:
		Promise(Job job) {
			this->state = 0;
			this->_promise = std::async(std::launch::async, [&](Job job, Promise* self) {
				Sleep(1);
				T result = job(self);
				if (&result != nullptr && self != nullptr) {
					self->resolve(result);
				}
				return result;
			}, job, this);
			this->fail(Log_Print);
		}

		void resolve(T result) {
			if (this->state != 0) {
				return;
			}
			this->state = 1;
			for (auto const& cb : this->_resultCBs) {
				cb(result);
			}
			this->callFinallyCBs();
		}

		void reject(std::string result) {
			if (this->state != 0) {
				return;
			}
			this->state = 2;
			for (auto const& cb : this->_failCBs) {
				cb(result.c_str());
			}
			this->callFinallyCBs();
		}

		void then(resultCallback cb) {
			if (this->state != 0) {
				return;
			}
			this->_resultCBs.push_back(cb);
		}

		void fail(failCallback cb) {
			if (this->state != 0) {
				return;
			}
			this->_failCBs.push_back(cb);
		}

		void finally(finallyCallback cb) {
			if (this->state != 0) {
				return;
			}
			this->_finallyCBs.push_back(cb);
		}

		bool wait(int timeout = 30000) {
			if (this->state != 0) {
				return this->state == 1;
			}
			std::future_status status;
			do {
				if (this->state != 0) {break;}
				status = this->_promise.wait_for(std::chrono::milliseconds(timeout));
				if (this->state != 0) { break; }
				if (status == std::future_status::deferred) {
					Sleep(100);
				} else if (status == std::future_status::timeout) {
					this->reject("Promise failed. Timed out on wait call.");
					return false;
				} else if (status == std::future_status::ready) {
					return this->state == 1;
				}
			} while (status != std::future_status::ready);
			return this->state == 1;
		}

		//Blocking!
		T get() {
			return this->_promise.get();
		}

		/*
		void free() {
			if (failureText) {
				delete failureText;
			}
			delete this;
		}
		*/
	};

}