package interfaces

import (
	"context"
)

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	// 健康检查
	Health(ctx context.Context) error
	
	// 关闭连接
	Close() error
}