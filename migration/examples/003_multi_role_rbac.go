package examples

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"Qingyu_backend/models/auth"
)

// MultiRoleRBAC 多角色RBAC迁移
// 将用户从单角色系统迁移到多角色系统，并添加VIP等级
type MultiRoleRBAC struct{}

func (m *MultiRoleRBAC) Version() string {
	return "003"
}

func (m *MultiRoleRBAC) Description() string {
	return "Migrate from single role to multi-role RBAC system and add VIP level"
}

// Up 执行迁移：将单角色改为多角色，添加VIP等级
func (m *MultiRoleRBAC) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 分批处理所有用户
	batchSize := int64(100)
	var lastID string

	for {
		// 查询一批用户
		filter := bson.M{}
		if lastID != "" {
			filter = bson.M{"_id": bson.M{"$gt": lastID}}
		}

		opts := options.Find().
			SetSort(bson.D{{Key: "_id", Value: 1}}).
			SetLimit(batchSize)

		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return fmt.Errorf("failed to find users: %w", err)
		}

		var users []bson.M
		if err = cursor.All(ctx, &users); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("failed to decode users: %w", err)
		}
		cursor.Close(ctx)

		if len(users) == 0 {
			break // 没有更多用户了
		}

		// 批量更新用户
		for _, user := range users {
			userID := user["_id"].(string)
			lastID = userID // 记录最后一个用户ID
			updateDoc, err := m.buildUpdateDoc(user)
			if err != nil {
				fmt.Printf("  Warning: Failed to build update for user %s: %v\n", userID, err)
				continue
			}

			update := bson.M{"$set": updateDoc}
			if _, err := collection.UpdateByID(ctx, userID, update); err != nil {
				fmt.Printf("  Warning: Failed to update user %s: %v\n", userID, err)
				continue
			}
		}

		fmt.Printf("  ✓ Processed batch of %d users\n", len(users))
	}

	fmt.Println("  ✓ Migration completed: Single role → Multi role + VIP level")
	return nil
}

// Down 回滚迁移：将多角色改为单角色，移除VIP等级
func (m *MultiRoleRBAC) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 分批处理所有用户
	batchSize := int64(100)

	for {
		// 查询一批拥有roles字段的用户
		filter := bson.M{"roles": bson.M{"$exists": true}}

		cursor, err := collection.Find(ctx, filter, options.Find().SetLimit(batchSize))
		if err != nil {
			return fmt.Errorf("failed to find users: %w", err)
		}

		var users []bson.M
		if err = cursor.All(ctx, &users); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("failed to decode users: %w", err)
		}
		cursor.Close(ctx)

		if len(users) == 0 {
			break
		}

		// 批量回滚用户
		for _, user := range users {
			userID := user["_id"].(string)
			rollbackDoc, err := m.buildRollbackDoc(user)
			if err != nil {
				fmt.Printf("  Warning: Failed to build rollback for user %s: %v\n", userID, err)
				continue
			}

			update := bson.M{"$set": rollbackDoc}
			if _, err := collection.UpdateByID(ctx, userID, update); err != nil {
				fmt.Printf("  Warning: Failed to rollback user %s: %v\n", userID, err)
				continue
			}
		}

		fmt.Printf("  ✓ Rolled back batch of %d users\n", len(users))
	}

	fmt.Println("  ✓ Rollback completed: Multi role → Single role")
	return nil
}

// buildUpdateDoc 根据用户当前数据构建更新文档
func (m *MultiRoleRBAC) buildUpdateDoc(user bson.M) (bson.M, error) {
	update := bson.M{}

	// 处理角色转换
	oldRole, hasRole := user["role"]
	if hasRole {
		roles := m.convertOldRoleToNewRoles(oldRole.(string))
		update["roles"] = roles
	} else {
		// 如果没有role字段，默认为reader
		update["roles"] = []string{auth.RoleReader}
	}

	// 处理VIP等级
	vipLevel, hasVIP := user["vip_level"]
	if hasVIP {
		// 已有vip_level字段，确保值在0-5范围内
		if level, ok := vipLevel.(int32); ok {
			update["vip_level"] = m.normalizeVIPLevel(int(level))
		} else if level, ok := vipLevel.(int); ok {
			update["vip_level"] = m.normalizeVIPLevel(level)
		}
	} else {
		// 没有vip_level字段，默认为0
		update["vip_level"] = 0
	}

	return update, nil
}

// buildRollbackDoc 构建回滚文档
func (m *MultiRoleRBAC) buildRollbackDoc(user bson.M) (bson.M, error) {
	update := bson.M{}

	// 从多角色转为单角色
	roles, hasRoles := user["roles"]
	if hasRoles {
		roleArray := roles.(bson.A)
		primaryRole := m.getPrimaryRoleFromRoles(roleArray)
		update["role"] = primaryRole
	} else {
		update["role"] = "user"
	}

	return update, nil
}

// convertOldRoleToNewRoles 将旧角色转换为新角色数组
func (m *MultiRoleRBAC) convertOldRoleToNewRoles(oldRole string) []string {
	switch oldRole {
	case "user":
		return []string{auth.RoleReader}
	case "author":
		return []string{auth.RoleReader, auth.RoleAuthor}
	case "admin":
		return []string{auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin}
	default:
		return []string{auth.RoleReader}
	}
}

// getPrimaryRoleFromRoles 从角色数组中获取主要角色
func (m *MultiRoleRBAC) getPrimaryRoleFromRoles(roles bson.A) string {
	roleSet := make(map[string]bool)
	for _, r := range roles {
		if str, ok := r.(string); ok {
			roleSet[str] = true
		}
	}

	// 如果有admin角色，返回admin
	if roleSet[auth.RoleAdmin] {
		return "admin"
	}

	// 如果有author角色，返回author
	if roleSet[auth.RoleAuthor] {
		return "author"
	}

	// 默认返回user
	return "user"
}

// normalizeVIPLevel 标准化VIP等级到0-5范围
func (m *MultiRoleRBAC) normalizeVIPLevel(level int) int {
	if level < 0 {
		return 0
	}
	if level > 5 {
		return 5
	}
	return level
}
