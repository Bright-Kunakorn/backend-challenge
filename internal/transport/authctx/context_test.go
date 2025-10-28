package authctx

import (
	"context"
	"testing"
)

func TestWithUserIDAndFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = WithUserID(ctx, "user-123")

	id, ok := UserIDFromContext(ctx)
	if !ok {
		t.Fatalf("expected user id in context")
	}
	if id != "user-123" {
		t.Fatalf("expected user-123 got %s", id)
	}

	if _, ok := UserIDFromContext(context.Background()); ok {
		t.Fatalf("expected no user id in fresh context")
	}
}
