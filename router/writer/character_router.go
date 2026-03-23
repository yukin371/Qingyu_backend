package writer

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitCharacterRoutes 初始化角色路由
func InitCharacterRoutes(router *gin.RouterGroup, characterService interfaces.CharacterService) {
	zap.L().Info("InitCharacterRoutes: 开始注册角色路由")

	api := writer.NewCharacterApi(characterService)

	// 项目级别的角色路由（需要认证）
	projectGroup := router.Group("/projects/:id")
	{
		// 创建和列表查询
		projectGroup.POST("/characters", api.CreateCharacter)
		projectGroup.GET("/characters", api.ListCharacters)

		// 关系相关
		projectGroup.GET("/characters/relations", api.ListCharacterRelations)
		projectGroup.GET("/characters/graph", api.GetCharacterGraph)

		zap.L().Info("InitCharacterRoutes: 项目级角色路由已注册到 /projects/:id/characters")
	}

	// 角色级别的路由（需要认证）
	characterGroup := router.Group("/characters")
	{
		// 单个角色操作（需要传递projectId作为查询参数进行权限验证）
		characterGroup.GET("/:characterId", api.GetCharacter)
		characterGroup.PUT("/:characterId", api.UpdateCharacter)
		characterGroup.DELETE("/:characterId", api.DeleteCharacter)

		// 关系管理
		characterGroup.POST("/relations", api.CreateCharacterRelation)
		characterGroup.DELETE("/relations/:relationId", api.DeleteCharacterRelation)

		zap.L().Info("InitCharacterRoutes: 角色级路由已注册到 /characters")
	}

	zap.L().Info("InitCharacterRoutes: 角色路由注册完成")
}
