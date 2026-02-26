package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validSearchIndicesYAML = `
indices:
  books:
    alias: books_search
    number_of_shards: 1
    number_of_replicas: 0
    mapping:
      properties:
        title:
          type: text
settings:
  analysis:
    analyzer: {}
`

func TestLoadSearchIndicesConfig_UsesExplicitPath(t *testing.T) {
	file := writeSearchIndicesFile(t, t.TempDir(), "custom_indices.yaml", validSearchIndicesYAML)

	cfg, err := LoadSearchIndicesConfig(file)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "books_search", cfg.Indices["books"].Alias)
}

func TestLoadSearchIndicesConfig_UsesEnvPath(t *testing.T) {
	temp := t.TempDir()
	file := writeSearchIndicesFile(t, temp, "env_indices.yaml", validSearchIndicesYAML)

	t.Setenv("SEARCH_INDICES_CONFIG", file)

	cfg, err := LoadSearchIndicesConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "books_search", cfg.Indices["books"].Alias)
}

func TestLoadSearchIndicesConfig_DefaultPrefersConfigs(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	writeSearchIndicesFile(t, temp, filepath.Join("configs", "search_indices.yaml"), validSearchIndicesYAML)
	writeSearchIndicesFile(t, temp, filepath.Join("config", "search_indices.yaml"), `
indices:
  books:
    alias: legacy_alias
    number_of_shards: 1
    number_of_replicas: 0
    mapping:
      properties:
        title:
          type: text
settings:
  analysis:
    analyzer: {}
`)

	cfg, err := LoadSearchIndicesConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "books_search", cfg.Indices["books"].Alias)
}

func TestLoadSearchIndicesConfig_FallsBackToLegacyConfigDir(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	writeSearchIndicesFile(t, temp, filepath.Join("config", "search_indices.yaml"), validSearchIndicesYAML)

	cfg, err := LoadSearchIndicesConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "books_search", cfg.Indices["books"].Alias)
}

func TestLoadSearchIndicesConfig_ReturnsErrorWhenNoFileFound(t *testing.T) {
	temp := t.TempDir()
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	cfg, err := LoadSearchIndicesConfig("")
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "tried:")
	assert.Contains(t, err.Error(), "configs/search_indices.yaml")
	assert.Contains(t, err.Error(), "config/search_indices.yaml")
}

func writeSearchIndicesFile(t *testing.T, base, rel, content string) string {
	t.Helper()

	path := filepath.Join(base, rel)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}
