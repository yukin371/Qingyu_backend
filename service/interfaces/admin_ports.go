package interfaces

import (
	"context"

	bookstoreModel "Qingyu_backend/models/bookstore"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
)

// AdminContentBookPort 是管理端内容导出读取书籍所需的最小端口。
// HTTP handler 不应依赖书籍仓储的写入、事务和批量管理能力。
type AdminContentBookPort interface {
	List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.Book, error)
}

// AdminContentChapterPort 是管理端内容导出读取章节所需的最小端口。
type AdminContentChapterPort interface {
	List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.Chapter, error)
}
