package main

import (
	"log"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

// @title           青羽写作平台 API
// @version         1.0
// @description     青羽写作平台后端服务API文档
// @host            localhost:9090
// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("[Main] Configuration loaded successfully")

	// 初始化服务器（包含数据库和路由）
	r, err := core.InitServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// 运行服务器
	if err := core.RunServer(r); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
