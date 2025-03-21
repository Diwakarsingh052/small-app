package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"rest-api/auth"
	"strings"
)

func (m *Mid) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Trace ID from the context
		ctx := r.Context()
		traceId, ok := ctx.Value(TraceIdKey).(string)
		if !ok {
			slog.Error("trace id not present in the context")
			traceId = "unknown"
		}

		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			slog.Error("invalid authorization header format",
				slog.String("Trace ID", traceId),
				slog.String("Error", err.Error()))

			// Respond with HTTP 401 Unauthorized if the header is invalid
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := m.a.ValidateToken(token)

		if err != nil {
			slog.Error("invalid token",
				slog.String("Trace ID", traceId),
				slog.String("Error", err.Error()))

			// Respond with HTTP 401 Unauthorized if the token is invalid
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Attach the claims to the context
		ctx = context.WithValue(ctx, auth.Key, claims)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
