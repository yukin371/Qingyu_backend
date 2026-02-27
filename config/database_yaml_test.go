package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConfig_LoadOldFormatYAML 测试从YAML加载旧格式配置
func TestDatabaseConfig_LoadOldFormatYAML(t *testing.T) {
	t.Setenv("MONGODB_URI", "")
	t.Setenv("MONGODB_DATABASE", "")
	t.Setenv("MONGODB_MAX_POOL_SIZE", "")
	t.Setenv("MONGODB_PROFILING_LEVEL", "")
	t.Setenv("MONGODB_SLOW_MS", "")
	t.Setenv("MONGODB_PROFILER_SIZE_MB", "")

	yamlContent := `
uri: "mongodb://localhost:27017"
name: "qingyu"
max_pool_size: 100
min_pool_size: 10
connect_timeout: 10s
type: mongodb
primary:
  type: ""
`
	// 创建临时文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config_old.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// 加载配置
	config, err := LoadDatabaseConfig(configPath)
	require.NoError(t, err)

	// 验证可以获取MongoDB配置
	mongoConfig, err := config.GetMongoConfig()
	require.NoError(t, err)
	assert.Equal(t, "mongodb://localhost:27017", mongoConfig.URI)
	assert.Equal(t, "qingyu", mongoConfig.Database)
	assert.Equal(t, uint64(100), mongoConfig.MaxPoolSize)
	assert.Equal(t, uint64(10), mongoConfig.MinPoolSize)
}

// TestDatabaseConfig_LoadNewFormatYAML 测试从YAML加载新格式配置
func TestDatabaseConfig_LoadNewFormatYAML(t *testing.T) {
	t.Setenv("MONGODB_URI", "")
	t.Setenv("MONGODB_DATABASE", "")
	t.Setenv("MONGODB_MAX_POOL_SIZE", "")
	t.Setenv("MONGODB_PROFILING_LEVEL", "")
	t.Setenv("MONGODB_SLOW_MS", "")
	t.Setenv("MONGODB_PROFILER_SIZE_MB", "")

	yamlContent := `
type: mongodb
primary:
  type: mongodb
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "qingyu"
    max_pool_size: 100
    min_pool_size: 10
    connect_timeout: 10s
    server_timeout: 30s
    profiling_level: 1
    slow_ms: 100
    profiler_size_mb: 100
`
	// 创建临时文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config_new.yaml")
	err := os.WriteFile(configPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// 加载配置
	config, err := LoadDatabaseConfig(configPath)
	require.NoError(t, err)

	// 验证可以获取MongoDB配置
	mongoConfig, err := config.GetMongoConfig()
	require.NoError(t, err)
	assert.Equal(t, "mongodb://localhost:27017", mongoConfig.URI)
	assert.Equal(t, "qingyu", mongoConfig.Database)
	assert.Equal(t, uint64(100), mongoConfig.MaxPoolSize)
	assert.Equal(t, uint64(10), mongoConfig.MinPoolSize)
	assert.Equal(t, 1, mongoConfig.ProfilingLevel)
	assert.Equal(t, int64(100), mongoConfig.SlowMS)
	assert.Equal(t, int64(100), mongoConfig.ProfilerSizeMB)
}

// TestDatabaseConfig_ValidateOldFormatYAML 测试旧格式YAML的验证
func TestDatabaseConfig_ValidateOldFormatYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
		errMsg      string
	}{
		{
			name: "valid_old_format",
			yamlContent: `
uri: "mongodb://localhost:27017"
name: "qingyu"
type: mongodb
`,
			wantErr: false,
		},
		{
			name: "invalid_old_format_missing_uri",
			yamlContent: `
name: "qingyu"
type: mongodb
`,
			wantErr: true, // 有name但没有uri，应该失败
			errMsg:  "invalid database configuration",
		},
		{
			name: "no_config_at_all",
			yamlContent: `
type: mongodb
`,
			wantErr: false, // 完全没有配置，验证通过
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configPath, []byte(tt.yamlContent), 0644)
			require.NoError(t, err)

			config, err := LoadDatabaseConfig(configPath)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.NoError(t, config.Validate())
			}
		})
	}
}
