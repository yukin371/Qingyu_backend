package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/messaging"
	mongoMessaging "Qingyu_backend/repository/mongodb/messaging"
	"Qingyu_backend/test/testutil"
)

// setupAnnouncementTest 设置测试环境
func setupAnnouncementTest(t *testing.T) context.Context {
	t.Helper()

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	// 清空announcements集合
	ctx := context.Background()
	_ = db.Collection("announcements").Drop(ctx)

	return ctx
}

// createTestAnnouncement 创建测试公告
func createTestAnnouncement(title, content string) *messaging.Announcement {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	return &messaging.Announcement{
		Title:       title,
		Content:     content,
		Type:        messaging.AnnouncementTypeSystem,
		Priority:    messaging.AnnouncementPriorityNormal,
		IsActive:    true,
		TargetRole:  "all",
		ViewCount:   0,
		CreatedBy:   "admin",
		StartTime:   &now,
		EndTime:     &later,
	}
}

// ==================== 基础CRUD测试 ====================

func TestAnnouncementRepository_Create(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("成功创建公告", func(t *testing.T) {
		announcement := createTestAnnouncement("测试公告", "这是一个测试公告内容")

		err := repo.Create(ctx, announcement)
		assert.NoError(t, err)
		assert.NotEmpty(t, announcement.ID)
		assert.False(t, announcement.CreatedAt.IsZero())
		assert.False(t, announcement.UpdatedAt.IsZero())
		assert.Equal(t, "测试公告", announcement.Title)
		assert.Equal(t, messaging.AnnouncementTypeSystem, announcement.Type)
	})

	t.Run("创建多个公告", func(t *testing.T) {
		_ = db.Collection("announcements").Drop(ctx)

		ann1 := createTestAnnouncement("公告1", "内容1")
		ann2 := createTestAnnouncement("公告2", "内容2")

		err := repo.Create(ctx, ann1)
		require.NoError(t, err)

		err = repo.Create(ctx, ann2)
		assert.NoError(t, err)
		assert.NotEqual(t, ann1.ID, ann2.ID)
	})
}

func TestAnnouncementRepository_GetByID(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("成功获取公告", func(t *testing.T) {
		announcement := createTestAnnouncement("获取测试", "获取测试内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, announcement.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, announcement.ID, found.ID)
		assert.Equal(t, announcement.Title, found.Title)
		assert.Equal(t, announcement.Content, found.Content)
	})

	t.Run("公告不存在", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent_id")
		assert.Error(t, err)
	})
}

func TestAnnouncementRepository_Update(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("成功更新公告", func(t *testing.T) {
		announcement := createTestAnnouncement("原标题", "原内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)

		// 等待一小段时间确保UpdatedAt会变化
		time.Sleep(10 * time.Millisecond)

		updated := &messaging.Announcement{
			Title:    "新标题",
			Content:  "新内容",
			Priority: messaging.AnnouncementPriorityHigh,
		}

		err = repo.Update(ctx, announcement.ID, updated)
		assert.NoError(t, err)

		// 验证更新
		found, err := repo.GetByID(ctx, announcement.ID)
		require.NoError(t, err)
		assert.Equal(t, "新标题", found.Title)
		assert.Equal(t, "新内容", found.Content)
		assert.Equal(t, messaging.AnnouncementPriorityHigh, found.Priority)
	})

	t.Run("更新不存在的公告", func(t *testing.T) {
		updated := &messaging.Announcement{
			Title: "标题",
		}

		err := repo.Update(ctx, "nonexistent_id", updated)
		assert.Error(t, err)
	})
}

func TestAnnouncementRepository_Delete(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("成功删除公告", func(t *testing.T) {
		announcement := createTestAnnouncement("待删除", "内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)

		err = repo.Delete(ctx, announcement.ID)
		assert.NoError(t, err)

		// 验证已删除
		_, err = repo.GetByID(ctx, announcement.ID)
		assert.Error(t, err)
	})

	t.Run("删除不存在的公告", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent_id")
		assert.Error(t, err)
	})
}

// ==================== 查询测试 ====================

func TestAnnouncementRepository_List(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("获取公告列表", func(t *testing.T) {
		// 创建测试数据
		for i := 1; i <= 5; i++ {
			ann := createTestAnnouncement("公告"+string(rune(i)), "内容"+string(rune(i)))
			_ = repo.Create(ctx, ann)
		}

		filter := &messaging.AnnouncementFilter{
			Limit: 10,
		}

		announcements, total, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(5))
		assert.GreaterOrEqual(t, len(announcements), 5)
	})

	t.Run("分页查询", func(t *testing.T) {
		_ = db.Collection("announcements").Drop(ctx)

		// 创建15个公告
		for i := 1; i <= 15; i++ {
			ann := createTestAnnouncement("公告"+string(rune(i)), "内容")
			_ = repo.Create(ctx, ann)
		}

		// 第一页
		filter1 := &messaging.AnnouncementFilter{
			Limit:  10,
			Offset: 0,
		}
		page1, total1, err := repo.List(ctx, filter1)
		require.NoError(t, err)
		assert.Equal(t, int64(15), total1)
		assert.LessOrEqual(t, len(page1), 10)

		// 第二页
		filter2 := &messaging.AnnouncementFilter{
			Limit:  10,
			Offset: 10,
		}
		page2, _, err := repo.List(ctx, filter2)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(page2), 5)
	})

	t.Run("按类型筛选", func(t *testing.T) {
		_ = db.Collection("announcements").Drop(ctx)

		// 创建不同类型的公告
		systemAnn := createTestAnnouncement("系统公告", "内容")
		systemAnn.Type = messaging.AnnouncementTypeSystem
		_ = repo.Create(ctx, systemAnn)

		eventAnn := createTestAnnouncement("活动公告", "内容")
		eventAnn.Type = messaging.AnnouncementTypeEvent
		_ = repo.Create(ctx, eventAnn)

		filter := &messaging.AnnouncementFilter{
			Type:  &messaging.AnnouncementTypeSystem,
			Limit: 10,
		}

		announcements, _, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		for _, ann := range announcements {
			assert.Equal(t, messaging.AnnouncementTypeSystem, ann.Type)
		}
	})
}

func TestAnnouncementRepository_GetEffectiveAnnouncements(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("获取有效公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)
		past := now.Add(-24 * time.Hour)

		// 有效公告
		validAnn := createTestAnnouncement("有效公告", "内容")
		validAnn.IsActive = true
		validAnn.StartTime = &past
		validAnn.EndTime = &later
		_ = repo.Create(ctx, validAnn)

		// 未激活公告
		inactiveAnn := createTestAnnouncement("未激活公告", "内容")
		inactiveAnn.IsActive = false
		_ = repo.Create(ctx, inactiveAnn)

		// 过期公告
		expiredAnn := createTestAnnouncement("过期公告", "内容")
		expiredEnd := now.Add(-1 * time.Hour)
		expiredAnn.IsActive = true
		expiredAnn.EndTime = &expiredEnd
		_ = repo.Create(ctx, expiredAnn)

		announcements, err := repo.GetEffectiveAnnouncements(ctx, "all", 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(announcements), 1)

		// 验证返回的都是有效公告
		for _, ann := range announcements {
			assert.True(t, ann.IsActive)
		}
	})

	t.Run("按角色筛选", func(t *testing.T) {
		_ = db.Collection("announcements").Drop(ctx)

		now := time.Now()
		later := now.Add(24 * time.Hour)

		// 全员公告
		allAnn := createTestAnnouncement("全员公告", "内容")
		allAnn.TargetRole = "all"
		allAnn.StartTime = &now
		allAnn.EndTime = &later
		_ = repo.Create(ctx, allAnn)

		// 作者公告
		writerAnn := createTestAnnouncement("作者公告", "内容")
		writerAnn.TargetRole = "writer"
		writerAnn.StartTime = &now
		writerAnn.EndTime = &later
		_ = repo.Create(ctx, writerAnn)

		// 获取全员公告
		announcements, err := repo.GetEffectiveAnnouncements(ctx, "all", 10)
		assert.NoError(t, err)

		// 应该包含全员公告
		hasAllAnn := false
		for _, ann := range announcements {
			if ann.Title == "全员公告" {
				hasAllAnn = true
				break
			}
		}
		assert.True(t, hasAllAnn)
	})
}

// ==================== 批量操作测试 ====================

func TestAnnouncementRepository_IncrementViewCount(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("增加查看次数", func(t *testing.T) {
		announcement := createTestAnnouncement("热门公告", "内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)

		initialCount := announcement.ViewCount

		// 增加5次查看
		for i := 0; i < 5; i++ {
			err = repo.IncrementViewCount(ctx, announcement.ID)
			assert.NoError(t, err)
		}

		// 验证
		found, err := repo.GetByID(ctx, announcement.ID)
		require.NoError(t, err)
		assert.Equal(t, initialCount+5, found.ViewCount)
	})
}

func TestAnnouncementRepository_BatchUpdateStatus(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("批量更新状态", func(t *testing.T) {
		var ids []string

		// 创建3个公告
		for i := 1; i <= 3; i++ {
			ann := createTestAnnouncement("公告"+string(rune(i)), "内容")
			ann.IsActive = true
			err := repo.Create(ctx, ann)
			require.NoError(t, err)
			ids = append(ids, ann.ID)
		}

		// 批量禁用
		err := repo.BatchUpdateStatus(ctx, ids, false)
		assert.NoError(t, err)

		// 验证
		for _, id := range ids {
			found, err := repo.GetByID(ctx, id)
			require.NoError(t, err)
			assert.False(t, found.IsActive)
		}
	})
}

func TestAnnouncementRepository_BatchDelete(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("批量删除", func(t *testing.T) {
		var ids []string

		// 创建3个公告
		for i := 1; i <= 3; i++ {
			ann := createTestAnnouncement("待删除"+string(rune(i)), "内容")
			err := repo.Create(ctx, ann)
			require.NoError(t, err)
			ids = append(ids, ann.ID)
		}

		// 批量删除
		err := repo.BatchDelete(ctx, ids)
		assert.NoError(t, err)

		// 验证已删除
		for _, id := range ids {
			_, err := repo.GetByID(ctx, id)
			assert.Error(t, err)
		}
	})
}

// ==================== 健康检查测试 ====================

func TestAnnouncementRepository_Health(t *testing.T) {
	ctx := setupAnnouncementTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("健康检查成功", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)
	})
}
