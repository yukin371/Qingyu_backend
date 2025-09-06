package config

import (
	"time"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoURI        string
	DBName          string
	ConnectTimeout  time.Duration
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
	RetryWrites     bool
	RetryReads      bool
}

// LoadDatabaseConfig 加载数据库配置
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		MongoURI:        getEnv("DB_URI", "mongodb://localhost:27017"), // URI
		DBName:          getEnv("DB_NAME", "Qingyu_writer"),            // 数据库名称
		ConnectTimeout:  time.Duration(getEnvAsInt("MONGO_CONNECT_TIMEOUT", 10)) * time.Second,
		MaxPoolSize:     uint64(getEnvAsInt("MONGO_MAX_POOL_SIZE", 100)),                          // 最大连接数
		MinPoolSize:     uint64(getEnvAsInt("MONGO_MIN_POOL_SIZE", 10)),                           // 最小连接数
		MaxConnIdleTime: time.Duration(getEnvAsInt("MONGO_MAX_CONN_IDLE_TIME", 60)) * time.Second, // 最大连接空闲时间
		RetryWrites:     getEnvAsBool("MONGO_RETRY_WRITES", true),                                 // 重试写入
		RetryReads:      getEnvAsBool("MONGO_RETRY_READS", true),                                  // 重试读取
	}
}
