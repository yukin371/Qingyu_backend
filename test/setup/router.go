package setup

import (
	"github.com/gin-gonic/gin"

	// 暂时不设置具体路由
	// 测试时会根据需要动态添加
)

// SetupRouter 设置测试路由
func SetupRouter() *gin.Engine {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// 注册基础中间件
	router.Use(gin.Recovery())

	return router
}
