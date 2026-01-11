package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/writing"
	"Qingyu_backend/test/testutil"
)

// TestProjectRepository_Create 测试创建项目
func TestProjectRepository_Create(t *testing.T) {
	// 1. 初始化测试环境
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 2. 测试正常创建
	t.Run("正常创建项目", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "测试项目",
			Summary:  "这是一个测试项目",
			Category: "玄幻",
			Tags:     []string{"热血", "升级流"},
		}

		err := repo.Create(ctx, project)
		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.Equal(t, writer.StatusDraft, project.Status)
		assert.Equal(t, writer.VisibilityPrivate, project.Visibility)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	})

	// 3. 测试统计初始化
	t.Run("统计信息初始化", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "测试项目2",
		}

		err := repo.Create(ctx, project)
		require.NoError(t, err)
		assert.Equal(t, 0, project.Statistics.TotalWords)
		assert.Equal(t, 0, project.Statistics.ChapterCount)
		assert.Equal(t, 0, project.Statistics.DocumentCount)
		assert.NotZero(t, project.Statistics.LastUpdateAt)
	})

	// 4. 测试设置初始化
	t.Run("设置信息初始化", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "测试项目3",
		}

		err := repo.Create(ctx, project)
		require.NoError(t, err)
		assert.True(t, project.Settings.AutoBackup)
		assert.Equal(t, 24, project.Settings.BackupInterval)
	})

	// 5. 测试空对象创建
	t.Run("空对象创建应该失败", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能为空")
	})

	// 6. 测试必填字段验证
	t.Run("缺少标题应该失败", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "",
		}

		err := repo.Create(ctx, project)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "验证失败")
	})
}

// TestProjectRepository_GetByID 测试根据ID获取项目
func TestProjectRepository_GetByID(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 1. 创建测试数据
	project := &writer.Project{
		AuthorID: "user123",
		Title:    "测试项目",
		Summary:  "测试简介",
	}
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 2. 测试查询存在的项目
	t.Run("查询存在的项目", func(t *testing.T) {
		found, err := repo.GetByID(ctx, project.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, project.Title, found.Title)
		assert.Equal(t, project.AuthorID, found.AuthorID)
	})

	// 3. 测试查询不存在的项目
	t.Run("查询不存在的项目", func(t *testing.T) {
		found, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	// 4. 测试无效ID格式
	t.Run("无效的项目ID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, "invalid_id")
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	// 5. 测试软删除后查询
	t.Run("软删除后不可查询", func(t *testing.T) {
		err := repo.SoftDelete(ctx, project.ID, project.AuthorID)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

// TestProjectRepository_Update 测试更新项目
func TestProjectRepository_Update(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 创建测试数据
	project := &writer.Project{
		AuthorID: "user123",
		Title:    "原标题",
		Summary:  "原简介",
	}
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 测试更新
	t.Run("更新项目信息", func(t *testing.T) {
		updates := map[string]interface{}{
			"title":   "新标题",
			"summary": "新简介",
		}

		err := repo.Update(ctx, project.ID, updates)
		require.NoError(t, err)

		// 验证更新
		updated, err := repo.GetByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "新标题", updated.Title)
		assert.Equal(t, "新简介", updated.Summary)
		assert.True(t, updated.UpdatedAt.After(project.UpdatedAt))
	})

	// 测试更新不存在的项目
	t.Run("更新不存在的项目", func(t *testing.T) {
		updates := map[string]interface{}{
			"title": "新标题",
		}

		err := repo.Update(ctx, "507f1f77bcf86cd799439011", updates)
		assert.Error(t, err)
	})
}

// TestProjectRepository_GetListByOwnerID 测试获取作者的项目列表
func TestProjectRepository_GetListByOwnerID(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	authorID := "user123"

	// 创建多个项目
	for i := 0; i < 5; i++ {
		project := &writer.Project{
			AuthorID: authorID,
			Title:    fmt.Sprintf("项目%d", i+1),
		}
		err := repo.Create(ctx, project)
		require.NoError(t, err)

		// 稍微延迟以确保不同的updated_at
		time.Sleep(10 * time.Millisecond)
	}

	// 创建其他作者的项目
	otherProject := &writer.Project{
		AuthorID: "other_user",
		Title:    "其他项目",
	}
	err := repo.Create(ctx, otherProject)
	require.NoError(t, err)

	// 测试查询作者的所有项目
	t.Run("查询作者的所有项目", func(t *testing.T) {
		projects, err := repo.GetListByOwnerID(ctx, authorID, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 5, len(projects))

		// 验证按updated_at降序排列
		for i := 0; i < len(projects)-1; i++ {
			assert.True(t, projects[i].UpdatedAt.After(projects[i+1].UpdatedAt) || projects[i].UpdatedAt.Equal(projects[i+1].UpdatedAt))
		}
	})

	// 测试分页查询
	t.Run("分页查询", func(t *testing.T) {
		// 第一页
		page1, err := repo.GetListByOwnerID(ctx, authorID, 2, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(page1))

		// 第二页
		page2, err := repo.GetListByOwnerID(ctx, authorID, 2, 2)
		require.NoError(t, err)
		assert.Equal(t, 2, len(page2))

		// 确保不同页的数据不重复
		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	})
}

// TestProjectRepository_GetByOwnerAndStatus 测试按状态查询
func TestProjectRepository_GetByOwnerAndStatus(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	authorID := "user123"

	// 创建不同状态的项目
	statuses := []writer.ProjectStatus{
		writer.StatusDraft,
		writer.StatusDraft,
		writer.StatusSerializing,
		writer.StatusCompleted,
	}

	for i, status := range statuses {
		project := &writer.Project{
			AuthorID: authorID,
			Title:    fmt.Sprintf("项目%d", i+1),
			Status:   status,
		}
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 测试查询草稿状态的项目
	t.Run("查询草稿状态项目", func(t *testing.T) {
		projects, err := repo.GetByOwnerAndStatus(ctx, authorID, string(writer.StatusDraft), 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(projects))

		for _, p := range projects {
			assert.Equal(t, writer.StatusDraft, p.Status)
		}
	})

	// 测试查询连载中的项目
	t.Run("查询连载中项目", func(t *testing.T) {
		projects, err := repo.GetByOwnerAndStatus(ctx, authorID, string(writer.StatusSerializing), 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(projects))
	})
}

// TestProjectRepository_UpdateByOwner 测试所有者更新
func TestProjectRepository_UpdateByOwner(t *testing.T) {
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 创建测试项目
	project := &writer.Project{
		AuthorID: "user123",
		Title:    "原标题",
	}
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 测试所有者更新
	t.Run("所有者可以更新", func(t *testing.T) {
		updates := map[string]interface{}{
			"title": "新标题",
		}

		err := repo.UpdateByOwner(ctx, project.ID, "user123", updates)
		require.NoError(t, err)

		// 验证更新
		updated, err := repo.GetByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "新标题", updated.Title)
	})

	// 测试非所有者更新
	t.Run("非所有者不能更新", func(t *testing.T) {
		updates := map[string]interface{}{
			"title": "黑客标题",
		}

		err := repo.UpdateByOwner(ctx, project.ID, "hacker", updates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在或无权限")
	})
}

// TestProjectRepository_SoftDelete 测试软删除
func TestProjectRepository_SoftDelete(t *testing.T) {
	db, cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 创建测试项目
	project := &writer.Project{
		AuthorID: "user123",
		Title:    "测试项目",
	}
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 测试软删除
	t.Run("软删除成功", func(t *testing.T) {
		err := repo.SoftDelete(ctx, project.ID, "user123")
		require.NoError(t, err)

		// 验证已删除
		found, err := repo.GetByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	// 测试恢复
	t.Run("恢复已删除的项目", func(t *testing.T) {
		err := repo.Restore(ctx, project.ID, "user123")
		require.NoError(t, err)

		// 验证已恢复
		found, err := repo.GetByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "测试项目", found.Title)
	})
}

// TestProjectRepository_IsOwner 测试所有者检查
func TestProjectRepository_IsOwner(t *testing.T) {
	db, cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 创建测试项目
	project := &writer.Project{
		AuthorID: "user123",
		Title:    "测试项目",
	}
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// 测试真实所有者
	t.Run("真实所有者", func(t *testing.T) {
		isOwner, err := repo.IsOwner(ctx, project.ID, "user123")
		require.NoError(t, err)
		assert.True(t, isOwner)
	})

	// 测试非所有者
	t.Run("非所有者", func(t *testing.T) {
		isOwner, err := repo.IsOwner(ctx, project.ID, "other_user")
		require.NoError(t, err)
		assert.False(t, isOwner)
	})
}

// TestProjectRepository_Count 测试统计功能
func TestProjectRepository_Count(t *testing.T) {
	db, cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	authorID := "user123"

	// 创建多个项目
	for i := 0; i < 3; i++ {
		project := &writer.Project{
			AuthorID: authorID,
			Title:    fmt.Sprintf("项目%d", i+1),
		}
		err := repo.Create(ctx, project)
		require.NoError(t, err)
	}

	// 创建一个软删除的项目
	deletedProject := &writer.Project{
		AuthorID: authorID,
		Title:    "已删除项目",
	}
	err := repo.Create(ctx, deletedProject)
	require.NoError(t, err)
	err = repo.SoftDelete(ctx, deletedProject.ID, authorID)
	require.NoError(t, err)

	// 测试统计
	t.Run("统计作者的项目数（不包含已删除）", func(t *testing.T) {
		count, err := repo.CountByOwner(ctx, authorID)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count) // 不包含已删除的
	})

	// 测试按状态统计
	t.Run("统计草稿状态项目", func(t *testing.T) {
		count, err := repo.CountByStatus(ctx, string(writer.StatusDraft))
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

// TestProjectRepository_Transaction 测试事务支持
func TestProjectRepository_Transaction(t *testing.T) {
	db, cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	// 测试成功的事务
	t.Run("事务成功创建", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "事务测试项目",
		}

		err := repo.CreateWithTransaction(ctx, project, func(ctx context.Context) error {
			// 模拟其他操作
			return nil
		})

		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)

		// 验证项目已创建
		found, err := repo.GetByID(ctx, project.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	})

	// 测试失败的事务（应该回滚）
	t.Run("事务失败回滚", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
			Title:    "事务回滚项目",
		}

		err := repo.CreateWithTransaction(ctx, project, func(ctx context.Context) error {
			// 模拟操作失败
			return fmt.Errorf("模拟失败")
		})

		assert.Error(t, err)

		// 验证项目未创建（已回滚）
		if project.ID != "" {
			found, err := repo.GetByID(ctx, project.ID)
			assert.NoError(t, err)
			assert.Nil(t, found)
		}
	})
}

// TestProjectRepository_Health 测试健康检查
func TestProjectRepository_Health(t *testing.T) {
	db, cleanup := testutil.SetupTestDatabase(t)
	defer cleanup()

	repo := writing.NewMongoProjectRepository(db)
	ctx := context.Background()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}

// TestProject_BusinessMethods 测试Project模型的业务方法
func TestProject_BusinessMethods(t *testing.T) {
	// 测试IsOwner
	t.Run("IsOwner方法", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
		}

		assert.True(t, project.IsOwner("user123"))
		assert.False(t, project.IsOwner("other_user"))
	})

	// 测试CanEdit
	t.Run("CanEdit方法", func(t *testing.T) {
		now := time.Now()
		project := &writer.Project{
			AuthorID: "user123",
			Collaborators: []writer.Collaborator{
				{
					UserID:     "editor1",
					Role:       writer.RoleEditor,
					InvitedAt:  now,
					AcceptedAt: &now,
				},
				{
					UserID:     "viewer1",
					Role:       writer.RoleViewer,
					InvitedAt:  now,
					AcceptedAt: &now,
				},
			},
		}

		// 所有者可以编辑
		assert.True(t, project.CanEdit("user123"))

		// 编辑者可以编辑
		assert.True(t, project.CanEdit("editor1"))

		// 查看者不能编辑
		assert.False(t, project.CanEdit("viewer1"))

		// 非协作者不能编辑
		assert.False(t, project.CanEdit("stranger"))
	})

	// 测试CanView
	t.Run("CanView方法", func(t *testing.T) {
		now := time.Now()
		privateProject := &writer.Project{
			AuthorID:   "user123",
			Visibility: writer.VisibilityPrivate,
			Collaborators: []writer.Collaborator{
				{
					UserID:     "viewer1",
					Role:       writer.RoleViewer,
					InvitedAt:  now,
					AcceptedAt: &now,
				},
			},
		}

		// 所有者可以查看
		assert.True(t, privateProject.CanView("user123"))

		// 协作者可以查看
		assert.True(t, privateProject.CanView("viewer1"))

		// 非协作者不能查看私密项目
		assert.False(t, privateProject.CanView("stranger"))

		// 公开项目任何人都可以查看
		publicProject := &writer.Project{
			AuthorID:   "user123",
			Visibility: writer.VisibilityPublic,
		}
		assert.True(t, publicProject.CanView("stranger"))
	})

	// 测试UpdateStatistics
	t.Run("UpdateStatistics方法", func(t *testing.T) {
		project := &writer.Project{
			AuthorID: "user123",
		}

		oldTime := project.UpdatedAt

		stats := writer.ProjectStats{
			TotalWords:   10000,
			ChapterCount: 10,
		}

		project.UpdateStatistics(stats)

		assert.Equal(t, 10000, project.Statistics.TotalWords)
		assert.Equal(t, 10, project.Statistics.ChapterCount)
		assert.True(t, project.UpdatedAt.After(oldTime))
	})

	// 测试Validate
	t.Run("Validate方法", func(t *testing.T) {
		// 正常的项目
		validProject := &writer.Project{
			AuthorID:   "user123",
			Title:      "测试项目",
			Status:     writer.StatusDraft,
			Visibility: writer.VisibilityPrivate,
		}
		assert.NoError(t, validProject.Validate())

		// 缺少作者ID
		noAuthor := &writer.Project{
			Title: "测试项目",
		}
		assert.Error(t, noAuthor.Validate())

		// 缺少标题
		noTitle := &writer.Project{
			AuthorID: "user123",
		}
		assert.Error(t, noTitle.Validate())

		// 标题过长
		longTitle := &writer.Project{
			AuthorID:   "user123",
			Title:      string(make([]byte, 101)),
			Status:     writer.StatusDraft,
			Visibility: writer.VisibilityPrivate,
		}
		assert.Error(t, longTitle.Validate())

		// 无效状态
		invalidStatus := &writer.Project{
			AuthorID:   "user123",
			Title:      "测试项目",
			Status:     "invalid_status",
			Visibility: writer.VisibilityPrivate,
		}
		assert.Error(t, invalidStatus.Validate())
	})
}
