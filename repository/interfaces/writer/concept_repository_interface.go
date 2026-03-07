package writer

import (
	"context"

	"Qingyu_backend/models/writer"
)

// ConceptRepository 设定百科Repository接口
//
// 提供Concept模型的数据库操作接口，支持AI写作助手功能中的设定百科管理喵~
//
// 功能包括：
//   - 基础CRUD操作
//   - 按项目查询设定
//   - 按分类和关键词搜索
//   - 批量查询操作
type ConceptRepository interface {
	// Create 创建设定
	//
	// 参数：
	//   - ctx: 上下文
	//   - concept: 设定对象指针
	//
	// 返回：
	//   - error: 创建错误
	Create(ctx context.Context, concept *writer.Concept) error

	// GetByID 根据ID获取设定
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 设定ID
	//
	// 返回：
	//   - *writer.Concept: 设定对象
	//   - error: 查询错误
	GetByID(ctx context.Context, id string) (*writer.Concept, error)

	// ListByProject 获取项目的设定列表
	//
	// 参数：
	//   - ctx: 上下文
	//   - projectID: 项目ID
	//
	// 返回：
	//   - []*writer.Concept: 设定列表
	//   - error: 查询错误
	ListByProject(ctx context.Context, projectID string) ([]*writer.Concept, error)

	// Search 搜索设定
	//
	// 支持按分类和关键词搜索，可选参数喵~
	//
	// 参数：
	//   - ctx: 上下文
	//   - projectID: 项目ID
	//   - category: 分类（可选）
	//   - keyword: 关键词（可选）
	//
	// 返回：
	//   - []*writer.Concept: 设定列表
	//   - error: 查询错误
	Search(ctx context.Context, projectID, category, keyword string) ([]*writer.Concept, error)

	// Update 更新设定
	//
	// 参数：
	//   - ctx: 上下文
	//   - concept: 设定对象指针
	//
	// 返回：
	//   - error: 更新错误
	Update(ctx context.Context, concept *writer.Concept) error

	// Delete 删除设定
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 设定ID
	//
	// 返回：
	//   - error: 删除错误
	Delete(ctx context.Context, id string) error

	// BatchGetByIDs 批量获取设定
	//
	// 参数：
	//   - ctx: 上下文
	//   - ids: 设定ID列表
	//
	// 返回：
	//   - []*writer.Concept: 设定列表
	//   - error: 查询错误
	BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.Concept, error)
}
