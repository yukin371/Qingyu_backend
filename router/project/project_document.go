package document

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文档相关路由
// TODO: 重构说明 - api/v1/document/ 已合并到 api/v1/writer/
// 待服务层实现完成后，使用 router/writer/writer.go 中的完整实现
// 当前为空实现，避免路由冲突
func RegisterRoutes(r *gin.RouterGroup) {
	// 临时禁用，避免与其他路由冲突
	// 完整的写作端路由请参考 router/writer/writer.go
}
