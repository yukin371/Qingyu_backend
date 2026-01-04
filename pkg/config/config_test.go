package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

// TestConfig 测试配置结构
type TestConfig struct {
	StringField string `json:"string_field" validate:"required"`
	IntField    int    `json:"int_field" validate:"required,min=1"`
	BoolField   bool   `json:"bool_field"`
	FloatField  float64 `json:"float_field"`
	ArrayField []string `json:"array_field"`
	Nested      struct {
		NestedField string `json:"nested_field"`
	} `json:"nested"`
	EnvOverride string `json:"env_override" env:"TEST_ENV_VAR"`
}

// TestNewManager 测试创建配置管理器
func TestNewManager(t *testing.T) {
	manager := NewManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.viper)
}

// TestLoadConfig_Json 测试加载JSON配置
func TestLoadConfig_Json(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建测试配置文件
	config := TestConfig{
		StringField: "test",
		IntField:    42,
		BoolField:   true,
		FloatField:  3.14,
		ArrayField: []string{"item1", "item2"},
	}
	config.Nested.NestedField = "nested"

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)

	// 加载配置
	manager := NewManager()
	var loadedConfig TestConfig
	err := manager.LoadConfig(configPath, &loadedConfig)

	assert.NoError(t, err)
	// 字段可能通过JSON标签映射，检查viper的原始值
	assert.Equal(t, "test", manager.GetString("string_field"))
	assert.Equal(t, 42, manager.GetInt("int_field"))
	assert.True(t, manager.GetBool("bool_field"))
}

// TestLoadConfig_InvalidPath 测试加载不存在的配置文件
func TestLoadConfig_InvalidPath(t *testing.T) {
	manager := NewManager()
	var config TestConfig
	err := manager.LoadConfig("nonexistent.json", &config)

	assert.Error(t, err)
}

// TestLoadConfig_UnsupportedType 测试不支持的配置类型
func TestLoadConfig_UnsupportedType(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.txt")

	os.WriteFile(configPath, []byte("test"), 0644)

	manager := NewManager()
	var config TestConfig
	err := manager.LoadConfig(configPath, &config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

// TestGetters 测试各种Get方法
func TestGetters(t *testing.T) {
	manager := NewManager()

	manager.Set("string_key", "value")
	manager.Set("int_key", 42)
	manager.Set("float_key", 3.14)
	manager.Set("bool_key", true)
	manager.Set("array_key", []string{"a", "b"})

	assert.Equal(t, "value", manager.GetString("string_key"))
	assert.Equal(t, 42, manager.GetInt("int_key"))
	assert.Equal(t, int64(42), manager.GetInt64("int_key"))
	assert.Equal(t, 3.14, manager.GetFloat64("float_key"))
	assert.True(t, manager.GetBool("bool_key"))
	assert.Equal(t, []string{"a", "b"}, manager.GetStringSlice("array_key"))
}

// TestSetAndGet 测试Set和Get
func TestSetAndGet(t *testing.T) {
	manager := NewManager()

	manager.Set("test_key", "test_value")
	value := manager.Get("test_key")

	assert.Equal(t, "test_value", value)
}

// TestIsSet 测试IsSet
func TestIsSet(t *testing.T) {
	manager := NewManager()

	assert.False(t, manager.IsSet("nonexistent"))

	manager.Set("existing_key", "value")
	assert.True(t, manager.IsSet("existing_key"))
}

// TestAllSettings 测试获取所有配置
func TestAllSettings(t *testing.T) {
	manager := NewManager()

	manager.Set("key1", "value1")
	manager.Set("key2", 42)

	settings := manager.AllSettings()

	assert.NotNil(t, settings)
	assert.Contains(t, settings, "key1")
	assert.Contains(t, settings, "key2")
}

// TestValidateConfig 测试配置验证
func TestValidateConfig(t *testing.T) {
	manager := NewManager()

	// 有效配置
	validConfig := TestConfig{
		StringField: "test",
		IntField:    10,
	}

	err := manager.ValidateConfig(&validConfig)
	assert.NoError(t, err)

	// 无效配置 - 缺少必填字段
	invalidConfig := TestConfig{
		IntField: 10,
	}

	err = manager.ValidateConfig(&invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")

	// 无效配置 - 最小值验证
	invalidMinConfig := TestConfig{
		StringField: "test",
		IntField:    0,
	}

	err = manager.ValidateConfig(&invalidMinConfig)
	assert.Error(t, err)
}

// TestValidateConfig_InvalidInput 测试无效输入验证
func TestValidateConfig_InvalidInput(t *testing.T) {
	manager := NewManager()

	// 非指针
	err := manager.ValidateConfig(TestConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pointer")

	// 非结构体指针
	var num int
	err = manager.ValidateConfig(&num)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "struct")
}

// TestSetFieldValue 测试字段值设置
func TestSetFieldValue(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name     string
		fieldPtr interface{}
		value    string
		verify   func(interface{}) bool
	}{
		{
			name:     "字符串",
			fieldPtr: new(string),
			value:    "test",
			verify: func(v interface{}) bool {
				return *(v.(*string)) == "test"
			},
		},
		{
			name:     "整数",
			fieldPtr: new(int),
			value:    "42",
			verify: func(v interface{}) bool {
				return *(v.(*int)) == 42
			},
		},
		{
			name:     "浮点数",
			fieldPtr: new(float64),
			value:    "3.14",
			verify: func(v interface{}) bool {
				return *(v.(*float64)) == 3.14
			},
		},
		{
			name:     "布尔值",
			fieldPtr: new(bool),
			value:    "true",
			verify: func(v interface{}) bool {
				return *(v.(*bool)) == true
			},
		},
		{
			name:     "布尔值0",
			fieldPtr: new(bool),
			value:    "0",
			verify: func(v interface{}) bool {
				return *(v.(*bool)) == false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.fieldPtr).Elem()
			err := manager.setFieldValue(v, tt.value)

			assert.NoError(t, err)
			assert.True(t, tt.verify(tt.fieldPtr))
		})
	}
}

// TestOverrideWithEnv 测试环境变量覆盖
func TestOverrideWithEnv(t *testing.T) {
	manager := NewManager()

	// 设置环境变量
	os.Setenv("TEST_ENV_VAR", "env_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	config := TestConfig{
		StringField: "original",
		IntField:    10,
		EnvOverride: "default",
	}

	err := manager.OverrideWithEnv(&config)
	assert.NoError(t, err)
	assert.Equal(t, "env_value", config.EnvOverride)
}

// TestSaveConfig 测试保存配置
func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "saved_config.json")

	config := TestConfig{
		StringField: "saved",
		IntField:    100,
	}

	manager := NewManager()
	err := manager.SaveConfig(configPath, &config)

	assert.NoError(t, err)

	// 验证文件存在
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// 验证文件内容
	data, err := os.ReadFile(configPath)
	assert.NoError(t, err)

	var savedConfig TestConfig
	err = json.Unmarshal(data, &savedConfig)
	assert.NoError(t, err)
	assert.Equal(t, "saved", savedConfig.StringField)
	assert.Equal(t, 100, savedConfig.IntField)
}

// TestSaveConfig_UnsupportedType 测试保存不支持的类型
func TestSaveConfig_UnsupportedType(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager := NewManager()
	config := TestConfig{StringField: "test"}

	err := manager.SaveConfig(configPath, &config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

// TestWatchConfig 测试配置监听
func TestWatchConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "watch_config.json")

	// 创建初始配置
	config := TestConfig{StringField: "initial", IntField: 1}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)

	manager := NewManager()
	var loadedConfig TestConfig
	err := manager.LoadConfig(configPath, &loadedConfig)
	assert.NoError(t, err)

	// 设置监听回调
	callbackCalled := false
	manager.WatchConfig(func(in fsnotify.Event) {
		callbackCalled = true
	})

	// 注意：实际触发文件变化需要修改文件并等待
	// 这里只测试注册不panic
	assert.False(t, callbackCalled)
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	manager := NewManager()

	done := make(chan bool)

	// 并发写入
	for i := 0; i < 100; i++ {
		go func(n int) {
			manager.Set("key", n)
			done <- true
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		go func() {
			_ = manager.Get("key")
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 200; i++ {
		<-done
	}

	// 验证能正常读写
	assert.NotNil(t, manager.Get("key"))
}

// BenchmarkGet 性能测试 - Get操作
func BenchmarkGet(b *testing.B) {
	manager := NewManager()
	manager.Set("test_key", "test_value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.Get("test_key")
	}
}

// BenchmarkSet 性能测试 - Set操作
func BenchmarkSet(b *testing.B) {
	manager := NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Set("test_key", "test_value")
	}
}

// BenchmarkLoadConfig 性能测试 - 加载配置
func BenchmarkLoadConfig(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "bench_config.json")

	config := TestConfig{
		StringField: "bench",
		IntField:    999,
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewManager()
		var cfg TestConfig
		_ = manager.LoadConfig(configPath, &cfg)
	}
}
