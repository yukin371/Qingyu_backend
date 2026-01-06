package messaging

import (
	"Qingyu_backend/models/messaging"
	"Qingyu_backend/pkg/errors"
	repo "Qingyu_backend/repository/interfaces/messaging"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnnouncementService 公告服务接口
type AnnouncementService interface {
	// 基础CRUD
	GetAnnouncementByID(ctx context.Context, id string) (*messaging.Announcement, error)
	GetAnnouncements(ctx context.Context, req *GetAnnouncementsRequest) (*GetAnnouncementsResponse, error)
	GetEffectiveAnnouncements(ctx context.Context, targetUsers string, limit int) ([]*messaging.Announcement, error)
	CreateAnnouncement(ctx context.Context, req *CreateAnnouncementRequest) (*messaging.Announcement, error)
	UpdateAnnouncement(ctx context.Context, id string, req *UpdateAnnouncementRequest) error
	DeleteAnnouncement(ctx context.Context, id string) error

	// 批量操作
	BatchUpdateStatus(ctx context.Context, req *BatchUpdateAnnouncementStatusRequest) error
	BatchDelete(ctx context.Context, ids []string) error

	// 统计
	IncrementViewCount(ctx context.Context, id string) error
}

// announcementServiceImpl 公告服务实现
type announcementServiceImpl struct {
	announcementRepo repo.AnnouncementRepository
}

// NewAnnouncementService 创建公告服务实例
func NewAnnouncementService(announcementRepo repo.AnnouncementRepository) AnnouncementService {
	return &announcementServiceImpl{
		announcementRepo: announcementRepo,
	}
}

// GetAnnouncementsRequest 获取公告列表请求
type GetAnnouncementsRequest struct {
	IsActive   *bool   `json:"isActive"`
	Type       *string `json:"type"`
	TargetRole *string `json:"targetRole"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
	SortBy     string  `json:"sortBy"`    // priority, created_at, view_count
	SortOrder  string  `json:"sortOrder"` // asc, desc
}

// GetAnnouncementsResponse 获取公告列表响应
type GetAnnouncementsResponse struct {
	Announcements []*messaging.Announcement `json:"announcements"`
	Total         int64                     `json:"total"`
}

// CreateAnnouncementRequest 创建公告请求
type CreateAnnouncementRequest struct {
	Title      string     `json:"title" validate:"required,min=1,max=200"`
	Content    string     `json:"content" validate:"required,min=1"`
	Type       string     `json:"type" validate:"required,oneof=info warning notice"`
	Priority   int        `json:"priority"`
	IsActive   bool       `json:"isActive"`
	StartTime  *time.Time `json:"startTime"`
	EndTime    *time.Time `json:"endTime"`
	TargetRole string     `json:"targetRole" validate:"required,oneof=all reader writer admin"`
	CreatedBy  string     `json:"createdBy"`
}

// UpdateAnnouncementRequest 更新公告请求
type UpdateAnnouncementRequest struct {
	Title      *string    `json:"title" validate:"omitempty,min=1,max=200"`
	Content    *string    `json:"content" validate:"omitempty,min=1"`
	Type       *string    `json:"type" validate:"omitempty,oneof=info warning notice"`
	Priority   *int       `json:"priority"`
	IsActive   *bool      `json:"isActive"`
	StartTime  *time.Time `json:"startTime"`
	EndTime    *time.Time `json:"endTime"`
	TargetRole *string    `json:"targetRole" validate:"omitempty,oneof=all reader writer admin"`
}

// BatchUpdateAnnouncementStatusRequest 批量更新公告状态请求
type BatchUpdateAnnouncementStatusRequest struct {
	AnnouncementIDs []string `json:"announcementIds" validate:"required,min=1"`
	IsActive        bool     `json:"isActive"`
}

// GetAnnouncementByID 获取公告详情
func (s *announcementServiceImpl) GetAnnouncementByID(ctx context.Context, id string) (*messaging.Announcement, error) {
	announcementID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", id)
	}

	announcement, err := s.announcementRepo.GetByID(ctx, announcementID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_GET_FAILED", "获取公告失败", err)
	}

	if announcement == nil {
		return nil, errors.BookstoreServiceFactory.NotFoundError("Announcement", id)
	}

	return announcement, nil
}

// GetAnnouncements 获取公告列表
func (s *announcementServiceImpl) GetAnnouncements(ctx context.Context, req *GetAnnouncementsRequest) (*GetAnnouncementsResponse, error) {
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 构建Filter
	var announcementType *messaging.AnnouncementType
	if req.Type != nil {
		t := messaging.AnnouncementType(*req.Type)
		announcementType = &t
	}

	filter := &messaging.AnnouncementFilter{
		IsActive:   req.IsActive,
		Type:       announcementType,
		TargetRole: req.TargetRole,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	announcements, err := s.announcementRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_LIST_FAILED", "获取公告列表失败", err)
	}

	// 获取总数
	total, err := s.announcementRepo.Count(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_COUNT_FAILED", "获取公告总数失败", err)
	}

	return &GetAnnouncementsResponse{
		Announcements: announcements,
		Total:         total,
	}, nil
}

// GetEffectiveAnnouncements 获取当前有效的公告
func (s *announcementServiceImpl) GetEffectiveAnnouncements(ctx context.Context, targetUsers string, limit int) ([]*messaging.Announcement, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	announcements, err := s.announcementRepo.GetEffective(ctx, targetUsers, limit)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_GET_FAILED", "获取有效公告失败", err)
	}

	return announcements, nil
}

// CreateAnnouncement 创建公告
func (s *announcementServiceImpl) CreateAnnouncement(ctx context.Context, req *CreateAnnouncementRequest) (*messaging.Announcement, error) {
	// 验证时间范围
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return nil, errors.BookstoreServiceFactory.ValidationError("INVALID_TIME_RANGE", "开始时间不能晚于结束时间")
	}

	announcement := &messaging.Announcement{
		Content:    req.Content,
		Type:       messaging.AnnouncementType(req.Type),
		Priority:   req.Priority,
		IsActive:   req.IsActive,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		TargetRole: req.TargetRole,
		ViewCount:  0,
		CreatedBy:  req.CreatedBy,
	}

	// 设置嵌入字段
	announcement.Title = req.Title
	announcement.CreatedAt = time.Now()
	announcement.UpdatedAt = time.Now()

	if err := s.announcementRepo.Create(ctx, announcement); err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_CREATE_FAILED", "创建公告失败", err)
	}

	return announcement, nil
}

// UpdateAnnouncement 更新公告
func (s *announcementServiceImpl) UpdateAnnouncement(ctx context.Context, id string, req *UpdateAnnouncementRequest) error {
	announcementID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", id)
	}

	// 检查公告是否存在
	exists, err := s.announcementRepo.Exists(ctx, announcementID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_CHECK_FAILED", "检查公告是否存在失败", err)
	}
	if !exists {
		return errors.BookstoreServiceFactory.NotFoundError("Announcement", id)
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.StartTime != nil {
		updates["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = *req.EndTime
	}
	if req.TargetRole != nil {
		updates["target_role"] = *req.TargetRole
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}

	// 验证时间范围
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_TIME_RANGE", "开始时间不能晚于结束时间")
	}

	if len(updates) == 0 {
		return nil // 没有更新
	}

	if err := s.announcementRepo.Update(ctx, announcementID, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_UPDATE_FAILED", "更新公告失败", err)
	}

	return nil
}

// DeleteAnnouncement 删除公告
func (s *announcementServiceImpl) DeleteAnnouncement(ctx context.Context, id string) error {
	announcementID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", id)
	}

	if err := s.announcementRepo.Delete(ctx, announcementID); err != nil {
		return errors.BookstoreServiceFactory.InternalError("ANNOUNCEMENT_DELETE_FAILED", "删除公告失败", err)
	}

	return nil
}

// BatchUpdateStatus 批量更新状态
func (s *announcementServiceImpl) BatchUpdateStatus(ctx context.Context, req *BatchUpdateAnnouncementStatusRequest) error {
	if len(req.AnnouncementIDs) == 0 {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_INPUT", "公告ID列表不能为空")
	}

	// 转换ID
	announcementIDs := make([]primitive.ObjectID, 0, len(req.AnnouncementIDs))
	for _, idStr := range req.AnnouncementIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", idStr)
		}
		announcementIDs = append(announcementIDs, id)
	}

	if err := s.announcementRepo.BatchUpdateStatus(ctx, announcementIDs, req.IsActive); err != nil {
		return errors.BookstoreServiceFactory.InternalError("BATCH_UPDATE_FAILED", "批量更新状态失败", err)
	}

	return nil
}

// BatchDelete 批量删除
func (s *announcementServiceImpl) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_INPUT", "公告ID列表不能为空")
	}

	// 转换ID
	announcementIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, idStr := range ids {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", idStr)
		}
		announcementIDs = append(announcementIDs, id)
	}

	if err := s.announcementRepo.BatchDelete(ctx, announcementIDs); err != nil {
		return errors.BookstoreServiceFactory.InternalError("BATCH_DELETE_FAILED", "批量删除失败", err)
	}

	return nil
}

// IncrementViewCount 增加查看次数
func (s *announcementServiceImpl) IncrementViewCount(ctx context.Context, id string) error {
	announcementID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的公告ID", id)
	}

	if err := s.announcementRepo.IncrementViewCount(ctx, announcementID); err != nil {
		return errors.BookstoreServiceFactory.InternalError("INCREMENT_FAILED", "增加查看次数失败", err)
	}

	return nil
}
