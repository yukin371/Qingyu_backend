package main

import (
	"log"

	_ "Qingyu_backend/docs" // Import generated docs
	serverapp "Qingyu_backend/internal/app/server"
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

// @host      localhost:9090
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
	if err := serverapp.Run(); err != nil {
		log.Fatal(err)
	}
}
