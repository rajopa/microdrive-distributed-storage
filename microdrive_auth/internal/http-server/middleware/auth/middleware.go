package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"microdrive_auth/internal/lib/logger/sl"
	"microdrive_auth/internal/lib/logger/sl/jwt"
)

type PermissionProvider interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type ctxKey string

const (
	uidKey     ctxKey = "uid"
	isAdminKey ctxKey = "isAdmin"
	errorKey   ctxKey = "error"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

func New(
	log *slog.Logger,
	appSecret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt.Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UID)
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), uidKey, claims.UID)
			ctx = context.WithValue(r.Context(), isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}
