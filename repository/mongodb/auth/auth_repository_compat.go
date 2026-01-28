package auth

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	authModel "Qingyu_backend/models/auth"
	authInterface "Qingyu_backend/repository/interfaces/auth"
)

// AuthRepository 认证相关Repository (别名，指向 RoleRepositoryImpl)
// Deprecated: 使用 RoleRepository 替代
type AuthRepository = RoleRepositoryImpl

// NewAuthRepository 创建AuthRepository (兼容函数)
// Deprecated: 使用 NewRoleRepository 替代
func NewAuthRepository(db *mongo.Database, logger *zap.Logger) *AuthRepository {
	repo := NewRoleRepository(db)
	return (*AuthRepository)(repo)
}

// 确保类型断言有效
var _ authInterface.RoleRepository = (*AuthRepository)(nil)
