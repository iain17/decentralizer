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
		const char* failureText;

		void callFinallyCBs() {
			for (auto const& cb : this->_finallyCBs) {
				cb();
			}
			//free();
		}
	public:
		Promise(Job job) {
			failureText = nullptr;
			this->_promise = std::async(std::launch::async, [&](Job job, Promise* self) {
				Sleep(1);
				T result = job(self);
				if (self != nullptr) {
					self->resolve(result);
				}
				return result;
			}, job, this);
		}

		void resolve(T result) {
			if (failureText) {
				return;
			}
			for (auto const& cb : this->_resultCBs) {
				cb(result);
			}
			this->callFinallyCBs();
		}

		void reject(const char* result) {
			this->failureText = result;
			for (auto const& cb : this->_failCBs) {
				cb(result);
			}
			//this->_promise._Set_value(NULL);
			this->callFinallyCBs();
		}

		Promise* then(resultCallback cb) {
			this->_resultCBs.push_back(cb);
			return this;
		}

		Promise* fail(failCallback cb) {
			this->_failCBs.push_back(cb);
			return this;
		}

		Promise* finally(finallyCallback cb) {
			this->_finallyCBs.push_back(cb);
			return this;
		}

		bool wait(int timeout = 7500) {
			std::future_status status;
			do {
				if (failureText) {return false;}
				status = this->_promise.wait_for(std::chrono::milliseconds(timeout));
				if (failureText) { return false; }
				if (status == std::future_status::deferred) {
					Sleep(100);
				} else if (status == std::future_status::timeout) {
					this->reject("Promise failed. Timed out on wait call.");
					return false;
				} else if (status == std::future_status::ready) {
					return failureText == 0;
				}
			} while (status != std::future_status::ready);
			return failureText == 0;
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