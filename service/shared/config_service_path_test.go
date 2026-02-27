package shared

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigService_UsesProvidedPath(t *testing.T) {
	svc := NewConfigService("./custom/config.yaml")
	assert.Equal(t, "./custom/config.yaml", svc.configPath)
}

func TestNewConfigService_DefaultPrefersConfigs(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	writeConfigFile(t, filepath.Join("configs", "config.yaml"))
	writeConfigFile(t, filepath.Join("config", "config.yaml"))

	svc := NewConfigService("")
	assert.Equal(t, "./configs/config.yaml", svc.configPath)
}

func TestNewConfigService_FallsBackToLegacyPath(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	writeConfigFile(t, filepath.Join("config", "config.yaml"))

	svc := NewConfigService("")
	assert.Equal(t, "./config/config.yaml", svc.configPath)
}

func TestNewConfigService_DefaultPathWhenNoFileExists(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	svc := NewConfigService("")
	assert.Equal(t, "./configs/config.yaml", svc.configPath)
}

func writeConfigFile(t *testing.T, rel string) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Dir(rel), 0755))
	require.NoError(t, os.WriteFile(rel, []byte("server:\n  port: \":9090\"\n"), 0644))
}
