package adapters

import (
	"context"

	"Qingyu_backend/models/system"
	"Qingyu_backend/repository"
	"Qingyu_backend/repository/interfaces"
)

// UserRepositoryAdapter 用户Repository适配器
// 将新的UserRepository接口适配到旧的接口，确保API兼容性
type UserRepositoryAdapter struct {
	newRepo interfaces.UserRepository
}

// NewUserRepositoryAdapter 创建用户Repository适配器
func NewUserRepositoryAdapter(newRepo interfaces.UserRepository) repository.UserRepository {
	return &UserRepositoryAdapter{
		newRepo: newRepo,
	}
}

// Create 创建用户
func (a *UserRepositoryAdapter) Create(ctx context.Context, user *system.User) error {
	return a.newRepo.Create(ctx, &user)
}

// GetByID 根据ID获取用户
func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id string) (*system.User, error) {
	// 使用新的接口，通过ID过滤
	filter := interfaces.UserFilter{ID: id}
	user, err := a.newRepo.GetByID(ctx, filter)
	if err != nil {
		return nil, err
	}
	return *user, nil
}

// GetByUsername 根据用户名获取用户
func (a *UserRepositoryAdapter) GetByUsername(ctx context.Context, username string) (*system.User, error) {
	return a.newRepo.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户
func (a *UserRepositoryAdapter) GetByEmail(ctx context.Context, email string) (*system.User, error) {
	return a.newRepo.GetByEmail(ctx, email)
}

// Update 更新用户信息
func (a *UserRepositoryAdapter) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	// 使用新的接口，通过ID过滤
	filter := interfaces.UserFilter{ID: id}
	return a.newRepo.Update(ctx, filter, updates)
}

// Delete 删除用户（软删除）
func (a *UserRepositoryAdapter) Delete(ctx context.Context, id string) error {
	// 使用新的接口，通过ID过滤
	filter := interfaces.UserFilter{ID: id}
	return a.newRepo.Delete(ctx, filter)
}

// HardDelete 硬删除用户
func (a *UserRepositoryAdapter) HardDelete(ctx context.Context, id string) error {
	// 由于新的接口没有HardDelete方法，我们使用Delete方法
	// TODO：这是一个临时解决方案，直到接口设计被完善
	filter := interfaces.UserFilter{ID: id}
	return a.newRepo.Delete(ctx, filter)
}

// List 获取用户列表
func (a *UserRepositoryAdapter) List(ctx context.Context, filter repository.UserFilter) ([]*system.User, error) {
	// 转换过滤器格式
	newFilter := interfaces.UserFilter{
		Username: filter.Username,
		Email:    filter.Email,
		Status:   filter.Status,
		FromDate: filter.FromDate,
		ToDate:   filter.ToDate,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}

	// 使用新接口的List方法
	users, err := a.newRepo.List(ctx, newFilter)
	if err != nil {
		return nil, err
	}

	// 转换结果格式
	var result []*system.User
	for _, user := range users {
		if user != nil && *user != nil {
			result = append(result, *user)
		}
	}

	return result, nil
}

// Count 统计用户数量
func (a *UserRepositoryAdapter) Count(ctx context.Context, filter repository.UserFilter) (int64, error) {
	// 转换过滤器格式
	newFilter := interfaces.UserFilter{
		Username: filter.Username,
		Email:    filter.Email,
		Status:   filter.Status,
		FromDate: filter.FromDate,
		ToDate:   filter.ToDate,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}

	return a.newRepo.Count(ctx, newFilter)
}

// Exists 检查用户是否存在
func (a *UserRepositoryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	// 使用新的接口，通过ID过滤
	filter := interfaces.UserFilter{ID: id}
	return a.newRepo.Exists(ctx, filter)
}

// ExistsByUsername 检查用户名是否存在
func (a *UserRepositoryAdapter) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return a.newRepo.ExistsByUsername(ctx, username)
}

// ExistsByEmail 检查邮箱是否存在
func (a *UserRepositoryAdapter) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return a.newRepo.ExistsByEmail(ctx, email)
}

// UpdateLastLogin 更新最后登录时间
func (a *UserRepositoryAdapter) UpdateLastLogin(ctx context.Context, id string) error {
	return a.newRepo.UpdateLastLogin(ctx, id)
}

// UpdatePassword 更新密码
func (a *UserRepositoryAdapter) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	return a.newRepo.UpdatePassword(ctx, id, hashedPassword)
}

// GetActiveUsers 获取活跃用户
func (a *UserRepositoryAdapter) GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error) {
	return a.newRepo.GetActiveUsers(ctx, limit)
}

// BatchUpdate 批量更新用户
func (a *UserRepositoryAdapter) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	// 由于新的接口BatchUpdate参数类型不匹配，我们暂时返回错误
	// TODO：这是一个临时解决方案，直到接口设计被完善
	return interfaces.NewRepositoryError(
		interfaces.ErrorTypeInternal,
		"BatchUpdate方法需要重新设计以支持新的接口",
		nil,
	)
}

// Transaction 执行事务操作
func (a *UserRepositoryAdapter) Transaction(ctx context.Context, fn func(ctx context.Context, repo repository.UserRepository) error) error {
	return a.newRepo.Transaction(ctx, func(ctx context.Context, repo interfaces.UserRepository) error {
		// 创建适配器包装新的repo，传递给旧的回调函数
		adapter := NewUserRepositoryAdapter(repo)
		return fn(ctx, adapter)
	})
}

// ProjectRepositoryAdapter 项目Repository适配器
type ProjectRepositoryAdapter struct {
	newRepo interfaces.ProjectRepository
}

// NewProjectRepositoryAdapter 创建项目Repository适配器
func NewProjectRepositoryAdapter(newRepo interfaces.ProjectRepository) repository.ProjectRepository {
	return &ProjectRepositoryAdapter{
		newRepo: newRepo,
	}
}

// Create 创建项目
func (a *ProjectRepositoryAdapter) Create(ctx context.Context, project interface{}) error {
	return a.newRepo.Create(ctx, &project)
}

// GetByID 根据ID获取项目
func (a *ProjectRepositoryAdapter) GetByID(ctx context.Context, id string) (interface{}, error) {
	return a.newRepo.GetByID(ctx, id)
}

// GetByCreatorID 根据创建者ID获取项目列表
func (a *ProjectRepositoryAdapter) GetByCreatorID(ctx context.Context, creatorID string) ([]interface{}, error) {
	return a.newRepo.GetByCreatorID(ctx, creatorID)
}

// Update 更新项目
func (a *ProjectRepositoryAdapter) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return a.newRepo.Update(ctx, id, updates)
}

// Delete 删除项目
func (a *ProjectRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.newRepo.Delete(ctx, id)
}

// List 获取项目列表
func (a *ProjectRepositoryAdapter) List(ctx context.Context, filter interface{}) ([]interface{}, error) {
	// 由于新的接口设计问题，我们需要使用其他方法
	// TODO：这是一个临时解决方案，直到接口设计被修复
	return nil, interfaces.NewRepositoryError(
		interfaces.ErrorTypeInternal,
		"Project List方法需要重新设计",
		nil,
	)
}

// Count 统计项目数量
func (a *ProjectRepositoryAdapter) Count(ctx context.Context, filter interface{}) (int64, error) {
	// 由于新的接口设计问题，我们需要使用其他方法
	// TODO：这是一个临时解决方案，直到接口设计被修复
	return 0, interfaces.NewRepositoryError(
		interfaces.ErrorTypeInternal,
		"Project Count方法需要重新设计",
		nil,
	)
}

// ArchiveByCreatorID 根据创建者ID归档项目
func (a *ProjectRepositoryAdapter) ArchiveByCreatorID(ctx context.Context, creatorID string) error {
	return a.newRepo.ArchiveByCreatorID(ctx, creatorID)
}

// 角色相关

// RoleRepositoryAdapter 角色Repository适配器
type RoleRepositoryAdapter struct {
	newRepo interfaces.RoleRepository
}

// NewRoleRepositoryAdapter 创建角色Repository适配器
func NewRoleRepositoryAdapter(newRepo interfaces.RoleRepository) repository.RoleRepository {
	return &RoleRepositoryAdapter{
		newRepo: newRepo,
	}
}

// Create 创建角色
func (a *RoleRepositoryAdapter) Create(ctx context.Context, role interface{}) error {
	return a.newRepo.Create(ctx, &role)
}

// GetByID 根据ID获取角色
func (a *RoleRepositoryAdapter) GetByID(ctx context.Context, id string) (interface{}, error) {
	return a.newRepo.GetByID(ctx, id)
}

// GetByName 根据名称获取角色
func (a *RoleRepositoryAdapter) GetByName(ctx context.Context, name string) (interface{}, error) {
	return a.newRepo.GetByName(ctx, name)
}

// Update 更新角色
func (a *RoleRepositoryAdapter) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return a.newRepo.Update(ctx, id, updates)
}

// Delete 删除角色
func (a *RoleRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.newRepo.Delete(ctx, id)
}

// List 获取角色列表
func (a *RoleRepositoryAdapter) List(ctx context.Context) ([]interface{}, error) {
	// 由于新的接口设计问题，我们需要使用其他方法
	// 这是一个临时解决方案，直到接口设计被修复
	return nil, interfaces.NewRepositoryError(
		interfaces.ErrorTypeInternal,
		"Role List方法需要重新设计",
		nil,
	)
}

// GetDefaultRole 获取默认角色
func (a *RoleRepositoryAdapter) GetDefaultRole(ctx context.Context) (interface{}, error) {
	return a.newRepo.GetDefaultRole(ctx)
}

// GetUserRoles 获取用户角色
func (a *RoleRepositoryAdapter) GetUserRoles(ctx context.Context, userID string) ([]interface{}, error) {
	return a.newRepo.GetUserRoles(ctx, userID)
}

// AssignRole 分配角色
func (a *RoleRepositoryAdapter) AssignRole(ctx context.Context, userID, roleID string) error {
	return a.newRepo.AssignRole(ctx, userID, roleID)
}

// RemoveRole 移除角色
func (a *RoleRepositoryAdapter) RemoveRole(ctx context.Context, userID, roleID string) error {
	return a.newRepo.RemoveRole(ctx, userID, roleID)
}

// GetUserPermissions 获取用户权限
func (a *RoleRepositoryAdapter) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return a.newRepo.GetUserPermissions(ctx, userID)
}

// RepositoryFactoryAdapter 仓储工厂适配器
type RepositoryFactoryAdapter struct {
	newFactory interfaces.RepositoryFactory
}

// NewRepositoryFactoryAdapter 创建仓储工厂适配器
func NewRepositoryFactoryAdapter(newFactory interfaces.RepositoryFactory) repository.RepositoryFactory {
	return &RepositoryFactoryAdapter{
		newFactory: newFactory,
	}
}

// CreateUserRepository 创建用户Repository
func (a *RepositoryFactoryAdapter) CreateUserRepository() repository.UserRepository {
	newRepo := a.newFactory.CreateUserRepository()
	return NewUserRepositoryAdapter(newRepo)
}

// CreateProjectRepository 创建项目Repository
func (a *RepositoryFactoryAdapter) CreateProjectRepository() repository.ProjectRepository {
	newRepo := a.newFactory.CreateProjectRepository()
	return NewProjectRepositoryAdapter(newRepo)
}

// CreateRoleRepository 创建角色Repository
func (a *RepositoryFactoryAdapter) CreateRoleRepository() repository.RoleRepository {
	newRepo := a.newFactory.CreateRoleRepository()
	return NewRoleRepositoryAdapter(newRepo)
}

// 数据库相关
// Close 关闭连接
func (a *RepositoryFactoryAdapter) Close() error {
	return a.newFactory.Close()
}

// Health 健康检查
func (a *RepositoryFactoryAdapter) Health(ctx context.Context) error {
	return a.newFactory.Health(ctx)
}

// GetDatabaseType 获取数据库类型
func (a *RepositoryFactoryAdapter) GetDatabaseType() string {
	return a.newFactory.GetDatabaseType()
}
