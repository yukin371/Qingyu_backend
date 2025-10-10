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
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	UserInterface "Qingyu_backend/repository/interfaces/user"
)

// MongoUserRepository MongoDB用户仓储实现
// 基于新的架构设计，实现BaseRepository和UserRepository接口
type MongoUserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoUserRepository 创建新的MongoDB用户仓储实例
func NewMongoUserRepository(db *mongo.Database) UserInterface.UserRepository {
	// 创建新的MongoDB用户仓储实例
	return &MongoUserRepository{
		db:         db,
		collection: db.Collection("users"),
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

// Create 创建用户
func (r *MongoUserRepository) Create(ctx context.Context, user *usersModel.User) error {
	if user == nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"用户对象不能为空",
			nil,
		)
	}

	actualUser := *user

	// 验证用户数据
	if err := r.ValidateUser(actualUser); err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"用户数据验证失败",
			err,
		)
	}

	// 设置创建时间
	now := time.Now()
	actualUser.CreatedAt = now
	actualUser.UpdatedAt = now

	// 生成ObjectID
	if actualUser.ID == "" {
		actualUser.ID = primitive.NewObjectID().Hex()
	}

	// 插入文档
	_, err := r.collection.InsertOne(ctx, actualUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeDuplicate,
				"用户名或邮箱已存在",
				err,
			)
		}
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"创建用户失败",
			err,
		)
	}

	return nil
}

// GetByID 根据ID获取用户
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*usersModel.User, error) {
	var user usersModel.User

	// 构建查询条件
	mongoFilter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}

	err := r.collection.FindOne(ctx, mongoFilter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				"用户不存在",
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

// Update 更新用户信息
func (r *MongoUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 构建更新条件
	mongoFilter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, mongoFilter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新用户失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			"用户不存在",
			nil,
		)
	}

	return nil
}

// Delete 软删除用户
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	// 构建查询条件
	mongoFilter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.collection.UpdateOne(ctx, mongoFilter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"删除用户失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			"没有找到匹配的用户",
			nil,
		)
	}

	return nil
}

// HardDelete 硬删除用户
func (r *MongoUserRepository) HardDelete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
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

// List 获取用户列表
func (r *MongoUserRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*usersModel.User, error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// 尝试将 Filter 转换为 UserFilter 以使用特定字段
	var userFilter UserInterface.UserFilter
	if f, ok := filter.(UserInterface.UserFilter); ok {
		userFilter = f
	}

	// 如果是UserFilter类型，使用自定义过滤条件
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

	// 构建查询选项
	opts := options.Find()
	if userFilter.Limit > 0 {
		opts.SetLimit(userFilter.Limit)
	}
	if userFilter.Offset > 0 {
		opts.SetSkip(userFilter.Offset)
	}
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
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
		userPtr := &user
		users = append(users, userPtr)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	return users, nil
}

// ListWithPagination 分页获取用户列表
func (r *MongoUserRepository) ListWithPagination(ctx context.Context, filter UserInterface.UserFilter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[usersModel.User], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// 添加过滤条件
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
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
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
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
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
		userPtr := &user
		users = append(users, userPtr)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
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

// Count 统计用户数量
func (r *MongoUserRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// 尝试将 Filter 转换为 UserFilter 以使用特定字段
	var userFilter UserInterface.UserFilter
	if f, ok := filter.(UserInterface.UserFilter); ok {
		userFilter = f
	}

	// 添加过滤条件
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
	if !userFilter.FromDate.IsZero() {
		mongoFilter["created_at"] = bson.M{"$gte": userFilter.FromDate}
	}
	if !userFilter.ToDate.IsZero() {
		if mongoFilter["created_at"] == nil {
			mongoFilter["created_at"] = bson.M{}
		}
		mongoFilter["created_at"].(bson.M)["$lte"] = userFilter.ToDate
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计用户数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查用户是否存在
func (r *MongoUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查用户存在性失败",
			err,
		)
	}

	return count > 0, nil
}

// GetByUsername 根据用户名获取用户
func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*usersModel.User, error) {
	var user usersModel.User

	filter := bson.M{
		"username":   username,
		"deleted_at": bson.M{"$exists": false},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				fmt.Sprintf("用户名 %s 不存在", username),
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

// GetByEmail 根据邮箱获取用户
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	var user usersModel.User

	filter := bson.M{
		"email":      email,
		"deleted_at": bson.M{"$exists": false},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeNotFound,
				fmt.Sprintf("邮箱 %s 不存在", email),
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

// ExistsByUsername 检查用户名是否存在
func (r *MongoUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{
		"username":   username,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查用户名存在性失败",
			err,
		)
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *MongoUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{
		"email":      email,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查邮箱存在性失败",
			err,
		)
	}

	return count > 0, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *MongoUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"last_login_at": now, "updated_at": now}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新最后登录时间失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}

	return nil
}

// UpdatePassword 更新密码
func (r *MongoUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updated_at": now}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新密码失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}

	return nil
}

// GetActiveUsers 获取活跃用户列表
func (r *MongoUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*usersModel.User, error) {
	// 构建查询条件：未删除且状态为活跃
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"status":     "active",
	}

	// 按最后登录时间排序
	opts := options.Find().
		SetSort(bson.M{"last_login_at": -1}).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
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

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	return users, nil
}

// BatchUpdate 批量更新用户
func (r *MongoUserRepository) BatchUpdate(ctx context.Context, filters []UserInterface.UserFilter, updates map[string]interface{}) error {
	if len(filters) == 0 {
		return nil
	}

	// 构建批量更新的过滤条件
	var orConditions []bson.M
	for _, filter := range filters {
		condition := bson.M{}

		// 根据过滤器构建查询条件
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
	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeInternal,
			"批量更新用户失败", err)
	}

	return nil
}

// BatchDelete 批量删除用户
func (r *MongoUserRepository) BatchDelete(ctx context.Context, filters []UserInterface.UserFilter) error {
	if len(filters) == 0 {
		return nil
	}

	// 构建删除条件
	var deleteFilters []bson.M
	for _, filter := range filters {
		mongoFilter := bson.M{}

		// 根据过滤器构建查询条件
		if filter.Username != "" {
			mongoFilter["username"] = filter.Username
		}
		if filter.Email != "" {
			mongoFilter["email"] = filter.Email
		}
		if filter.Status != "" {
			mongoFilter["status"] = filter.Status
		}

		if len(mongoFilter) > 0 {
			deleteFilters = append(deleteFilters, mongoFilter)
		}
	}

	if len(deleteFilters) == 0 {
		return nil
	}

	// 批量删除
	filter := bson.M{"$or": deleteFilters}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeInternal,
			"批量删除用户失败", err)
	}

	return nil
}

// Transaction 执行事务操作
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
		// 创建事务内的Repository实例
		txRepo := &MongoUserRepository{
			db:         r.db,
			collection: r.collection,
		}

		return nil, fn(sessCtx, txRepo)
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

// BatchCreate 批量创建用户
// 建议仅当测试使用，生产环境中应避免批量创建（通过Excel文件导入用户数据来验证功能）
func (r *MongoUserRepository) BatchCreate(ctx context.Context, users []*usersModel.User) error {
	if len(users) == 0 {
		return nil
	}

	// 验证用户对象
	for _, user := range users {
		if user == nil {
			return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeValidation,
				"用户对象不能为空", nil)
		}
	}

	// 验证所有用户数据
	for i, user := range users {
		if err := r.ValidateUser(*user); err != nil {
			return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeValidation,
				fmt.Sprintf("用户 #%d 数据验证失败: %v", i+1, err), err)
		}
	}

	// 准备批量插入的文档
	var documents []interface{}
	now := time.Now()
	for _, user := range users {
		(*user).CreatedAt = now
		(*user).UpdatedAt = now
		if (*user).ID == "" {
			(*user).ID = primitive.NewObjectID().Hex()
		}
		documents = append(documents, *user)
	}

	// 执行批量插入
	_, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeDuplicate,
				"批量创建用户失败，存在重复数据", err)
		}
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeInternal,
			"批量创建用户失败", err)
	}

	return nil
}

// 删除buildFilter方法，已在各方法中直接实现过滤器构建

// ValidateUser 验证用户数据
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

// FindWithPagination 通用分页查询
func (r *MongoUserRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[usersModel.User], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

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
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
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
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
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
		userPtr := &user
		users = append(users, userPtr)
	}

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
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

// 删除重复的Health方法
// 已在文件开头定义了Health方法
