#pragma once
#include <future>

namespace libdn {

	template <class T>
	class Promise {
		typedef std::function<T(Promise*)> Job;
		typedef std::function<void(T)> resultCallback;
		typedef std::function<void(std::string)> failCallback;
		typedef std::function<void()> finallyCallback;
	private:
		std::future<T> _promise;
		std::vector<resultCallback> _resultCBs;
		std::vector<failCallback> _failCBs;
		std::vector<finallyCallback> _finallyCBs;

		int referenced = 1;
		void callFinallyCBs() {
			for (auto const& cb : this->_finallyCBs) {
				cb();
			}
		}
		void free() {
			referenced--;
			if (referenced > 0) {
				return;
			}
			delete this;
		}
	public:
		std::string failureText;
		Promise(Job job) {
			this->_promise = std::async(std::launch::async, [&](Job job, Promise* self) {
				Sleep(1);
				T result = job(self);
				if (self != nullptr && !self->failureText.empty()) {
					self->resolve(result);
				}
				return result;
			}, job, this);
		}

		void resolve(T result) {
			for (auto const& cb : this->_resultCBs) {
				cb(result);
			}
			this->callFinallyCBs();
			free();
		}

		void reject(std::string result) {
			this->failureText = result;
			for (auto const& cb : this->_failCBs) {
				cb(result);
			}
			this->callFinallyCBs();
			free();
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

		//Blocking!
		T get(int timeout = 7500) {
			referenced++;
			std::future_status status;
			do {
				status = this->_promise.wait_for(std::chrono::milliseconds(timeout));
				if (status == std::future_status::deferred) {
					Sleep(100);
				} else if (status == std::future_status::timeout) {
					this->reject("timed out");
					return NULL;
				} else if (status == std::future_status::ready) {
					return this->_promise.get();
				}
			} while (status != std::future_status::ready);
			return this->_promise.get();
		}

	};

}