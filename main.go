package main

import (
	"log"
	"os"
	"path/filepath"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

func main() {
	// 获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认使用当前目录
		configPath = "."
	}
	configPath = filepath.Clean(configPath)

	// 加载配置
	_, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 注册配置重载处理器
	config.RegisterReloadHandler("database", func() {
		if err := core.InitDB(); err != nil {
			log.Printf("Failed to reload database configuration: %v", err)
		}
	})

	config.RegisterReloadHandler("server", func() {
		// 服务器配置变更时的处理逻辑
		log.Println("Server configuration reloaded")
	})

	// 启用配置热重载
	config.EnableHotReload()

	// 初始化数据库
	if err := core.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化服务器
	r, err := core.InitServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// 运行服务器
	if err := core.RunServer(r); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
