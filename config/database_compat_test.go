package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDatabaseConfig_OldFormat 测试旧格式（扁平）配置加载
func TestDatabaseConfig_OldFormat(t *testing.T) {
	// 模拟旧格式配置
	cfg := &DatabaseConfig{
		// 旧格式（扁平）字段
		URI:         "mongodb://localhost:27017",
		Name:        "qingyu",
		MaxPoolSize: 100,
		MinPoolSize: 10,
	}

	mongoConfig, err := cfg.GetMongoConfig()

	assert.NoError(t, err)
	assert.Equal(t, "mongodb://localhost:27017", mongoConfig.URI)
	assert.Equal(t, "qingyu", mongoConfig.Database)
	assert.Equal(t, uint64(100), mongoConfig.MaxPoolSize)
	assert.Equal(t, uint64(10), mongoConfig.MinPoolSize)
}

// TestDatabaseConfig_NewFormat 测试新格式（嵌套）配置加载
func TestDatabaseConfig_NewFormat(t *testing.T) {
	// 模拟新格式配置
	cfg := &DatabaseConfig{
		Primary: DatabaseConnection{
			Type: DatabaseTypeMongoDB,
			MongoDB: &MongoDBConfig{
				URI:         "mongodb://localhost:27017",
				Database:    "qingyu",
				MaxPoolSize: 100,
			},
		},
	}

	mongoConfig, err := cfg.GetMongoConfig()

	assert.NoError(t, err)
	assert.Equal(t, "mongodb://localhost:27017", mongoConfig.URI)
	assert.Equal(t, "qingyu", mongoConfig.Database)
}

// TestDatabaseConfig_NewFormatWithDefaults 测试新格式配置也会填充默认值
func TestDatabaseConfig_NewFormatWithDefaults(t *testing.T) {
	cfg := &DatabaseConfig{
		Primary: DatabaseConnection{
			Type: DatabaseTypeMongoDB,
			MongoDB: &MongoDBConfig{
				URI:      "mongodb://localhost:27017",
				Database: "qingyu",
			},
		},
	}

	mongoConfig, err := cfg.GetMongoConfig()

	assert.NoError(t, err)
	assert.Equal(t, uint64(100), mongoConfig.MaxPoolSize)
	assert.Equal(t, uint64(5), mongoConfig.MinPoolSize)
	assert.Equal(t, 10*time.Second, mongoConfig.ConnectTimeout)
	assert.Equal(t, 30*time.Second, mongoConfig.ServerTimeout)
	assert.Equal(t, int64(100), mongoConfig.ProfilerSizeMB)
}

// TestDatabaseConfig_Priority 测试新格式优先级高于旧格式
func TestDatabaseConfig_Priority(t *testing.T) {
	// 新格式优先
	cfg := &DatabaseConfig{
		// 旧格式
		URI:  "mongodb://old-format:27017",
		Name: "old_db",
		// 新格式
		Primary: DatabaseConnection{
			Type: DatabaseTypeMongoDB,
			MongoDB: &MongoDBConfig{
				URI:      "mongodb://new-format:27017",
				Database: "new_db",
			},
		},
	}

	mongoConfig, err := cfg.GetMongoConfig()

	assert.NoError(t, err)
	// 应该使用新格式
	assert.Equal(t, "mongodb://new-format:27017", mongoConfig.URI)
	assert.Equal(t, "new_db", mongoConfig.Database)
}

// TestDatabaseConfig_OldFormatWithDefaults 测试旧格式配置的默认值
func TestDatabaseConfig_OldFormatWithDefaults(t *testing.T) {
	// 只提供必需字段，其他使用默认值
	cfg := &DatabaseConfig{
		URI:  "mongodb://localhost:27017",
		Name: "qingyu",
		// MaxPoolSize, MinPoolSize, ConnectTimeout 未设置，应使用默认值
	}

	mongoConfig, err := cfg.GetMongoConfig()

	assert.NoError(t, err)
	assert.Equal(t, "mongodb://localhost:27017", mongoConfig.URI)
	assert.Equal(t, "qingyu", mongoConfig.Database)
	assert.Equal(t, uint64(100), mongoConfig.MaxPoolSize)   // 默认值
	assert.Equal(t, uint64(10), mongoConfig.MinPoolSize)    // 默认值（注意：旧文档说5，但应该与MongoDBConfig一致）
	assert.Equal(t, 10*1000000000, int(mongoConfig.ConnectTimeout)) // 10秒
}

// TestDatabaseConfig_NoConfig 测试无配置的错误情况
func TestDatabaseConfig_NoConfig(t *testing.T) {
	cfg := &DatabaseConfig{}

	_, err := cfg.GetMongoConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database configuration")
}

// TestDatabaseConfig_ValidateWithOldFormat 测试旧格式配置的验证
func TestDatabaseConfig_ValidateWithOldFormat(t *testing.T) {
	cfg := &DatabaseConfig{
		URI:  "mongodb://localhost:27017",
		Name: "qingyu",
	}

	err := cfg.Validate()

	assert.NoError(t, err)
}

// TestDatabaseConfig_ValidateWithNewFormat 测试新格式配置的验证
func TestDatabaseConfig_ValidateWithNewFormat(t *testing.T) {
	cfg := &DatabaseConfig{
		Primary: DatabaseConnection{
			Type: DatabaseTypeMongoDB,
			MongoDB: &MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "qingyu",
				ProfilingLevel: 1,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
		},
	}

	err := cfg.Validate()

	assert.NoError(t, err)
}

// TestDatabaseConfig_ValidateEmptyOldFormat 测试旧格式空配置验证失败
func TestDatabaseConfig_ValidateEmptyOldFormat(t *testing.T) {
	cfg := &DatabaseConfig{
		URI:  "", // 空
		Name: "",
		Type: "mongodb",
	}

	// 尝试获取MongoDB配置应该失败
	_, err := cfg.GetMongoConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database configuration")
}

// TestDatabaseConfig_MultipleCalls 测试多次调用GetMongoConfig返回相同结果
func TestDatabaseConfig_MultipleCalls(t *testing.T) {
	cfg := &DatabaseConfig{
		URI:  "mongodb://localhost:27017",
		Name: "qingyu",
	}

	config1, err1 := cfg.GetMongoConfig()
	config2, err2 := cfg.GetMongoConfig()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// 应该返回相同的配置对象（已缓存）
	assert.Same(t, config1, config2)
}
