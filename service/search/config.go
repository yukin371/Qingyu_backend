package search

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// SearchConfig 搜索服务配置（从配置文件加载）
type SearchConfig struct {
	// 缓存配置
	Cache CacheConfig `yaml:"cache"`
	// 限流配置
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	// 优化器配置
	Optimizer OptimizerConfig `yaml:"optimizer"`
	// 书籍搜索配置
	Books BookSearchConfig `yaml:"books"`
	// 项目搜索配置
	Projects ProjectSearchConfig `yaml:"projects"`
	// 文档搜索配置
	Documents DocumentSearchConfig `yaml:"documents"`
	// 用户搜索配置
	Users UserSearchConfig `yaml:"users"`

	// Elasticsearch 配置
	ES ESConfig `yaml:"elasticsearch"`
}

// SearchIndicesConfig 索引配置（从 search_indices.yaml 加载）
type SearchIndicesConfig struct {
	Indices map[string]IndexConfig `yaml:"indices"`
	Settings SettingsConfig         `yaml:"settings"`
}

// IndexConfig 单个索引配置
type IndexConfig struct {
	Alias            string                 `yaml:"alias"`
	NumberOfShards   int                    `yaml:"number_of_shards"`
	NumberOfReplicas int                    `yaml:"number_of_replicas"`
	Mapping          map[string]interface{} `yaml:"mapping"`
}

// SettingsConfig 全局设置配置
type SettingsConfig struct {
	Analysis AnalysisConfig `yaml:"analysis"`
}

// AnalysisConfig 分词器配置
type AnalysisConfig struct {
	Analyzer map[string]AnalyzerConfig `yaml:"analyzer"`
}

// AnalyzerConfig 分词器定义
type AnalyzerConfig struct {
	Type      string `yaml:"type"`
	Tokenizer string `yaml:"tokenizer"`
}

// LoadSearchIndicesConfig 加载索引配置文件
func LoadSearchIndicesConfig(configPath string) (*SearchIndicesConfig, error) {
	// 支持环境变量覆盖配置路径
	if configPath == "" {
		configPath = os.Getenv("SEARCH_INDICES_CONFIG")
		if configPath == "" {
			configPath = "config/search_indices.yaml"
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read search indices config file: %w", err)
	}

	// 解析 YAML
	var config SearchIndicesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse search indices config: %w", err)
	}

	// 验证配置
	if err := validateIndicesConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid search indices config: %w", err)
	}

	return &config, nil
}

// GetDefaultIndicesConfig 获取默认索引配置（当配置文件不存在时）
func GetDefaultIndicesConfig() *SearchIndicesConfig {
	return &SearchIndicesConfig{
		Indices: map[string]IndexConfig{
			"books": {
				Alias:            "books_search",
				NumberOfShards:   3,
				NumberOfReplicas: 1,
				Mapping: map[string]interface{}{
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":            "text",
							"analyzer":        "ik_max_word",
							"search_analyzer": "ik_smart",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"content": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
						},
						"author": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"status": map[string]interface{}{
							"type": "keyword",
						},
						"created_at": map[string]interface{}{
							"type": "date",
						},
						"updated_at": map[string]interface{}{
							"type": "date",
						},
					},
				},
			},
			"projects": {
				Alias:            "projects_search",
				NumberOfShards:   2,
				NumberOfReplicas: 1,
				Mapping: map[string]interface{}{
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
						},
						"description": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
						},
						"user_id": map[string]interface{}{
							"type": "keyword",
						},
						"status": map[string]interface{}{
							"type": "keyword",
						},
						"created_at": map[string]interface{}{
							"type": "date",
						},
					},
				},
			},
			"documents": {
				Alias:            "documents_search",
				NumberOfShards:   2,
				NumberOfReplicas: 1,
				Mapping: map[string]interface{}{
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
						},
						"content": map[string]interface{}{
							"type":     "text",
							"analyzer": "ik_max_word",
						},
						"project_id": map[string]interface{}{
							"type": "keyword",
						},
						"user_id": map[string]interface{}{
							"type": "keyword",
						},
						"created_at": map[string]interface{}{
							"type": "date",
						},
					},
				},
			},
		},
		Settings: SettingsConfig{
			Analysis: AnalysisConfig{
				Analyzer: map[string]AnalyzerConfig{
					"ik_max_word": {
						Type:      "custom",
						Tokenizer: "ik_max_word",
					},
					"ik_smart": {
						Type:      "custom",
						Tokenizer: "ik_smart",
					},
				},
			},
		},
	}
}

// validateIndicesConfig 验证索引配置
func validateIndicesConfig(config *SearchIndicesConfig) error {
	if len(config.Indices) == 0 {
		return fmt.Errorf("no indices defined")
	}

	// 验证每个索引配置
	for name, indexConfig := range config.Indices {
		if indexConfig.Alias == "" {
			return fmt.Errorf("index '%s': alias is required", name)
		}
		if indexConfig.NumberOfShards <= 0 {
			return fmt.Errorf("index '%s': number_of_shards must be positive", name)
		}
		if indexConfig.NumberOfReplicas < 0 {
			return fmt.Errorf("index '%s': number_of_replicas cannot be negative", name)
		}
		if indexConfig.Mapping == nil {
			return fmt.Errorf("index '%s': mapping is required", name)
		}
	}

	return nil
}

// BuildIndexMapping 构建索引映射（包含 settings 和 mappings）
func (c *IndexConfig) BuildIndexMapping(globalSettings SettingsConfig) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   c.NumberOfShards,
			"number_of_replicas": c.NumberOfReplicas,
			"analysis":           globalSettings.Analysis,
		},
		"mappings": c.Mapping,
	}
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled    bool          `yaml:"enabled"`
	DefaultTTL time.Duration `yaml:"default_ttl"`
	HotTTL     time.Duration `yaml:"hot_ttl"`
	KeyPrefix  string        `yaml:"key_prefix"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
}

// OptimizerConfig 查询优化器配置
type OptimizerConfig struct {
	MaxResults      int  `yaml:"max_results"`
	MaxPageSize     int  `yaml:"max_page_size"`
	MinQueryLength  int  `yaml:"min_query_length"`
	EnableFuzziness bool `yaml:"enable_fuzziness"`
}

// BookSearchConfig 书籍搜索配置
type BookSearchConfig struct {
	// 可搜索的状态列表
	AllowedStatuses []string `yaml:"allowed_statuses"`
	// 可搜索的隐私设置
	AllowedPrivacy []bool `yaml:"allowed_privacy"`
}

// ProjectSearchConfig 项目搜索配置
type ProjectSearchConfig struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`
}

// DocumentSearchConfig 文档搜索配置
type DocumentSearchConfig struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`
}

// UserSearchConfig 用户搜索配置
type UserSearchConfig struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`
}

// ESConfig Elasticsearch 配置
type ESConfig struct {
	// 是否启用 ES
	Enabled bool `yaml:"enabled"`
	// ES 地址
	URL string `yaml:"url"`
	// 索引前缀
	IndexPrefix string `yaml:"index_prefix"`
	// 灰度发布配置
	GrayScale GrayScaleConfig `yaml:"grayscale"`
}

// GrayScaleConfig 灰度发布配置
type GrayScaleConfig struct {
	// 是否启用灰度
	Enabled bool `yaml:"enabled"`
	// 灰度流量百分比(0-100)
	Percent int `yaml:"percent"`
}
