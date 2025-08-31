package config

import (
    "os"
)

// Config 存储应用配置
type Config struct {
    MongoURI string
    DBName   string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
    return &Config{
        MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
        DBName:   getEnv("DB_NAME", "Qingyu_writer"),
    }
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
    }
    return value
}