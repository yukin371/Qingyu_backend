package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/mongodb/base"
	usersModel "Qingyu_backend/models/users"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	UserInterface "Qingyu_backend/repository/interfaces/user"
)

// MongoUserRepository MongoDB用户仓储实现
// 基于新的架构设计，实现BaseRepository和UserRepository接口
type MongoUserRepository struct {
	*base.BaseMongoRepository  // 嵌入基类，继承ID转换和通用CRUD方法喵~
	db *mongo.Database
}

// NewMongoUserRepository 创建新的MongoDB用户仓储实例
func NewMongoUserRepository(db *mongo.Database) UserInterface.UserRepository {
	return &MongoUserRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "users"),
		db:                 db,
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
	if actualUser.ID.IsZero() {
		actualUser.ID = primitive.NewObjectID()
	}

	// 插入文档
	_, err := r.GetCollection().InsertOne(ctx, actualUser)
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

	// 将生成的ID和时间戳设置回原始对象
	user.ID = actualUser.ID
	user.CreatedAt = actualUser.CreatedAt
	user.UpdatedAt = actualUser.UpdatedAt

	return nil
}

// GetByID 根据ID获取用户
func (r *MongoUserRepository) GetByID(ctx context.Context, id string) (*usersModel.User, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	var user usersModel.User

	// 构建查询条件
	mongoFilter := bson.M{
		"_id":        objectID,
		"deleted_at": bson.M{"$exists": false},
	}

	err = r.GetCollection().FindOne(ctx, mongoFilter).Decode(&user)
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 构建更新条件
	mongoFilter := bson.M{
		"_id":        objectID,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": updates}

	result, err := r.GetCollection().UpdateOne(ctx, mongoFilter, update)
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	// 构建查询条件
	mongoFilter := bson.M{
		"_id":        objectID,
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, mongoFilter, update)
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
	objectID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	filter := bson.M{"_id": objectID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
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
		objID, err := r.ParseID(userFilter.ID)
		if err != nil {
			return nil, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeValidation,
				"无效的用户ID",
				err,
			)
		}
		mongoFilter["_id"] = objID
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
	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
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
	total, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
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
	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
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
		objID, err := r.ParseID(userFilter.ID)
		if err != nil {
			return 0, UserInterface.NewUserRepositoryError(
				UserInterface.ErrorTypeValidation,
				"无效的用户ID",
				err,
			)
		}
		mongoFilter["_id"] = objID
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

	count, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
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

	count, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
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
		"username": username,
		// 注意：User 模型使用 status 字段而不是 deleted_at
		// 软删除检查应该在业务层处理，这里只查询用户是否存在
	}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&user)
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
		"email": email,
		// 注意：User 模型使用 status 字段而不是 deleted_at
		// 软删除检查应该在业务层处理，这里只查询用户是否存在
	}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&user)
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
		"username": username,
		// 注意：User 模型使用 status 字段，软删除检查在业务层处理
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
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
		"email": email,
		// 注意：User 模型使用 status 字段，软删除检查在业务层处理
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查邮箱存在性失败",
			err,
		)
	}

	return count > 0, nil
}

// UpdateLastLogin 更新最后登录时间和IP
func (r *MongoUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	objID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{
		"last_login_at": now,
		"last_login_ip": ip,
		"updated_at":    now,
	}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
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
	objID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
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

// UpdatePasswordByEmail 根据邮箱更新密码
func (r *MongoUserRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	now := time.Now()
	filter := bson.M{"email": email, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"根据邮箱更新密码失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeNotFound,
			fmt.Sprintf("邮箱 %s 不存在", email),
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

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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
	_, err := r.GetCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(UserInterface.ErrorTypeInternal,
			"批量更新用户失败", err)
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
			BaseMongoRepository: r.BaseMongoRepository,
			db:                 r.db,
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
		if (*user).ID.IsZero() {
			(*user).ID = primitive.NewObjectID()
		}
		documents = append(documents, *user)
	}

	// 执行批量插入
	_, err := r.GetCollection().InsertMany(ctx, documents)
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
			objID, err := r.ParseID(userFilter.ID)
			if err != nil {
				return nil, UserInterface.NewUserRepositoryError(
					UserInterface.ErrorTypeValidation,
					"无效的用户ID",
					err,
				)
			}
			mongoFilter["_id"] = objID
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
	total, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
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
	cursor, err := r.GetCollection().Find(ctx, mongoFilter, findOptions)
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

// GetByPhone 根据手机号获取用户
func (r *MongoUserRepository) GetByPhone(ctx context.Context, phone string) (*usersModel.User, error) {
	var user usersModel.User

	filter := bson.M{
		"phone":      phone,
		"deleted_at": bson.M{"$exists": false},
	}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&user)
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

// ExistsByPhone 检查手机号是否存在
func (r *MongoUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	filter := bson.M{
		"phone":      phone,
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"检查手机号存在性失败",
			err,
		)
	}

	return count > 0, nil
}

// UpdateStatus 更新用户状态
func (r *MongoUserRepository) UpdateStatus(ctx context.Context, id string, status usersModel.UserStatus) error {
	objID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"更新用户状态失败",
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

// GetUsersByRole 根据角色获取用户列表
func (r *MongoUserRepository) GetUsersByRole(ctx context.Context, role string, limit int64) ([]*usersModel.User, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"role":       role,
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit)

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	return users, nil
}

// SetEmailVerified 设置邮箱验证状态
func (r *MongoUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	objID, err := r.ParseID(id)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"email_verified": verified, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"设置邮箱验证状态失败",
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

// SetPhoneVerified 设置手机号验证状态
func (r *MongoUserRepository) SetPhoneVerified(ctx context.Context, id string, verified bool) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"phone_verified": verified, "updated_at": now}}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"设置手机验证状态失败",
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

// BatchUpdateStatus 批量更新用户状态
func (r *MongoUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status usersModel.UserStatus) error {
	if len(ids) == 0 {
		return nil
	}

	// 将string IDs转换为ObjectIDs
	objectIDs, err := r.ParseIDs(ids)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeValidation,
			"无效的用户ID",
			err,
		)
	}

	now := time.Now()
	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": now}}

	_, err = r.GetCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"批量更新用户状态失败",
			err,
		)
	}

	return nil
}

// BatchDelete 批量删除用户（软删除）
func (r *MongoUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	now := time.Now()
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"deleted_at": bson.M{"$exists": false},
	}
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}

	_, err := r.GetCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"批量删除用户失败",
			err,
		)
	}

	return nil
}

// FindWithFilter 使用UserFilter进行高级查询
func (r *MongoUserRepository) FindWithFilter(ctx context.Context, filter *usersModel.UserFilter) ([]*usersModel.User, int64, error) {
	// 设置默认值
	filter.SetDefaults()

	// 构建MongoDB查询条件
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

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
		mongoFilter["role"] = filter.Role
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

	// 关键词搜索（用户名、昵称、邮箱）
	if filter.SearchKeyword != "" {
		mongoFilter["$or"] = []bson.M{
			{"username": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
			{"email": bson.M{"$regex": filter.SearchKeyword, "$options": "i"}},
		}
	}

	// 获取总数
	total, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
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
	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
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

	if err := cursor.Err(); err != nil {
		return nil, 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	return users, total, nil
}

// SearchUsers 关键词搜索用户
func (r *MongoUserRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*usersModel.User, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"$or": []bson.M{
			{"username": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
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

	if err := cursor.Err(); err != nil {
		return nil, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	return users, nil
}

// CountByRole 按角色统计用户数量
func (r *MongoUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"role":       role,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计角色用户数量失败",
			err,
		)
	}

	return count, nil
}

// CountByStatus 按状态统计用户数量
func (r *MongoUserRepository) CountByStatus(ctx context.Context, status usersModel.UserStatus) (int64, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"status":     status,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"统计状态用户数量失败",
			err,
		)
	}

	return count, nil
}

// ==================== 新增方法：邮箱/手机/设备管理 ====================

// UnbindEmail 解绑邮箱（清空邮箱并标记为未验证）
func (r *MongoUserRepository) UnbindEmail(ctx context.Context, id string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{
		"$set": bson.M{
			"email":          "",
			"email_verified": false,
			"updated_at":     now,
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"解绑邮箱失败",
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

// UnbindPhone 解绑手机（清空手机号并标记为未验证）
func (r *MongoUserRepository) UnbindPhone(ctx context.Context, id string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{
		"$set": bson.M{
			"phone":          "",
			"phone_verified": false,
			"updated_at":     now,
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return UserInterface.NewUserRepositoryError(
			UserInterface.ErrorTypeInternal,
			"解绑手机失败",
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

// DeleteDevice 删除设备
func (r *MongoUserRepository) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	// TODO: 实现设备删除逻辑
	// 需要创建DeviceRepository和Device模型
	// 暂时返回"设备不存在"错误
	return UserInterface.NewUserRepositoryError(
		UserInterface.ErrorTypeNotFound,
		"设备不存在或已删除",
		nil,
	)
}

// GetDevices 获取用户设备列表
func (r *MongoUserRepository) GetDevices(ctx context.Context, userID string) ([]interface{}, error) {
	// TODO: 实现获取设备列表逻辑
	// 需要创建DeviceRepository和Device模型
	// 暂时返回空列表
	return []interface{}{}, nil
}
