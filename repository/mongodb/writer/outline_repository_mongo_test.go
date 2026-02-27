package writer_test

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	writerInterfaces "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// setupOutlineRepo 设置大纲Repository测试环境
func setupOutlineRepo(t *testing.T) (writerInterfaces.OutlineRepository, *mongo.Database, context.Context, func()) {
	t.Helper()
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewOutlineRepository(db)
	ctx := context.Background()
	return repo, db, ctx, func() {
		// 清理outlines集合
		_ = db.Collection("outlines").Drop(ctx)
		cleanup()
	}
}

// createTestOutline 创建测试大纲节点
func createTestOutline(projectID, title, parentID string, order int) *writer.OutlineNode {
	outline := &writer.OutlineNode{}
	outline.ProjectID = projectID
	outline.Title = title
	outline.ParentID = parentID
	outline.Order = order
	outline.Summary = "测试大纲摘要"
	outline.Type = "chapter"
	outline.Tension = 5
	outline.Characters = []string{}
	outline.Items = []string{}
	return outline
}

// TestOutlineRepository_Create 测试创建大纲节点
func TestOutlineRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "第一章", "", 0)

	err := repo.Create(ctx, outline)
	require.NoError(t, err)
	assert.False(t, outline.ID.IsZero())
	assert.NotZero(t, outline.CreatedAt)
	assert.NotZero(t, outline.UpdatedAt)
	assert.Equal(t, "第一章", outline.Title)
}

// TestOutlineRepository_Create_WithParent 测试创建子节点
func TestOutlineRepository_Create_WithParent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建父节点
	parent := createTestOutline(projectID, "第一卷", "", 0)
	err := repo.Create(ctx, parent)
	require.NoError(t, err)

	// 创建子节点
	child := createTestOutline(projectID, "第一章", parent.ID.Hex(), 0)
	err = repo.Create(ctx, child)
	require.NoError(t, err)

	assert.Equal(t, parent.ID.Hex(), child.ParentID)
	assert.Equal(t, 0, child.Order)
}

// TestOutlineRepository_FindByID 测试根据ID查询
func TestOutlineRepository_FindByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "第一章", "", 0)
	err := repo.Create(ctx, outline)
	require.NoError(t, err)

	// 查询
	retrieved, err := repo.FindByID(ctx, outline.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, outline.ID, retrieved.ID)
	assert.Equal(t, outline.Title, retrieved.Title)
	assert.Equal(t, outline.ProjectID, retrieved.ProjectID)
}

// TestOutlineRepository_FindByID_NotFound 测试查询不存在的节点
func TestOutlineRepository_FindByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	_, err := repo.FindByID(ctx, primitive.NewObjectID().Hex())
	assert.Error(t, err)
}

// TestOutlineRepository_FindByProjectID 测试查询项目所有大纲
func TestOutlineRepository_FindByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()

	// 创建当前项目的大纲
	for i := 1; i <= 3; i++ {
		outline := createTestOutline(projectID, "第"+string(rune('0'+i))+"章", "", i-1)
		err := repo.Create(ctx, outline)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 创建其他项目的大纲
	otherOutline := createTestOutline(otherProjectID, "其他章节", "", 0)
	err := repo.Create(ctx, otherOutline)
	require.NoError(t, err)

	// 查询当前项目的大纲
	outlines, err := repo.FindByProjectID(ctx, projectID)
	require.NoError(t, err)
	assert.Len(t, outlines, 3)

	// 验证都属于当前项目
	for _, o := range outlines {
		assert.Equal(t, projectID, o.ProjectID)
	}
}

// TestOutlineRepository_Update 测试更新大纲
func TestOutlineRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "原标题", "", 0)
	err := repo.Create(ctx, outline)
	require.NoError(t, err)

	// 更新
	outline.Title = "新标题"
	outline.Summary = "新摘要"
	outline.Tension = 8

	oldUpdatedAt := outline.UpdatedAt
	time.Sleep(10 * time.Millisecond) // 确保时间差异

	err = repo.Update(ctx, outline)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.FindByID(ctx, outline.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "新标题", retrieved.Title)
	assert.Equal(t, "新摘要", retrieved.Summary)
	assert.Equal(t, 8, retrieved.Tension)
	assert.True(t, retrieved.UpdatedAt.After(oldUpdatedAt))
}

// TestOutlineRepository_Update_NotFound 测试更新不存在的节点
func TestOutlineRepository_Update_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "测试", "", 0)
	outline.ID = primitive.NewObjectID() // 不存在的ID

	err := repo.Update(ctx, outline)
	assert.Error(t, err)
}

// TestOutlineRepository_Delete 测试删除大纲
func TestOutlineRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "待删除", "", 0)
	err := repo.Create(ctx, outline)
	require.NoError(t, err)

	// 删除
	err = repo.Delete(ctx, outline.ID.Hex())
	require.NoError(t, err)

	// 验证已删除
	_, err = repo.FindByID(ctx, outline.ID.Hex())
	assert.Error(t, err)
}

// TestOutlineRepository_Delete_NotFound 测试删除不存在的节点
func TestOutlineRepository_Delete_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	err := repo.Delete(ctx, primitive.NewObjectID().Hex())
	assert.Error(t, err)
}

// TestOutlineRepository_FindByParentID 测试查询子节点
func TestOutlineRepository_FindByParentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建父节点
	parent := createTestOutline(projectID, "第一卷", "", 0)
	err := repo.Create(ctx, parent)
	require.NoError(t, err)

	// 创建多个子节点
	for i := 1; i <= 3; i++ {
		child := createTestOutline(projectID, "第"+string(rune('0'+i))+"章", parent.ID.Hex(), i-1)
		err := repo.Create(ctx, child)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 查询子节点
	children, err := repo.FindByParentID(ctx, projectID, parent.ID.Hex())
	require.NoError(t, err)
	assert.Len(t, children, 3)

	// 验证都是子节点
	for _, child := range children {
		assert.Equal(t, parent.ID.Hex(), child.ParentID)
		assert.Equal(t, projectID, child.ProjectID)
	}
}

// TestOutlineRepository_FindByParentID_Empty 测试查询空子节点列表
func TestOutlineRepository_FindByParentID_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	// 查询不存在的父节点的子节点
	children, err := repo.FindByParentID(ctx, projectID, parentID)
	require.NoError(t, err)
	assert.Len(t, children, 0)
}

// TestOutlineRepository_FindRoots 测试查询根节点
func TestOutlineRepository_FindRoots(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建根节点
	for i := 1; i <= 2; i++ {
		root := createTestOutline(projectID, "第"+string(rune('0'+i))+"卷", "", i-1)
		err := repo.Create(ctx, root)
		require.NoError(t, err)
	}

	// 创建一个子节点
	parent := createTestOutline(projectID, "第一卷", "", 0)
	err := repo.Create(ctx, parent)
	require.NoError(t, err)
	child := createTestOutline(projectID, "第一章", parent.ID.Hex(), 0)
	err = repo.Create(ctx, child)
	require.NoError(t, err)

	// 查询根节点
	roots, err := repo.FindRoots(ctx, projectID)
	require.NoError(t, err)
	assert.Len(t, roots, 3) // 2个显式根节点 + 1个父节点

	// 验证都是根节点
	for _, root := range roots {
		assert.Equal(t, "", root.ParentID)
		assert.Equal(t, projectID, root.ProjectID)
	}
}

// TestOutlineRepository_ExistsByID 测试检查存在性
func TestOutlineRepository_ExistsByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	outline := createTestOutline(projectID, "测试章节", "", 0)
	err := repo.Create(ctx, outline)
	require.NoError(t, err)

	// 检查存在的节点
	exists, err := repo.ExistsByID(ctx, outline.ID.Hex())
	require.NoError(t, err)
	assert.True(t, exists)

	// 检查不存在的节点
	exists, err = repo.ExistsByID(ctx, primitive.NewObjectID().Hex())
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestOutlineRepository_CountByProjectID 测试统计项目大纲数量
func TestOutlineRepository_CountByProjectID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()

	// 创建当前项目的大纲
	for i := 1; i <= 5; i++ {
		outline := createTestOutline(projectID, "第"+string(rune('0'+i))+"章", "", i-1)
		err := repo.Create(ctx, outline)
		require.NoError(t, err)
	}

	// 创建其他项目的大纲
	otherOutline := createTestOutline(otherProjectID, "其他章节", "", 0)
	err := repo.Create(ctx, otherOutline)
	require.NoError(t, err)

	// 统计当前项目
	count, err := repo.CountByProjectID(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// 统计其他项目
	count, err = repo.CountByProjectID(ctx, otherProjectID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestOutlineRepository_CountByParentID 测试统计子节点数量
func TestOutlineRepository_CountByParentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建父节点
	parent := createTestOutline(projectID, "第一卷", "", 0)
	err := repo.Create(ctx, parent)
	require.NoError(t, err)

	// 创建子节点
	for i := 1; i <= 3; i++ {
		child := createTestOutline(projectID, "第"+string(rune('0'+i))+"章", parent.ID.Hex(), i-1)
		err := repo.Create(ctx, child)
		require.NoError(t, err)
	}

	// 创建根节点
	root := createTestOutline(projectID, "第二卷", "", 1)
	err = repo.Create(ctx, root)
	require.NoError(t, err)

	// 统计子节点数量
	count, err := repo.CountByParentID(ctx, projectID, parent.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 统计根节点数量（parent_id为空）
	count, err = repo.CountByParentID(ctx, projectID, "")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count) // 1个parent + 1个root
}

// TestOutlineRepository_Ordering 测试排序
func TestOutlineRepository_Ordering(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建不同order的节点
	for i := 3; i >= 1; i-- {
		outline := createTestOutline(projectID, "第"+string(rune('0'+i))+"章", "", i)
		err := repo.Create(ctx, outline)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 查询根节点（应该按order排序）
	roots, err := repo.FindRoots(ctx, projectID)
	require.NoError(t, err)
	assert.Len(t, roots, 3)

	// 验证order顺序
	assert.Equal(t, 1, roots[0].Order)
	assert.Equal(t, 2, roots[1].Order)
	assert.Equal(t, 3, roots[2].Order)
}

// TestOutlineRepository_HierarchicalStructure 测试层级结构
func TestOutlineRepository_HierarchicalStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOutlineRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID().Hex()

	// 创建三层结构
	volume := createTestOutline(projectID, "第一卷", "", 0)
	err := repo.Create(ctx, volume)
	require.NoError(t, err)

	chapter1 := createTestOutline(projectID, "第一章", volume.ID.Hex(), 0)
	err = repo.Create(ctx, chapter1)
	require.NoError(t, err)

	section1 := createTestOutline(projectID, "第一节", chapter1.ID.Hex(), 0)
	err = repo.Create(ctx, section1)
	require.NoError(t, err)

	// 验证层级关系
	// 查询卷的子节点（章节）
	chapters, err := repo.FindByParentID(ctx, projectID, volume.ID.Hex())
	require.NoError(t, err)
	assert.Len(t, chapters, 1)
	assert.Equal(t, "第一章", chapters[0].Title)

	// 查询章的子节点（节）
	sections, err := repo.FindByParentID(ctx, projectID, chapter1.ID.Hex())
	require.NoError(t, err)
	assert.Len(t, sections, 1)
	assert.Equal(t, "第一节", sections[0].Title)
}
