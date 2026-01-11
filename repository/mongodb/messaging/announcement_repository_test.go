package messaging_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	messagingModel "Qingyu_backend/models/messaging"
	"Qingyu_backend/repository/mongodb/messaging"
	"Qingyu_backend/test/testutil"
)

// 测试辅助函数
func setupAnnouncementTest(t *testing.T) (interface {
	Create(ctx context.Context, announcement *messagingModel.Announcement) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*messagingModel.Announcement, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter interface{}) ([]*messagingModel.Announcement, error)
	GetActive(ctx context.Context, targetUsers string, limit, offset int) ([]*messagingModel.Announcement, error)
	GetByType(ctx context.Context, announcementType string, limit, offset int) ([]*messagingModel.Announcement, error)
	GetEffective(ctx context.Context, targetUsers string, limit int) ([]*messagingModel.Announcement, error)
	IncrementViewCount(ctx context.Context, announcementID primitive.ObjectID) error
	Count(ctx context.Context, filter interface{}) (int64, error)
	Exists(ctx context.Context, id primitive.ObjectID) (bool, error)
	BatchUpdateStatus(ctx context.Context, announcementIDs []primitive.ObjectID, isActive bool) error
	BatchDelete(ctx context.Context, announcementIDs []primitive.ObjectID) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
	GetViewStats(ctx context.Context, announcementID primitive.ObjectID) (int64, error)
	GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*messagingModel.Announcement, error)
}, context.Context) {
	t.Helper()

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	repo := messaging.NewMongoAnnouncementRepository(db.Client(), db.Name())

	return repo, ctx
}

func createTestAnnouncement(title, content string) *messagingModel.Announcement {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	ann := &messagingModel.Announcement{
		Content:     content,
		Type:        messagingModel.AnnouncementTypeInfo,
		Priority:    1,
		IsActive:    true,
		TargetRole:  "all",
		ViewCount:   0,
		CreatedBy:   "admin_" + primitive.NewObjectID().Hex(),
		StartTime:   &now,
		EndTime:     &later,
	}

	// 设置嵌入的字段
	ann.Title = title
	ann.ID = ""
	ann.IsPinned = false
	ann.ExpiresAt = &later

	return ann
}

// TestAnnouncementRepository_Create 测试创建公告
func TestAnnouncementRepository_Create(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("测试公告", "这是一个测试公告内容")

	err := repo.Create(ctx, announcement)
	assert.NoError(t, err)
	assert.NotEmpty(t, announcement.ID)
	assert.False(t, announcement.CreatedAt.IsZero())
	assert.False(t, announcement.UpdatedAt.IsZero())
}

// TestAnnouncementRepository_GetByID 测试根据ID获取公告
func TestAnnouncementRepository_GetByID(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("获取测试", "获取测试内容")
	err := repo.Create(ctx, announcement)
	require.NoError(t, err)

	objectID, err := primitive.ObjectIDFromHex(announcement.ID)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, objectID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, announcement.ID, found.ID)
}

// TestAnnouncementRepository_GetByID_NotFound 测试获取不存在的公告
func TestAnnouncementRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	objectID := primitive.NewObjectID()

	found, err := repo.GetByID(ctx, objectID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestAnnouncementRepository_Update 测试更新公告
func TestAnnouncementRepository_Update(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("原标题", "原内容")
	err := repo.Create(ctx, announcement)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	objectID, err := primitive.ObjectIDFromHex(announcement.ID)
	require.NoError(t, err)

	updates := map[string]interface{}{
		"title":    "新标题",
		"content":  "新内容",
		"priority": 5,
	}

	err = repo.Update(ctx, objectID, updates)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, objectID)
	require.NoError(t, err)
	assert.Equal(t, "新标题", found.Title)
}

// TestAnnouncementRepository_Delete 测试删除公告
func TestAnnouncementRepository_Delete(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("删除测试", "删除测试内容")
	err := repo.Create(ctx, announcement)
	require.NoError(t, err)

	objectID, err := primitive.ObjectIDFromHex(announcement.ID)
	require.NoError(t, err)

	err = repo.Delete(ctx, objectID)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, objectID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestAnnouncementRepository_List 测试列表查询
func TestAnnouncementRepository_List(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	for i := 0; i < 5; i++ {
		announcement := createTestAnnouncement("公告"+string(rune('A'+i)), "内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)
	}

	// 测试列表查询
	filter := &messagingModel.AnnouncementFilter{}

	announcements, err := repo.List(ctx, filter)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(announcements), 5)
}

// TestAnnouncementRepository_GetActive 测试获取活跃公告
func TestAnnouncementRepository_GetActive(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	now := time.Now()
	later := now.Add(24 * time.Hour)

	ann := &messagingModel.Announcement{
		Content:    "活跃内容",
		Type:       messagingModel.AnnouncementTypeInfo,
		Priority:   1,
		IsActive:   true,
		TargetRole: "all",
		CreatedBy:  "admin",
		StartTime:  &now,
		EndTime:    &later,
	}
	ann.Title = "活跃公告"
	ann.ExpiresAt = &later

	err := repo.Create(ctx, ann)
	require.NoError(t, err)

	announcements, err := repo.GetActive(ctx, "all", 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(announcements), 1)
}

// TestAnnouncementRepository_GetByType 测试根据类型获取公告
func TestAnnouncementRepository_GetByType(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	now := time.Now()
	later := now.Add(24 * time.Hour)

	ann := &messagingModel.Announcement{
		Content:    "警告内容",
		Type:       messagingModel.AnnouncementTypeWarning,
		Priority:   1,
		IsActive:   true,
		TargetRole: "all",
		CreatedBy:  "admin",
		StartTime:  &now,
		EndTime:    &later,
	}
	ann.Title = "警告公告"
	ann.ExpiresAt = &later

	err := repo.Create(ctx, ann)
	require.NoError(t, err)

	announcements, err := repo.GetByType(ctx, string(messagingModel.AnnouncementTypeWarning), 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(announcements), 1)
}

// TestAnnouncementRepository_IncrementViewCount 测试增加浏览次数
func TestAnnouncementRepository_IncrementViewCount(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("浏览测试", "浏览内容")
	err := repo.Create(ctx, announcement)
	require.NoError(t, err)

	objectID, err := primitive.ObjectIDFromHex(announcement.ID)
	require.NoError(t, err)

	err = repo.IncrementViewCount(ctx, objectID)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, objectID)
	require.NoError(t, err)
	assert.Greater(t, found.ViewCount, int64(0))
}

// TestAnnouncementRepository_Count 测试统计公告数
func TestAnnouncementRepository_Count(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	for i := 0; i < 3; i++ {
		announcement := createTestAnnouncement("统计测试"+string(rune('A'+i)), "内容")
		err := repo.Create(ctx, announcement)
		require.NoError(t, err)
	}

	filter := &messagingModel.AnnouncementFilter{}

	count, err := repo.Count(ctx, filter)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}

// TestAnnouncementRepository_Exists 测试检查公告是否存在
func TestAnnouncementRepository_Exists(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	announcement := createTestAnnouncement("存在测试", "内容")
	err := repo.Create(ctx, announcement)
	require.NoError(t, err)

	objectID, err := primitive.ObjectIDFromHex(announcement.ID)
	require.NoError(t, err)

	exists, err := repo.Exists(ctx, objectID)
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestAnnouncementRepository_Exists_NotFound 测试检查不存在的公告
func TestAnnouncementRepository_Exists_NotFound(t *testing.T) {
	repo, ctx := setupAnnouncementTest(t)

	objectID := primitive.NewObjectID()

	exists, err := repo.Exists(ctx, objectID)
	require.NoError(t, err)
	assert.False(t, exists)
}
