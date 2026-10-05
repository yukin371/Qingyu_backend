package serverapp

import (
	"path/filepath"
	"testing"
)

func TestResolveConfigPath_PrefersConfigFile(t *testing.T) {
	t.Setenv("CONFIG_FILE", filepath.Join("configs", "file.yaml"))
	t.Setenv("CONFIG_PATH", filepath.Join("configs", "directory"))

	got := resolveConfigPath()
	want := filepath.Clean(filepath.Join("configs", "file.yaml"))
	if got != want {
		t.Fatalf("resolveConfigPath() = %q, want %q", got, want)
	}
}

func TestResolveConfigPathUsesConfigPathWhenFileIsUnset(t *testing.T) {
	t.Setenv("CONFIG_FILE", "")
	t.Setenv("CONFIG_PATH", filepath.Join("configs", "directory"))

	got := resolveConfigPath()
	want := filepath.Clean(filepath.Join("configs", "directory"))
	if got != want {
		t.Fatalf("resolveConfigPath() = %q, want %q", got, want)
	}
}
