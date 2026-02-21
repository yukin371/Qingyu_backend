package config

import (
	"os"
	"testing"
	"time"
)

// TestMongoDBProfilingDefaults 测试MongoDB Profiling默认配置
func TestMongoDBProfilingDefaults(t *testing.T) {
	config := getDefaultDatabaseConfig()

	if config.Primary.MongoDB.ProfilingLevel != 1 {
		t.Errorf("Expected ProfilingLevel=1, got %d", config.Primary.MongoDB.ProfilingLevel)
	}
	if config.Primary.MongoDB.SlowMS != 100 {
		t.Errorf("Expected SlowMS=100, got %d", config.Primary.MongoDB.SlowMS)
	}
	if config.Primary.MongoDB.ProfilerSizeMB != 100 {
		t.Errorf("Expected ProfilerSizeMB=100, got %d", config.Primary.MongoDB.ProfilerSizeMB)
	}
}

// TestMongoDBProfilingValidation 测试MongoDB Profiling验证逻辑
func TestMongoDBProfilingValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  MongoDBConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with profiling",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 1,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
			wantErr: false,
		},
		{
			name: "valid config with profiling disabled",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 0,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
			wantErr: false,
		},
		{
			name: "valid config with full profiling",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 2,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid profiling level - negative",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: -1,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
			wantErr: true,
			errMsg:  "ProfilingLevel必须在0-2之间",
		},
		{
			name: "invalid profiling level - too high",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 5,
				SlowMS:         100,
				ProfilerSizeMB: 100,
			},
			wantErr: true,
			errMsg:  "ProfilingLevel必须在0-2之间",
		},
		{
			name: "invalid slow ms - negative",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 1,
				SlowMS:         -10,
				ProfilerSizeMB: 100,
			},
			wantErr: true,
			errMsg:  "SlowMS必须非负",
		},
		{
			name: "invalid profiler size - zero",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 1,
				SlowMS:         100,
				ProfilerSizeMB: 0,
			},
			wantErr: true,
			errMsg:  "ProfilerSizeMB必须至少为1MB",
		},
		{
			name: "invalid profiler size - negative",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 1,
				SlowMS:         100,
				ProfilerSizeMB: -10,
			},
			wantErr: true,
			errMsg:  "ProfilerSizeMB必须至少为1MB",
		},
		{
			name: "valid zero slow ms",
			config: MongoDBConfig{
				URI:            "mongodb://localhost:27017",
				Database:       "test",
				ProfilingLevel: 1,
				SlowMS:         0,
				ProfilerSizeMB: 100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestMongoDBConfigProfilingDefaults 测试MongoDB配置的默认值
func TestMongoDBConfigProfilingDefaults(t *testing.T) {
	config := &MongoDBConfig{
		URI:            "mongodb://localhost:27017",
		Database:       "test",
		ProfilingLevel: 0,
		SlowMS:         0,
		ProfilerSizeMB: 1, // 设置为最小合法值
	}

	// 验证时会填充默认值
	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() unexpectedly failed: %v", err)
	}

	// Profiling字段没有默认值，应该保持0值
	if config.ProfilingLevel != 0 {
		t.Errorf("Expected default ProfilingLevel=0, got %d", config.ProfilingLevel)
	}

	// 检查其他字段的默认值
	if config.MaxPoolSize != 100 {
		t.Errorf("Expected default MaxPoolSize=100, got %d", config.MaxPoolSize)
	}
	if config.MinPoolSize != 5 {
		t.Errorf("Expected default MinPoolSize=5, got %d", config.MinPoolSize)
	}
	if config.ConnectTimeout != 10*time.Second {
		t.Errorf("Expected default ConnectTimeout=10s, got %v", config.ConnectTimeout)
	}
	if config.ServerTimeout != 30*time.Second {
		t.Errorf("Expected default ServerTimeout=30s, got %v", config.ServerTimeout)
	}
}

// TestMongoDBConfigToRepositoryConfig 测试配置转换包含Profiling字段
func TestMongoDBConfigToRepositoryConfig(t *testing.T) {
	config := &MongoDBConfig{
		URI:            "mongodb://localhost:27017",
		Database:       "testdb",
		MaxPoolSize:    200,
		MinPoolSize:    10,
		ConnectTimeout: 20 * time.Second,
		ServerTimeout:  60 * time.Second,
		ProfilingLevel: 2,
		SlowMS:         200,
		ProfilerSizeMB: 200,
	}

	repoConfig := config.ToRepositoryConfig()

	// 验证基本字段
	if repoConfig["type"] != "mongodb" {
		t.Errorf("Expected type=mongodb, got %v", repoConfig["type"])
	}
	if repoConfig["uri"] != "mongodb://localhost:27017" {
		t.Errorf("Expected uri=mongodb://localhost:27017, got %v", repoConfig["uri"])
	}
	if repoConfig["database"] != "testdb" {
		t.Errorf("Expected database=testdb, got %v", repoConfig["database"])
	}

	// 验证Profiling配置正确传递
	if repoConfig["profiling_level"] != 2 {
		t.Errorf("Expected profiling_level=2, got %v", repoConfig["profiling_level"])
	}
	if repoConfig["slow_ms"] != int64(200) {
		t.Errorf("Expected slow_ms=200, got %v", repoConfig["slow_ms"])
	}
	if repoConfig["profiler_size_mb"] != int64(200) {
		t.Errorf("Expected profiler_size_mb=200, got %v", repoConfig["profiler_size_mb"])
	}
}

// TestMongoDBProfilingEnvOverrides 测试MongoDB Profiling环境变量覆盖
func TestMongoDBProfilingEnvOverrides(t *testing.T) {
	// 保存原始环境变量
	originalLevel := os.Getenv("MONGODB_PROFILING_LEVEL")
	originalSlowMS := os.Getenv("MONGODB_SLOW_MS")
	originalSizeMB := os.Getenv("MONGODB_PROFILER_SIZE_MB")

	defer func() {
		// 恢复原始值
		if originalLevel != "" {
			_ = os.Setenv("MONGODB_PROFILING_LEVEL", originalLevel)
		} else {
			_ = os.Unsetenv("MONGODB_PROFILING_LEVEL")
		}
		if originalSlowMS != "" {
			_ = os.Setenv("MONGODB_SLOW_MS", originalSlowMS)
		} else {
			_ = os.Unsetenv("MONGODB_SLOW_MS")
		}
		if originalSizeMB != "" {
			_ = os.Setenv("MONGODB_PROFILER_SIZE_MB", originalSizeMB)
		} else {
			_ = os.Unsetenv("MONGODB_PROFILER_SIZE_MB")
		}
	}()

	// 设置测试环境变量
	_ = os.Setenv("MONGODB_PROFILING_LEVEL", "2")
	_ = os.Setenv("MONGODB_SLOW_MS", "200")
	_ = os.Setenv("MONGODB_PROFILER_SIZE_MB", "200")

	// 加载配置
	config, err := LoadDatabaseConfig("")
	if err != nil {
		t.Fatalf("LoadDatabaseConfig failed: %v", err)
	}

	// 验证环境变量覆盖默认值
	if config.Primary.MongoDB.ProfilingLevel != 2 {
		t.Errorf("Expected ProfilingLevel=2 from env, got %d", config.Primary.MongoDB.ProfilingLevel)
	}
	if config.Primary.MongoDB.SlowMS != 200 {
		t.Errorf("Expected SlowMS=200 from env, got %d", config.Primary.MongoDB.SlowMS)
	}
	if config.Primary.MongoDB.ProfilerSizeMB != 200 {
		t.Errorf("Expected ProfilerSizeMB=200 from env, got %d", config.Primary.MongoDB.ProfilerSizeMB)
	}
}
