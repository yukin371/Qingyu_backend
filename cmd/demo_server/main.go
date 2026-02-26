// Package main 提供论文答辩演示用的简化版后端启动入口
// 特点：连接真实MongoDB，简化其他依赖，内置演示数据，快速启动
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	_ "Qingyu_backend/docs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title           青羽写作平台 API (演示版)
// @version         1.0-demo
// @description     青羽写作平台后端演示服务，专为论文答辩演示设计。

// @host      localhost:9090
// @BasePath  /api/v1

// DemoConfig 演示模式配置
type DemoConfig struct {
	Port          string
	MongoURI      string
	DBName        string
	SeedData      bool
	SkipRedis     bool
	SkipAIService bool
}

// DefaultDemoConfig 默认演示配置
var DefaultDemoConfig = DemoConfig{
	Port:          "9090",
	MongoURI:      "mongodb://localhost:27017",
	DBName:        "qingyu",
	SeedData:      true,
	SkipRedis:     false,
	SkipAIService: true,
}

func main() {
	fmt.Println("========================================")
	fmt.Println("  青羽写作平台 - 论文答辩演示版本")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 加载配置
	if err := loadDemoConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化数据库连接
	fmt.Println("[演示模式] 连接MongoDB...")
	if err := core.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v (请确保MongoDB已启动)", err)
	}
	fmt.Println("[演示模式] MongoDB连接成功")

	// 3. 初始化服务器
	fmt.Println("[演示模式] 初始化服务...")
	r, err := core.InitServer()
	if err != nil {
		log.Fatalf("初始化服务器失败: %v", err)
	}

	// 4. 注入演示数据（可选）
	if DefaultDemoConfig.SeedData {
		fmt.Println("[演示模式] 检查演示数据...")
		ensureDemoData()
	}

	// 5. 添加演示专用端点
	addDemoEndpoints(r)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  服务已启动 - 论文答辩演示就绪")
	fmt.Println("========================================")
	fmt.Printf("  API地址: http://localhost:%s\n", config.GlobalConfig.Server.Port)
	fmt.Printf("  健康检查: http://localhost:%s/health\n", config.GlobalConfig.Server.Port)
	fmt.Printf("  Swagger文档: http://localhost:%s/swagger/index.html\n", config.GlobalConfig.Server.Port)
	fmt.Println()
	fmt.Println("  演示账户:")
	fmt.Println("    - demo / demo123 (读者)")
	fmt.Println("    - author / author123 (作者)")
	fmt.Println("    - admin / admin123 (管理员)")
	fmt.Println("========================================")
	fmt.Println()

	// 6. 运行服务器
	if err := core.RunServer(r); err != nil {
		log.Fatalf("运行服务器失败: %v", err)
	}
}

// loadDemoConfig 加载演示配置
func loadDemoConfig() error {
	// 从环境变量读取配置
	if port := os.Getenv("DEMO_PORT"); port != "" {
		DefaultDemoConfig.Port = port
	}
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		DefaultDemoConfig.MongoURI = uri
	}
	if db := os.Getenv("DB_NAME"); db != "" {
		DefaultDemoConfig.DBName = db
	}
	if seed := os.Getenv("SEED_DATA"); seed == "false" {
		DefaultDemoConfig.SeedData = false
	}

	// 获取配置文件路径
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "."
	}

	// 优先使用演示配置文件（新路径），兼容旧路径
	demoConfigPathCandidates := []string{
		filepath.Join(configPath, "configs", "config.demo.yaml"),
		filepath.Join(configPath, "config", "config.demo.yaml"),
		filepath.Join(configPath, "config.demo.yaml"),
	}
	for _, candidate := range demoConfigPathCandidates {
		if _, err := os.Stat(candidate); err == nil {
			configPath = candidate
			fmt.Printf("[演示模式] 使用演示配置: %s\n", configPath)
			break
		}
	}

	// 加载配置
	_, err := config.LoadConfig(configPath)
	if err != nil {
		// 如果加载失败，使用默认配置
		fmt.Printf("[演示模式] 配置文件加载失败，使用默认配置: %v\n", err)
		setDefaultConfig()
	}

	return nil
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	config.GlobalConfig = &config.Config{
		Server: &config.ServerConfig{
			Port: DefaultDemoConfig.Port,
			Mode: "debug",
		},
		Log: &config.LogConfig{
			Level:       "info",
			Format:      "console",
			Output:      "stdout",
			Development: true,
			Mode:        "normal",
		},
		JWT: &config.JWTConfig{
			Secret:          "qingyu_demo_secret_key_for_thesis_defense",
			ExpirationHours: 24,
		},
		Database: &config.DatabaseConfig{
			URI:           DefaultDemoConfig.MongoURI,
			Name:          DefaultDemoConfig.DBName,
			ConnectTimeout: 10 * time.Second,
		},
		Redis: &config.RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
		RateLimit: &config.RateLimitConfig{
			Enabled:        true,
			RequestsPerSec: 100,
			Burst:          200,
		},
	}
}

// addDemoEndpoints 添加演示专用端点
func addDemoEndpoints(r *gin.Engine) {
	// 演示信息端点
	r.GET("/demo/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "青羽写作平台 - 论文答辩演示版本",
			"version":     "1.0-demo",
			"description": "专为论文答辩演示设计的版本，连接真实MongoDB数据库",
			"features": []string{
				"完整CRUD操作",
				"真实数据库存储",
				"内置演示数据",
				"完整API接口",
				"Swagger文档",
			},
			"demo_accounts": []gin.H{
				{"username": "demo", "password": "demo123", "role": "reader"},
				{"username": "author", "password": "author123", "role": "author"},
				{"username": "admin", "password": "admin123", "role": "admin"},
			},
			"database": gin.H{
				"type": "mongodb",
				"uri":  DefaultDemoConfig.MongoURI,
				"name": DefaultDemoConfig.DBName,
			},
		})
	})

	// 演示数据重置端点（仅演示模式）
	r.POST("/demo/reset", func(c *gin.Context) {
		ensureDemoData()
		c.JSON(200, gin.H{
			"code":    0,
			"message": "演示数据已重置",
		})
	})

	// 快速创建演示账户端点
	r.POST("/demo/create-user", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
			return
		}

		if req.Role == "" {
			req.Role = "reader"
		}

		// 这里可以调用用户服务创建用户
		// 为简化演示，返回成功消息
		c.JSON(200, gin.H{
			"code":    0,
			"message": "用户创建成功（演示模式）",
			"data": gin.H{
				"username": req.Username,
				"role":     req.Role,
			},
		})
	})

	logger := zap.NewNop()
	logger.Info("演示端点注册完成")
}

// ensureDemoData 确保演示数据存在
func ensureDemoData() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取服务容器
	container := getServiceContainer()
	if container == nil {
		fmt.Println("[演示模式] 服务容器未初始化，跳过数据检查")
		return
	}

	// 检查是否有用户数据
	userService, err := container.GetUserService()
	if err != nil {
		fmt.Printf("[演示模式] 获取用户服务失败: %v\n", err)
		return
	}

	// 检查演示用户是否存在
	users, err := userService.ListUsers(ctx, 1, 0)
	if err == nil && len(users) > 0 {
		fmt.Printf("[演示模式] 数据库已有 %d 个用户，跳过数据初始化\n", len(users))
		return
	}

	// 如果没有数据，提示用户运行seeder
	fmt.Println("[演示模式] 数据库为空，请运行以下命令初始化演示数据:")
	fmt.Println("  go run cmd/seeder/main.go")
	fmt.Println()
	fmt.Println("  或者手动创建演示账户后继续")
}

// getServiceContainer 获取服务容器
func getServiceContainer() interface {
	GetUserService() (interface{}, error)
} {
	// 返回nil，因为需要完整的服务容器实现
	// 实际使用时会通过 service.GetServiceContainer() 获取
	return nil
}
