package admin

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	authModel "Qingyu_backend/models/auth"
	"Qingyu_backend/models/users"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	sharedrepo "Qingyu_backend/repository/mongodb/shared"
)

const (
	UserCollection         = "users"
	UserActivityCollection = "user_activities"
)

type MongoUserAdminRepository struct {
	// 组合使用统一的MongoUserRepository
	baseRepo *sharedrepo.MongoUserRepository
	db       *mongo.Database
}

type userEmailFilter struct {
	Email string `bson:"email"`
}

// NewMongoUserAdminRepository 创建用户管理仓储
func NewMongoUserAdminRepository(db *mongo.Database) adminrepo.UserAdminRepository {
	return &MongoUserAdminRepository{
		baseRepo: sharedrepo.NewMongoUserRepository(db),
		db:       db,
	}
}

// === 委托共享方法到 baseRepo (使用string ID) ===

// Create 创建用户
func (r *MongoUserAdminRepository) Create(ctx context.Context, user *users.User) error {
	return r.baseRepo.Create(ctx, user)
}

// GetByID 根据ID获取用户 (string ID version)
func (r *MongoUserAdminRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	return r.baseRepo.GetByID(ctx, id)
}

// Update 更新用户信息 (string ID version)
func (r *MongoUserAdminRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.baseRepo.Update(ctx, id, updates)
}

// Delete 删除用户 (string ID version)
func (r *MongoUserAdminRepository) Delete(ctx context.Context, id string) error {
	return r.baseRepo.Delete(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (r *MongoUserAdminRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	return r.baseRepo.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户
func (r *MongoUserAdminRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return r.baseRepo.GetByEmail(ctx, email)
}

// UpdateStatus 更新用户状态 (string ID version)
func (r *MongoUserAdminRepository) UpdateStatus(ctx context.Context, id string, status users.UserStatus) error {
	return r.baseRepo.UpdateStatus(ctx, id, status)
}

// UpdatePassword 更新用户密码 (string ID version)
func (r *MongoUserAdminRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	return r.baseRepo.UpdatePassword(ctx, id, hashedPassword)
}

// SetEmailVerified 设置邮箱验证状态 (string ID version)
func (r *MongoUserAdminRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	return r.baseRepo.SetEmailVerified(ctx, id, verified)
}

// BatchUpdateStatus 批量更新用户状态 (string IDs version)
func (r *MongoUserAdminRepository) BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error {
	return r.baseRepo.BatchUpdateStatus(ctx, ids, status)
}

// BatchDelete 批量删除用户 (string IDs version)
func (r *MongoUserAdminRepository) BatchDelete(ctx context.Context, ids []string) error {
	return r.baseRepo.BatchDelete(ctx, ids)
}

// CountByStatus 按状态统计用户数量 (返回int64)
func (r *MongoUserAdminRepository) CountByStatus(ctx context.Context, status users.UserStatus) (int64, error) {
	return r.baseRepo.CountByStatus(ctx, status)
}

// CountByRole 按角色统计用户数量
func (r *MongoUserAdminRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	return r.baseRepo.CountByRole(ctx, role)
}

// List 获取用户列表 (基础版本)
func (r *MongoUserAdminRepository) List(ctx context.Context, filter base.Filter) ([]*users.User, error) {
	// 直接使用 MongoDB 查询，绕过 baseRepo 以避免 Filter 接口转换问题
	var bsonFilter bson.M
	if filter != nil {
		conditions := filter.GetConditions()
		bsonFilter = bson.M{}
		for k, v := range conditions {
			bsonFilter[k] = v
		}
	} else {
		bsonFilter = bson.M{}
	}

	cursor, err := r.db.Collection(UserCollection).Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*users.User
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// Count 统计用户数量
func (r *MongoUserAdminRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	// 直接使用 MongoDB 查询，绕过 baseRepo 以避免 Filter 接口转换问题
	var bsonFilter bson.M
	if filter != nil {
		conditions := filter.GetConditions()
		bsonFilter = bson.M{}
		for k, v := range conditions {
			bsonFilter[k] = v
		}
	} else {
		bsonFilter = bson.M{}
	}

	count, err := r.db.Collection(UserCollection).CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Exists 检查用户是否存在
func (r *MongoUserAdminRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.baseRepo.Exists(ctx, id)
}

// === admin特有方法的实现 ===

// ListWithPagination 获取用户列表（带分页和筛选，admin特有）
func (r *MongoUserAdminRepository) ListWithPagination(ctx context.Context, filter *adminrepo.UserFilter, page, pageSize int) ([]*users.User, int64, error) {
	// Use a fixed DB query and perform filtering in-memory to avoid dynamic query construction.
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Collection(UserCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var userList []*users.User
	if err = cursor.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	filtered := make([]*users.User, 0, len(userList))
	for _, user := range userList {
		if userMatchesFilter(user, filter) {
			filtered = append(filtered, user)
		}
	}

	total := int64(len(filtered))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return []*users.User{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

// BatchCreate 批量创建用户（admin特有）
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

// HardDelete 硬删除用户（admin特有）
func (r *MongoUserAdminRepository) HardDelete(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	_, err := r.db.Collection(UserCollection).DeleteOne(ctx, filter)
	return err
}

// GetActivities 获取用户活动记录（admin特有）
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

// GetStatistics 获取用户统计信息（admin特有）
func (r *MongoUserAdminRepository) GetStatistics(ctx context.Context, userID primitive.ObjectID) (*users.UserStatistics, error) {
	// 简化实现，实际应该从各个集合统计
	user, err := r.GetByID(ctx, userID.Hex())
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

// ResetPassword 重置用户密码（admin特有）
func (r *MongoUserAdminRepository) ResetPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error {
	user, err := r.GetByID(ctx, userID.Hex())
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

// UpdateRoles 更新用户角色（admin特有）
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

// SearchUsers 搜索用户（admin特有）
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

// GetUsersByRole 根据角色获取用户列表（admin特有）
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

// CountByStatusMap 按状态统计用户数量（admin特有，返回map）
func (r *MongoUserAdminRepository) CountByStatusMap(ctx context.Context) (map[string]int64, error) {
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

// GetRecentUsers 获取最近注册的用户（admin特有）
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

// GetActiveUsers 获取活跃用户（admin特有）
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

// === 辅助方法 ===

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

func exactUserMatchRegex(value string) primitive.Regex {
	return primitive.Regex{
		Pattern: "^" + regexp.QuoteMeta(value) + "$",
		Options: "",
	}
}

func normalizeUserStatus(status users.UserStatus) (string, bool) {
	switch status {
	case users.UserStatusActive:
		return string(users.UserStatusActive), true
	case users.UserStatusInactive:
		return string(users.UserStatusInactive), true
	case users.UserStatusBanned:
		return string(users.UserStatusBanned), true
	case users.UserStatusDeleted:
		return string(users.UserStatusDeleted), true
	default:
		return "", false
	}
}

func normalizeUserRole(role string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case authModel.RoleReader:
		return authModel.RoleReader, true
	case authModel.RoleAuthor:
		return authModel.RoleAuthor, true
	case authModel.RoleAdmin:
		return authModel.RoleAdmin, true
	default:
		return "", false
	}
}

func userMatchesFilter(user *users.User, filter *adminrepo.UserFilter) bool {
	if user == nil || filter == nil {
		return true
	}

	if filter.Keyword != "" {
		keyword := strings.ToLower(strings.TrimSpace(filter.Keyword))
		if keyword != "" {
			username := strings.ToLower(user.Username)
			email := strings.ToLower(user.Email)
			nickname := strings.ToLower(user.Nickname)
			if !strings.Contains(username, keyword) &&
				!strings.Contains(email, keyword) &&
				!strings.Contains(nickname, keyword) {
				return false
			}
		}
	}

	if filter.Status != "" {
		if status, ok := normalizeUserStatus(filter.Status); !ok || string(user.Status) != status {
			return false
		}
	}

	if filter.Role != "" {
		if role, ok := normalizeUserRole(filter.Role); !ok || !userHasRole(user, role) {
			return false
		}
	}

	if filter.DateFrom != nil && user.CreatedAt.Before(*filter.DateFrom) {
		return false
	}
	if filter.DateTo != nil && user.CreatedAt.After(*filter.DateTo) {
		return false
	}
	if filter.LastActive != nil && user.LastLoginAt.Before(*filter.LastActive) {
		return false
	}

	return true
}

func userHasRole(user *users.User, role string) bool {
	for _, r := range user.Roles {
		if strings.EqualFold(strings.TrimSpace(r), role) {
			return true
		}
	}
	return false
}
