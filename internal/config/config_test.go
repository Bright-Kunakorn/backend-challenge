package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("PORT", "9000")
	t.Setenv("GRPC_PORT", "9100")
	t.Setenv("MONGO_URI", "mongodb://example:27017")
	t.Setenv("MONGO_DB", "db")
	t.Setenv("JWT_ISSUER", "issuer")
	t.Setenv("JWT_EXPIRY", "2h")
	t.Setenv("USER_COUNT_TICK", "30s")
	t.Setenv("ENVIRONMENT", "test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if cfg.Port != "9000" || cfg.GRPCPort != "9100" || cfg.MongoURI != "mongodb://example:27017" {
		t.Fatalf("unexpected config %+v", cfg)
	}
	if cfg.JWTExpiry != 2*time.Hour {
		t.Fatalf("expected expiry 2h got %v", cfg.JWTExpiry)
	}
	if cfg.BackgroundTick != 30*time.Second {
		t.Fatalf("expected tick 30s got %v", cfg.BackgroundTick)
	}
}

func TestLoadMissingSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	if _, err := Load(); err == nil {
		t.Fatal("expected error when JWT_SECRET missing")
	}
}

func TestParseDurationFallback(t *testing.T) {
	if d := parseDuration("bad", time.Minute); d != time.Minute {
		t.Fatalf("expected fallback duration got %v", d)
	}
}

func TestMustParseInt(t *testing.T) {
	if MustParseInt("MISSING", 42) != 42 {
		t.Fatal("expected fallback value")
	}
	t.Setenv("TEST_INT", "not-an-int")
	if MustParseInt("TEST_INT", 5) != 5 {
		t.Fatal("expected fallback when parse fails")
	}
	t.Setenv("TEST_INT", "17")
	if MustParseInt("TEST_INT", 5) != 17 {
		t.Fatal("expected parsed value")
	}
}
