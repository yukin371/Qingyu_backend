package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"Qingyu_backend/config"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ConfigService 配置管理服务
type ConfigService struct {
	configPath string
	mu         sync.RWMutex
	logger     *zap.Logger
}

// NewConfigService 创建配置管理服务
func NewConfigService(configPath string) *ConfigService {
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	return &ConfigService{
		configPath: configPath,
		logger:     zap.L(),
	}
}

// ConfigItem 配置项
type ConfigItem struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`        // string, number, boolean, object
	Description string      `json:"description"` // 配置说明
	Editable    bool        `json:"editable"`    // 是否可编辑
	Sensitive   bool        `json:"sensitive"`   // 是否敏感信息（如密码）
}

// ConfigGroup 配置组
type ConfigGroup struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Items       []*ConfigItem `json:"items"`
}

// GetAllConfigs 获取所有配置（分组显示）
func (s *ConfigService) GetAllConfigs(ctx context.Context) ([]*ConfigGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, fmt.Errorf("配置未加载")
	}

	groups := []*ConfigGroup{
		{
			Name:        "server",
			Description: "服务器配置",
			Items: []*ConfigItem{
				{
					Key:         "server.port",
					Value:       cfg.Server.Port,
					Type:        "string",
					Description: "服务器端口",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "server.mode",
					Value:       cfg.Server.Mode,
					Type:        "string",
					Description: "运行模式 (debug/release)",
					Editable:    true,
					Sensitive:   false,
				},
			},
		},
		{
			Name:        "database",
			Description: "数据库配置",
			Items: []*ConfigItem{
				{
					Key:         "database.uri",
					Value:       s.maskSensitiveValue(cfg.Database.Primary.MongoDB.URI),
					Type:        "string",
					Description: "MongoDB连接URI",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "database.name",
					Value:       cfg.Database.Primary.MongoDB.Database,
					Type:        "string",
					Description: "数据库名称",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "database.max_pool_size",
					Value:       cfg.Database.Primary.MongoDB.MaxPoolSize,
					Type:        "number",
					Description: "最大连接池大小",
					Editable:    true,
					Sensitive:   false,
				},
			},
		},
		{
			Name:        "jwt",
			Description: "JWT配置",
			Items: []*ConfigItem{
				{
					Key:         "jwt.secret",
					Value:       s.maskSensitiveValue(cfg.JWT.Secret),
					Type:        "string",
					Description: "JWT密钥",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "jwt.expiration_hours",
					Value:       cfg.JWT.ExpirationHours,
					Type:        "number",
					Description: "Token过期时间（小时）",
					Editable:    true,
					Sensitive:   false,
				},
			},
		},
	}

	// 如果有Redis配置
	if cfg.Redis != nil {
		groups = append(groups, &ConfigGroup{
			Name:        "redis",
			Description: "Redis配置",
			Items: []*ConfigItem{
				{
					Key:         "redis.host",
					Value:       cfg.Redis.Host,
					Type:        "string",
					Description: "Redis主机地址",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "redis.port",
					Value:       cfg.Redis.Port,
					Type:        "number",
					Description: "Redis端口",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "redis.password",
					Value:       s.maskSensitiveValue(cfg.Redis.Password),
					Type:        "string",
					Description: "Redis密码",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "redis.db",
					Value:       cfg.Redis.DB,
					Type:        "number",
					Description: "Redis数据库编号",
					Editable:    true,
					Sensitive:   false,
				},
			},
		})
	}

	// 如果有AI配置
	if cfg.AI != nil {
		groups = append(groups, &ConfigGroup{
			Name:        "ai",
			Description: "AI服务配置",
			Items: []*ConfigItem{
				{
					Key:         "ai.api_key",
					Value:       s.maskSensitiveValue(cfg.AI.APIKey),
					Type:        "string",
					Description: "AI API密钥",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "ai.base_url",
					Value:       cfg.AI.BaseURL,
					Type:        "string",
					Description: "AI API基础URL",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "ai.max_tokens",
					Value:       cfg.AI.MaxTokens,
					Type:        "number",
					Description: "最大Token数",
					Editable:    true,
					Sensitive:   false,
				},
			},
		})
	}

	// 邮件配置组
	if cfg.Email != nil {
		groups = append(groups, &ConfigGroup{
			Name:        "email",
			Description: "邮件配置",
			Items: []*ConfigItem{
				{
					Key:         "email.enabled",
					Value:       cfg.Email.Enabled,
					Type:        "boolean",
					Description: "是否启用邮件服务",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "email.smtp_host",
					Value:       cfg.Email.SMTPHost,
					Type:        "string",
					Description: "SMTP服务器地址",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "email.smtp_port",
					Value:       cfg.Email.SMTPPort,
					Type:        "number",
					Description: "SMTP端口",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "email.username",
					Value:       s.maskSensitiveValue(cfg.Email.Username),
					Type:        "string",
					Description: "SMTP用户名",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "email.password",
					Value:       s.maskSensitiveValue(cfg.Email.Password),
					Type:        "string",
					Description: "SMTP密码",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "email.from_address",
					Value:       cfg.Email.FromAddress,
					Type:        "string",
					Description: "发件人邮箱地址",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "email.from_name",
					Value:       cfg.Email.FromName,
					Type:        "string",
					Description: "发件人名称",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "email.use_tls",
					Value:       cfg.Email.UseTLS,
					Type:        "boolean",
					Description: "是否使用TLS",
					Editable:    true,
					Sensitive:   false,
				},
			},
		})
	}

	// 支付配置组
	if cfg.Payment != nil {
		paymentItems := []*ConfigItem{
			{
				Key:         "payment.enabled",
				Value:       cfg.Payment.Enabled,
				Type:        "boolean",
				Description: "是否启用支付功能",
				Editable:    true,
				Sensitive:   false,
			},
			{
				Key:         "payment.default_provider",
				Value:       cfg.Payment.DefaultProvider,
				Type:        "string",
				Description: "默认支付提供商 (alipay/wechat)",
				Editable:    true,
				Sensitive:   false,
			},
		}

		// 支付宝配置
		if cfg.Payment.Alipay != nil {
			paymentItems = append(paymentItems, []*ConfigItem{
				{
					Key:         "payment.alipay.enabled",
					Value:       cfg.Payment.Alipay.Enabled,
					Type:        "boolean",
					Description: "是否启用支付宝",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "payment.alipay.app_id",
					Value:       s.maskSensitiveValue(cfg.Payment.Alipay.AppID),
					Type:        "string",
					Description: "支付宝应用ID",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "payment.alipay.sandbox",
					Value:       cfg.Payment.Alipay.Sandbox,
					Type:        "boolean",
					Description: "是否使用沙箱环境",
					Editable:    true,
					Sensitive:   false,
				},
			}...)
		}

		// 微信支付配置
		if cfg.Payment.Wechat != nil {
			paymentItems = append(paymentItems, []*ConfigItem{
				{
					Key:         "payment.wechat.enabled",
					Value:       cfg.Payment.Wechat.Enabled,
					Type:        "boolean",
					Description: "是否启用微信支付",
					Editable:    true,
					Sensitive:   false,
				},
				{
					Key:         "payment.wechat.app_id",
					Value:       s.maskSensitiveValue(cfg.Payment.Wechat.AppID),
					Type:        "string",
					Description: "微信应用ID",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "payment.wechat.mch_id",
					Value:       s.maskSensitiveValue(cfg.Payment.Wechat.MchID),
					Type:        "string",
					Description: "微信商户ID",
					Editable:    true,
					Sensitive:   true,
				},
				{
					Key:         "payment.wechat.sandbox",
					Value:       cfg.Payment.Wechat.Sandbox,
					Type:        "boolean",
					Description: "是否使用沙箱环境",
					Editable:    true,
					Sensitive:   false,
				},
			}...)
		}

		groups = append(groups, &ConfigGroup{
			Name:        "payment",
			Description: "支付配置",
			Items:       paymentItems,
		})
	}

	return groups, nil
}

// GetConfigByKey 根据Key获取单个配置
func (s *ConfigService) GetConfigByKey(ctx context.Context, key string) (*ConfigItem, error) {
	groups, err := s.GetAllConfigs(ctx)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		for _, item := range group.Items {
			if item.Key == key {
				return item, nil
			}
		}
	}

	return nil, fmt.Errorf("配置项不存在: %s", key)
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(ctx context.Context, req *UpdateConfigRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 验证配置项是否可编辑
	item, err := s.GetConfigByKey(ctx, req.Key)
	if err != nil {
		return err
	}
	if !item.Editable {
		return fmt.Errorf("配置项 %s 不可编辑", req.Key)
	}

	// 2. 读取当前配置文件
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 3. 解析YAML
	var yamlConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 4. 更新配置值
	if err := s.setNestedValue(yamlConfig, req.Key, req.Value); err != nil {
		return fmt.Errorf("更新配置值失败: %w", err)
	}

	// 5. 备份原配置文件
	backupPath := s.configPath + ".backup"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		s.logger.Warn("创建备份文件失败", zap.Error(err))
	}

	// 6. 写入新配置
	newData, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(s.configPath, newData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	// 7. 重新加载配置
	if err := s.reloadConfig(); err != nil {
		// 如果重新加载失败，恢复备份
		s.logger.Error("重新加载配置失败，恢复备份", zap.Error(err))
		if restoreErr := os.Rename(backupPath, s.configPath); restoreErr != nil {
			s.logger.Error("恢复备份失败", zap.Error(restoreErr))
		}
		return fmt.Errorf("重新加载配置失败: %w", err)
	}

	s.logger.Info("配置更新成功",
		zap.String("key", req.Key),
		zap.Any("value", req.Value),
	)

	return nil
}

// BatchUpdateConfig 批量更新配置
func (s *ConfigService) BatchUpdateConfig(ctx context.Context, requests []*UpdateConfigRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 读取当前配置文件
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 2. 解析YAML
	var yamlConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 3. 批量更新配置值
	for _, req := range requests {
		if err := s.setNestedValue(yamlConfig, req.Key, req.Value); err != nil {
			return fmt.Errorf("更新配置值 %s 失败: %w", req.Key, err)
		}
	}

	// 4. 备份原配置文件
	backupPath := s.configPath + ".backup"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		s.logger.Warn("创建备份文件失败", zap.Error(err))
	}

	// 5. 写入新配置
	newData, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(s.configPath, newData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	// 6. 重新加载配置
	if err := s.reloadConfig(); err != nil {
		s.logger.Error("重新加载配置失败，恢复备份", zap.Error(err))
		if restoreErr := os.Rename(backupPath, s.configPath); restoreErr != nil {
			s.logger.Error("恢复备份失败", zap.Error(restoreErr))
		}
		return fmt.Errorf("重新加载配置失败: %w", err)
	}

	s.logger.Info("批量配置更新成功", zap.Int("count", len(requests)))

	return nil
}

// ValidateConfig 验证配置
func (s *ConfigService) ValidateConfig(ctx context.Context, yamlContent string) error {
	var cfg config.Config
	if err := yaml.Unmarshal([]byte(yamlContent), &cfg); err != nil {
		return fmt.Errorf("YAML格式错误: %w", err)
	}

	return config.ValidateConfig(&cfg)
}

// GetConfigBackups 获取配置备份列表
func (s *ConfigService) GetConfigBackups(ctx context.Context) ([]string, error) {
	dir := filepath.Dir(s.configPath)
	fileName := filepath.Base(s.configPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Base(entry.Name()) == fileName+".backup" {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// RestoreConfigBackup 恢复配置备份
func (s *ConfigService) RestoreConfigBackup(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	backupPath := s.configPath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在")
	}

	// 备份当前文件
	currentBackup := s.configPath + ".current"
	if err := os.Rename(s.configPath, currentBackup); err != nil {
		return fmt.Errorf("备份当前配置失败: %w", err)
	}

	// 恢复备份
	if err := os.Rename(backupPath, s.configPath); err != nil {
		// 恢复失败，还原当前配置
		os.Rename(currentBackup, s.configPath)
		return fmt.Errorf("恢复备份失败: %w", err)
	}

	// 重新加载配置
	if err := s.reloadConfig(); err != nil {
		s.logger.Error("重新加载配置失败", zap.Error(err))
		return fmt.Errorf("重新加载配置失败: %w", err)
	}

	s.logger.Info("配置恢复成功")
	return nil
}

// 私有方法

// maskSensitiveValue 掩码敏感值
func (s *ConfigService) maskSensitiveValue(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}

// setNestedValue 设置嵌套的配置值
func (s *ConfigService) setNestedValue(data map[string]interface{}, key string, value interface{}) error {
	// 分割key，如 "server.port" -> ["server", "port"]
	parts := splitKey(key)
	if len(parts) == 0 {
		return fmt.Errorf("无效的配置键: %s", key)
	}

	// 遍历到最后一层
	current := data
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			// 创建中间层
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}

	// 设置最后一层的值
	current[parts[len(parts)-1]] = value
	return nil
}

// reloadConfig 重新加载配置
func (s *ConfigService) reloadConfig() error {
	_, err := config.LoadConfig(filepath.Dir(s.configPath))
	return err
}

// splitKey 分割配置键
func splitKey(key string) []string {
	var parts []string
	current := ""
	for _, ch := range key {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
