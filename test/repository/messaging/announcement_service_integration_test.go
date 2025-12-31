package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/messaging"
	mongoMessaging "Qingyu_backend/repository/mongodb/messaging"
	"Qingyu_backend/test/testutil"
)

func setupServiceTest(t *testing.T) context.Context {
	t.Helper()

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	_ = db.Collection("announcements").Drop(ctx)

	return ctx
}

func TestAnnouncementService_Integration(t *testing.T) {
	ctx := setupServiceTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("创建并获取公告", func(t *testing.T) {
		// 创建公告
		announcement := &messaging.Announcement{
			Content:    "这是一个测试公告",
			Type:       messaging.AnnouncementTypeSystem,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			CreatedBy:  "test-admin",
		}
		announcement.Title = "测试公告"
		announcement.CreatedAt = time.Now()
		announcement.UpdatedAt = time.Now()

		err := repo.Create(ctx, announcement)
		require.NoError(t, err)
		assert.NotEmpty(t, announcement.ID)

		// 获取公告
		objectID, err := primitive.ObjectIDFromHex(announcement.ID)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, objectID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, announcement.ID, found.ID)
		assert.Equal(t, "测试公告", found.Title)
		assert.Equal(t, "这是一个测试公告", found.Content)
	})

	t.Run("获取有效公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)
		past := now.Add(-1 * time.Hour)

		// 创建有效公告
		validAnn := &messaging.Announcement{
			Type:       messaging.AnnouncementTypeSystem,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			StartTime:  &past,
			EndTime:    &later,
			CreatedBy:  "test-admin",
		}
		validAnn.Title = "有效公告"
		validAnn.Content = "这是有效公告内容"

		err := repo.Create(ctx, validAnn)
		require.NoError(t, err)

		// 获取有效公告
		announcements, err := repo.GetEffective(ctx, "all", 10)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(announcements), 1)

		// 验证包含刚创建的公告
		found := false
		for _, ann := range announcements {
			if ann.ID == validAnn.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "应该找到刚创建的有效公告")
	})

	t.Run("增加查看次数", func(t *testing.T) {
		announcement := &messaging.Announcement{
			Type:       messaging.AnnouncementTypeNotice,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			ViewCount:  0,
			CreatedBy:  "test-admin",
		}
		announcement.Title = "热门公告"
		announcement.Content = "公告内容"

		err := repo.Create(ctx, announcement)
		require.NoError(t, err)

		objectID, err := primitive.ObjectIDFromHex(announcement.ID)
		require.NoError(t, err)

		// 增加查看次数
		err = repo.IncrementViewCount(ctx, objectID)
		assert.NoError(t, err)

		// 验证
		stats, err := repo.GetViewStats(ctx, objectID)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats)
	})

	t.Run("按类型筛选公告", func(t *testing.T) {
		_ = db.Collection("announcements").Drop(ctx)

		// 创建不同类型的公告
		systemAnn := &messaging.Announcement{
			Type:       messaging.AnnouncementTypeSystem,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			CreatedBy:  "test-admin",
		}
		systemAnn.Title = "系统公告"
		systemAnn.Content = "系统公告内容"

		err := repo.Create(ctx, systemAnn)
		require.NoError(t, err)

		eventAnn := &messaging.Announcement{
			Type:       messaging.AnnouncementTypeEvent,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			CreatedBy:  "test-admin",
		}
		eventAnn.Title = "活动公告"
		eventAnn.Content = "活动公告内容"

		err = repo.Create(ctx, eventAnn)
		require.NoError(t, err)

		// 按类型筛选
		anns, err := repo.GetByType(ctx, string(messaging.AnnouncementTypeSystem), 10, 0)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(anns), 1)

		// 验证返回的都是系统公告
		for _, ann := range anns {
			assert.Equal(t, messaging.AnnouncementTypeSystem, ann.Type)
		}
	})

	t.Run("批量更新状态", func(t *testing.T) {
		var ids []primitive.ObjectID

		// 创建3个公告
		for i := 0; i < 3; i++ {
			ann := &messaging.Announcement{
				Type:       messaging.AnnouncementTypeNotice,
				Priority:   messaging.AnnouncementPriorityNormal,
				IsActive:   true,
				TargetRole: "all",
				CreatedBy:  "test-admin",
			}
			ann.Title = "批量测试公告"
			ann.Content = "内容"

			err := repo.Create(ctx, ann)
			require.NoError(t, err)

			objectID, err := primitive.ObjectIDFromHex(ann.ID)
			require.NoError(t, err)
			ids = append(ids, objectID)
		}

		// 批量禁用
		err := repo.BatchUpdateStatus(ctx, ids, false)
		assert.NoError(t, err)

		// 验证
		for _, id := range ids {
			ann, err := repo.GetByID(ctx, id)
			require.NoError(t, err)
			assert.False(t, ann.IsActive)
		}
	})

	t.Run("批量删除", func(t *testing.T) {
		var ids []primitive.ObjectID

		// 创建3个公告
		for i := 0; i < 3; i++ {
			ann := &messaging.Announcement{
				Type:       messaging.AnnouncementTypeNotice,
				Priority:   messaging.AnnouncementPriorityNormal,
				IsActive:   true,
				TargetRole: "all",
				CreatedBy:  "test-admin",
			}
			ann.Title = "待删除公告"
			ann.Content = "内容"

			err := repo.Create(ctx, ann)
			require.NoError(t, err)

			objectID, err := primitive.ObjectIDFromHex(ann.ID)
			require.NoError(t, err)
			ids = append(ids, objectID)
		}

		// 批量删除
		err := repo.BatchDelete(ctx, ids)
		assert.NoError(t, err)

		// 验证已删除
		for _, id := range ids {
			exists, err := repo.Exists(ctx, id)
			require.NoError(t, err)
			assert.False(t, exists)
		}
	})
}

func TestAnnouncementService_Filter(t *testing.T) {
	ctx := setupServiceTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("使用Filter查询公告", func(t *testing.T) {
		// 创建测试数据
		isActive := true
		typeValue := messaging.AnnouncementTypeSystem
		targetRole := "reader"

		filter := &messaging.AnnouncementFilter{
			IsActive:   &isActive,
			Type:       &typeValue,
			TargetRole: &targetRole,
			Limit:      10,
		}

		announcements, err := repo.List(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, announcements)
	})

	t.Run("使用Filter统计公告数量", func(t *testing.T) {
		isActive := true

		filter := &messaging.AnnouncementFilter{
			IsActive: &isActive,
			Limit:    10,
		}

		count, err := repo.Count(ctx, filter)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(0))
	})
}

func TestAnnouncementService_TimeRange(t *testing.T) {
	ctx := setupServiceTest(t)

	db, _ := testutil.SetupTestDB(t)
	repo := mongoMessaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	t.Run("按时间范围查询", func(t *testing.T) {
		now := time.Now()
		startTime := now.Add(-24 * time.Hour)
		endTime := now.Add(24 * time.Hour)

		// 创建在时间范围内的公告
		ann := &messaging.Announcement{
			Type:       messaging.AnnouncementTypeSystem,
			Priority:   messaging.AnnouncementPriorityNormal,
			IsActive:   true,
			TargetRole: "all",
			CreatedBy:  "test-admin",
		}
		ann.Title = "时间范围测试公告"
		ann.Content = "内容"
		ann.StartTime = &startTime
		ann.EndTime = &endTime

		err := repo.Create(ctx, ann)
		require.NoError(t, err)

		// 查询时间范围内的公告
		anns, err := repo.GetByTimeRange(ctx, &startTime, &endTime, 10, 0)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(anns), 1)
	})
}
