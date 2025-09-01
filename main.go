package main

import (
	"Qingyu_backend/api"
	"Qingyu_backend/config"
	"Qingyu_backend/database"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 连接数据库
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("无法连接到数据库: ", err)
	}

	// 确保在程序结束时断开数据库连接
	defer database.DisconnectDB()

	// 3. 创建 Gin 引擎
	router := gin.Default()

	// 4. 注册路由
	api.RegisterRoutes(router)

	// 5. 关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的goroutine中启动服务
	go func() {
		// 在生产环境中，使用 os.Getenv("PORT") 来动态设置端口
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("正在关闭服务器...")

	// TODO: 可以添加其他清理工作
	log.Println("服务器已关闭")
}
