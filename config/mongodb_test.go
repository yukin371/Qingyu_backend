package config

import (
	"os"
	"testing"
	"time"
)

// TestMongoDBConfigProfilingDefaults 测试MongoDB配置的profiling默认值
func TestMongoDBConfigProfilingDefaults(t *testing.T) {
	config := getDefaultDatabaseConfig()

	if config.Primary.MongoDB == nil {
		t.Fatal("MongoDB配置不应为空")
	}

	mongoConfig := config.Primary.MongoDB

	// 测试ProfilingLevel默认值应该是1（仅记录慢查询）
	if mongoConfig.ProfilingLevel != 1 {
		t.Errorf("期望ProfilingLevel=1, 实际=%d", mongoConfig.ProfilingLevel)
	}

	// 测试SlowMS默认值应该是100ms
	if mongoConfig.SlowMS != 100 {
		t.Errorf("期望SlowMS=100, 实际=%d", mongoConfig.SlowMS)
	}

	// 测试ProfilerSizeMB默认值应该是100MB
	if mongoConfig.ProfilerSizeMB != 100 {
		t.Errorf("期望ProfilerSizeMB=100, 实际=%d", mongoConfig.ProfilerSizeMB)
	}
}

// TestMongoDBConfigProfilingEnvOverrides 测试从环境变量覆盖profiling配置
func TestMongoDBConfigProfilingEnvOverrides(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("MONGODB_PROFILING_LEVEL", "2")
	os.Setenv("MONGODB_SLOW_MS", "200")
	os.Setenv("MONGODB_PROFILER_SIZE_MB", "200")
	defer func() {
		os.Unsetenv("MONGODB_PROFILING_LEVEL")
		os.Unsetenv("MONGODB_SLOW_MS")
		os.Unsetenv("MONGODB_PROFILER_SIZE_MB")
	}()

	config := getDefaultDatabaseConfig()
	mongoConfig := config.Primary.MongoDB

	// 验证环境变量覆盖
	if mongoConfig.ProfilingLevel != 2 {
		t.Errorf("期望从环境变量覆盖ProfilingLevel=2, 实际=%d", mongoConfig.ProfilingLevel)
	}

	if mongoConfig.SlowMS != 200 {
		t.Errorf("期望从环境变量覆盖SlowMS=200, 实际=%d", mongoConfig.SlowMS)
	}

	if mongoConfig.ProfilerSizeMB != 200 {
		t.Errorf("期望从环境变量覆盖ProfilerSizeMB=200, 实际=%d", mongoConfig.ProfilerSizeMB)
	}
}

// TestMongoDBConfigProfilingValidation 测试profiling配置验证
func TestMongoDBConfigProfilingValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      MongoDBConfig
		expectError bool
	}{
		{
			name: "有效的profiling配置",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  1,
				SlowMS:          100,
				ProfilerSizeMB:  100,
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: false,
		},
		{
			name: "ProfilingLevel=0是有效的",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  0,
				SlowMS:          100,
				ProfilerSizeMB:  100,
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: false,
		},
		{
			name: "ProfilingLevel=2是有效的",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  2,
				SlowMS:          100,
				ProfilerSizeMB:  100,
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: false,
		},
		{
			name: "无效的ProfilingLevel应该被拒绝",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  3, // 无效值
				SlowMS:          100,
				ProfilerSizeMB:  100,
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: true,
		},
		{
			name: "负数的SlowMS应该被拒绝",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  1,
				SlowMS:          -100, // 无效值
				ProfilerSizeMB:  100,
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: true,
		},
		{
			name: "过小的ProfilerSizeMB应该被拒绝",
			config: MongoDBConfig{
				URI:             "mongodb://localhost:27017",
				Database:        "test",
				ProfilingLevel:  1,
				SlowMS:          100,
				ProfilerSizeMB:  0, // 无效值
				ConnectTimeout:  10 * time.Second,
				ServerTimeout:   30 * time.Second,
				MaxPoolSize:     100,
				MinPoolSize:     5,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("期望验证失败但验证成功")
			}
			if !tt.expectError && err != nil {
				t.Errorf("期望验证成功但验证失败: %v", err)
			}
		})
	}
}

// TestMongoDBConfigToRepositoryConfigWithProfiling 测试转换配置包含profiling设置
func TestMongoDBConfigToRepositoryConfigWithProfiling(t *testing.T) {
	config := MongoDBConfig{
		URI:             "mongodb://localhost:27017",
		Database:        "test",
		ProfilingLevel:  1,
		SlowMS:          150,
		ProfilerSizeMB:  200,
		ConnectTimeout:  10 * time.Second,
		ServerTimeout:   30 * time.Second,
		MaxPoolSize:     100,
		MinPoolSize:     5,
	}

	repoConfig := config.ToRepositoryConfig()

	// 验证profiling配置被包含在转换后的配置中
	if level, ok := repoConfig["profiling_level"].(int); !ok || level != 1 {
		t.Errorf("期望profiling_level=1, 实际=%v", repoConfig["profiling_level"])
	}

	if slowMS, ok := repoConfig["slow_ms"].(int64); !ok || slowMS != 150 {
		t.Errorf("期望slow_ms=150, 实际=%v", repoConfig["slow_ms"])
	}

	if sizeMB, ok := repoConfig["profiler_size_mb"].(int64); !ok || sizeMB != 200 {
		t.Errorf("期望profiler_size_mb=200, 实际=%v", repoConfig["profiler_size_mb"])
	}
}
