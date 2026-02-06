// Package config provides configuration management for the seeder
package config

// Config holds the main configuration for the seeder
type Config struct {
	MongoDB   MongoDBConfig
	Scale     string
	Modules   []string
	Clean     bool
	BatchSize int
}

// MongoDBConfig holds MongoDB connection configuration
type MongoDBConfig struct {
	URI      string
	Database string
}

// ScaleConfig defines the data scale for generation
type ScaleConfig struct {
	Scale         string // 数据规模标识：small, medium, large
	Users         int
	Books         int
	Authors       int
	MinChapters   int // 最少章节数
	MaxChapters   int // 最多章节数
	// 角色比例配置
	ReaderPercent float64 // reader 角色比例 (默认 0.70)
	AuthorPercent float64 // author 角色比例 (默认 0.25)
	AdminPercent  float64 // admin 角色比例 (默认 0.05)
}

// DefaultConfig is the default configuration
var DefaultConfig = &Config{
	MongoDB: MongoDBConfig{
		URI:      "mongodb://localhost:27017",
		Database: "qingyu",
	},
	Scale:     "medium",
	BatchSize: 100,
}

// GetScaleConfig returns the scale configuration for the given scale name
func GetScaleConfig(scale string) ScaleConfig {
	scales := map[string]ScaleConfig{
		"small": {
			Scale:         "small",
			Users:         50,
			Books:         100,
			Authors:       10,
			MinChapters:   1,
			MaxChapters:   2,
			ReaderPercent: 0.70,
			AuthorPercent: 0.25,
			AdminPercent:  0.05,
		},
		"medium": {
			Scale:         "medium",
			Users:         500,
			Books:         500,
			Authors:       100,
			MinChapters:   3,
			MaxChapters:   5,
			ReaderPercent: 0.70,
			AuthorPercent: 0.25,
			AdminPercent:  0.05,
		},
		"large": {
			Scale:         "large",
			Users:         2000,
			Books:         1200,
			Authors:       400,
			MinChapters:   10,
			MaxChapters:   30,
			ReaderPercent: 0.70,
			AuthorPercent: 0.25,
			AdminPercent:  0.05,
		},
	}

	if cfg, ok := scales[scale]; ok {
		return cfg
	}
	return scales["medium"]
}
