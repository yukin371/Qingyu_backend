package writing_test

import (
	"Qingyu_backend/models/writer"
	writingInterface "Qingyu_backend/repository/interfaces/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/test/testutil"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 测试辅助函数
func setupProjectRepo(t *testing.T) (writingInterface.ProjectRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewMongoProjectRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

func createTestProject(authorID, title string) *writer.Project {
	project := &writer.Project{
		Summary:    "Test project summary",
		Status:     writer.StatusDraft,
		Visibility: writer.VisibilityPrivate,
		Tags:       []string{"test", "novel"},
		Category:   "fantasy",
	}
	project.AuthorID = authorID
	project.Title = title
	return project
}

// 1. 测试项目创建
func TestProjectRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "My First Novel")

	err := repo.Create(ctx, project)
	require.NoError(t, err)
	assert.NotEmpty(t, project.ID)
	assert.NotZero(t, project.CreatedAt)
	assert.NotZero(t, project.UpdatedAt)
	assert.Equal(t, writer.StatusDraft, project.Status)
	assert.Equal(t, writer.VisibilityPrivate, project.Visibility)
	assert.True(t, project.Settings.AutoBackup)
	assert.Equal(t, 24, project.Settings.BackupInterval)
}

// 2. 测试创建空项目（参数验证）
func TestProjectRepository_Create_NilProject(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	err := repo.Create(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目对象不能为空")
}

// 3. 测试创建缺少必需字段的项目
func TestProjectRepository_Create_MissingFields(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	tests := []struct {
		name    string
		project *writer.Project
		errMsg  string
	}{
		{
			name: "Missing AuthorID",
			project: func() *writer.Project {
				p := &writer.Project{}
				p.Title = "Test Project"
				return p
			}(),
			errMsg: "作者ID不能为空",
		},
		{
			name: "Missing Title",
			project: func() *writer.Project {
				p := &writer.Project{}
				p.AuthorID = "author123"
				return p
			}(),
			errMsg: "项目标题不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.project)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

// 4. 测试根据ID获取项目
func TestProjectRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建项目
	project := createTestProject("author123", "My Novel")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 获取项目
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, project.ID, retrieved.ID)
	assert.Equal(t, project.Title, retrieved.Title)
	assert.Equal(t, project.AuthorID, retrieved.AuthorID)
}

// 5. 测试获取不存在的项目
func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	retrieved, err := repo.GetByID(ctx, primitive.NewObjectID().Hex())
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 6. 测试获取无效ID（字符串ID无需验证，直接查询）
func TestProjectRepository_GetByID_InvalidID(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 使用不存在的ID，应该返回nil而不是错误
	retrieved, err := repo.GetByID(ctx, "non-existent-id")
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 7. 测试更新项目
func TestProjectRepository_Update(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建项目
	project := createTestProject("author123", "Original Title")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 更新项目
	updates := map[string]interface{}{
		"title":   "Updated Title",
		"summary": "Updated Summary",
		"status":  writer.StatusSerializing,
	}

	err = repo.Update(ctx, project.ID, updates)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", retrieved.Title)
	assert.Equal(t, "Updated Summary", retrieved.Summary)
	assert.Equal(t, writer.StatusSerializing, retrieved.Status)
}

// 8. 测试更新不存在的项目
func TestProjectRepository_Update_NotFound(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	updates := map[string]interface{}{"title": "New Title"}
	err := repo.Update(ctx, primitive.NewObjectID().Hex(), updates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目不存在或已删除")
}

// 9. 测试删除项目
func TestProjectRepository_Delete(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建项目
	project := createTestProject("author123", "To Be Deleted")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 删除项目
	err = repo.Delete(ctx, project.ID)
	require.NoError(t, err)

	// 验证已删除
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 10. 测试根据所有者获取项目列表
func TestProjectRepository_GetListByOwnerID(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建多个项目
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "Project "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // 确保时间差异
	}

	// 创建其他作者的项目
	otherProject := createTestProject("other_author", "Other Project")
	err := repo.Create(ctx, otherProject)
	require.NoError(t, err)

	// 查询
	projects, err := repo.GetListByOwnerID(ctx, authorID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, projects, 3)
	assert.Equal(t, authorID, projects[0].AuthorID)
}

// 11. 测试根据所有者和状态获取项目
func TestProjectRepository_GetByOwnerAndStatus(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建不同状态的项目
	draftProject := createTestProject(authorID, "Draft Project")
	draftProject.Status = writer.StatusDraft
	err := repo.Create(ctx, draftProject)
	require.NoError(t, err)

	serializingProject := createTestProject(authorID, "Serializing Project")
	serializingProject.Status = writer.StatusSerializing
	err = repo.Create(ctx, serializingProject)
	require.NoError(t, err)

	// 查询草稿状态项目
	projects, err := repo.GetByOwnerAndStatus(ctx, authorID, string(writer.StatusDraft), 10, 0)
	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, writer.StatusDraft, projects[0].Status)
}

// 12. 测试更新项目（根据所有者）
func TestProjectRepository_UpdateByOwner(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"
	project := createTestProject(authorID, "Owner Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 所有者更新
	updates := map[string]interface{}{"title": "Owner Updated"}
	err = repo.UpdateByOwner(ctx, project.ID, authorID, updates)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, "Owner Updated", retrieved.Title)
}

// 13. 测试非所有者更新失败
func TestProjectRepository_UpdateByOwner_NotOwner(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "Owner Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 非所有者尝试更新
	updates := map[string]interface{}{"title": "Hacked"}
	err = repo.UpdateByOwner(ctx, project.ID, "hacker", updates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目不存在或无权限")
}

// 14. 测试检查所有者
func TestProjectRepository_IsOwner(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"
	project := createTestProject(authorID, "Owner Check Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 验证所有者
	isOwner, err := repo.IsOwner(ctx, project.ID, authorID)
	require.NoError(t, err)
	assert.True(t, isOwner)

	// 验证非所有者
	isOwner, err = repo.IsOwner(ctx, project.ID, "other_user")
	require.NoError(t, err)
	assert.False(t, isOwner)
}

// 15. 测试软删除
func TestProjectRepository_SoftDelete(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"
	project := createTestProject(authorID, "Soft Delete Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 软删除
	err = repo.SoftDelete(ctx, project.ID, authorID)
	require.NoError(t, err)

	// 普通查询应找不到
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 16. 测试硬删除
func TestProjectRepository_HardDelete(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "Hard Delete Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 硬删除
	err = repo.HardDelete(ctx, project.ID)
	require.NoError(t, err)

	// 验证完全不存在
	exists, err := repo.Exists(ctx, project.ID)
	require.NoError(t, err)
	assert.False(t, exists)
}

// 17. 测试恢复已删除项目
func TestProjectRepository_Restore(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"
	project := createTestProject(authorID, "Restore Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 软删除
	err = repo.SoftDelete(ctx, project.ID, authorID)
	require.NoError(t, err)

	// 恢复
	err = repo.Restore(ctx, project.ID, authorID)
	require.NoError(t, err)

	// 验证可以查询到
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, project.Title, retrieved.Title)
}

// 18. 测试根据所有者统计项目数
func TestProjectRepository_CountByOwner(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建3个项目
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "Count Project "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 统计
	count, err := repo.CountByOwner(ctx, authorID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// 19. 测试根据状态统计项目数
func TestProjectRepository_CountByStatus(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建不同状态的项目
	for i := 0; i < 2; i++ {
		draftProject := createTestProject("author"+string(rune('1'+i)), "Draft "+string(rune('A'+i)))
		draftProject.Status = writer.StatusDraft
		err := repo.Create(ctx, draftProject)
		require.NoError(t, err)
	}

	completedProject := createTestProject("author3", "Completed")
	completedProject.Status = writer.StatusCompleted
	err := repo.Create(ctx, completedProject)
	require.NoError(t, err)

	// 统计草稿
	count, err := repo.CountByStatus(ctx, string(writer.StatusDraft))
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// 统计完成
	count, err = repo.CountByStatus(ctx, string(writer.StatusCompleted))
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// 20. 测试列表查询
func TestProjectRepository_List(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建多个项目
	for i := 1; i <= 5; i++ {
		project := createTestProject("author"+string(rune('0'+i)), "List Project "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 查询所有
	projects, err := repo.List(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, projects, 5)
}

// 21. 测试带筛选条件的列表查询
func TestProjectRepository_List_WithFilter(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建多个项目
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "Filter Project "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 创建其他作者的项目
	otherProject := createTestProject("other_author", "Other Project")
	err := repo.Create(ctx, otherProject)
	require.NoError(t, err)

	// 使用筛选查询
	filter := &testutil.SimpleFilter{
		Conditions: map[string]interface{}{
			"author_id": authorID,
		},
	}

	projects, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.Len(t, projects, 3)
}

// 22. 测试项目存在检查
func TestProjectRepository_Exists(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "Exists Check Project")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 检查存在
	exists, err := repo.Exists(ctx, project.ID)
	require.NoError(t, err)
	assert.True(t, exists)

	// 检查不存在
	exists, err = repo.Exists(ctx, primitive.NewObjectID().Hex())
	require.NoError(t, err)
	assert.False(t, exists)
}

// 23. 测试统计总数
func TestProjectRepository_Count(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	// 创建5个项目
	for i := 1; i <= 5; i++ {
		project := createTestProject("author"+string(rune('0'+i)), "Count All "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 统计全部
	count, err := repo.Count(ctx, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

// 24. 测试带筛选的统计
func TestProjectRepository_Count_WithFilter(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建3个项目
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "Count Filter "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 创建其他作者的2个项目
	for i := 1; i <= 2; i++ {
		project := createTestProject("other_author", "Other "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 统计特定作者
	filter := &testutil.SimpleFilter{
		Conditions: map[string]interface{}{
			"author_id": authorID,
		},
	}

	count, err := repo.Count(ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// 25. 测试事务创建（需要MongoDB副本集支持）
func TestProjectRepository_CreateWithTransaction(t *testing.T) {
	t.Skip("事务功能需要MongoDB副本集支持，跳过测试")

	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "Transaction Project")

	// 使用事务创建
	err := repo.CreateWithTransaction(ctx, project, func(txCtx context.Context) error {
		// 在事务中执行额外操作（这里只是验证回调被调用）
		return nil
	})

	require.NoError(t, err)
	assert.NotEmpty(t, project.ID)

	// 验证项目已创建
	retrieved, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
}

// 26. 测试事务回滚（需要MongoDB副本集支持）
func TestProjectRepository_CreateWithTransaction_Rollback(t *testing.T) {
	t.Skip("事务功能需要MongoDB副本集支持，跳过测试")

	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	project := createTestProject("author123", "Transaction Rollback Project")

	// 使用事务创建，但回调失败
	err := repo.CreateWithTransaction(ctx, project, func(txCtx context.Context) error {
		return assert.AnError // 模拟错误
	})

	assert.Error(t, err)

	// 验证项目未创建（事务已回滚）
	if project.ID != "" {
		retrieved, err := repo.GetByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Nil(t, retrieved)
	}
}

// 27. 测试健康检查
func TestProjectRepository_Health(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}

// 28. 测试分页查询
func TestProjectRepository_GetListByOwnerID_WithPagination(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建10个项目
	for i := 1; i <= 10; i++ {
		project := createTestProject(authorID, "Page Project "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}

	// 获取第一页（5个）
	page1, err := repo.GetListByOwnerID(ctx, authorID, 5, 0)
	require.NoError(t, err)
	assert.Len(t, page1, 5)

	// 获取第二页（5个）
	page2, err := repo.GetListByOwnerID(ctx, authorID, 5, 5)
	require.NoError(t, err)
	assert.Len(t, page2, 5)

	// 确保两页数据不重复
	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

// 29. 测试软删除后的所有者统计
func TestProjectRepository_CountByOwner_AfterSoftDelete(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建3个项目
	var projectID string
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "Count After Delete "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
		if i == 1 {
			projectID = project.ID
		}
	}

	// 软删除一个项目
	err := repo.SoftDelete(ctx, projectID, authorID)
	require.NoError(t, err)

	// 统计应该只剩2个
	count, err := repo.CountByOwner(ctx, authorID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// 30. 测试List不返回已删除项目
func TestProjectRepository_List_ExcludesDeleted(t *testing.T) {
	repo, ctx, cleanup := setupProjectRepo(t)
	defer cleanup()

	authorID := "author123"

	// 创建3个项目
	var projectID string
	for i := 1; i <= 3; i++ {
		project := createTestProject(authorID, "List Deleted "+string(rune('A'+i-1)))
		err := repo.Create(ctx, project)
		require.NoError(t, err)
		if i == 1 {
			projectID = project.ID
		}
	}

	// 软删除一个项目
	err := repo.SoftDelete(ctx, projectID, authorID)
	require.NoError(t, err)

	// List应该只返回2个
	projects, err := repo.List(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, projects, 2)

	// 确保已删除的项目不在列表中
	for _, p := range projects {
		assert.NotEqual(t, projectID, p.ID)
	}
}
