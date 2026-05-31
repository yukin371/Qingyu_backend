package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestValidateBootstrapOptions(t *testing.T) {
	t.Parallel()

	valid := bootstrapOptions{
		Username:    "admin",
		Email:       "admin@example.com",
		PasswordEnv: defaultBootstrapPasswordEnv,
	}
	if err := validateBootstrapOptions(valid); err != nil {
		t.Fatalf("expected valid options, got error: %v", err)
	}

	invalidCases := []bootstrapOptions{
		{Email: "admin@example.com", PasswordEnv: defaultBootstrapPasswordEnv},
		{Username: "admin", PasswordEnv: defaultBootstrapPasswordEnv},
		{Username: "admin", Email: "invalid", PasswordEnv: defaultBootstrapPasswordEnv},
		{Username: "admin", Email: "admin@example.com"},
	}
	for _, item := range invalidCases {
		if err := validateBootstrapOptions(item); err == nil {
			t.Fatalf("expected validation error for %#v", item)
		}
	}
}

func TestResolveBootstrapPassword(t *testing.T) {
	t.Parallel()

	getenv := func(key string) string {
		switch key {
		case "GOOD":
			return "12345678"
		case "SHORT":
			return "123"
		default:
			return ""
		}
	}

	if _, err := resolveBootstrapPassword("MISSING", getenv); err == nil {
		t.Fatal("expected missing env error")
	}
	if _, err := resolveBootstrapPassword("SHORT", getenv); err == nil {
		t.Fatal("expected short password error")
	}
	if value, err := resolveBootstrapPassword("GOOD", getenv); err != nil || value != "12345678" {
		t.Fatalf("expected valid password, got value=%q err=%v", value, err)
	}
}

func TestMergeAdminRoles(t *testing.T) {
	t.Parallel()

	roles := mergeAdminRoles([]string{"reader", "author", "reader"})
	if len(roles) != 3 {
		t.Fatalf("expected 3 roles, got %d (%v)", len(roles), roles)
	}

	foundAdmin := false
	for _, role := range roles {
		if role == "admin" {
			foundAdmin = true
		}
	}
	if !foundAdmin {
		t.Fatalf("expected admin role in %v", roles)
	}
}

func TestDetectBootstrapConfigPath(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	configPath := filepath.Join(tempDir, "configs", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(configPath, []byte("server:\n  port: \"8080\"\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	resolved, err := detectBootstrapConfigPath("")
	if err != nil {
		t.Fatalf("expected config path, got error: %v", err)
	}
	if resolved != filepath.Clean(configPath) {
		t.Fatalf("expected %s, got %s", filepath.Clean(configPath), resolved)
	}
}

func TestShouldCreateBootstrapAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		findErr     error
		adminCount  int64
		wantCreate  bool
		wantErr     bool
		wantErrText string
	}{
		{
			name:       "empty database creates admin",
			findErr:    mongo.ErrNoDocuments,
			adminCount: 0,
			wantCreate: true,
		},
		{
			name:        "existing admin blocks duplicate bootstrap",
			findErr:     mongo.ErrNoDocuments,
			adminCount:  1,
			wantErr:     true,
			wantErrText: "系统中已存在 1 个管理员",
		},
		{
			name:        "lookup error bubbles up",
			findErr:     errors.New("lookup failed"),
			wantErr:     true,
			wantErrText: "查询管理员账号失败",
		},
		{
			name:       "matched user falls through to update branch",
			adminCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCreate, err := shouldCreateBootstrapAdmin(tt.findErr, tt.adminCount)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrText != "" && !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrText, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotCreate != tt.wantCreate {
				t.Fatalf("expected wantCreate=%v, got %v", tt.wantCreate, gotCreate)
			}
		})
	}
}
