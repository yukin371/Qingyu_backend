package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	base "Qingyu_backend/repository/interfaces"
)

// MongoRepositoryFactoryNew 新的MongoDB仓储工厂实现
type MongoRepositoryFactoryNew struct {
	client   *mongo.Client
	database *mongo.Database
	config   *base.MongoConfig
}

// NewMongoRepositoryFactoryNew 创建新的MongoDB仓储工厂
func NewMongoRepositoryFactoryNew(config *base.MongoConfig) (base.RepositoryFactory, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("MongoDB配置验证失败: %w", err)
	}

	// 创建客户端选项
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetConnectTimeout(config.ConnectTimeout).
		SetServerSelectionTimeout(config.ServerTimeout)

	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("连接MongoDB失败: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("MongoDB连接测试失败: %w", err)
	}

	database := client.Database(config.Database)

	return &MongoRepositoryFactoryNew{
		client:   client,
		database: database,
		config:   config,
	}, nil
}

// CreateUserRepository 创建用户Repository
func (f *MongoRepositoryFactoryNew) CreateUserRepository() base.UserRepository {
	return NewMongoUserRepositoryNew(f.database)
}

// CreateProjectRepository 创建项目Repository
func (f *MongoRepositoryFactoryNew) CreateProjectRepository() base.ProjectRepository {
	return NewMongoProjectRepositoryNew(f.database)
}

// CreateRoleRepository 创建角色Repository
func (f *MongoRepositoryFactoryNew) CreateRoleRepository() base.RoleRepository {
	return NewMongoRoleRepositoryNew(f.database)
}

// Close 关闭数据库连接
func (f *MongoRepositoryFactoryNew) Close() error {
	if f.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return f.client.Disconnect(ctx)
	}
	return nil
}

// Health 健康检查
func (f *MongoRepositoryFactoryNew) Health(ctx context.Context) error {
	if f.client == nil {
		return fmt.Errorf("MongoDB客户端未初始化")
	}

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return f.client.Ping(pingCtx, nil)
}

// GetDatabaseType 获取数据库类型
func (f *MongoRepositoryFactoryNew) GetDatabaseType() string {
	return base.DatabaseTypeMongoDB
}

// MongoProjectRepositoryNew 新的MongoDB项目仓储实现
type MongoProjectRepositoryNew struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoProjectRepositoryNew 创建新的MongoDB项目仓储实例
func NewMongoProjectRepositoryNew(db *mongo.Database) base.ProjectRepository {
	return &MongoProjectRepositoryNew{
		db:         db,
		collection: db.Collection("projects"),
	}
}

// Create 创建项目
func (r *MongoProjectRepositoryNew) Create(ctx context.Context, project *interface{}) error {
	if project == nil {
		return base.NewRepositoryError(
			base.ErrorTypeValidation,
			"项目对象不能为空",
			nil,
		)
	}

	_, err := r.collection.InsertOne(ctx, *project)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"创建项目失败",
			err,
		)
	}
	return nil
}

// GetByID 根据ID获取项目
func (r *MongoProjectRepositoryNew) GetByID(ctx context.Context, id interface{}) (*interface{}, error) {
	var project interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": id}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, base.NewRepositoryError(
				base.ErrorTypeNotFound,
				fmt.Sprintf("项目ID %v 不存在", id),
				err,
			)
		}
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询项目失败",
			err,
		)
	}
	return &project, nil
}

// Update 更新项目
func (r *MongoProjectRepositoryNew) Update(ctx context.Context, id interface{}, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, map[string]interface{}{"_id": id}, map[string]interface{}{"$set": updates})
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"更新项目失败",
			err,
		)
	}
	return nil
}

// Delete 删除项目
func (r *MongoProjectRepositoryNew) Delete(ctx context.Context, id interface{}) error {
	_, err := r.collection.UpdateOne(ctx,
		map[string]interface{}{"_id": id},
		map[string]interface{}{"$set": map[string]interface{}{"deleted_at": time.Now()}})
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"删除项目失败",
			err,
		)
	}
	return nil
}

// HardDelete 硬删除项目
func (r *MongoProjectRepositoryNew) HardDelete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, map[string]interface{}{"_id": id})
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"硬删除项目失败",
			err,
		)
	}
	return nil
}

// List 获取项目列表
func (r *MongoProjectRepositoryNew) List(ctx context.Context, filter base.Filter) ([]*interface{}, error) {
	mongoFilter := bson.M{}

	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询项目列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var projects []*interface{}
	for cursor.Next(ctx) {
		var project interface{}
		if err := cursor.Decode(&project); err != nil {
			return nil, base.NewRepositoryError(
				base.ErrorTypeInternal,
				"解析项目数据失败",
				err,
			)
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

// ListWithPagination 分页获取项目列表
func (r *MongoProjectRepositoryNew) ListWithPagination(ctx context.Context, filter interface{}, pagination base.Pagination) (*base.PagedResult[interface{}], error) {
	// 实现分页逻辑
	return nil, fmt.Errorf("项目分页查询功能待实现")
}

// FindWithPagination 分页查询项目
func (r *MongoProjectRepositoryNew) FindWithPagination(ctx context.Context, filter base.Filter, pagination base.Pagination) (*base.PagedResult[interface{}], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

	if filter != nil {
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
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"统计项目总数失败",
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
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询项目失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var projects []*interface{}
	for cursor.Next(ctx) {
		var project interface{}
		if err := cursor.Decode(&project); err != nil {
			return nil, base.NewRepositoryError(
				base.ErrorTypeInternal,
				"解析项目数据失败",
				err,
			)
		}
		projects = append(projects, &project)
	}

	if err := cursor.Err(); err != nil {
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"遍历项目数据失败",
			err,
		)
	}

	// 创建分页结果
	return base.NewPagedResult[interface{}](projects, total, pagination), nil
}

// Count 统计项目数量
func (r *MongoProjectRepositoryNew) Count(ctx context.Context, filter base.Filter) (int64, error) {
	mongoFilter := bson.M{}

	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"统计项目数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查项目是否存在
func (r *MongoProjectRepositoryNew) Exists(ctx context.Context, id interface{}) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, map[string]interface{}{"_id": id})
	if err != nil {
		return false, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"检查项目存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// GetByCreatorID 根据创建者ID获取项目列表
func (r *MongoProjectRepositoryNew) GetByCreatorID(ctx context.Context, creatorID string) ([]interface{}, error) {
	// 创建一个简单的Filter实现
	filter := &base.BaseFilter{
		Conditions: map[string]interface{}{
			"creator_id": creatorID,
			"deleted_at": map[string]interface{}{"$exists": false},
		},
	}

	// 调用List方法并转换返回类型
	projectPointers, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 将[]*interface{}转换为[]interface{}
	projects := make([]interface{}, len(projectPointers))
	for i, projectPtr := range projectPointers {
		if projectPtr != nil {
			projects[i] = *projectPtr
		}
	}

	return projects, nil
}

// ArchiveByCreatorID 根据创建者ID归档项目
func (r *MongoProjectRepositoryNew) ArchiveByCreatorID(ctx context.Context, creatorID string) error {
	filter := map[string]interface{}{
		"creator_id": creatorID,
		"deleted_at": map[string]interface{}{"$exists": false},
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"archived_at": time.Now(),
			"updated_at":  time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"归档项目失败",
			err,
		)
	}
	return nil
}

// BatchUpdate 批量更新项目
func (r *MongoProjectRepositoryNew) BatchUpdate(ctx context.Context, ids []interface{}, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	// 转换ID为ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		idStr, ok := id.(string)
		if !ok {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"项目ID必须是字符串类型", nil)
		}
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"无效的项目ID", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{"$set": updates}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"批量更新项目失败",
			err,
		)
	}
	return nil
}

// BatchCreate 批量创建项目
func (r *MongoProjectRepositoryNew) BatchCreate(ctx context.Context, projects []*interface{}) error {
	if len(projects) == 0 {
		return nil
	}

	// 转换为实际的项目对象
	actualProjects := make([]interface{}, len(projects))
	for i, project := range projects {
		if project == nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"项目对象不能为空", nil)
		}
		actualProjects[i] = *project
	}

	// 批量插入
	_, err := r.collection.InsertMany(ctx, actualProjects)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return base.NewRepositoryError(base.ErrorTypeDuplicate,
				"项目已存在", err)
		}
		return base.NewRepositoryError(base.ErrorTypeInternal,
			"批量创建项目失败", err)
	}

	return nil
}

// MongoRoleRepositoryNew 新的MongoDB角色仓储实现
type MongoRoleRepositoryNew struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoRoleRepositoryNew 创建新的MongoDB角色仓储实例
func NewMongoRoleRepositoryNew(db *mongo.Database) base.RoleRepository {
	return &MongoRoleRepositoryNew{
		db:         db,
		collection: db.Collection("roles"),
	}
}

// Create 创建角色
func (r *MongoRoleRepositoryNew) Create(ctx context.Context, role *interface{}) error {
	if role == nil {
		return base.NewRepositoryError(
			base.ErrorTypeValidation,
			"角色对象不能为空",
			nil,
		)
	}

	_, err := r.collection.InsertOne(ctx, *role)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"创建角色失败",
			err,
		)
	}
	return nil
}

// GetByID 根据ID获取角色
func (r *MongoRoleRepositoryNew) GetByID(ctx context.Context, id interface{}) (*interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, base.NewRepositoryError(
				base.ErrorTypeNotFound,
				fmt.Sprintf("角色ID %v 不存在", id),
				err,
			)
		}
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	return &role, nil
}

// Update 更新角色
func (r *MongoRoleRepositoryNew) Update(ctx context.Context, id interface{}, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updates}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete 删除角色
func (r *MongoRoleRepositoryNew) Delete(ctx context.Context, id interface{}) error {
	_, err := r.collection.UpdateOne(ctx,
		map[string]interface{}{"_id": id},
		map[string]interface{}{"$set": map[string]interface{}{"deleted_at": time.Now()}})
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"删除角色失败",
			err,
		)
	}
	return nil
}

// HardDelete 硬删除角色
func (r *MongoRoleRepositoryNew) HardDelete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, map[string]interface{}{"_id": id})
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"硬删除角色失败",
			err,
		)
	}
	return nil
}

// List 获取角色列表
func (r *MongoRoleRepositoryNew) List(ctx context.Context, filter base.Filter) ([]*interface{}, error) {
	mongoFilter := bson.M{}

	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询角色列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roles []*interface{}
	for cursor.Next(ctx) {
		var role interface{}
		if err := cursor.Decode(&role); err != nil {
			return nil, base.NewRepositoryError(
				base.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

// ListWithPagination 分页获取角色列表
func (r *MongoRoleRepositoryNew) ListWithPagination(ctx context.Context, filter interface{}, pagination base.Pagination) (*base.PagedResult[interface{}], error) {
	// 实现分页逻辑
	return nil, fmt.Errorf("角色分页查询功能待实现")
}

// FindWithPagination 分页查询角色
func (r *MongoRoleRepositoryNew) FindWithPagination(ctx context.Context, filter base.Filter, pagination base.Pagination) (*base.PagedResult[interface{}], error) {
	// 构建MongoDB过滤器
	mongoFilter := bson.M{"deleted_at": bson.M{"$exists": false}}

	if filter != nil {
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
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"统计角色总数失败",
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
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var roles []*interface{}
	for cursor.Next(ctx) {
		var role interface{}
		if err := cursor.Decode(&role); err != nil {
			return nil, base.NewRepositoryError(
				base.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}

	if err := cursor.Err(); err != nil {
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"遍历角色数据失败",
			err,
		)
	}

	// 创建分页结果
	return base.NewPagedResult[interface{}](roles, total, pagination), nil
}

// Count 统计角色数量
func (r *MongoRoleRepositoryNew) Count(ctx context.Context, filter base.Filter) (int64, error) {
	mongoFilter := bson.M{}

	if filter != nil {
		conditions := filter.GetConditions()
		for k, v := range conditions {
			mongoFilter[k] = v
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"统计角色数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查角色是否存在
func (r *MongoRoleRepositoryNew) Exists(ctx context.Context, id interface{}) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, map[string]interface{}{"_id": id})
	if err != nil {
		return false, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"检查角色存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// GetByName 根据名称获取角色
func (r *MongoRoleRepositoryNew) GetByName(ctx context.Context, name string) (interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, base.NewRepositoryError(
				base.ErrorTypeNotFound,
				fmt.Sprintf("角色名称 %s 不存在", name),
				err,
			)
		}
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	return role, nil
}

// GetDefaultRole 获取默认角色
func (r *MongoRoleRepositoryNew) GetDefaultRole(ctx context.Context) (interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"is_default": true}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, base.NewRepositoryError(
				base.ErrorTypeNotFound,
				"默认角色不存在",
				err,
			)
		}
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询默认角色失败",
			err,
		)
	}
	return role, nil
}

// GetUserRoles 获取用户角色
func (r *MongoRoleRepositoryNew) GetUserRoles(ctx context.Context, userID string) ([]interface{}, error) {
	// 这里需要查询用户角色关联表
	userRolesCollection := r.db.Collection("user_roles")
	cursor, err := userRolesCollection.Find(ctx, map[string]interface{}{"user_id": userID})
	if err != nil {
		return nil, base.NewRepositoryError(
			base.ErrorTypeInternal,
			"查询用户角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roleIDs []string
	for cursor.Next(ctx) {
		var userRole map[string]interface{}
		if err := cursor.Decode(&userRole); err != nil {
			return nil, base.NewRepositoryError(
				base.ErrorTypeInternal,
				"解析用户角色数据失败",
				err,
			)
		}
		if roleID, ok := userRole["role_id"].(string); ok {
			roleIDs = append(roleIDs, roleID)
		}
	}

	// 根据角色ID查询角色详情
	if len(roleIDs) == 0 {
		return []interface{}{}, nil
	}

	// 创建一个简单的Filter实现
	filter := &base.BaseFilter{
		Conditions: map[string]interface{}{"_id": map[string]interface{}{"$in": roleIDs}},
	}

	// 调用List方法并转换返回类型
	rolePointers, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 将[]*interface{}转换为[]interface{}
	roles := make([]interface{}, len(rolePointers))
	for i, rolePtr := range rolePointers {
		if rolePtr != nil {
			roles[i] = *rolePtr
		}
	}

	return roles, nil
}

// AssignRole 分配角色
func (r *MongoRoleRepositoryNew) AssignRole(ctx context.Context, userID, roleID string) error {
	userRolesCollection := r.db.Collection("user_roles")
	userRole := map[string]interface{}{
		"user_id":    userID,
		"role_id":    roleID,
		"created_at": time.Now(),
	}

	_, err := userRolesCollection.InsertOne(ctx, userRole)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"分配角色失败",
			err,
		)
	}
	return nil
}

// RemoveRole 移除角色
func (r *MongoRoleRepositoryNew) RemoveRole(ctx context.Context, userID, roleID string) error {
	userRolesCollection := r.db.Collection("user_roles")
	filter := map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	}

	_, err := userRolesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"移除角色失败",
			err,
		)
	}
	return nil
}

// GetUserPermissions 获取用户权限
func (r *MongoRoleRepositoryNew) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	// 获取用户角色
	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 收集所有权限
	permissionSet := make(map[string]bool)
	for _, role := range roles {
		if roleMap, ok := role.(map[string]interface{}); ok {
			if permissions, ok := roleMap["permissions"].([]interface{}); ok {
				for _, perm := range permissions {
					if permStr, ok := perm.(string); ok {
						permissionSet[permStr] = true
					}
				}
			}
		}
	}

	// 转换为切片
	var permissions []string
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// BatchUpdate 批量更新角色
func (r *MongoRoleRepositoryNew) BatchUpdate(ctx context.Context, ids []interface{}, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	// 转换ID为ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		idStr, ok := id.(string)
		if !ok {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"角色ID必须是字符串类型", nil)
		}
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"无效的角色ID", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{"$set": updates}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return base.NewRepositoryError(
			base.ErrorTypeInternal,
			"批量更新角色失败",
			err,
		)
	}
	return nil
}

// BatchCreate 批量创建角色
func (r *MongoRoleRepositoryNew) BatchCreate(ctx context.Context, roles []*interface{}) error {
	if len(roles) == 0 {
		return nil
	}

	// 转换为实际的角色对象
	actualRoles := make([]interface{}, len(roles))
	for i, role := range roles {
		if role == nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"角色对象不能为空", nil)
		}
		actualRoles[i] = *role
	}

	// 批量插入
	_, err := r.collection.InsertMany(ctx, actualRoles)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return base.NewRepositoryError(base.ErrorTypeDuplicate,
				"角色已存在", err)
		}
		return base.NewRepositoryError(base.ErrorTypeInternal,
			"批量创建角色失败", err)
	}

	return nil
}

// BatchDelete 批量删除项目
func (r *MongoProjectRepositoryNew) BatchDelete(ctx context.Context, ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	// 转换ID为ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		idStr, ok := id.(string)
		if !ok {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"项目ID必须是字符串类型", nil)
		}
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"无效的项目ID", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	// 批量删除
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return base.NewRepositoryError(base.ErrorTypeInternal,
			"批量删除项目失败", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoProjectRepositoryNew) Health(ctx context.Context) error {
	// 执行简单的ping操作来检查连接
	return r.db.Client().Ping(ctx, nil)
}

// BatchDelete 批量删除角色
func (r *MongoRoleRepositoryNew) BatchDelete(ctx context.Context, ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	// 转换ID为ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		idStr, ok := id.(string)
		if !ok {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"角色ID必须是字符串类型", nil)
		}
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return base.NewRepositoryError(base.ErrorTypeValidation,
				"无效的角色ID", err)
		}
		objectIDs = append(objectIDs, objectID)
	}

	// 批量删除
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return base.NewRepositoryError(base.ErrorTypeInternal,
			"批量删除角色失败", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoRoleRepositoryNew) Health(ctx context.Context) error {
	// 执行简单的ping操作来检查连接
	return r.db.Client().Ping(ctx, nil)
}
