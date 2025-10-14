package rest

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ctxUserID contextKey = "userID"
	ctxRole   contextKey = "role"
)

func (s *Server) jwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		userID, role, err := s.jwt.ParseToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxUserID, userID)
		ctx = context.WithValue(ctx, ctxRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := roleFromContext(r.Context())
			if _, ok := allowed[role]; !ok {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func roleFromContext(ctx context.Context) string {
	if v := ctx.Value(ctxRole); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func userIDFromContext(ctx context.Context) int64 {
	if v := ctx.Value(ctxUserID); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func (s *Server) requireManager(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := roleFromContext(r.Context())
		if role != "manager" {
			writeError(w, http.StatusForbidden, "manager role required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
