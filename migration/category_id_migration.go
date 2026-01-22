package migration

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CategoryIDMigration_v1 分类模型 ID 类型迁移
//
// 迁移目标：
// - 将 Category.ID 从 primitive.ObjectID 迁移到 string
// - 将 Category.ParentID 从 *primitive.ObjectID 迁移到 *string
//
// 迁移策略：
// 1. 读取所有现有分类数据
// 2. 将 ObjectID 转换为十六进制字符串
// 3. 更新数据库中的文档
// 4. 验证迁移结果
//
// 回滚策略：
// - 保留原始 _id 字段
// - 可通过 Hex() -> ObjectID 反向转换

type CategoryIDMigration_v1 struct {
	db         *mongo.Database
	categories *mongo.Collection
}

// NewCategoryIDMigration_v1 创建迁移实例
func NewCategoryIDMigration_v1(db *mongo.Database) *CategoryIDMigration_v1 {
	return &CategoryIDMigration_v1{
		db:         db,
		categories: db.Collection("categories"),
	}
}

// Name 返回迁移名称
func (m *CategoryIDMigration_v1) Name() string {
	return "category_id_type_to_string_v1"
}

// Description 返回迁移描述
func (m *CategoryIDMigration_v1) Description() string {
	return "将 Category 模型的 ID 字段类型从 ObjectID 转换为 string (Hex)"
}

// Up 执行迁移
func (m *CategoryIDMigration_v1) Up(ctx context.Context) error {
	log.Println("[Migration] 开始执行 Category ID 类型迁移...")

	// 1. 统计需要迁移的数据量
	total, err := m.categories.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计分类数量失败: %w", err)
	}

	if total == 0 {
		log.Println("[Migration] 没有需要迁移的分类数据")
		return nil
	}

	log.Printf("[Migration] 找到 %d 个分类需要迁移\n", total)

	// 2. 批量读取所有分类
	cursor, err := m.categories.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("读取分类数据失败: %w", err)
	}
	defer cursor.Close(ctx)

	type OldCategory struct {
		ID       primitive.ObjectID  `bson:"_id"`
		ParentID *primitive.ObjectID `bson:"parent_id,omitempty"`
		Name     string              `bson:"name"`
	}

	var oldCategories []OldCategory
	if err = cursor.All(ctx, &oldCategories); err != nil {
		return fmt.Errorf("解析分类数据失败: %w", err)
	}

	log.Printf("[Migration] 成功读取 %d 个分类\n", len(oldCategories))

	// 3. 批量更新（使用 BulkWrite）
	var writes []mongo.WriteModel
	for _, cat := range oldCategories {
		// 构建更新操作
		update := bson.M{
			"$set": bson.M{
				"_id":        cat.ID.Hex(), // 将 ObjectID 转换为 Hex 字符串
				"parent_id":  convertParentIDToString(cat.ParentID),
				"migrated":   true, // 标记为已迁移
				"migrated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		}

		// 使用原始 _id 作为查询条件
		filter := bson.M{"_id": cat.ID}
		writes = append(writes, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	// 4. 执行批量更新
	if len(writes) > 0 {
		result, err := m.categories.BulkWrite(ctx, writes)
		if err != nil {
			return fmt.Errorf("批量更新失败: %w", err)
		}

		log.Printf("[Migration] 成功迁移 %d 个分类\n", result.MatchedCount)
	}

	// 5. 验证迁移结果
	newTotal, err := m.categories.CountDocuments(ctx, bson.M{"migrated": true})
	if err != nil {
		return fmt.Errorf("验证迁移结果失败: %w", err)
	}

	if newTotal != total {
		return fmt.Errorf("迁移验证失败: 预期 %d 条，实际 %d 条", total, newTotal)
	}

	log.Println("[Migration] Category ID 类型迁移完成！")
	return nil
}

// Down 回滚迁移
func (m *CategoryIDMigration_v1) Down(ctx context.Context) error {
	log.Println("[Migration] 开始回滚 Category ID 类型迁移...")

	// 1. 查找所有已迁移的分类
	cursor, err := m.categories.Find(ctx, bson.M{"migrated": true})
	if err != nil {
		return fmt.Errorf("查找已迁移分类失败: %w", err)
	}
	defer cursor.Close(ctx)

	type MigratedCategory struct {
		ID       string `bson:"_id"`
		ParentID *string `bson:"parent_id,omitempty"`
	}

	var migratedCategories []MigratedCategory
	if err = cursor.All(ctx, &migratedCategories); err != nil {
		return fmt.Errorf("解析已迁移分类失败: %w", err)
	}

	log.Printf("[Migration] 找到 %d 个已迁移分类需要回滚\n", len(migratedCategories))

	// 2. 批量回滚
	var writes []mongo.WriteModel
	for _, cat := range migratedCategories {
		// 将字符串 ID 转换回 ObjectID
		objectID, err := primitive.ObjectIDFromHex(cat.ID)
		if err != nil {
			log.Printf("[Migration] 警告: 无法转换 ID %s: %v，跳过\n", cat.ID, err)
			continue
		}

		update := bson.M{
			"$set": bson.M{
				"_id":        objectID,
				"parent_id":  convertParentIDToObjectID(cat.ParentID),
			},
			"$unset": bson.M{
				"migrated":   "",
				"migrated_at": "",
			},
		}

		// 使用字符串 _id 作为查询条件
		filter := bson.M{"_id": cat.ID}
		writes = append(writes, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}

	// 3. 执行批量回滚
	if len(writes) > 0 {
		result, err := m.categories.BulkWrite(ctx, writes)
		if err != nil {
			return fmt.Errorf("批量回滚失败: %w", err)
		}

		log.Printf("[Migration] 成功回滚 %d 个分类\n", result.MatchedCount)
	}

	log.Println("[Migration] Category ID 类型迁移回滚完成！")
	return nil
}

// convertParentIDToString 将 *ObjectID 转换为 *string
func convertParentIDToString(parentID *primitive.ObjectID) *string {
	if parentID == nil {
		return nil
	}
	hex := parentID.Hex()
	return &hex
}

// convertParentIDToObjectID 将 *string 转换为 *ObjectID
func convertParentIDToObjectID(parentID *string) *primitive.ObjectID {
	if parentID == nil || *parentID == "" {
		return nil
	}
	objectID, err := primitive.ObjectIDFromHex(*parentID)
	if err != nil {
		log.Printf("[Migration] 警告: 无法转换 ParentID %s: %v\n", *parentID, err)
		return nil
	}
	return &objectID
}

// Validate 验证迁移状态
func (m *CategoryIDMigration_v1) Validate(ctx context.Context) error {
	// 检查是否有未迁移的数据
	unmigrated, err := m.categories.CountDocuments(ctx, bson.M{"migrated": bson.M{"$ne": true}})
	if err != nil {
		return fmt.Errorf("验证失败: %w", err)
	}

	if unmigrated > 0 {
		return fmt.Errorf("发现 %d 个未迁移的分类", unmigrated)
	}

	// 检查 _id 格式
	cursor, err := m.categories.Find(ctx, bson.M{"migrated": true})
	if err != nil {
		return fmt.Errorf("查找已迁移分类失败: %w", err)
	}
	defer cursor.Close(ctx)

	var invalidCount int
	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			continue
		}

		if id, ok := result["_id"].(string); ok {
			if _, err = primitive.ObjectIDFromHex(id); err != nil {
				invalidCount++
			}
		}
	}

	if invalidCount > 0 {
		return fmt.Errorf("发现 %d 个无效的 ID 格式", invalidCount)
	}

	log.Println("[Migration] 验证通过，所有分类 ID 格式正确")
	return nil
}
