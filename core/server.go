package core

import (
	"fmt"

	"Qingyu_backend/config"
	"Qingyu_backend/middleware"
	"Qingyu_backend/router"

	"github.com/gin-gonic/gin"
)

// InitServer 初始化服务器
func InitServer() (*gin.Engine, error) {
	cfg := config.GlobalConfig.Server
	if cfg == nil {
		return nil, fmt.Errorf("server configuration is missing")
	}

	// 设置gin模式
	gin.SetMode(cfg.Mode)

	// 创建gin实例
	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORSMiddleware())

	// 注册路由
	router.RegisterRoutes(r)

	return r, nil
}

// RunServer 运行服务器
func RunServer(r *gin.Engine) error {
	cfg := config.GlobalConfig.Server
	if cfg == nil {
		return fmt.Errorf("server configuration is missing")
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("Server is running on port %s in %s mode\n", cfg.Port, cfg.Mode)
	return r.Run(addr)
}
