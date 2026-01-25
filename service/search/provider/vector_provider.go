package provider

import (
	"context"
	"fmt"

	"Qingyu_backend/models/search"
)

// VectorProvider 向量搜索提供者
type VectorProvider struct {
	// TODO: 注入 Milvus Engine
	typeVal search.SearchType
}

// NewVectorProvider 创建向量搜索提供者
func NewVectorProvider() (*VectorProvider, error) {
	return &VectorProvider{
		typeVal: search.SearchTypeVector,
	}, nil
}

// Search 执行向量搜索
func (p *VectorProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// TODO: 实现向量搜索逻辑
	// 1. 验证用户认证
	// 2. 调用 Milvus 进行语义搜索
	// 3. 支持混合检索（后续迭代）
	// 4. 返回结果
	return nil, fmt.Errorf("vector provider search not implemented")
}

// Type 获取搜索类型
func (p *VectorProvider) Type() search.SearchType {
	return p.typeVal
}

// Validate 验证搜索参数
func (p *VectorProvider) Validate(req *search.SearchRequest) error {
	// TODO: 实现参数验证
	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *VectorProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	// TODO: 实现获取逻辑
	return nil, fmt.Errorf("get by id not implemented")
}

// GetBatch 批量获取文档
func (p *VectorProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	// TODO: 实现批量获取逻辑
	return nil, fmt.Errorf("get batch not implemented")
}
