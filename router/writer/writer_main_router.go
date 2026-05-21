package writer

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/service/interfaces"
)

// InitWriterRoutes 初始化所有写作相关路由
func InitWriterRoutes(
	router *gin.RouterGroup,
	characterService interfaces.CharacterService,
	locationService interfaces.LocationService,
	timelineService interfaces.TimelineService,
	outlineService interfaces.OutlineService,
	entityService interfaces.EntityService,
) {
	zap.L().Info("InitWriterRoutes: 开始注册设定百科路由")

	// 创建 /writer 子组，与 InitWriterRouter 保持一致
	writerGroup := router.Group("/writer")
	writerGroup.Use(auth.JWTAuth())
	{
		// 角色管理路由
		if characterService != nil {
			InitCharacterRoutes(writerGroup, characterService)
			zap.L().Info("InitWriterRoutes: 角色路由注册完成")
		} else {
			zap.L().Warn("InitWriterRoutes: CharacterService为nil，跳过角色路由注册")
		}

		// 地点管理路由
		if locationService != nil {
			InitLocationRoutes(writerGroup, locationService)
			zap.L().Info("InitWriterRoutes: 地点路由注册完成")
		}

		// 时间线管理路由
		if timelineService != nil {
			InitTimelineRoutes(writerGroup, timelineService)
			zap.L().Info("InitWriterRoutes: 时间线路由注册完成")
		}

		// 大纲管理路由
		if outlineService != nil {
			InitOutlineRoutes(writerGroup, outlineService)
			zap.L().Info("InitWriterRoutes: 大纲路由注册完成")
		}

		// 统一实体路由
		if entityService != nil {
			InitEntityRoutes(writerGroup, entityService)
			zap.L().Info("InitWriterRoutes: 实体路由注册完成")
		}
	}

	zap.L().Info("InitWriterRoutes: 设定百科路由注册完成")
}
