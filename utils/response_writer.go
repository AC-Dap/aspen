package utils

import "net/http"

type TrackingResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type trackingResponseWriter struct {
	http.ResponseWriter
	statusCode int

	// Go says that any number of 1xx status codes can be written, before
	// a final 2xx-5xx status code. We only want to track the final code.
	// Additionally, we only take the **first** final code written, even if
	// WriteHeader is called multiple times.
	wroteStatusCode bool
}

func NewTrackingResponseWriter(w http.ResponseWriter) TrackingResponseWriter {
	return trackingResponseWriter{
		ResponseWriter:  w,
		statusCode:      http.StatusOK,
		wroteStatusCode: false,
	}
}

func (t trackingResponseWriter) Write(b []byte) (int, error) {
	// Implicit 200 only if a final status hasn't been sent yet
	if !t.wroteStatusCode {
		t.WriteHeader(http.StatusOK)
	}
	return t.ResponseWriter.Write(b)
}

func (t trackingResponseWriter) WriteHeader(statusCode int) {
	if !t.wroteStatusCode && statusCode >= 200 {
		t.statusCode = statusCode
		t.wroteStatusCode = true
	}

	t.ResponseWriter.WriteHeader(statusCode)
}

func (t trackingResponseWriter) Status() int {
	return t.statusCode
}
