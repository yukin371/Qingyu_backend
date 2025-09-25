package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/system"
	"Qingyu_backend/repository/base"
	"Qingyu_backend/repository/interfaces"
)

// MongoUserRepositoryNew 新的MongoDB用户仓储实现
// 基于新的架构设计，实现BaseRepository和UserRepository接口
type MongoUserRepositoryNew struct {
	db         *mongo.Database
	collection *mongo.Collection
	queryBuilder base.QueryBuilder
}

// NewMongoUserRepositoryNew 创建新的MongoDB用户仓储实例
func NewMongoUserRepositoryNew(db *mongo.Database) interfaces.UserRepository {
	return &MongoUserRepositoryNew{
		db:         db,
		collection: db.Collection("users"),
		queryBuilder: base.NewMongoQueryBuilder(),
	}
}

// Create 创建用户
func (r *MongoUserRepositoryNew) Create(ctx context.Context, user **system.User) error {
	if user == nil || *user == nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeValidation,
			"用户对象不能为空",
			nil,
		)
	}

	actualUser := *user
	
	// 验证用户数据
	if err := r.validateUser(actualUser); err != nil {
		return err
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
			return interfaces.NewRepositoryError(
				interfaces.ErrorTypeDuplicate,
				"用户名或邮箱已存在",
				err,
			)
		}
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"创建用户失败",
			err,
		)
	}

	return nil
}

// GetByID 根据ID获取用户
func (r *MongoUserRepositoryNew) GetByID(ctx context.Context, filter interfaces.UserFilter) (**system.User, error) {
	var user system.User
	
	// 构建查询条件
	mongoFilter := r.buildFilter(filter)
	
	err := r.collection.FindOne(ctx, mongoFilter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeNotFound,
				"用户不存在",
				err,
			)
		}
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}

	userPtr := &user
	return &userPtr, nil
}

// Update 更新用户信息
func (r *MongoUserRepositoryNew) Update(ctx context.Context, filter interfaces.UserFilter, updates map[string]interface{}) error {
	// 添加更新时间
	updates["updated_at"] = time.Now()
	
	// 构建更新条件
	mongoFilter := r.buildFilter(filter)
	update := bson.M{"$set": updates}
	
	result, err := r.collection.UpdateOne(ctx, mongoFilter, update)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"更新用户失败",
			err,
		)
	}
	
	if result.MatchedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			"用户不存在",
			nil,
		)
	}
	
	return nil
}

// Delete 软删除用户
func (r *MongoUserRepositoryNew) Delete(ctx context.Context, filter interfaces.UserFilter) error {
	now := time.Now()
	mongoFilter := r.buildFilter(filter)
	update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}
	
	result, err := r.collection.UpdateMany(ctx, mongoFilter, update)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"删除用户失败",
			err,
		)
	}
	
	if result.MatchedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			"没有找到匹配的用户",
			nil,
		)
	}
	
	return nil
}

// HardDelete 硬删除用户
func (r *MongoUserRepositoryNew) HardDelete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"硬删除用户失败",
			err,
		)
	}
	
	if result.DeletedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}
	
	return nil
}

// List 获取用户列表
func (r *MongoUserRepositoryNew) List(ctx context.Context, filter base.Filter) ([]**system.User, error) {
	// 构建查询条件
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}
	
	// 如果是UserFilter类型，使用buildFilter方法
	if userFilter, ok := filter.(interfaces.UserFilter); ok {
		mongoFilter = r.buildFilter(userFilter)
	} else if filter != nil {
		// 使用通用Filter接口的条件
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	// 构建查询选项
	opts := options.Find()
	if userFilter, ok := filter.(interfaces.UserFilter); ok {
		if userFilter.Limit > 0 {
			opts.SetLimit(userFilter.Limit)
		}
		if userFilter.Offset > 0 {
			opts.SetSkip(userFilter.Offset)
		}
	}
	opts.SetSort(bson.D{{"created_at", -1}})

	// 执行查询
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)
	
	// 解析结果
	var users []**system.User
	for cursor.Next(ctx) {
		var user system.User
		if err := cursor.Decode(&user); err != nil {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		userPtr := &user
		users = append(users, &userPtr)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}
	
	return users, nil
}

// ListWithPagination 分页获取用户列表
func (r *MongoUserRepositoryNew) ListWithPagination(ctx context.Context, filter interfaces.UserFilter, pagination base.Pagination) (*base.PagedResult[*system.User], error) {
	// 构建查询条件
	mongoFilter := r.buildFilter(filter)
	
	// 计算分页参数
	pagination.CalculatePagination()
	
	// 获取总数
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
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
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)
	
	// 解析结果
	var users []**system.User
	for cursor.Next(ctx) {
		var user system.User
		if err := cursor.Decode(&user); err != nil {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		userPtr := &user
		users = append(users, &userPtr)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}
	
	// 创建分页结果
	return base.NewPagedResult[*system.User](users, total, pagination), nil
}

// Count 统计用户数量
func (r *MongoUserRepositoryNew) Count(ctx context.Context, filter base.Filter) (int64, error) {
	mongoFilter := bson.M{}
	
	// 如果是UserFilter类型，使用buildFilter方法
	if userFilter, ok := filter.(interfaces.UserFilter); ok {
		mongoFilter = r.buildFilter(userFilter)
	} else if filter != nil {
		// 使用通用Filter接口的条件
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"统计用户数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查用户是否存在
func (r *MongoUserRepositoryNew) Exists(ctx context.Context, filter interfaces.UserFilter) (bool, error) {
	mongoFilter := r.buildFilter(filter)
	
	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return false, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"检查用户存在性失败",
			err,
		)
	}
	
	return count > 0, nil
}

// GetByUsername 根据用户名获取用户
func (r *MongoUserRepositoryNew) GetByUsername(ctx context.Context, username string) (*system.User, error) {
	var user system.User
	
	filter := bson.M{
		"username":   username,
		"deleted_at": bson.M{"$exists": false},
	}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeNotFound,
				fmt.Sprintf("用户名 %s 不存在", username),
				err,
			)
		}
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}

	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *MongoUserRepositoryNew) GetByEmail(ctx context.Context, email string) (*system.User, error) {
	var user system.User
	
	filter := bson.M{
		"email":      email,
		"deleted_at": bson.M{"$exists": false},
	}
	
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeNotFound,
				fmt.Sprintf("邮箱 %s 不存在", email),
				err,
			)
		}
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}

	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *MongoUserRepositoryNew) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{
		"username":   username,
		"deleted_at": bson.M{"$exists": false},
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"检查用户名存在性失败",
			err,
		)
	}
	
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *MongoUserRepositoryNew) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{
		"email":      email,
		"deleted_at": bson.M{"$exists": false},
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"检查邮箱存在性失败",
			err,
		)
	}
	
	return count > 0, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *MongoUserRepositoryNew) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"last_login_at": now, "updated_at": now}}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"更新最后登录时间失败",
			err,
		)
	}
	
	if result.MatchedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}
	
	return nil
}

// UpdatePassword 更新密码
func (r *MongoUserRepositoryNew) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	now := time.Now()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"password": hashedPassword, "updated_at": now}}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"更新密码失败",
			err,
		)
	}
	
	if result.MatchedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			fmt.Sprintf("用户ID %s 不存在", id),
			nil,
		)
	}
	
	return nil
}

// GetActiveUsers 获取活跃用户
func (r *MongoUserRepositoryNew) GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error) {
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
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询活跃用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)
	
	var users []*system.User
	for cursor.Next(ctx) {
		var user system.User
		if err := cursor.Decode(&user); err != nil {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		users = append(users, &user)
	}
	
	if err := cursor.Err(); err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}
	
	return users, nil
}

// BatchUpdate 批量更新用户
func (r *MongoUserRepositoryNew) BatchUpdate(ctx context.Context, filters []interfaces.UserFilter, updates map[string]interface{}) error {
	if len(filters) == 0 {
		return nil
	}

	// 构建批量更新的过滤条件
	var orConditions []bson.M
	for _, filter := range filters {
		condition := r.buildFilter(filter)
		if len(condition) > 0 {
			orConditions = append(orConditions, condition)
		}
	}

	if len(orConditions) == 0 {
		return interfaces.NewRepositoryError(interfaces.ErrorTypeValidation, 
			"没有有效的过滤条件", nil)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 批量更新
	filter := bson.M{"$or": orConditions}
	update := bson.M{"$set": updates}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return interfaces.NewRepositoryError(interfaces.ErrorTypeInternal, 
			"批量更新用户失败", err)
	}

	return nil
}

// BatchDelete 批量删除用户
func (r *MongoUserRepositoryNew) BatchDelete(ctx context.Context, filters []interfaces.UserFilter) error {
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
		return interfaces.NewRepositoryError(interfaces.ErrorTypeInternal, 
			"批量删除用户失败", err)
	}

	return nil
}

// Transaction 执行事务操作
func (r *MongoUserRepositoryNew) Transaction(ctx context.Context, fn func(ctx context.Context, repo interfaces.UserRepository) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"启动事务会话失败",
			err,
		)
	}
	defer session.EndSession(ctx)
	
	// 执行事务
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 创建事务内的Repository实例
		txRepo := &MongoUserRepositoryNew{
			db:         r.db,
			collection: r.collection,
			queryBuilder: r.queryBuilder,
		}
		
		return nil, fn(sessCtx, txRepo)
	})
	
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"执行事务失败",
			err,
		)
	}
	
	return nil
}

// BatchCreate 批量创建用户
func (r *MongoUserRepositoryNew) BatchCreate(ctx context.Context, users []**system.User) error {
	if len(users) == 0 {
		return nil
	}

	// 验证用户对象
	for _, user := range users {
		if user == nil || *user == nil {
			return interfaces.NewRepositoryError(interfaces.ErrorTypeValidation, 
				"用户对象不能为空", nil)
		}
	}

	// 验证所有用户数据
	for _, user := range users {
		if err := r.validateUser(*user); err != nil {
			return interfaces.NewRepositoryError(interfaces.ErrorTypeValidation, 
				fmt.Sprintf("用户数据验证失败: %v", err), err)
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
			return interfaces.NewRepositoryError(interfaces.ErrorTypeDuplicate, 
				"批量创建用户失败，存在重复数据", err)
		}
		return interfaces.NewRepositoryError(interfaces.ErrorTypeInternal, 
			"批量创建用户失败", err)
	}

	return nil
}

// buildFilter 构建MongoDB查询过滤器
func (r *MongoUserRepositoryNew) buildFilter(filter interfaces.UserFilter) bson.M {
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}
	
	if filter.Username != "" {
		mongoFilter["username"] = bson.M{"$regex": filter.Username, "$options": "i"}
	}
	
	if filter.Email != "" {
		mongoFilter["email"] = bson.M{"$regex": filter.Email, "$options": "i"}
	}
	
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}
	
	// 时间范围过滤
	if !filter.FromDate.IsZero() || !filter.ToDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.FromDate.IsZero() {
			dateFilter["$gte"] = filter.FromDate
		}
		if !filter.ToDate.IsZero() {
			dateFilter["$lte"] = filter.ToDate
		}
		mongoFilter["created_at"] = dateFilter
	}
	
	return mongoFilter
}

// validateUser 验证用户数据
func (r *MongoUserRepositoryNew) validateUser(user *system.User) error {
	if user.Username == "" {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeValidation,
			"用户名不能为空",
			nil,
		)
	}
	
	if user.Email == "" {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeValidation,
			"邮箱不能为空",
			nil,
		)
	}
	
	if user.Password == "" {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeValidation,
			"密码不能为空",
			nil,
		)
	}
	
	return nil
}

// FindWithPagination 分页查询用户
func (r *MongoUserRepositoryNew) FindWithPagination(ctx context.Context, filter base.Filter, pagination base.Pagination) (*base.PagedResult[*system.User], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}
	
	// 如果是UserFilter类型，使用buildFilter方法
	if userFilter, ok := filter.(interfaces.UserFilter); ok {
		mongoFilter = r.buildFilter(userFilter)
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
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
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
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"查询用户失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var users []**system.User
	for cursor.Next(ctx) {
		var user system.User
		if err := cursor.Decode(&user); err != nil {
			return nil, interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"解析用户数据失败",
				err,
			)
		}
		userPtr := &user
		users = append(users, &userPtr)
	}

	if err := cursor.Err(); err != nil {
		return nil, interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"遍历用户数据失败",
			err,
		)
	}

	// 创建分页结果
	return base.NewPagedResult[*system.User](users, total, pagination), nil
}

// Health 健康检查
func (r *MongoUserRepositoryNew) Health(ctx context.Context) error {
	// 执行简单的ping操作来检查数据库连接
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"数据库连接健康检查失败",
			err,
		)
	}
	return nil
}