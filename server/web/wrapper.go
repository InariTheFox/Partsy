package web

import "net/http"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode        int
	statusCodeWritten bool
}

func NewWrappedWriter(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter:    w,
		statusCodeWritten: false,
	}
}

func (w *ResponseWriterWrapper) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.statusCodeWritten = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriterWrapper) Write(data []byte) (int, error) {
	if !w.statusCodeWritten {
		w.statusCode = http.StatusOK
	}

	return w.ResponseWriter.Write(data)
}
