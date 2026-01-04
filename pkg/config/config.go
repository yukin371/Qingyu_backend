package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Manager 配置管理器
type Manager struct {
	viper  *viper.Viper
	mu     sync.RWMutex
	config interface{}
}

// NewManager 创建配置管理器
func NewManager() *Manager {
	return &Manager{
		viper: viper.New(),
	}
}

// LoadConfig 从文件加载配置
func (m *Manager) LoadConfig(configPath string, config interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 设置配置文件路径
	m.viper.SetConfigFile(configPath)

	// 获取文件扩展名
	ext := filepath.Ext(configPath)

	// 根据扩展名设置配置类型
	switch ext {
	case ".yaml", ".yml":
		m.viper.SetConfigType("yaml")
	case ".json":
		m.viper.SetConfigType("json")
	case ".toml":
		m.viper.SetConfigType("toml")
	case ".env":
		m.viper.SetConfigType("env")
	default:
		return fmt.Errorf("unsupported config file type: %s", ext)
	}

	// 启用环境变量支持
	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("APP")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := m.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置到结构体
	if err := m.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = config
	return nil
}

// LoadConfigFromDirectory 从目录加载配置
func (m *Manager) LoadConfigFromDirectory(configDir, configName string, config interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 设置配置名称和路径
	m.viper.SetConfigName(configName)
	m.viper.AddConfigPath(configDir)
	m.viper.AddConfigPath(".")

	// 支持的配置类型
	m.viper.SetConfigType("yaml")

	// 启用环境变量
	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("APP")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置
	if err := m.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// 解析配置
	if err := m.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = config
	return nil
}

// LoadConfigFromEnv 从环境变量加载配置
func (m *Manager) LoadConfigFromEnv(config interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.viper.AutomaticEnv()
	m.viper.SetEnvPrefix("APP")
	m.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := m.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = config
	return nil
}

// Get 获取配置值
func (m *Manager) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.Get(key)
}

// GetString 获取字符串配置
func (m *Manager) GetString(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetString(key)
}

// GetInt 获取整数配置
func (m *Manager) GetInt(key string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetInt(key)
}

// GetInt64 获取int64配置
func (m *Manager) GetInt64(key string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetInt64(key)
}

// GetFloat64 获取float64配置
func (m *Manager) GetFloat64(key string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetFloat64(key)
}

// GetBool 获取布尔配置
func (m *Manager) GetBool(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetBool(key)
}

// GetStringSlice 获取字符串数组配置
func (m *Manager) GetStringSlice(key string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.GetStringSlice(key)
}

// Set 设置配置值
func (m *Manager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.viper.Set(key, value)
}

// IsSet 检查配置是否设置
func (m *Manager) IsSet(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.IsSet(key)
}

// AllSettings 获取所有配置
func (m *Manager) AllSettings() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.viper.AllSettings()
}

// WatchConfig 监听配置文件变化
func (m *Manager) WatchConfig(callback func(in fsnotify.Event)) {
	m.viper.OnConfigChange(func(in fsnotify.Event) {
		m.mu.Lock()
		callback(in)
		m.mu.Unlock()
	})
	m.viper.WatchConfig()
}

// SaveConfig 保存配置到文件
func (m *Manager) SaveConfig(configPath string, config interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 根据文件扩展名选择编码器
	ext := filepath.Ext(configPath)
	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config to JSON: %w", err)
		}
	case ".yaml", ".yml":
		// YAML需要yaml包，这里简化处理
		return fmt.Errorf("YAML encoding not implemented")
	default:
		return fmt.Errorf("unsupported config file type: %s", ext)
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateConfig 验证配置
func (m *Manager) ValidateConfig(config interface{}) error {
	v := reflect.ValueOf(config)

	// 检查是否是指针
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("config must be a pointer")
	}

	// 解引用指针
	v = v.Elem()

	// 检查是否是结构体
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct")
	}

	// 遍历结构体字段
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取validate标签
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// 执行验证
		if err := m.validateField(value, validateTag); err != nil {
			return fmt.Errorf("field %s validation failed: %w", field.Name, err)
		}
	}

	return nil
}

// validateField 验证单个字段
func (m *Manager) validateField(value reflect.Value, validateTag string) error {
	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		parts := strings.Split(rule, "=")
		ruleName := parts[0]

		switch ruleName {
		case "required":
			if value.IsZero() {
				return fmt.Errorf("field is required")
			}
		case "min":
			if len(parts) < 2 {
				continue
			}
			minValue := parts[1]
			// 这里简化处理，实际应该根据类型比较
			if value.Kind() == reflect.String && value.String() < minValue {
				return fmt.Errorf("value is less than minimum: %s", minValue)
			}
		case "max":
			if len(parts) < 2 {
				continue
			}
			maxValue := parts[1]
			if value.Kind() == reflect.String && value.String() > maxValue {
				return fmt.Errorf("value is greater than maximum: %s", maxValue)
			}
		}
	}

	return nil
}

// OverrideWithEnv 使用环境变量覆盖配置
func (m *Manager) OverrideWithEnv(config interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取env标签
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// 从环境变量读取
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		// 设置值
		if err := m.setFieldValue(value, envValue); err != nil {
			return fmt.Errorf("failed to set field %s from env: %w", field.Name, err)
		}
	}

	return nil
}

// setFieldValue 设置字段值
func (m *Manager) setFieldValue(value reflect.Value, envValue string) error {
	if !value.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(envValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var num int64
		if _, err := fmt.Sscanf(envValue, "%d", &num); err != nil {
			return err
		}
		value.SetInt(num)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var num uint64
		if _, err := fmt.Sscanf(envValue, "%d", &num); err != nil {
			return err
		}
		value.SetUint(num)
	case reflect.Float32, reflect.Float64:
		var num float64
		if _, err := fmt.Sscanf(envValue, "%f", &num); err != nil {
			return err
		}
		value.SetFloat(num)
	case reflect.Bool:
		b := strings.ToLower(envValue) == "true" || envValue == "1"
		value.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %s", value.Kind())
	}

	return nil
}

// GetConfig 获取当前配置
func (m *Manager) GetConfig() interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}
