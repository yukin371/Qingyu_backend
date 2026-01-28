package auth

import (
	"go.mongodb.org/mongo-driver/mongo"

	authInterface "Qingyu_backend/repository/interfaces/auth"
)

// AuthRepository 认证相关Repository (别名，指向接口)
// Deprecated: 使用 RoleRepository 替代
type AuthRepository = authInterface.RoleRepository

// NewAuthRepository 创建AuthRepository (兼容函数)
// Deprecated: 使用 NewRoleRepository 替代
func NewAuthRepository(db *mongo.Database, logger interface{}) AuthRepository {
	return NewRoleRepository(db)
}
