package provider

import (
	"context"
	"fmt"

	"Qingyu_backend/models/search"
)

// DocumentProvider 文档搜索提供者
type DocumentProvider struct {
	// TODO: 注入 Engine
	typeVal search.SearchType
}

// NewDocumentProvider 创建文档搜索提供者
func NewDocumentProvider() (*DocumentProvider, error) {
	return &DocumentProvider{
		typeVal: search.SearchTypeDocuments,
	}, nil
}

// Search 搜索文档
func (p *DocumentProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// TODO: 实现文档搜索逻辑
	// 1. 验证用户认证
	// 2. 强制过滤：只能搜索自己项目的文档
	// 3. 支持按项目 ID 过滤
	// 4. 构建查询：标题、内容搜索
	// 5. 执行搜索
	// 6. 返回结果
	return nil, fmt.Errorf("document provider search not implemented")
}

// Type 获取搜索类型
func (p *DocumentProvider) Type() search.SearchType {
	return p.typeVal
}

// Validate 验证搜索参数
func (p *DocumentProvider) Validate(req *search.SearchRequest) error {
	// TODO: 实现参数验证
	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *DocumentProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	// TODO: 实现获取逻辑
	return nil, fmt.Errorf("get by id not implemented")
}

// GetBatch 批量获取文档
func (p *DocumentProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	// TODO: 实现批量获取逻辑
	return nil, fmt.Errorf("get batch not implemented")
}
