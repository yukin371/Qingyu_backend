package admin

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/users"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
)

const (
	UserCollection         = "users"
	UserActivityCollection = "user_activities"
)

type MongoUserAdminRepository struct {
	db *mongo.Database
}

// NewMongoUserAdminRepository 创建用户管理仓储
func NewMongoUserAdminRepository(db *mongo.Database) adminrepo.UserAdminRepository {
	return &MongoUserAdminRepository{db: db}
}

// List 获取用户列表（分页、筛选）
func (r *MongoUserAdminRepository) List(ctx context.Context, filter *adminrepo.UserFilter, page, pageSize int) ([]*users.User, int64, error) {
	mongoFilter := r.buildFilter(filter)

	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

// Create 创建用户
func (r *MongoUserAdminRepository) Create(ctx context.Context, user *users.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Collection(UserCollection).InsertOne(ctx, user)
	return err
}

// BatchCreate 批量创建用户
func (r *MongoUserAdminRepository) BatchCreate(ctx context.Context, usersList []*users.User) error {
	if len(usersList) == 0 {
		return nil
	}

	documents := make([]interface{}, len(usersList))
	for i, user := range usersList {
		user.ID = primitive.NewObjectID()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		documents[i] = user
	}

	_, err := r.db.Collection(UserCollection).InsertMany(ctx, documents)
	return err
}

// GetByID 根据ID获取用户
func (r *MongoUserAdminRepository) GetByID(ctx context.Context, userID primitive.ObjectID) (*users.User, error) {
	var user users.User
	filter := bson.M{"_id": userID}

	err := r.db.Collection(UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *MongoUserAdminRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	canonicalEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	var user users.User
	filter := bson.M{"email": canonicalEmail}

	err = r.db.Collection(UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func normalizeEmail(email string) (string, error) {
	trimmed := strings.TrimSpace(email)
	parsed, err := mail.ParseAddress(trimmed)
	if err != nil || parsed == nil {
		return "", fmt.Errorf("invalid email format")
	}
	if parsed.Address != trimmed {
		return "", fmt.Errorf("invalid email format")
	}
	return strings.ToLower(parsed.Address), nil
}

// Update 更新用户信息
func (r *MongoUserAdminRepository) Update(ctx context.Context, userID primitive.ObjectID, user *users.User) error {
	user.UpdatedAt = time.Now()
	user.Touch()

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": user}

	_, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	return err
}

// UpdateStatus 更新用户状态
func (r *MongoUserAdminRepository) UpdateStatus(ctx context.Context, userID primitive.ObjectID, status users.UserStatus) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	return err
}

// Delete 删除用户（软删除）
func (r *MongoUserAdminRepository) Delete(ctx context.Context, userID primitive.ObjectID) error {
	return r.UpdateStatus(ctx, userID, users.UserStatusDeleted)
}

// HardDelete 硬删除用户
func (r *MongoUserAdminRepository) HardDelete(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	_, err := r.db.Collection(UserCollection).DeleteOne(ctx, filter)
	return err
}

// GetActivities 获取用户活动记录
func (r *MongoUserAdminRepository) GetActivities(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*users.UserActivity, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.db.Collection(UserActivityCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserActivityCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var activities []*users.UserActivity
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

// UpdateRoles 更新用户角色
func (r *MongoUserAdminRepository) UpdateRoles(ctx context.Context, userID primitive.ObjectID, role string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"role":       role,
			"updated_at": time.Now(),
		},
	}

	_, err := r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	return err
}

// GetStatistics 获取用户统计信息
func (r *MongoUserAdminRepository) GetStatistics(ctx context.Context, userID primitive.ObjectID) (*users.UserStatistics, error) {
	// 简化实现，实际应该从各个集合统计
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	registerDays := int(time.Since(user.CreatedAt).Hours() / 24)

	stats := &users.UserStatistics{
		UserID:       user.ID.Hex(),
		LastActiveAt: user.LastLoginAt,
		RegisterDays: registerDays,
	}

	// TODO: 从其他集合统计书籍、章节、评论等数据
	return stats, nil
}

// ResetPassword 重置用户密码
func (r *MongoUserAdminRepository) ResetPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"password":   user.Password,
			"updated_at": time.Now(),
		},
	}

	_, err = r.db.Collection(UserCollection).UpdateOne(ctx, filter, update)
	return err
}

// BatchUpdateStatus 批量更新用户状态
func (r *MongoUserAdminRepository) BatchUpdateStatus(ctx context.Context, userIDs []primitive.ObjectID, status users.UserStatus) error {
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.db.Collection(UserCollection).UpdateMany(ctx, filter, update)
	return err
}

// BatchDelete 批量删除用户
func (r *MongoUserAdminRepository) BatchDelete(ctx context.Context, userIDs []primitive.ObjectID) error {
	return r.BatchUpdateStatus(ctx, userIDs, users.UserStatusDeleted)
}

// GetUsersByRole 根据角色获取用户列表
func (r *MongoUserAdminRepository) GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error) {
	filter := bson.M{"role": role}

	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

// SearchUsers 搜索用户
func (r *MongoUserAdminRepository) SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
			{"nickname": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	total, err := r.db.Collection(UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

// CountByStatus 按状态统计用户数量
func (r *MongoUserAdminRepository) CountByStatus(ctx context.Context) (map[string]int64, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":   "$status",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.db.Collection(UserCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]int64)
	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	for _, doc := range docs {
		status := doc["_id"].(string)
		count := doc["count"].(int32)
		results[status] = int64(count)
	}

	// 确保所有状态都有值
	allStatuses := []string{"active", "inactive", "banned", "deleted"}
	for _, status := range allStatuses {
		if _, exists := results[status]; !exists {
			results[status] = 0
		}
	}

	return results, nil
}

// GetRecentUsers 获取最近注册的用户
func (r *MongoUserAdminRepository) GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, err
	}

	return userList, nil
}

// GetActiveUsers 获取活跃用户
func (r *MongoUserAdminRepository) GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	filter := bson.M{
		"last_login_at": bson.M{"$gte": cutoffDate},
		"status":        bson.M{"$ne": "deleted"},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "last_login_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, err
	}

	return userList, nil
}

// buildFilter 构建查询过滤器
func (r *MongoUserAdminRepository) buildFilter(filter *adminrepo.UserFilter) bson.M {
	mongoFilter := bson.M{}

	if filter != nil {
		if filter.Keyword != "" {
			mongoFilter["$or"] = []bson.M{
				{"username": bson.M{"$regex": filter.Keyword, "$options": "i"}},
				{"email": bson.M{"$regex": filter.Keyword, "$options": "i"}},
				{"nickname": bson.M{"$regex": filter.Keyword, "$options": "i"}},
			}
		}
		if filter.Status != "" {
			mongoFilter["status"] = filter.Status
		}
		if filter.Role != "" {
			mongoFilter["role"] = filter.Role
		}
		if filter.DateFrom != nil {
			if _, exists := mongoFilter["created_at"]; !exists {
				mongoFilter["created_at"] = bson.M{}
			}
			mongoFilter["created_at"].(bson.M)["$gte"] = *filter.DateFrom
		}
		if filter.DateTo != nil {
			if _, exists := mongoFilter["created_at"]; !exists {
				mongoFilter["created_at"] = bson.M{}
			}
			mongoFilter["created_at"].(bson.M)["$lte"] = *filter.DateTo
		}
		if filter.LastActive != nil {
			mongoFilter["last_login_at"] = bson.M{"$gte": *filter.LastActive}
		}
	}

	return mongoFilter
}
