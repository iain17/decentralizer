package web

import "net/http"

type recoverableResponseWriter struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	fallback       http.HandlerFunc
	recovered      bool
}

func (rrw *recoverableResponseWriter) WriteHeader(c int) {
	switch c {
	case 404:
		rrw.fallback.ServeHTTP(rrw.responseWriter, rrw.request)
		rrw.recovered = true
	default:
		rrw.responseWriter.WriteHeader(c)
	}
}

func (rrw *recoverableResponseWriter) Write(b []byte) (int, error) {
	if rrw.recovered {
		return 0, nil
	}
	return rrw.responseWriter.Write(b)
}

func (rrw *recoverableResponseWriter) Header() http.Header {
	return rrw.responseWriter.Header()
}
