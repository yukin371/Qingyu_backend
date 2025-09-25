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

// MongoUserRepository MongoDB用户仓储实现
type MongoUserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	queryBuilder base.QueryBuilder
}

// NewMongoUserRepository 创建MongoDB用户仓储实例
func NewMongoUserRepository(db *mongo.Database) interfaces.UserRepository {
	return &MongoUserRepository{
		db:         db,
		collection: db.Collection("users"),
		queryBuilder: base.NewMongoQueryBuilder(),
	}
}

// Create 创建用户
func (r *MongoUserRepository) Create(ctx context.Context, user **system.User) error {
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
func (r *MongoUserRepository) GetByID(ctx context.Context, filter interfaces.UserFilter) (**system.User, error) {
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

// Update 更新用户
func (r *MongoUserRepository) Update(ctx context.Context, filter interfaces.UserFilter, updates map[string]interface{}) error {
	// 构建查询条件
	mongoFilter := r.buildFilter(filter)
	
	// 添加更新时间
	updates["updated_at"] = time.Now()
	
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
func (r *MongoUserRepository) Delete(ctx context.Context, filter interfaces.UserFilter) error {
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
func (r *MongoUserRepository) HardDelete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"删除用户失败",
			err,
		)
	}
	
	if result.DeletedCount == 0 {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeNotFound,
			"用户不存在",
			nil,
		)
	}
	
	return nil
}

// List 获取用户列表
func (r *MongoUserRepository) List(ctx context.Context, filter base.Filter) ([]**system.User, error) {
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
func (r *MongoUserRepository) ListWithPagination(ctx context.Context, filter interfaces.UserFilter, pagination base.Pagination) (*base.PagedResult[*system.User], error) {
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
func (r *MongoUserRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	// 构建查询条件
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}
	
	if userFilter, ok := filter.(interfaces.UserFilter); ok {
		mongoFilter = r.buildFilter(userFilter)
	} else if filter != nil {
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
func (r *MongoUserRepository) Exists(ctx context.Context, filter interfaces.UserFilter) (bool, error) {
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
func (r *MongoUserRepository) GetByUsername(ctx context.Context, username string) (*system.User, error) {
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
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*system.User, error) {
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
func (r *MongoUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
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
func (r *MongoUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
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
func (r *MongoUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	filter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	
	update := bson.M{
		"$set": bson.M{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		},
	}
	
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
			"用户不存在",
			nil,
		)
	}
	
	return nil
}

// UpdatePassword 更新用户密码
func (r *MongoUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	filter := bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}
	
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		},
	}
	
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
			"用户不存在",
			nil,
		)
	}
	
	return nil
}

// GetActiveUsers 获取活跃用户
func (r *MongoUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
		"status":     "active",
	}
	
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{"last_login_at", -1}})
	
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
func (r *MongoUserRepository) BatchUpdate(ctx context.Context, filters []interfaces.UserFilter, updates map[string]interface{}) error {
	// 添加更新时间
	updates["updated_at"] = time.Now()
	
	for _, filter := range filters {
		mongoFilter := r.buildFilter(filter)
		update := bson.M{"$set": updates}
		
		_, err := r.collection.UpdateMany(ctx, mongoFilter, update)
		if err != nil {
			return interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"批量更新用户失败",
				err,
			)
		}
	}
	
	return nil
}

// BatchDelete 批量删除用户
func (r *MongoUserRepository) BatchDelete(ctx context.Context, filters []interfaces.UserFilter) error {
	now := time.Now()
	
	for _, filter := range filters {
		mongoFilter := r.buildFilter(filter)
		update := bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}}
		
		_, err := r.collection.UpdateMany(ctx, mongoFilter, update)
		if err != nil {
			return interfaces.NewRepositoryError(
				interfaces.ErrorTypeInternal,
				"批量删除用户失败",
				err,
			)
		}
	}
	
	return nil
}

// Transaction 执行事务
func (r *MongoUserRepository) Transaction(ctx context.Context, fn func(ctx context.Context, repo interfaces.UserRepository) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"启动事务失败",
			err,
		)
	}
	defer session.EndSession(ctx)
	
	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err := session.WithTransaction(sc, func(sc mongo.SessionContext) (interface{}, error) {
			err := fn(sc, r)
			return nil, err
		})
		return err
	})
}

// BatchCreate 批量创建用户
func (r *MongoUserRepository) BatchCreate(ctx context.Context, users []**system.User) error {
	if len(users) == 0 {
		return nil
	}
	
	var docs []interface{}
	now := time.Now()
	
	for _, user := range users {
		if user == nil || *user == nil {
			continue
		}
		
		actualUser := *user
		
		// 验证用户数据
		if err := r.validateUser(actualUser); err != nil {
			return err
		}
		
		// 设置时间戳
		actualUser.CreatedAt = now
		actualUser.UpdatedAt = now
		
		// 生成ID
		if actualUser.ID == "" {
			actualUser.ID = primitive.NewObjectID().Hex()
		}
		
		docs = append(docs, actualUser)
	}
	
	if len(docs) == 0 {
		return nil
	}
	
	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return interfaces.NewRepositoryError(
			interfaces.ErrorTypeInternal,
			"批量创建用户失败",
			err,
		)
	}
	
	return nil
}

// buildFilter 构建MongoDB过滤器
func (r *MongoUserRepository) buildFilter(filter interfaces.UserFilter) bson.M {
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}
	
	if filter.Username != "" {
		mongoFilter["username"] = filter.Username
	}
	if filter.Email != "" {
		mongoFilter["email"] = filter.Email
	}
	if filter.Status != "" {
		mongoFilter["status"] = filter.Status
	}
	
	return mongoFilter
}

// validateUser 验证用户数据
func (r *MongoUserRepository) validateUser(user *system.User) error {
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
	return nil
}

// FindWithPagination 分页查询用户
func (r *MongoUserRepository) FindWithPagination(ctx context.Context, filter base.Filter, pagination base.Pagination) (*base.PagedResult[*system.User], error) {
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
func (r *MongoUserRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}