package messaging

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnnouncementModel(t *testing.T) {
	t.Run("创建公告模型", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)

		announcement := &Announcement{
			Content:    "这是一个测试公告",
			Type:       AnnouncementTypeInfo,
			Priority:   1,
			IsActive:   true,
			TargetRole: "all",
			ViewCount:  0,
			CreatedBy:  "admin",
			StartTime:  &now,
			EndTime:    &later,
		}

		// 设置嵌入字段
		announcement.Title = "系统公告"
		announcement.CreatedAt = now
		announcement.UpdatedAt = now

		assert.Equal(t, "系统公告", announcement.Title)
		assert.Equal(t, "这是一个测试公告", announcement.Content)
		assert.Equal(t, AnnouncementTypeInfo, announcement.Type)
		assert.Equal(t, 1, announcement.Priority)
		assert.True(t, announcement.IsActive)
		assert.Equal(t, "all", announcement.TargetRole)
	})

	t.Run("公告类型验证", func(t *testing.T) {
		validTypes := []AnnouncementType{
			AnnouncementTypeInfo,
			AnnouncementTypeWarning,
			AnnouncementTypeNotice,
		}

		for _, annType := range validTypes {
			assert.NotEmpty(t, string(annType), "类型不应该为空: %v", annType)
		}
	})
}

func TestAnnouncementFilter(t *testing.T) {
	t.Run("创建过滤条件", func(t *testing.T) {
		isActive := true
		typeValue := AnnouncementTypeInfo
		targetRole := "reader"

		filter := &AnnouncementFilter{
			IsActive:   &isActive,
			Type:       &typeValue,
			TargetRole: &targetRole,
			Limit:      10,
			Offset:     0,
			SortBy:     "priority",
			SortOrder:  "desc",
		}

		conditions := filter.GetConditions()

		assert.NotNil(t, conditions)
		assert.Equal(t, true, conditions["is_active"])
		assert.Equal(t, AnnouncementTypeInfo, conditions["type"])
	})

	t.Run("空过滤条件", func(t *testing.T) {
		filter := &AnnouncementFilter{
			Limit:     20,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "desc",
		}

		conditions := filter.GetConditions()
		assert.NotNil(t, conditions)
	})

	t.Run("获取排序", func(t *testing.T) {
		filter := &AnnouncementFilter{
			SortBy:    "priority",
			SortOrder: "desc",
		}

		sort := filter.GetSort()
		assert.Equal(t, -1, sort["priority"])
	})

	t.Run("默认排序", func(t *testing.T) {
		filter := &AnnouncementFilter{}

		sort := filter.GetSort()
		assert.Equal(t, -1, sort["is_pinned"])
		assert.Equal(t, -1, sort["priority"])
		assert.Equal(t, -1, sort["created_at"])
	})
}

func TestAnnouncementTypes(t *testing.T) {
	t.Run("所有公告类型都应该有值", func(t *testing.T) {
		assert.Equal(t, AnnouncementType("info"), AnnouncementTypeInfo)
		assert.Equal(t, AnnouncementType("warning"), AnnouncementTypeWarning)
		assert.Equal(t, AnnouncementType("notice"), AnnouncementTypeNotice)
	})
}

func TestAnnouncementMethods(t *testing.T) {
	t.Run("IsEffective - 有效公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)
		past := now.Add(-1 * time.Hour)

		announcement := &Announcement{
			IsActive:   true,
			StartTime:  &past,
			EndTime:    &later,
			TargetRole: "all",
		}
		announcement.ExpiresAt = &later

		assert.True(t, announcement.IsEffective())
	})

	t.Run("IsEffective - 未激活公告", func(t *testing.T) {
		announcement := &Announcement{
			IsActive: false,
		}

		assert.False(t, announcement.IsEffective())
	})

	t.Run("ShouldShow - 全员公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)

		announcement := &Announcement{
			IsActive:   true,
			TargetRole: "all",
		}
		announcement.ExpiresAt = &later

		assert.True(t, announcement.ShouldShow("reader"))
		assert.True(t, announcement.ShouldShow("writer"))
		assert.True(t, announcement.ShouldShow("admin"))
	})

	t.Run("ShouldShow - 特定角色公告", func(t *testing.T) {
		now := time.Now()
		later := now.Add(24 * time.Hour)

		announcement := &Announcement{
			IsActive:   true,
			TargetRole: "reader",
		}
		announcement.ExpiresAt = &later

		assert.True(t, announcement.ShouldShow("reader"))
		assert.False(t, announcement.ShouldShow("writer"))
	})

	t.Run("Publish和Unpublish", func(t *testing.T) {
		announcement := &Announcement{
			IsActive: false,
		}

		announcement.Publish()
		assert.True(t, announcement.IsActive)

		announcement.Unpublish()
		assert.False(t, announcement.IsActive)
	})

	t.Run("IncrementView", func(t *testing.T) {
		announcement := &Announcement{
			ViewCount: 0,
		}

		announcement.IncrementView()
		assert.Equal(t, int64(1), announcement.ViewCount)

		announcement.IncrementView()
		announcement.IncrementView()
		assert.Equal(t, int64(3), announcement.ViewCount)
	})
}

func TestAnnouncementValidation(t *testing.T) {
	t.Run("Filter验证成功", func(t *testing.T) {
		filter := &AnnouncementFilter{
			Limit:  10,
			Offset: 0,
		}

		err := filter.Validate()
		assert.NoError(t, err)
	})

	t.Run("Filter验证失败 - 负数limit", func(t *testing.T) {
		filter := &AnnouncementFilter{
			Limit: -1,
		}

		err := filter.Validate()
		assert.Error(t, err)
	})

	t.Run("Filter验证失败 - 负数offset", func(t *testing.T) {
		filter := &AnnouncementFilter{
			Limit:  10,
			Offset: -1,
		}

		err := filter.Validate()
		assert.Error(t, err)
	})
}
