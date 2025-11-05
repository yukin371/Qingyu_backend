package main

import (
	"log"
	"os"
	"path/filepath"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	_ "Qingyu_backend/docs" // Import generated docs
)

// @title           青羽写作平台 API
// @version         1.0
// @description     青羽写作平台后端服务API文档，提供AI辅助写作、阅读社区、书城管理等核心功能。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name 书城
// @tag.description 书城相关接口，包括首页、书籍列表、分类等
// @tag.name 书籍
// @tag.description 书籍详情、搜索、评分等功能
// @tag.name 用户
// @tag.description 用户注册、登录、个人信息管理
// @tag.name 项目
// @tag.description 写作项目管理
// @tag.name 文档
// @tag.description 文档编辑、版本控制
// @tag.name AI辅助
// @tag.description AI写作辅助功能
// @tag.name 钱包
// @tag.description 钱包、充值、提现功能

func main() {
	// 获取当前工作目录（用于调试）
	wd, _ := os.Getwd()

	// 获取配置文件路径
	// Docker环境使用绝对路径，本地开发使用环境变量或当前目录
	log.Printf("[Main] Working directory: %s", wd)
	log.Printf("[Main] CONFIG_FILE env: %s", os.Getenv("CONFIG_FILE"))
	log.Printf("[Main] CONFIG_PATH env: %s", os.Getenv("CONFIG_PATH"))

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if configPath == "" {
		// 默认尝试 /app/config/config.yaml（Docker环境）
		if _, err := os.Stat("/app/config/config.yaml"); err == nil {
			configPath = "/app/config/config.yaml"
			log.Printf("[Main] Found Docker config at: %s", configPath)
		} else {
			// 否则使用当前目录（本地环境）
			configPath = "."
			log.Printf("[Main] Using current directory for config: %s", configPath)
		}
	}
	configPath = filepath.Clean(configPath)

	log.Printf("[Main] Final config path: %s", configPath)

	// 加载配置
	_, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("[Main] Configuration loaded successfully")

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
