package interfaces

import (
	"context"

	bookstoreModel "Qingyu_backend/models/bookstore"
	writerModel "Qingyu_backend/models/writer"
)

// ConceptCatalogPort 是写作百科查询场景所需的最小概念目录接口。
// 具体仓储由组合根注入，HTTP handler 不依赖完整 ConceptRepository 能力集。
type ConceptCatalogPort interface {
	GetByID(ctx context.Context, id string) (*writerModel.Concept, error)
	ListByProject(ctx context.Context, projectID string) ([]*writerModel.Concept, error)
	Search(ctx context.Context, projectID, category, keyword string) ([]*writerModel.Concept, error)
}

// WriterStatsBookPort 是作者统计查询场景所需的最小书籍目录接口。
// 统计 API 只解析作者作品、项目映射和作品基础信息，不应持有完整书籍仓储的写入能力。
type WriterStatsBookPort interface {
	GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstoreModel.Book, error)
	CountByAuthorID(ctx context.Context, authorID string) (int64, error)
	GetByID(ctx context.Context, id string) (*bookstoreModel.Book, error)
	GetByProjectID(ctx context.Context, projectID string) (*bookstoreModel.Book, error)
}

// DocumentSceneStatePort 是 AI 场景状态更新所需的最小文档写入接口。
type DocumentSceneStatePort interface {
	Update(ctx context.Context, id string, updates map[string]interface{}) error
}
