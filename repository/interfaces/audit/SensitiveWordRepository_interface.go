package audit

import (
	"context"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/repository/interfaces/infrastructure"
)

// SensitiveWordRepository 敏感词Repository接口
type SensitiveWordRepository interface {
	// 基础CRUD
	Create(ctx context.Context, word *audit.SensitiveWord) error
	GetByID(ctx context.Context, id string) (*audit.SensitiveWord, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询方法
	GetByWord(ctx context.Context, word string) (*audit.SensitiveWord, error)
	List(ctx context.Context, filter infrastructure.Filter) ([]*audit.SensitiveWord, error)
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[audit.SensitiveWord], error)

	// 业务方法
	GetEnabledWords(ctx context.Context) ([]*audit.SensitiveWord, error)
	GetByCategory(ctx context.Context, category string) ([]*audit.SensitiveWord, error)
	GetByLevel(ctx context.Context, minLevel int) ([]*audit.SensitiveWord, error)
	BatchCreate(ctx context.Context, words []*audit.SensitiveWord) error
	BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error
	BatchDelete(ctx context.Context, ids []string) error

	// 统计方法
	CountByCategory(ctx context.Context) (map[string]int64, error)
	CountByLevel(ctx context.Context) (map[int]int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
