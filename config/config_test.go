package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  port: "8080"
  mode: "debug"

jwt:
  secret: "test-secret"
  expiration_hours: 24

ai:
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
  max_tokens: 1000
  temperature: 1

database:
  type: "mongodb"
  primary:
    mongodb:
      uri: "mongodb://localhost:27017"
      database: "test_db"
      connect_timeout: "10s"
      max_pool_size: 100
      min_pool_size: 5
      server_timeout: "30s"
  replicas: []
  indexing:
    auto_create: true
    background: true
  validation:
    enabled: true
    strict_mode: false
  sync:
    enabled: false
    interval: "5m"
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	// 测试加载配置 - 直接传递临时文件路径
	config, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)
	assert.NotNil(t, config)

	// 验证服务器配置
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "debug", config.Server.Mode)

	// 验证JWT配置
	assert.Equal(t, "test-secret", config.JWT.Secret)
	assert.Equal(t, 24, config.JWT.ExpirationHours)

	// 验证AI配置
	assert.Equal(t, "test-key", config.AI.APIKey)
	assert.Equal(t, "https://api.openai.com/v1", config.AI.BaseURL)
	assert.Equal(t, 1000, config.AI.MaxTokens)
	assert.Equal(t, 1, config.AI.Temperature)

	// 验证数据库配置
	assert.Equal(t, "mongodb", config.Database.Type)
	assert.Equal(t, "mongodb://localhost:27017", config.Database.Primary.MongoDB.URI)
	assert.Equal(t, "test_db", config.Database.Primary.MongoDB.Database)
	assert.Equal(t, uint64(100), config.Database.Primary.MongoDB.MaxPoolSize)
	assert.Equal(t, uint64(5), config.Database.Primary.MongoDB.MinPoolSize)
	assert.Equal(t, 10*time.Second, config.Database.Primary.MongoDB.ConnectTimeout)

	// 验证索引配置
	assert.True(t, config.Database.Indexing.AutoCreate)
	assert.True(t, config.Database.Indexing.Background)

	// 验证验证配置
	assert.True(t, config.Database.Validation.Enabled)
	assert.False(t, config.Database.Validation.StrictMode)

	// 验证数据库同步配置
	assert.False(t, config.Database.Sync.Enabled)
	assert.Equal(t, 5*time.Minute, config.Database.Sync.Interval)
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv("QINGYU_SERVER_PORT", "9090")
	os.Setenv("QINGYU_JWT_SECRET", "env-secret")
	os.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI", "mongodb://env-host:27017")
	os.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE", "env_db")
	defer func() {
		os.Unsetenv("QINGYU_SERVER_PORT")
		os.Unsetenv("QINGYU_JWT_SECRET")
		os.Unsetenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI")
		os.Unsetenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE")
	}()

	// 创建基本配置文件
	configContent := `
server:
  port: "8080"
  mode: "debug"

jwt:
  secret: "test-secret"
  expiration_hours: 24

ai:
  api_key: "test-ai-key"
  base_url: "https://api.openai.com/v1"
  max_tokens: 1000
  temperature: 1

database:
  type: "mongodb"
  primary:
    mongodb:
      uri: "mongodb://localhost:27017"
      database: "test_db"
`

	tmpFile, err := os.CreateTemp("", "config_env_test_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	viper.SetConfigFile(tmpFile.Name())

	// 测试加载配置
	config, err := LoadConfig(".")
	require.NoError(t, err)

	// 验证环境变量覆盖
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "env-secret", config.JWT.Secret)
	assert.Equal(t, "mongodb://env-host:27017", config.Database.Primary.MongoDB.URI)
	assert.Equal(t, "env_db", config.Database.Primary.MongoDB.Database)
}

func TestViperConfigManager(t *testing.T) {
	manager := NewViperConfigManager()
	assert.NotNil(t, manager)

	// 测试设置默认值
	manager.setViperDefaults()

	// 验证一些默认值
	assert.Equal(t, "8080", manager.viper.GetString("server.port"))
	assert.Equal(t, "debug", manager.viper.GetString("server.mode"))
	assert.Equal(t, "mongodb://localhost:27017", manager.viper.GetString("database.primary.mongodb.uri"))
	assert.Equal(t, "qingyu", manager.viper.GetString("database.primary.mongodb.database"))
}

func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *DatabaseConfig
		wantErr bool
	}{
		{
			name: "valid mongodb config",
			config: &DatabaseConfig{
				Type: "mongodb",
				Primary: DatabaseConnection{
					Type: DatabaseTypeMongoDB,
					MongoDB: &MongoDBConfig{
						URI:            "mongodb://localhost:27017",
						Database:       "test_db",
						ConnectTimeout: 10 * time.Second,
						MaxPoolSize:    100,
						MinPoolSize:    10,
						ServerTimeout:  30 * time.Second,
					},
				},
				Replicas: []DatabaseConnection{},
				Indexing: IndexingConfig{
					AutoCreate: true,
					Background: true,
				},
				Validation: ValidationConfig{
					Enabled:    true,
					StrictMode: false,
				},
				Sync: SynchronizationConfig{
					Enabled:  false,
					Interval: time.Hour,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid mongodb config - empty URI",
			config: &DatabaseConfig{
				Type: "mongodb",
				Primary: DatabaseConnection{
					Type: DatabaseTypeMongoDB,
					MongoDB: &MongoDBConfig{
						URI:            "",
						Database:       "test_db",
						ConnectTimeout: 10 * time.Second,
						MaxPoolSize:    100,
						MinPoolSize:    10,
						ServerTimeout:  30 * time.Second,
					},
				},
				Replicas: []DatabaseConnection{},
			},
			wantErr: true,
		},
		{
			name: "invalid mongodb config - empty database name",
			config: &DatabaseConfig{
				Primary: DatabaseConnection{
					Type: DatabaseTypeMongoDB,
					MongoDB: &MongoDBConfig{
						URI:            "mongodb://localhost:27017",
						Database:       "",
						ConnectTimeout: 10 * time.Second,
						MaxPoolSize:    100,
						MinPoolSize:    10,
						ServerTimeout:  30 * time.Second,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
