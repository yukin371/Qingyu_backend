package api

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(router *gin.Engine) {
    router.GET("/ping", PingHandler)
}