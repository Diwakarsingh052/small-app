package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type key string

const TraceIdKey key = "1"

func Logger(next http.Handler) http.Handler {
	// http.Handler is an interface, so we can't return a function directly to it
	// we have wrapped the call in http.HandlerFunc which is a type that implements the interface
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := uuid.NewString()
		requestStartTime := time.Now()

		// taking the context out from the request object
		ctx := r.Context()
		// putting the traceId in the context
		ctx = context.WithValue(ctx, TraceIdKey, traceId)

		slog.Info("started", slog.String("Trace ID", traceId),
			slog.String("Method", r.Method), slog.String("URL Path", r.URL.Path),
		)
		// r.WithContext would update the request to use the updated context
		rw := &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx)) // call the next thing in the chain

		slog.Info("completed", slog.String("Trace ID", traceId),
			slog.String("Method", r.Method), slog.String("URL Path", r.URL.Path),
			slog.Int("Status Code", rw.status),
			slog.Int64("duration μs",
				time.Since(requestStartTime).Microseconds()),
		)
	})
}

// Creating a struct to provide custom functionality for the WriteHeader by providing a
// custom implementation
type responseWriterWrapper struct {
	http.ResponseWriter
	status int // we have added status code field to capture the status separately
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	//fmt.Println("custom write header functionality called")
	w.status = statusCode // setting the status code in the struct variable
	// writing the header to the request
	w.ResponseWriter.WriteHeader(statusCode)
}
