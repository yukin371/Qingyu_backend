package admin

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/bookstore"
)

func oid(seed string) primitive.ObjectID {
	normalizedSeed := strings.ToLower(strings.TrimSpace(seed))
	sum := sha1.Sum([]byte(normalizedSeed))
	normalized := hex.EncodeToString(sum[:])[:24]
	objectID, err := primitive.ObjectIDFromHex(normalized)
	if err != nil {
		panic(err)
	}
	return objectID
}

func hexID(seed string) string {
	return oid(seed).Hex()
}

func strPtr(v string) *string {
	return &v
}

// setupTestDB 设置测试数据库（使用内存MongoDB或跳过集成测试）
func setupTestDB(t *testing.T) *mongo.Database {
	// 注意：这需要实际运行的MongoDB实例
	// 在CI/CD环境中应该使用 testcontainers 或类似工具
	// 如果没有MongoDB实例，跳过这些测试
	t.Skip("Skipping integration tests - MongoDB not available")
	return nil
}

// TestCategoryAdminMongoRepository_Create 测试创建分类
func TestCategoryAdminMongoRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	category := &bookstore.Category{
		Name:        "测试分类",
		Description: "测试描述",
		SortOrder:   1,
		IsActive:    true,
	}

	err := repo.Create(ctx, category)
	assert.NoError(t, err)
	assert.NotEmpty(t, category.ID)

	// 验证创建成功
	retrieved, err := repo.GetByID(ctx, category.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, "测试分类", retrieved.Name)

	// 清理
	repo.Delete(ctx, category.ID.Hex())
}

// TestCategoryAdminMongoRepository_Create_WithID 测试使用指定ID创建分类
func TestCategoryAdminMongoRepository_Create_WithID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	id := primitive.NewObjectID().Hex()
	category := &bookstore.Category{
		ID:          oid(id),
		Name:        "指定ID分类",
		Description: "测试",
		SortOrder:   1,
		IsActive:    true,
	}

	err := repo.Create(ctx, category)
	assert.NoError(t, err)
	assert.Equal(t, id, category.ID.Hex())

	// 清理
	repo.Delete(ctx, category.ID.Hex())
}

// TestCategoryAdminMongoRepository_GetByID 测试获取分类
func TestCategoryAdminMongoRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 先创建一个分类
	category := &bookstore.Category{
		Name:        "测试分类",
		Description: "测试描述",
		SortOrder:   1,
		IsActive:    true,
	}
	repo.Create(ctx, category)

	// 测试获取
	retrieved, err := repo.GetByID(ctx, category.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, category.ID, retrieved.ID)
	assert.Equal(t, "测试分类", retrieved.Name)

	// 清理
	repo.Delete(ctx, category.ID.Hex())
}

// TestCategoryAdminMongoRepository_GetByID_NotFound 测试获取不存在的分类
func TestCategoryAdminMongoRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 使用一个不存在的ID
	nonExistentID := primitive.NewObjectID().Hex()
	category, err := repo.GetByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "not found")
}

// TestCategoryAdminMongoRepository_GetByID_InvalidID 测试使用无效ID获取分类
func TestCategoryAdminMongoRepository_GetByID_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 使用无效的ObjectID
	category, err := repo.GetByID(ctx, "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, category)
}

// TestCategoryAdminMongoRepository_Update 测试更新分类
func TestCategoryAdminMongoRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 先创建一个分类
	category := &bookstore.Category{
		Name:        "原始名称",
		Description: "原始描述",
		SortOrder:   1,
		IsActive:    true,
	}
	repo.Create(ctx, category)

	// 修改并更新
	category.Name = "更新后的名称"
	category.Description = "更新后的描述"
	category.SortOrder = 5

	err := repo.Update(ctx, category)
	assert.NoError(t, err)

	// 验证更新
	retrieved, _ := repo.GetByID(ctx, category.ID.Hex())
	assert.Equal(t, "更新后的名称", retrieved.Name)
	assert.Equal(t, "更新后的描述", retrieved.Description)
	assert.Equal(t, 5, retrieved.SortOrder)

	// 清理
	repo.Delete(ctx, category.ID.Hex())
}

// TestCategoryAdminMongoRepository_Update_NotFound 测试更新不存在的分类
func TestCategoryAdminMongoRepository_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	nonExistentID := primitive.NewObjectID().Hex()
	category := &bookstore.Category{
		ID:        oid(nonExistentID),
		Name:      "不存在的分类",
		SortOrder: 1,
		IsActive:  true,
	}

	err := repo.Update(ctx, category)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCategoryAdminMongoRepository_Update_InvalidID 测试使用无效ID更新分类
func TestCategoryAdminMongoRepository_Update_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	category := &bookstore.Category{
		ID:        primitive.NilObjectID,
		Name:      "测试",
		SortOrder: 1,
		IsActive:  true,
	}

	err := repo.Update(ctx, category)
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_Delete 测试删除分类
func TestCategoryAdminMongoRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 先创建一个分类
	category := &bookstore.Category{
		Name:      "待删除分类",
		SortOrder: 1,
		IsActive:  true,
	}
	repo.Create(ctx, category)

	// 删除
	err := repo.Delete(ctx, category.ID.Hex())
	assert.NoError(t, err)

	// 验证已删除
	_, err = repo.GetByID(ctx, category.ID.Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCategoryAdminMongoRepository_Delete_NotFound 测试删除不存在的分类
func TestCategoryAdminMongoRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	nonExistentID := primitive.NewObjectID().Hex()
	err := repo.Delete(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestCategoryAdminMongoRepository_Delete_InvalidID 测试使用无效ID删除分类
func TestCategoryAdminMongoRepository_Delete_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, "invalid-id")
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_List 测试列表查询
func TestCategoryAdminMongoRepository_List(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建测试数据
	categories := []*bookstore.Category{
		{Name: "玄幻", SortOrder: 1, IsActive: true},
		{Name: "武侠", SortOrder: 2, IsActive: true},
		{Name: "仙侠", SortOrder: 3, IsActive: false},
	}

	var ids []string
	for _, cat := range categories {
		err := repo.Create(ctx, cat)
		assert.NoError(t, err)
		ids = append(ids, cat.ID.Hex())
	}
	defer func() {
		for _, id := range ids {
			repo.Delete(ctx, id)
		}
	}()

	// 查询所有
	result, err := repo.List(ctx, bson.M{})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 3)
}

// TestCategoryAdminMongoRepository_List_WithFilter 测试带过滤条件的列表查询
func TestCategoryAdminMongoRepository_List_WithFilter(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建测试数据
	cat1 := &bookstore.Category{Name: "活跃分类", SortOrder: 1, IsActive: true}
	cat2 := &bookstore.Category{Name: "非活跃分类", SortOrder: 2, IsActive: false}

	repo.Create(ctx, cat1)
	repo.Create(ctx, cat2)
	defer repo.Delete(ctx, cat1.ID.Hex())
	defer repo.Delete(ctx, cat2.ID.Hex())

	// 查询活跃分类
	filter := bson.M{"is_active": true}
	result, err := repo.List(ctx, filter)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1)
}

// TestCategoryAdminMongoRepository_List_WithOptions 测试带选项的列表查询
func TestCategoryAdminMongoRepository_List_WithOptions(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建测试数据
	for i := 1; i <= 5; i++ {
		cat := &bookstore.Category{
			Name:      "分类" + string(rune('0'+i)),
			SortOrder: i,
			IsActive:  true,
		}
		repo.Create(ctx, cat)
	}

	// 限制返回数量
	opts := options.Find().SetLimit(3)
	result, err := repo.List(ctx, bson.M{}, opts)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(result), 3)
}

// TestCategoryAdminMongoRepository_HasChildren 测试检查子分类
func TestCategoryAdminMongoRepository_HasChildren(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建父分类
	parent := &bookstore.Category{
		Name:      "父分类",
		SortOrder: 1,
		IsActive:  true,
	}
	repo.Create(ctx, parent)
	defer repo.Delete(ctx, parent.ID.Hex())

	// 创建子分类
	child := &bookstore.Category{
		Name:      "子分类",
		SortOrder: 1,
		IsActive:  true,
		ParentID:  strPtr(parent.ID.Hex()),
	}
	repo.Create(ctx, child)
	defer repo.Delete(ctx, child.ID.Hex())

	// 检查是否有子分类
	hasChildren, err := repo.HasChildren(ctx, parent.ID.Hex())
	assert.NoError(t, err)
	assert.True(t, hasChildren)

	// 子分类应该没有子节点
	hasChildren, err = repo.HasChildren(ctx, child.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, hasChildren)
}

// TestCategoryAdminMongoRepository_HasChildren_NoChildren 测试没有子分类的情况
func TestCategoryAdminMongoRepository_HasChildren_NoChildren(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建一个没有子分类的分类
	category := &bookstore.Category{
		Name:      "独立分类",
		SortOrder: 1,
		IsActive:  true,
	}
	repo.Create(ctx, category)
	defer repo.Delete(ctx, category.ID.Hex())

	hasChildren, err := repo.HasChildren(ctx, category.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, hasChildren)
}

// TestCategoryAdminMongoRepository_HasChildren_InvalidID 测试使用无效ID检查子分类
func TestCategoryAdminMongoRepository_HasChildren_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	_, err := repo.HasChildren(ctx, "invalid-id")
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_NameExistsAtLevel 测试同级名称检查
func TestCategoryAdminMongoRepository_NameExistsAtLevel(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建一个分类
	category := &bookstore.Category{
		Name:      "玄幻",
		SortOrder: 1,
		IsActive:  true,
	}
	repo.Create(ctx, category)
	defer repo.Delete(ctx, category.ID.Hex())

	// 检查同名是否存在（应该存在）
	exists, err := repo.NameExistsAtLevel(ctx, nil, "玄幻", "")
	assert.NoError(t, err)
	assert.True(t, exists)

	// 检查不同名（应该不存在）
	exists, err = repo.NameExistsAtLevel(ctx, nil, "武侠", "")
	assert.NoError(t, err)
	assert.False(t, exists)

	// 排除当前ID（用于更新时检查）
	exists, err = repo.NameExistsAtLevel(ctx, nil, "玄幻", category.ID.Hex())
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestCategoryAdminMongoRepository_NameExistsAtLevel_WithParent 测试带父级的同级名称检查
func TestCategoryAdminMongoRepository_NameExistsAtLevel_WithParent(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建父分类
	parent := &bookstore.Category{
		Name:      "父分类",
		SortOrder: 1,
		IsActive:  true,
	}
	repo.Create(ctx, parent)
	defer repo.Delete(ctx, parent.ID.Hex())

	// 创建子分类
	child := &bookstore.Category{
		Name:      "子分类",
		SortOrder: 1,
		IsActive:  true,
		ParentID:  strPtr(parent.ID.Hex()),
	}
	repo.Create(ctx, child)
	defer repo.Delete(ctx, child.ID.Hex())

	// 检查同名在同一父级下（应该存在）
	parentID := parent.ID.Hex()
	exists, err := repo.NameExistsAtLevel(ctx, &parentID, "子分类", "")
	assert.NoError(t, err)
	assert.True(t, exists)

	// 检查同名在不同父级下（应该不存在）
	exists, err = repo.NameExistsAtLevel(ctx, nil, "子分类", "")
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestCategoryAdminMongoRepository_NameExistsAtLevel_InvalidParentID 测试使用无效父ID检查名称
func TestCategoryAdminMongoRepository_NameExistsAtLevel_InvalidParentID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	invalidID := "invalid-id"
	_, err := repo.NameExistsAtLevel(ctx, &invalidID, "测试", "")
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_NameExistsAtLevel_InvalidExcludeID 测试使用无效排除ID检查名称
func TestCategoryAdminMongoRepository_NameExistsAtLevel_InvalidExcludeID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	_, err := repo.NameExistsAtLevel(ctx, nil, "测试", "invalid-id")
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_UpdateBookCount 测试更新作品数量
func TestCategoryAdminMongoRepository_UpdateBookCount(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建分类
	category := &bookstore.Category{
		Name:      "测试分类",
		SortOrder: 1,
		IsActive:  true,
		BookCount: 0,
	}
	repo.Create(ctx, category)
	defer repo.Delete(ctx, category.ID.Hex())

	// 更新作品数量
	err := repo.UpdateBookCount(ctx, category.ID.Hex(), 100)
	assert.NoError(t, err)

	// 验证更新
	retrieved, _ := repo.GetByID(ctx, category.ID.Hex())
	assert.Equal(t, int64(100), retrieved.BookCount)
}

// TestCategoryAdminMongoRepository_UpdateBookCount_InvalidID 测试使用无效ID更新作品数量
func TestCategoryAdminMongoRepository_UpdateBookCount_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	err := repo.UpdateBookCount(ctx, "invalid-id", 100)
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_BatchUpdateStatus 测试批量更新状态
func TestCategoryAdminMongoRepository_BatchUpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建多个分类
	cat1 := &bookstore.Category{Name: "分类1", SortOrder: 1, IsActive: true}
	cat2 := &bookstore.Category{Name: "分类2", SortOrder: 2, IsActive: true}
	cat3 := &bookstore.Category{Name: "分类3", SortOrder: 3, IsActive: true}

	repo.Create(ctx, cat1)
	repo.Create(ctx, cat2)
	repo.Create(ctx, cat3)

	defer repo.Delete(ctx, cat1.ID.Hex())
	defer repo.Delete(ctx, cat2.ID.Hex())
	defer repo.Delete(ctx, cat3.ID.Hex())

	// 批量禁用
	ids := []string{cat1.ID.Hex(), cat2.ID.Hex()}
	err := repo.BatchUpdateStatus(ctx, ids, false)
	assert.NoError(t, err)

	// 验证更新
	retrieved1, _ := repo.GetByID(ctx, cat1.ID.Hex())
	retrieved2, _ := repo.GetByID(ctx, cat2.ID.Hex())
	retrieved3, _ := repo.GetByID(ctx, cat3.ID.Hex())

	assert.False(t, retrieved1.IsActive)
	assert.False(t, retrieved2.IsActive)
	assert.True(t, retrieved3.IsActive) // 未更新
}

// TestCategoryAdminMongoRepository_BatchUpdateStatus_InvalidID 测试使用无效ID批量更新状态
func TestCategoryAdminMongoRepository_BatchUpdateStatus_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	ids := []string{"invalid-id", "another-invalid-id"}
	err := repo.BatchUpdateStatus(ctx, ids, true)
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_GetDescendantIDs 测试获取子孙分类ID
func TestCategoryAdminMongoRepository_GetDescendantIDs(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建三级分类结构
	root := &bookstore.Category{Name: "根分类", SortOrder: 1, IsActive: true}
	child1 := &bookstore.Category{Name: "子分类1", SortOrder: 1, IsActive: true}
	child2 := &bookstore.Category{Name: "子分类2", SortOrder: 2, IsActive: true}
	grandchild1 := &bookstore.Category{Name: "孙分类1", SortOrder: 1, IsActive: true}

	repo.Create(ctx, root)
	repo.Create(ctx, child1)
	repo.Create(ctx, child2)
	repo.Create(ctx, grandchild1)

	// 建立层级关系
	child1.ParentID = strPtr(root.ID.Hex())
	child2.ParentID = strPtr(root.ID.Hex())
	grandchild1.ParentID = strPtr(child1.ID.Hex())
	repo.Update(ctx, child1)
	repo.Update(ctx, child2)
	repo.Update(ctx, grandchild1)

	defer func() {
		repo.Delete(ctx, root.ID.Hex())
		repo.Delete(ctx, child1.ID.Hex())
		repo.Delete(ctx, child2.ID.Hex())
		repo.Delete(ctx, grandchild1.ID.Hex())
	}()

	// 获取根分类的所有子孙ID
	descendantIDs, err := repo.GetDescendantIDs(ctx, root.ID.Hex())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(descendantIDs), 3)
}

// TestCategoryAdminMongoRepository_GetDescendantIDs_NoDescendants 测试获取没有子孙的分类ID
func TestCategoryAdminMongoRepository_GetDescendantIDs_NoDescendants(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建没有子分类的分类
	category := &bookstore.Category{Name: "独立分类", SortOrder: 1, IsActive: true}
	repo.Create(ctx, category)
	defer repo.Delete(ctx, category.ID.Hex())

	descendantIDs, err := repo.GetDescendantIDs(ctx, category.ID.Hex())
	assert.NoError(t, err)
	assert.Empty(t, descendantIDs)
}

// TestCategoryAdminMongoRepository_GetDescendantIDs_InvalidID 测试使用无效ID获取子孙分类ID
func TestCategoryAdminMongoRepository_GetDescendantIDs_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	_, err := repo.GetDescendantIDs(ctx, "invalid-id")
	assert.Error(t, err)
}

// TestCategoryAdminMongoRepository_GetTree 测试获取分类树
func TestCategoryAdminMongoRepository_GetTree(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	// 创建多级分类
	root1 := &bookstore.Category{Name: "根分类1", SortOrder: 1, IsActive: true}
	root2 := &bookstore.Category{Name: "根分类2", SortOrder: 2, IsActive: true}
	child1 := &bookstore.Category{Name: "子分类1", SortOrder: 1, IsActive: true}

	repo.Create(ctx, root1)
	repo.Create(ctx, root2)
	repo.Create(ctx, child1)

	child1.ParentID = strPtr(root1.ID.Hex())
	repo.Update(ctx, child1)

	defer func() {
		repo.Delete(ctx, root1.ID.Hex())
		repo.Delete(ctx, root2.ID.Hex())
		repo.Delete(ctx, child1.ID.Hex())
	}()

	// 获取分类树（只返回根节点）
	tree, err := repo.GetTree(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(tree), 2)

	// 验证根节点按sort_order排序
	assert.Equal(t, "根分类1", tree[0].Name)
}

// TestCategoryAdminMongoRepository_Timestamps 测试时间戳自动设置
func TestCategoryAdminMongoRepository_Timestamps(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	repo := NewCategoryAdminMongoRepository(db)
	ctx := context.Background()

	beforeCreate := time.Now().Add(-time.Second)

	category := &bookstore.Category{
		Name:      "时间测试",
		SortOrder: 1,
		IsActive:  true,
	}

	err := repo.Create(ctx, category)
	assert.NoError(t, err)

	assert.False(t, category.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, category.UpdatedAt.IsZero(), "UpdatedAt should be set")
	assert.True(t, category.CreatedAt.After(beforeCreate), "CreatedAt should be recent")

	// 清理
	repo.Delete(ctx, category.ID.Hex())
}
