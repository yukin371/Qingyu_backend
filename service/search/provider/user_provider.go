package provider

import (
	"context"
	"fmt"

	"Qingyu_backend/models/search"
)

// UserProvider 用户搜索提供者
type UserProvider struct {
	// TODO: 注入 Engine
	typeVal search.SearchType
}

// NewUserProvider 创建用户搜索提供者
func NewUserProvider() (*UserProvider, error) {
	return &UserProvider{
		typeVal: search.SearchTypeUsers,
	}, nil
}

// Search 搜索用户
func (p *UserProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// TODO: 实现用户搜索逻辑
	// 1. 支持用户名、昵称、简介搜索
	// 2. 支持角色、认证状态过滤
	// 3. 查询权重：username(10.0), nickname(8.0), bio(1.0)
	// 4. 执行搜索
	// 5. 返回结果
	return nil, fmt.Errorf("user provider search not implemented")
}

// Type 获取搜索类型
func (p *UserProvider) Type() search.SearchType {
	return p.typeVal
}

// Validate 验证搜索参数
func (p *UserProvider) Validate(req *search.SearchRequest) error {
	// TODO: 实现参数验证
	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *UserProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	// TODO: 实现获取逻辑
	return nil, fmt.Errorf("get by id not implemented")
}

// GetBatch 批量获取文档
func (p *UserProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	// TODO: 实现批量获取逻辑
	return nil, fmt.Errorf("get batch not implemented")
}
