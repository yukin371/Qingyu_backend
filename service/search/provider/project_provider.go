package provider

import (
	"context"
	"fmt"

	"Qingyu_backend/models/search"
)

// ProjectProvider 创作项目搜索提供者
type ProjectProvider struct {
	// TODO: 注入 Engine
	typeVal search.SearchType
}

// NewProjectProvider 创建项目搜索提供者
func NewProjectProvider() (*ProjectProvider, error) {
	return &ProjectProvider{
		typeVal: search.SearchTypeProjects,
	}, nil
}

// Search 搜索项目
func (p *ProjectProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// TODO: 实现项目搜索逻辑
	// 1. 验证用户认证
	// 2. 强制过滤：只能搜索作者自己的项目
	// 3. 构建查询：标题、简介、类型等搜索
	// 4. 执行搜索
	// 5. 返回结果
	return nil, fmt.Errorf("project provider search not implemented")
}

// Type 获取搜索类型
func (p *ProjectProvider) Type() search.SearchType {
	return p.typeVal
}

// Validate 验证搜索参数
func (p *ProjectProvider) Validate(req *search.SearchRequest) error {
	// TODO: 实现参数验证
	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *ProjectProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	// TODO: 实现获取逻辑
	return nil, fmt.Errorf("get by id not implemented")
}

// GetBatch 批量获取文档
func (p *ProjectProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	// TODO: 实现批量获取逻辑
	return nil, fmt.Errorf("get batch not implemented")
}
