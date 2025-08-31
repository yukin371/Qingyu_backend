package main

import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/api"
    "Qingyu_backend/config"
    "Qingyu_backend/database"
)

func main() {
    // 1. 加载配置
    cfg := config.LoadConfig()

    // 2. 连接数据库
    database.ConnectDB(cfg)

    // 3. 创建 Gin 引擎
    router := gin.Default()

    // 4. 注册路由
    api.RegisterRoutes(router)

    // 5. 启动服务
    // 在生产环境中，使用 os.Getenv("PORT") 来动态设置端口
    err := router.Run(":8080")
    if err != nil {
        panic("Failed to run server: " + err.Error())
    }
}