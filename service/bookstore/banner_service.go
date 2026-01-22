package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/pkg/errors"
	repo "Qingyu_backend/repository/interfaces/bookstore"
	"context"
	"time"

)

// BannerService Banner服务接口
type BannerService interface {
	// 基础CRUD
	GetBannerByID(ctx context.Context, id string) (*bookstore.Banner, error)
	GetBanners(ctx context.Context, req *GetBannersRequest) (*GetBannersResponse, error)
	CreateBanner(ctx context.Context, req *CreateBannerRequest) (*bookstore.Banner, error)
	UpdateBanner(ctx context.Context, id string, req *UpdateBannerRequest) error
	DeleteBanner(ctx context.Context, id string) error

	// 批量操作
	BatchUpdateStatus(ctx context.Context, req *BatchUpdateStatusRequest) error
	BatchUpdateSort(ctx context.Context, req *BatchUpdateSortRequest) error

	// 统计
	IncrementClickCount(ctx context.Context, id string) error
}

// bannerServiceImpl Banner服务实现
type bannerServiceImpl struct {
	bannerRepo repo.BannerRepository
}

// NewBannerService 创建Banner服务实例
func NewBannerService(bannerRepo repo.BannerRepository) BannerService {
	return &bannerServiceImpl{
		bannerRepo: bannerRepo,
	}
}

// GetBannersRequest 获取Banner列表请求
type GetBannersRequest struct {
	IsActive   *bool   `json:"isActive"`
	TargetType *string `json:"targetType"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
	SortBy     string  `json:"sortBy"`    // sort_order, click_count, created_at
	SortOrder  string  `json:"sortOrder"` // asc, desc
}

// GetBannersResponse 获取Banner列表响应
type GetBannersResponse struct {
	Banners []*bookstore.Banner `json:"banners"`
	Total   int64               `json:"total"`
}

// CreateBannerRequest 创建Banner请求
type CreateBannerRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=100"`
	Description string     `json:"description" validate:"max=200"`
	Image       string     `json:"image" validate:"required,url"`
	Target      string     `json:"target" validate:"required"`
	TargetType  string     `json:"targetType" validate:"required,oneof=book category url"`
	SortOrder   int        `json:"sortOrder"`
	IsActive    bool       `json:"isActive"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

// UpdateBannerRequest 更新Banner请求
type UpdateBannerRequest struct {
	Title       *string    `json:"title" validate:"omitempty,min=1,max=100"`
	Description *string    `json:"description" validate:"omitempty,max=200"`
	Image       *string    `json:"image" validate:"omitempty,url"`
	Target      *string    `json:"target"`
	TargetType  *string    `json:"targetType" validate:"omitempty,oneof=book category url"`
	SortOrder   *int       `json:"sortOrder"`
	IsActive    *bool      `json:"isActive"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	BannerIDs []string `json:"bannerIds" validate:"required,min=1"`
	IsActive  bool     `json:"isActive"`
}

// BatchUpdateSortRequest 批量更新排序请求
type BatchUpdateSortRequest struct {
	Items []struct {
		ID        string `json:"id" validate:"required"`
		SortOrder int    `json:"sortOrder"`
	} `json:"items" validate:"required,min=1"`
}

// GetBannerByID 获取Banner详情
func (s *bannerServiceImpl) GetBannerByID(ctx context.Context, id string) (*bookstore.Banner, error) {
	banner, err := s.bannerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("BANNER_GET_FAILED", "获取Banner失败", err)
	}

	if banner == nil {
		return nil, errors.BookstoreServiceFactory.NotFoundError("Banner", id)
	}

	return banner, nil
}

// GetBanners 获取Banner列表
func (s *bannerServiceImpl) GetBanners(ctx context.Context, req *GetBannersRequest) (*GetBannersResponse, error) {
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 构建查询条件
	var banners []*bookstore.Banner
	var err error

	// 根据不同条件查询
	if req.IsActive != nil && *req.IsActive {
		banners, err = s.bannerRepo.GetActive(ctx, req.Limit, req.Offset)
	} else if req.TargetType != nil {
		banners, err = s.bannerRepo.GetByTargetType(ctx, *req.TargetType, req.Limit, req.Offset)
	} else {
		// 使用通用List方法，需要构建Filter
		filter := &bookstore.BannerFilter{
			IsActive:   req.IsActive,
			TargetType: req.TargetType,
			SortBy:     req.SortBy,
			SortOrder:  req.SortOrder,
			Limit:      req.Limit,
			Offset:     req.Offset,
		}
		banners, err = s.bannerRepo.List(ctx, filter)
	}

	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("BANNER_LIST_FAILED", "获取Banner列表失败", err)
	}

	// 获取总数
	filter := &bookstore.BannerFilter{
		IsActive:   req.IsActive,
		TargetType: req.TargetType,
	}
	total, err := s.bannerRepo.Count(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("BANNER_COUNT_FAILED", "获取Banner总数失败", err)
	}

	return &GetBannersResponse{
		Banners: banners,
		Total:   total,
	}, nil
}

// CreateBanner 创建Banner
func (s *bannerServiceImpl) CreateBanner(ctx context.Context, req *CreateBannerRequest) (*bookstore.Banner, error) {
	// 验证时间范围
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return nil, errors.BookstoreServiceFactory.ValidationError("INVALID_TIME_RANGE", "开始时间不能晚于结束时间")
	}

	banner := &bookstore.Banner{
		Title:       req.Title,
		Description: req.Description,
		Image:       req.Image,
		Target:      req.Target,
		TargetType:  req.TargetType,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.bannerRepo.Create(ctx, banner); err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("BANNER_CREATE_FAILED", "创建Banner失败", err)
	}

	return banner, nil
}

// UpdateBanner 更新Banner
func (s *bannerServiceImpl) UpdateBanner(ctx context.Context, id string, req *UpdateBannerRequest) error {
	// 验证ID格式
	if id == "" {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的Banner ID", id)
	}

	// 检查Banner是否存在
	exists, err := s.bannerRepo.Exists(ctx, id)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("BANNER_CHECK_FAILED", "检查Banner是否存在失败", err)
	}
	if !exists {
		return errors.BookstoreServiceFactory.NotFoundError("Banner", id)
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}
	if req.Target != nil {
		updates["target"] = *req.Target
	}
	if req.TargetType != nil {
		updates["target_type"] = *req.TargetType
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
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

	// 验证时间范围
	if req.StartTime != nil && req.EndTime != nil && req.StartTime.After(*req.EndTime) {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_TIME_RANGE", "开始时间不能晚于结束时间")
	}

	if len(updates) == 0 {
		return nil // 没有更新
	}

	if err := s.bannerRepo.Update(ctx, id, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("BANNER_UPDATE_FAILED", "更新Banner失败", err)
	}

	return nil
}

// DeleteBanner 删除Banner
func (s *bannerServiceImpl) DeleteBanner(ctx context.Context, id string) error {
	// 验证ID格式
	if id == "" {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的Banner ID", id)
	}

	if err := s.bannerRepo.Delete(ctx, id); err != nil {
		return errors.BookstoreServiceFactory.InternalError("BANNER_DELETE_FAILED", "删除Banner失败", err)
	}

	return nil
}

// BatchUpdateStatus 批量更新状态
func (s *bannerServiceImpl) BatchUpdateStatus(ctx context.Context, req *BatchUpdateStatusRequest) error {
	if len(req.BannerIDs) == 0 {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_INPUT", "Banner ID列表不能为空")
	}

	// 验证所有ID格式
	for _, idStr := range req.BannerIDs {
		if idStr == "" {
			return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的Banner ID", idStr)
		}
	}

	if err := s.bannerRepo.BatchUpdateStatus(ctx, req.BannerIDs, req.IsActive); err != nil {
		return errors.BookstoreServiceFactory.InternalError("BATCH_UPDATE_FAILED", "批量更新状态失败", err)
	}

	return nil
}

// BatchUpdateSort 批量更新排序
func (s *bannerServiceImpl) BatchUpdateSort(ctx context.Context, req *BatchUpdateSortRequest) error {
	if len(req.Items) == 0 {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_INPUT", "排序项列表不能为空")
	}

	// 遍历更新每个Banner的排序
	for _, item := range req.Items {
		// 验证ID格式
		if item.ID == "" {
			return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的Banner ID", item.ID)
		}

		updates := map[string]interface{}{
			"sort_order": item.SortOrder,
		}

		if err := s.bannerRepo.Update(ctx, item.ID, updates); err != nil {
			return errors.BookstoreServiceFactory.InternalError("BATCH_UPDATE_FAILED", "更新排序失败", err)
		}
	}

	return nil
}

// IncrementClickCount 增加点击次数
func (s *bannerServiceImpl) IncrementClickCount(ctx context.Context, id string) error {
	// 验证ID格式
	if id == "" {
		return errors.BookstoreServiceFactory.ValidationError("INVALID_ID", "无效的Banner ID", id)
	}

	if err := s.bannerRepo.IncrementClickCount(ctx, id); err != nil {
		return errors.BookstoreServiceFactory.InternalError("INCREMENT_FAILED", "增加点击次数失败", err)
	}

	return nil
}
