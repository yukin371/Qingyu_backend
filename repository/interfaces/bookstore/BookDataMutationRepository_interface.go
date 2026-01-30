package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
)

// BookDataMutationRepository 书籍数据变更接口 - 专注于书籍的创建、更新和删除操作
// 用于后台管理系统、作者工作台等场景
//
// 职责：
// - 创建新书籍
// - 更新书籍信息
// - 删除书籍
// - 事务支持
//
// 方法总数：4个
type BookDataMutationRepository interface {
	// Create 创建书籍
	Create(ctx context.Context, book *bookstore2.Book) error

	// Update 更新书籍
	// updates 是一个map，包含要更新的字段和新值
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除书籍
	Delete(ctx context.Context, id string) error

	// Transaction 执行事务
	// 用于需要原子性执行多个操作的场景
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
