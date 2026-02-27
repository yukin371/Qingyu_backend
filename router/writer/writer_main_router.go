package writer

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/interfaces"
)

// InitWriterRoutes 初始化所有写作相关路由
func InitWriterRoutes(
	router *gin.RouterGroup,
	characterService interfaces.CharacterService,
	locationService interfaces.LocationService,
	timelineService interfaces.TimelineService,
	outlineService interfaces.OutlineService,
) {
	// 角色管理路由
	if characterService != nil {
		InitCharacterRoutes(router, characterService)
	}

	// 地点管理路由
	if locationService != nil {
		InitLocationRoutes(router, locationService)
	}

	// 时间线管理路由
	if timelineService != nil {
		InitTimelineRoutes(router, timelineService)
	}

	// 大纲管理路由
	if outlineService != nil {
		InitOutlineRoutes(router, outlineService)
	}
}
