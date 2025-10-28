package authctx

import "context"

type contextKey string

const userIDKey contextKey = "userID"

// WithUserID injects the authenticated user ID into the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts the user ID from context if present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDKey).(string)
	return val, ok
}
