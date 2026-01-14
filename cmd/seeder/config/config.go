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
	Users   int
	Books   int
	Authors int
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
			Users:   50,
			Books:   100,
			Authors: 20,
		},
		"medium": {
			Users:   500,
			Books:   500,
			Authors: 100,
		},
		"large": {
			Users:   2000,
			Books:   1200,
			Authors: 400,
		},
	}

	if cfg, ok := scales[scale]; ok {
		return cfg
	}
	return scales["medium"]
}
