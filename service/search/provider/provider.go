package provider

import (
	"context"

	"Qingyu_backend/models/search"
)

// Provider 业务搜索提供者接口
type Provider interface {
	// Search 执行搜索
	Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error)

	// Type 获取搜索类型
	Type() search.SearchType

	// Validate 验证搜索参数
	Validate(req *search.SearchRequest) error

	// GetByID 根据 ID 获取单个文档
	GetByID(ctx context.Context, id string) (*search.SearchItem, error)

	// GetBatch 批量获取文档
	GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error)
}
