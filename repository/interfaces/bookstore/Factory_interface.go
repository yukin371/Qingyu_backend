package bookstore

import (
	"context"
)

// BookstoreRepositoryFactory 书城仓储工厂接口
// 负责创建和管理所有书城相关的Repository实例
type BookstoreRepositoryFactory interface {
	// 书籍相关Repository
	GetBookRepository() BookRepository

	// 分类相关Repository
	GetCategoryRepository() CategoryRepository

	// Banner相关Repository
	GetBannerRepository() BannerRepository

	// 榜单相关Repository
	GetRankingRepository() RankingRepository

	// Health 基础设施方法
	Health(ctx context.Context) error
	Close() error
}
