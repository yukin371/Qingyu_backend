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
		log.Fatalf("配置加载失败: %v", err)
	}
	log.Println("[主程序] 配置加载成功")

	// 注册配置重载处理器
	config.RegisterReloadHandler("database", func() {
		log.Println("数据库配置已变更，请重启服务以生效")
	})

	config.RegisterReloadHandler("server", func() {
		log.Println("服务器配置已重新加载")
	})

	// 启用配置热重载
	config.EnableHotReload()
	log.Println("[主程序] 配置热重载已启用")

	// 初始化服务器（包含数据库和路由）
	r, err := core.InitServer()
	if err != nil {
		log.Fatalf("服务器初始化失败: %v", err)
	}

	// 运行服务器
	if err := core.RunServer(r); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}
