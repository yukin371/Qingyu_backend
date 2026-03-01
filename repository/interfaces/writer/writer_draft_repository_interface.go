package writer

import (
	"context"

	"Qingyu_backend/models/writer"
)

// WriterDraftRepository AI写作草稿Repository接口
//
// 提供WriterDraft模型的数据库操作接口，支持AI写作助手功能中的文档草稿管理喵~
//
// 功能包括：
//   - 基础CRUD操作
//   - 按项目和章节查询
//   - 批量查询操作
type WriterDraftRepository interface {
	// Create 创建草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - doc: 草稿对象指针
	//
	// 返回：
	//   - error: 创建错误
	Create(ctx context.Context, doc *writer.WriterDraft) error

	// GetByID 根据ID获取草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 草稿ID
	//
	// 返回：
	//   - *writer.WriterDraft: 草稿对象
	//   - error: 查询错误
	GetByID(ctx context.Context, id string) (*writer.WriterDraft, error)

	// GetByProjectAndChapter 根据项目ID和章节号获取草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - projectID: 项目ID
	//   - chapterNum: 章节号
	//
	// 返回：
	//   - *writer.WriterDraft: 草稿对象
	//   - error: 查询错误
	GetByProjectAndChapter(ctx context.Context, projectID string, chapterNum int) (*writer.WriterDraft, error)

	// ListByProject 获取项目的草稿列表
	//
	// 参数：
	//   - ctx: 上下文
	//   - projectID: 项目ID
	//   - limit: 返回数量限制
	//
	// 返回：
	//   - []*writer.WriterDraft: 草稿列表
	//   - error: 查询错误
	ListByProject(ctx context.Context, projectID string, limit int) ([]*writer.WriterDraft, error)

	// Update 更新草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - doc: 草稿对象指针
	//
	// 返回：
	//   - error: 更新错误
	Update(ctx context.Context, doc *writer.WriterDraft) error

	// Delete 删除草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - id: 草稿ID
	//
	// 返回：
	//   - error: 删除错误
	Delete(ctx context.Context, id string) error

	// BatchGetByIDs 批量获取草稿
	//
	// 参数：
	//   - ctx: 上下文
	//   - ids: 草稿ID列表
	//
	// 返回：
	//   - []*writer.WriterDraft: 草稿列表
	//   - error: 查询错误
	BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.WriterDraft, error)
}
