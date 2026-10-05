package interfaces

import (
	"context"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
)

// AdminContentTimeRange 是管理端内容查询使用的时间范围。
// 该类型属于 service 端口，避免把 repository 的过滤器模型泄漏到 handler。
type AdminContentTimeRange struct {
	Start *time.Time
	End   *time.Time
}

// AdminContentFilter 是管理端内容导出的领域过滤条件。
// 具体持久化字段映射由组合根适配器负责。
type AdminContentFilter struct {
	Status      string
	Author      string
	BookID      string
	CreatedAt   *AdminContentTimeRange
	PublishedAt *AdminContentTimeRange
}

// AdminContentBookPort 是管理端内容导出读取书籍所需的最小端口。
// HTTP handler 不应依赖书籍仓储的写入、事务和批量管理能力。
type AdminContentBookPort interface {
	List(ctx context.Context, filter AdminContentFilter) ([]*bookstoreModel.Book, error)
}

// AdminContentChapterPort 是管理端内容导出读取章节所需的最小端口。
type AdminContentChapterPort interface {
	List(ctx context.Context, filter AdminContentFilter) ([]*bookstoreModel.Chapter, error)
}
