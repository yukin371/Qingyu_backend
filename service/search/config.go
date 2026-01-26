package search

import "time"

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
