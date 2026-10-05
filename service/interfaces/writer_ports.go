package interfaces

import (
	"context"

	writerModel "Qingyu_backend/models/writer"
)

// ConceptCatalogPort 是写作百科查询场景所需的最小概念目录接口。
// 具体仓储由组合根注入，HTTP handler 不依赖完整 ConceptRepository 能力集。
type ConceptCatalogPort interface {
	GetByID(ctx context.Context, id string) (*writerModel.Concept, error)
	ListByProject(ctx context.Context, projectID string) ([]*writerModel.Concept, error)
	Search(ctx context.Context, projectID, category, keyword string) ([]*writerModel.Concept, error)
}
