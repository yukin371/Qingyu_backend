package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	BookstoreInterfaces "Qingyu_backend/repository/interfaces/bookstore"
	ReadingInterfaces "Qingyu_backend/repository/interfaces/reading"
	UserInterface "Qingyu_backend/repository/interfaces/user"
	WritingInterfaces "Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/config"
	base "Qingyu_backend/repository/interfaces/infrastructure"
)

// MongoRepositoryFactory MongoDB仓储工厂实现
type MongoRepositoryFactory struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config.MongoDBConfig
}

// NewMongoRepositoryFactory 创建MongoDB仓储工厂
func NewMongoRepositoryFactory(config *config.MongoDBConfig) (infrastructure.RepositoryFactory, error) {
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

	return &MongoRepositoryFactory{
		client:   client,
		database: database,
		config:   config,
	}, nil
}

// CreateUserRepository 创建用户Repository
func (f *MongoRepositoryFactory) CreateUserRepository() UserInterface.UserRepository {
	return NewMongoUserRepository(f.database)
}

// CreateProjectRepository 创建项目Repository
func (f *MongoRepositoryFactory) CreateProjectRepository() UserInterface.ProjectRepository {
	return NewMongoProjectRepositoryNew(f.database)
}

// CreateReadingSettingsRepository 创建阅读设置Repository
func (f *MongoRepositoryFactory) CreateReadingSettingsRepository() ReadingInterfaces.ReadingSettingsRepository {
	return NewMongoReadingSettingsRepository(f.database)
}

// CreateRoleRepository 创建角色Repository
func (f *MongoRepositoryFactory) CreateRoleRepository() UserInterface.RoleRepository {
	return NewMongoRoleRepository(f.database)
}

// Close 关闭数据库连接
func (f *MongoRepositoryFactory) Close() error {
	if f.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return f.client.Disconnect(ctx)
	}
	return nil
}

// Health 健康检查
func (f *MongoRepositoryFactory) Health(ctx context.Context) error {
	if f.client == nil {
		return fmt.Errorf("MongoDB客户端未初始化")
	}

	// 设置超时
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return f.client.Ping(pingCtx, nil)
}

// GetDatabaseType 获取数据库类型
func (f *MongoRepositoryFactory) GetDatabaseType() string {
	return infrastructure.DatabaseTypeMongoDB
}

// MongoProjectRepositoryNew 新的MongoDB项目仓储实现
type MongoProjectRepositoryNew struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoProjectRepositoryNew 创建新的MongoDB项目仓储实例
func NewMongoProjectRepositoryNew(db *mongo.Database) UserInterface.ProjectRepository {
	return &MongoProjectRepositoryNew{
		db:         db,
		collection: db.Collection("projects"),
	}
}

// Create 创建项目
func (r *MongoProjectRepositoryNew) Create(ctx context.Context, project *interface{}) error {
	if project == nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeValidation,
			"项目对象不能为空",
			nil,
		)
	}

	_, err := r.collection.InsertOne(ctx, *project)
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"创建项目失败",
			err,
		)
	}
	return nil
}

// GetByID 根据ID获取项目
func (r *MongoProjectRepositoryNew) GetByID(ctx context.Context, id interface{}) (*interface{}, error) {
	var project interface{}
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeNotFound,
				fmt.Sprintf("项目ID %v 不存在", id),
				err,
			)
		}
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询项目失败",
			err,
		)
	}
	return &project, nil
}

// Update 更新项目
func (r *MongoProjectRepositoryNew) Update(ctx context.Context, id interface{}, updates map[string]interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"更新项目失败",
			err,
		)
	}
	return nil
}

// Delete 软删除项目
func (r *MongoProjectRepositoryNew) Delete(ctx context.Context, id interface{}) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"is_deleted": true,
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"删除项目失败",
			err,
		)
	}
	return nil
}

// HardDelete 硬删除项目
func (r *MongoProjectRepositoryNew) HardDelete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"删除项目失败",
			err,
		)
	}
	return nil
}

// List 获取项目列表
func (r *MongoProjectRepositoryNew) List(ctx context.Context, filter infrastructure.Filter) ([]*interface{}, error) {
	cursor, err := r.collection.Find(ctx, filter.GetConditions())
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询项目列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var projects []*interface{}
	for cursor.Next(ctx) {
		var project interface{}
		if err := cursor.Decode(&project); err != nil {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeInternal,
				"解析项目数据失败",
				err,
			)
		}
		projects = append(projects, &project)
	}
	return projects, nil
}

// ListWithPagination 分页获取项目列表
func (r *MongoProjectRepositoryNew) ListWithPagination(ctx context.Context, filter interface{}, pagination infrastructure.Pagination) (*infrastructure.PagedResult[interface{}], error) {
	// 实现分页逻辑
	return nil, nil
}

// FindWithPagination 分页查找项目
func (r *MongoProjectRepositoryNew) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[interface{}], error) {
	// 构建查询条件
	conditions := filter.GetConditions()
	
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, conditions)
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"统计项目总数失败",
			err,
		)
	}

	// 构建查询选项
	opts := options.Find()
	if pagination.Skip > 0 {
		opts.SetSkip(int64(pagination.Skip))
	}
	if pagination.PageSize > 0 {
		opts.SetLimit(int64(pagination.PageSize))
	}

	// 执行查询
	cursor, err := r.collection.Find(ctx, conditions, opts)
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询项目失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var projects []interface{}
	for cursor.Next(ctx) {
		var project interface{}
		if err := cursor.Decode(&project); err != nil {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeInternal,
				"解析项目数据失败",
				err,
			)
		}
		projects = append(projects, project)
	}

	if err := cursor.Err(); err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"遍历项目数据失败",
			err,
		)
	}

	// 创建分页结果
	return infrastructure.NewPagedResult[interface{}](projects, total, pagination), nil
}

// Count 统计项目数量
func (r *MongoProjectRepositoryNew) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter.GetConditions())
	if err != nil {
		return 0, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"统计项目数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查项目是否存在
func (r *MongoProjectRepositoryNew) Exists(ctx context.Context, id interface{}) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"检查项目存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// GetByCreatorID 根据创建者ID获取项目
func (r *MongoProjectRepositoryNew) GetByCreatorID(ctx context.Context, creatorID string) ([]interface{}, error) {
	// 创建一个简单的Filter实现
	filter := &infrastructure.BaseFilter{
		Conditions: map[string]interface{}{"creator_id": creatorID},
	}

	cursor, err := r.collection.Find(ctx, filter.GetConditions())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []interface{}
	for cursor.Next(ctx) {
		var project interface{}
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// ArchiveByCreatorID 根据创建者ID归档项目
func (r *MongoProjectRepositoryNew) ArchiveByCreatorID(ctx context.Context, creatorID string) error {
	updates := map[string]interface{}{
		"archived_at": time.Now(),
		"is_archived": true,
	}
	_, err := r.collection.UpdateMany(ctx, bson.M{"creator_id": creatorID}, bson.M{"$set": updates})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"归档项目失败",
			err,
		)
	}
	return nil
}

// BatchUpdate 批量更新项目
func (r *MongoProjectRepositoryNew) BatchUpdate(ctx context.Context, ids []interface{}, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeValidation,
			"项目ID列表不能为空", nil)
	}
	if len(updates) == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeValidation,
			"更新数据不能为空", nil)
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"批量更新项目失败",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeNotFound,
			"没有找到匹配的项目", nil)
	}

	return nil
}

// BatchCreate 批量创建项目
func (r *MongoProjectRepositoryNew) BatchCreate(ctx context.Context, projects []*interface{}) error {
	if len(projects) == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeValidation,
			"项目列表不能为空", nil)
	}

	// 转换为[]interface{}类型
	docs := make([]interface{}, len(projects))
	for i, project := range projects {
		if project == nil {
			return infrastructure.NewRepositoryError(infrastructure.ErrorTypeDuplicate,
				fmt.Sprintf("第%d个项目对象不能为空", i+1), nil)
		}
		docs[i] = *project
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeInternal,
			"批量创建项目失败", err)
	}

	return nil
}

// BatchDelete 批量删除项目
func (r *MongoProjectRepositoryNew) BatchDelete(ctx context.Context, ids []interface{}) error {
	if len(ids) == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeValidation,
			"项目ID列表不能为空", nil)
	}
	if len(ids) > 100 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeValidation,
			"批量删除数量不能超过100个", nil)
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	updates := bson.M{"$set": bson.M{
		"deleted_at": time.Now(),
		"is_deleted": true,
	}}

	result, err := r.collection.UpdateMany(ctx, filter, updates)
	if err != nil {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeInternal,
			"批量删除项目失败", err)
	}

	if result.MatchedCount == 0 {
		return infrastructure.NewRepositoryError(infrastructure.ErrorTypeNotFound,
			"没有找到匹配的项目", nil)
	}

	return nil
}

// Health 健康检查
func (r *MongoProjectRepositoryNew) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// MongoRoleRepository 新的MongoDB角色仓储实现
type MongoRoleRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoRoleRepository 创建新的MongoDB角色仓储实例
func NewMongoRoleRepository(db *mongo.Database) UserInterface.RoleRepository {
	return &MongoRoleRepository{
		db:         db,
		collection: db.Collection("roles"),
	}
}

// Create 创建角色
func (r *MongoRoleRepository) Create(ctx context.Context, role *interface{}) error {
	if role == nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeValidation,
			"角色对象不能为空",
			nil,
		)
	}

	_, err := r.collection.InsertOne(ctx, *role)
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"创建角色失败",
			err,
		)
	}
	return nil
}

// GetByID 根据ID获取角色
func (r *MongoRoleRepository) GetByID(ctx context.Context, id interface{}) (*interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeNotFound,
				fmt.Sprintf("角色ID %v 不存在", id),
				err,
			)
		}
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	return &role, nil
}

// Update 更新角色
func (r *MongoRoleRepository) Update(ctx context.Context, id interface{}, updates map[string]interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"更新角色失败",
			err,
		)
	}
	return nil
}

// Delete 软删除角色
func (r *MongoRoleRepository) Delete(ctx context.Context, id interface{}) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"is_deleted": true,
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"删除角色失败",
			err,
		)
	}
	return nil
}

// HardDelete 硬删除角色
func (r *MongoRoleRepository) HardDelete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"删除角色失败",
			err,
		)
	}
	return nil
}

// List 获取角色列表
func (r *MongoRoleRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*interface{}, error) {
	cursor, err := r.collection.Find(ctx, filter.GetConditions())
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询角色列表失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roles []*interface{}
	for cursor.Next(ctx) {
		var role interface{}
		if err := cursor.Decode(&role); err != nil {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, &role)
	}
	return roles, nil
}

// ListWithPagination 分页获取角色列表
func (r *MongoRoleRepository) ListWithPagination(ctx context.Context, filter interface{}, pagination infrastructure.Pagination) (*infrastructure.PagedResult[interface{}], error) {
	// 实现分页逻辑
	return nil, nil
}

// FindWithPagination 分页查找角色
func (r *MongoRoleRepository) FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[interface{}], error) {
	// 构建查询条件
	conditions := filter.GetConditions()
	
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, conditions)
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"统计角色总数失败",
			err,
		)
	}

	// 构建查询选项
	opts := options.Find()
	if pagination.Skip > 0 {
		opts.SetSkip(int64(pagination.Skip))
	}
	if pagination.PageSize > 0 {
		opts.SetLimit(int64(pagination.PageSize))
	}

	// 执行查询
	cursor, err := r.collection.Find(ctx, conditions, opts)
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var roles []interface{}
	for cursor.Next(ctx) {
		var role interface{}
		if err := cursor.Decode(&role); err != nil {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeInternal,
				"解析角色数据失败",
				err,
			)
		}
		roles = append(roles, role)
	}

	if err := cursor.Err(); err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"遍历角色数据失败",
			err,
		)
	}

	// 创建分页结果
	return infrastructure.NewPagedResult[interface{}](roles, total, pagination), nil
}

// Count 统计角色数量
func (r *MongoRoleRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter.GetConditions())
	if err != nil {
		return 0, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"统计角色数量失败",
			err,
		)
	}
	return count, nil
}

// Exists 检查角色是否存在
func (r *MongoRoleRepository) Exists(ctx context.Context, id interface{}) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"检查角色存在性失败",
			err,
		)
	}
	return count > 0, nil
}

// GetByName 根据名称获取角色
func (r *MongoRoleRepository) GetByName(ctx context.Context, name string) (interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeNotFound,
				fmt.Sprintf("角色名称 %s 不存在", name),
				err,
			)
		}
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询角色失败",
			err,
		)
	}
	return role, nil
}

// GetDefaultRole 获取默认角色
func (r *MongoRoleRepository) GetDefaultRole(ctx context.Context) (interface{}, error) {
	var role interface{}
	err := r.collection.FindOne(ctx, map[string]interface{}{"is_default": true}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeNotFound,
				"默认角色不存在",
				err,
			)
		}
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询默认角色失败",
			err,
		)
	}
	return role, nil
}

// GetUserRoles 获取用户角色
func (r *MongoRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]interface{}, error) {
	// 这里需要查询用户角色关联表
	userRolesCollection := r.db.Collection("user_roles")
	cursor, err := userRolesCollection.Find(ctx, map[string]interface{}{"user_id": userID})
	if err != nil {
		return nil, infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"查询用户角色失败",
			err,
		)
	}
	defer cursor.Close(ctx)

	var roleIDs []string
	for cursor.Next(ctx) {
		var userRole map[string]interface{}
		if err := cursor.Decode(&userRole); err != nil {
			return nil, infrastructure.NewRepositoryError(
				infrastructure.ErrorTypeInternal,
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
	filter := &infrastructure.BaseFilter{
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
func (r *MongoRoleRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	userRolesCollection := r.db.Collection("user_roles")
	userRole := map[string]interface{}{
		"user_id":    userID,
		"role_id":    roleID,
		"created_at": time.Now(),
	}

	_, err := userRolesCollection.InsertOne(ctx, userRole)
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"分配角色失败",
			err,
		)
	}
	return nil
}

// RemoveRole 移除角色
func (r *MongoRoleRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	userRolesCollection := r.db.Collection("user_roles")
	_, err := userRolesCollection.DeleteOne(ctx, map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	})
	if err != nil {
		return infrastructure.NewRepositoryError(
			infrastructure.ErrorTypeInternal,
			"移除角色失败",
			err,
		)
	}
	return nil
}

// GetUserPermissions 获取用户权限
func (r *MongoRoleRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
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
func (r *MongoRoleRepository) BatchUpdate(ctx context.Context, ids []interface{}, updates map[string]interface{}) error {
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
func (r *MongoRoleRepository) BatchCreate(ctx context.Context, roles []*interface{}) error {
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

// BatchDelete 批量删除角色
func (r *MongoRoleRepository) BatchDelete(ctx context.Context, ids []interface{}) error {
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
func (r *MongoRoleRepository) Health(ctx context.Context) error {
	// 执行简单的ping操作来检查连接
	return r.db.Client().Ping(ctx, nil)
}
