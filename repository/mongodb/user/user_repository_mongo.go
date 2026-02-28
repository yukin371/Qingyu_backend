package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/repository/interfaces/infrastructure"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	sharedrepo "Qingyu_backend/repository/mongodb/shared"
)

const UserCollection = "users"

// MongoUserRepository MongoDB用户仓储实现
// 组合使用shared.MongoUserRepository作为基础，保留user特有方法
type MongoUserRepository struct {
	// 组合使用统一的MongoUserRepository
	baseRepo *sharedrepo.MongoUserRepository
	db       *mongo.Database
}

// NewMongoUserRepository 创建新的MongoDB用户仓储实例
func NewMongoUserRepository(db *mongo.Database) UserInterface.UserRepository {
	return &MongoUserRepository{
		baseRepo: sharedrepo.NewMongoUserRepository(db),
		db:       db,
	}
}

// Health 健康检查
func (r *MongoUserRepository) Health(ctx context.Context) error {
	// 检查数据库连接
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeConnection,
			"数据库连接失败",
			err,
		)
	}
	return nil
}

// === 委托共享方法到 baseRepo (使用string ID) ===

// Create 创建用户
func (r *MongoUserRepository) Create(ctx context.Context, user *usersModel.User) error {
	// 验证用户数据
	if err := r.ValidateUser(*user); err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"用户数据验证失败",
			err,
		)
	}
	return r.baseRepo.Create(ctx, user)
}

// GetByID 根据ID获取用户 (string ID)
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*usersModel.User, error) {
	return r.baseRepo.GetByID(ctx, id)
}

// Update 更新用户信息 (string ID)
func (r *MongoUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.baseRepo.Update(ctx, id, updates)
}

// Delete 软删除用户 (string ID)
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	return r.baseRepo.Delete(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*usersModel.User, error) {
	return r.baseRepo.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	return r.baseRepo.GetByEmail(ctx, email)
}

// UpdateStatus 更新用户状态 (string ID)
func (r *MongoUserRepository) UpdateStatus(ctx context.Context, id string, status usersModel.UserStatus) error {
	return r.baseRepo.UpdateStatus(ctx, id, status)
}

// UpdatePassword 更新密码 (string ID)
func (r *MongoUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	return r.baseRepo.UpdatePassword(ctx, id, hashedPassword)
}

// SetEmailVerified 设置邮箱验证状态 (string ID)
func (r *MongoUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	return r.baseRepo.SetEmailVerified(ctx, id, verified)
}

// BatchUpdateStatus 批量更新用户状态 (string IDs)
func (r *MongoUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status usersModel.UserStatus) error {
	return r.baseRepo.BatchUpdateStatus(ctx, ids, status)
}

// BatchDelete 批量删除用户 (string IDs)
func (r *MongoUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	return r.baseRepo.BatchDelete(ctx, ids)
}

// CountByStatus 按状态统计用户数量
func (r *MongoUserRepository) CountByStatus(ctx context.Context, status usersModel.UserStatus) (int64, error) {
	return r.baseRepo.CountByStatus(ctx, status)
}

// CountByRole 按角色统计用户数量
func (r *MongoUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	return r.baseRepo.CountByRole(ctx, role)
}

// List 获取用户列表
func (r *MongoUserRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*usersModel.User, error) {
	return r.baseRepo.List(ctx, filter)
}

// Count 统计用户数量
func (r *MongoUserRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return r.baseRepo.Count(ctx, filter)
}

// Exists 检查用户是否存在
func (r *MongoUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.baseRepo.Exists(ctx, id)
}

// === user模块特有方法的实现 ===

// GetByPhone 根据手机号获取用户 (user特有)
func (r *MongoUserRepository) GetByPhone(ctx context.Context, phone string) (*usersModel.User, error) {
	var user usersModel.User
	filter := bson.M{"phone": phone}

	err := r.db.Collection(UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				fmt.Sprintf("手机号 %s 不存在", phone),
				err,
			)
		}
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}

	return &user, nil
}

// ExistsByUsername 检查用户名是否存在 (user特有)
func (r *MongoUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查用户名存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在 (user特有)
func (r *MongoUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查邮箱存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// ExistsByPhone 检查手机号是否存在 (user特有)
func (r *MongoUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	filter := bson.M{"phone": phone}
	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查手机号存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// UpdateLastLogin 更新最后登录时间和IP (user特有)
func (r *MongoUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	updates := map[string]interface{}{
		"last_login_at": time.Now(),
		"last_login_ip": ip,
	}
	return r.baseRepo.Update(ctx, id, updates)
}

// UpdatePasswordByEmail 根据邮箱更新密码 (user特有)
func (r *MongoUserRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	// 先获取用户ID
	user, err := r.GetByEmail(ctx, email)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("邮箱 %s 不存在", email),
			nil,
		)
	}
	return r.baseRepo.UpdatePassword(ctx, user.ID.Hex(), hashedPassword)
}

// GetActiveUsers 获取活跃用户列表 (user特有)
func (r *MongoUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*usersModel.User, error) {
	filter := bson.M{"status": "active"}
	opts := options.Find().
		SetSort(bson.M{"last_login_at": -1}).
		SetLimit(limit)

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询活跃用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetUsersByRole 根据角色获取用户列表 (user特有)
func (r *MongoUserRepository) GetUsersByRole(ctx context.Context, role string, limit int64) ([]*usersModel.User, error) {
	filter := bson.M{"roles": role} // 使用 roles 数组
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit)

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询角色用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	return users, nil
}

// SetPhoneVerified 设置手机号验证状态 (user特有)
func (r *MongoUserRepository) SetPhoneVerified(ctx context.Context, id string, verified bool) error {
	updates := map[string]interface{}{
		"phone_verified": verified,
	}
	return r.baseRepo.Update(ctx, id, updates)
}

// UnbindEmail 解绑邮箱 (user特有)
func (r *MongoUserRepository) UnbindEmail(ctx context.Context, id string) error {
	updates := map[string]interface{}{
		"email":          "",
		"email_verified": false,
	}
	return r.baseRepo.Update(ctx, id, updates)
}

// UnbindPhone 解绑手机 (user特有)
func (r *MongoUserRepository) UnbindPhone(ctx context.Context, id string) error {
	updates := map[string]interface{}{
		"phone":          "",
		"phone_verified": false,
	}
	return r.baseRepo.Update(ctx, id, updates)
}

// DeleteDevice 删除设备 (user特有，暂未实现)
func (r *MongoUserRepository) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	// TODO: 实现设备删除逻辑
	return UserInterface.NewUserRepositoryError(
		UserInterface.ErrorTypeNotFound,
		"设备不存在或已删除",
		nil,
	)
}

// GetDevices 获取用户设备列表 (user特有，暂未实现)
func (r *MongoUserRepository) GetDevices(ctx context.Context, userID string) ([]interface{}, error) {
	// TODO: 实现获取设备列表逻辑
	return []interface{}{}, nil
}

// FindWithFilter 使用UserFilter进行高级查询 (user特有)
func (r *MongoUserRepository) FindWithFilter(ctx context.Context, filter *usersModel.UserFilter) ([]*usersModel.User, int64, error) {
	// 设置默认值
	filter.SetDefaults()

	// 构建MongoDB查询条件
	mongoFilter := bson.M{}

	// 基础字段过滤
	if filter.ID != "" {
		mongoFilter["_id"] = filter.ID
	}
	if filter.Username != "" {
		mongoFilter["username"] = bson.M{"$regex": filter.Username, "$options": "i"}
	}
	if filter.Email != "" {
		mongoFilter["email"] = bson.M{"$regex": filter.Email, "$options": "i"}
	}
	if filter.Phone != "" {
		mongoFilter["phone"] = filter.Phone
	}

	// 角色和状态
	if filter.Role != "" {
		mongoFilter["roles"] = filter.Role // 使用 roles 数组
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}

	// 验证状态
	if filter.EmailVerified != nil {
		mongoFilter["email_verified"] = *filter.EmailVerified
	}
	if filter.PhoneVerified != nil {
		mongoFilter["phone_verified"] = *filter.PhoneVerified
	}

	// 时间范围
	if filter.CreatedAfter != nil {
		mongoFilter["created_at"] = bson.M{"$gte": *filter.CreatedAfter}
	}
	if filter.CreatedBefore != nil {
		if mongoFilter["created_at"] == nil {
			mongoFilter["created_at"] = bson.M{}
		}
		mongoFilter["created_at"].(bson.M)["$lte"] = *filter.CreatedBefore
	}

	// 关键词搜索
	if filter.SearchKeyword != "" {
		mongoFilter["$or"] = []bson.M{
			{"username": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
			{"email": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
		}
	}

	// 获取总数
	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计用户数量失败",
			err,
		)
	}

	// 构建查询选项
	opts := options.Find()
	opts.SetSkip(int64(filter.GetSkip()))
	opts.SetLimit(int64(filter.GetLimit()))

	// 排序
	sortOrder := 1
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}
	opts.SetSort(bson.M{filter.SortBy: sortOrder})

	// 执行查询
	cursor, err := r.db.Collection(UserCollection).Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	return users, total, nil
}

// SearchUsers 关键词搜索用户 (user特有)
func (r *MongoUserRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*usersModel.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"搜索用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	return users, nil
}

// Transaction 执行事务操作 (user特有)
func (r *MongoUserRepository) Transaction(ctx context.Context,
	user *usersModel.User,
	fn func(ctx context.Context, repo UserInterface.UserRepository) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"启动事务会话失败",
			err,
		)
	}
	defer session.EndSession(ctx)

	// 执行事务
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx, r)
	})

	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"执行事务失败",
			err,
		)
	}

	return nil
}

// BatchCreate 批量创建用户 (user特有)
func (r *MongoUserRepository) BatchCreate(ctx context.Context, users []*usersModel.User) error {
	if len(users) == 0 {
		return nil
	}

	// 验证所有用户数据
	for i, user := range users {
		if user == nil {
			return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeValidation,
				"用户对象不能为空", nil)
		}
		if err := r.ValidateUser(*user); err != nil {
			return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeValidation,
				fmt.Sprintf("用户 #%d 数据验证失败: %v", i+1, err), err)
		}
	}

	// 使用 shared 实现
	for _, user := range users {
		if err := r.baseRepo.Create(ctx, user); err != nil {
			if err.Error() == "用户数据验证失败" {
				return err
			}
			return UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"批量创建用户失败",
				err,
			)
		}
	}

	return nil
}

// === 其他 user 特有方法 ===

// ListWithPagination 分页获取用户列表 (user特有)
func (r *MongoUserRepository) ListWithPagination(ctx context.Context, filter UserInterface.UserFilter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[usersModel.User], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{}

	if filter.ID != "" {
		mongoFilter["_id"] = filter.ID
	}
	if filter.Username != "" {
		mongoFilter["username"] = bson.M{"$regex": filter.Username, "$options": "i"}
	}
	if filter.Email != "" {
		mongoFilter["email"] = bson.M{"$regex": filter.Email, "$options": "i"}
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}

	// 计算分页参数
	pagination.CalculatePagination()

	// 获取总数
	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"获取用户总数失败",
			err,
		)
	}

	// 构建查询选项
	opts := options.Find()
	opts.SetSkip(int64(pagination.Skip))
	opts.SetLimit(int64(pagination.PageSize))

	// 执行查询
	cursor, err := r.db.Collection(UserCollection).Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询用户列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	// 创建分页结果
	pagedResult := &infrastructure.PagedResult[usersModel.User]{
		Data:       users,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)),
		HasNext:    pagination.Page < int((total+int64(pagination.PageSize)-1)/int64(pagination.PageSize)),
		HasPrev:    pagination.Page > 1,
	}
	return pagedResult, nil
}

// FindWithPagination 通用分页查询 (user特有)
func (r *MongoUserRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[usersModel.User], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{}

	// 如果是UserFilter类型，添加过滤条件
	if userFilter, ok := filter.(UserInterface.UserFilter); ok {
		if userFilter.ID != "" {
			mongoFilter["_id"] = userFilter.ID
		}
		if userFilter.Username != "" {
			mongoFilter["username"] = bson.M{"$regex": userFilter.Username, "$options": "i"}
		}
		if userFilter.Email != "" {
			mongoFilter["email"] = bson.M{"$regex": userFilter.Email, "$options": "i"}
		}
		if userFilter.Status != "" {
			mongoFilter["status"] = userFilter.Status
		}
	} else if filter != nil {
		// 使用通用Filter接口的条件
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	// 计算分页参数
	pagination.CalculatePagination()

	// 获取总数
	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计用户总数失败",
			err,
		)
	}

	// 构建查询选项
	findOptions := options.Find()
	findOptions.SetSkip(int64(pagination.Skip))
	findOptions.SetLimit(int64(pagination.PageSize))

	// 设置排序
	if filter != nil {
		sort := filter.GetSort()
		if len(sort) > 0 {
			findOptions.SetSort(sort)
		}
	}

	// 执行查询
	cursor, err := r.db.Collection(UserCollection).Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []*usersModel.User
	for cursor.Next(ctx) {
		var user usersModel.User
		if err := cursor.Decode(&user); err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}

	// 创建分页结果
	pagedResult := &infrastructure.PagedResult[usersModel.User]{
		Data:       users,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)),
		HasNext:    pagination.Page < int((total+int64(pagination.PageSize)-1)/int64(pagination.PageSize)),
		HasPrev:    pagination.Page > 1,
	}
	return pagedResult, nil
}

// BatchUpdate 批量更新用户 (user特有)
func (r *MongoUserRepository) BatchUpdate(ctx context.Context, filters []UserInterface.UserFilter, updates map[string]interface{}) error {
	if len(filters) == 0 {
		return nil
	}

	// 构建批量更新的过滤条件
	var orConditions []bson.M
	for _, filter := range filters {
		condition := bson.M{}

		if filter.Username != "" {
			condition["username"] = filter.Username
		}
		if filter.Email != "" {
			condition["email"] = filter.Email
		}
		if filter.Status != "" {
			condition["status"] = filter.Status
		}

		if len(condition) > 0 {
			orConditions = append(orConditions, condition)
		}
	}

	if len(orConditions) == 0 {
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeValidation,
			"没有有效的过滤条件", nil)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 批量更新
	filter := bson.M{"$or": orConditions}
	update := bson.M{"$set": updates}
	_, err := r.db.Collection(UserCollection).UpdateMany(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeInternal,
			"批量更新用户失败", err)
	}

	return nil
}

// HardDelete 硬删除用户 (user特有)
func (r *MongoUserRepository) HardDelete(ctx context.Context, id string) error {
	// 将string ID转换为ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	filter := bson.M{"_id": objID}
	result, err := r.db.Collection(UserCollection).DeleteOne(ctx, filter)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"硬删除用户失败",
			err,
		)
	}

	if result.DeletedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}

	return nil
}

// ValidateUser 验证用户数据 (user特有)
func (r *MongoUserRepository) ValidateUser(user usersModel.User) error {
	if user.Username == "" {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"用户名不能为空",
			nil,
		)
	}

	if user.Email == "" && user.Phone == "" {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"邮箱或手机号不能为空",
			nil,
		)
	}

	if user.Password == "" {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"密码不能为空",
			nil,
		)
	}

	return nil
}
