package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/router"

	"github.com/gin-gonic/gin"
)

// InitServer 初始化服务器
func InitServer() *gin.Engine {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	app := gin.Default()

	// 注册路由
	router.RegisterRoutes(app)

	return app
}

// RunServer 运行服务器
func RunServer() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化服务器
	app := InitServer()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: app,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	CloseMongoDB()

	log.Println("Server exited")
}
