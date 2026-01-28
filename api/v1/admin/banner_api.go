package admin

import (
	"strconv"

	"Qingyu_backend/pkg/response"
	bookstoreModel "Qingyu_backend/models/bookstore" // Imported for Swagger annotations
	"Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
)

// BannerAPI Banner管理API
type BannerAPI struct {
	bannerService bookstore.BannerService
}

// NewBannerAPI 创建Banner管理API实例
func NewBannerAPI(bannerService bookstore.BannerService) *BannerAPI {
	return &BannerAPI{
		bannerService: bannerService,
	}
}

// GetBanners 获取Banner列表
// @Summary 获取Banner列表
// @Description 获取Banner列表，支持筛选和分页
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param isActive query boolean false "是否激活"
// @Param targetType query string false "目标类型(book/category/url)"
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Param sortBy query string false "排序字段(sort_order/click_count/created_at)" default(sort_order)
// @Param sortOrder query string false "排序方向(asc/desc)" default(asc)
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners [get]
func (api *BannerAPI) GetBanners(c *gin.Context) {
	// 解析查询参数
	req := &bookstore.GetBannersRequest{
		Limit:     20,
		Offset:    0,
		SortBy:    "sort_order",
		SortOrder: "asc",
	}

	// 解析isActive
	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			req.IsActive = &isActive
		}
	}

	// 解析targetType
	if targetType := c.Query("targetType"); targetType != "" {
		req.TargetType = &targetType
	}

	// 解析limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	// 解析offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	// 解析sortBy
	if sortBy := c.Query("sortBy"); sortBy != "" {
		req.SortBy = sortBy
	}

	// 解析sortOrder
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		req.SortOrder = sortOrder
	}

	// 调用Service层
	resp, err := api.bannerService.GetBanners(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetBannerByID 获取Banner详情
// @Summary 获取Banner详情
// @Description 根据ID获取Banner详情
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Banner ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners/{id} [get]
func (api *BannerAPI) GetBannerByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Banner ID不能为空", "")
		return
	}

	banner, err := api.bannerService.GetBannerByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, banner)
}

// CreateBanner 创建Banner
// @Summary 创建Banner
// @Description 创建新的Banner
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object true "创建Banner请求"
// @Success 201 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners [post]
func (api *BannerAPI) CreateBanner(c *gin.Context) {
	var req bookstore.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	banner, err := api.bannerService.CreateBanner(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, banner)
}

// UpdateBanner 更新Banner
// @Summary 更新Banner
// @Description 更新Banner信息
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Banner ID"
// @Param request body object true "更新Banner请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners/{id} [put]
func (api *BannerAPI) UpdateBanner(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Banner ID不能为空", "")
		return
	}

	var req bookstore.UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.bannerService.UpdateBanner(c.Request.Context(), id, &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteBanner 删除Banner
// @Summary 删除Banner
// @Description 删除指定的Banner
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Banner ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners/{id} [delete]
func (api *BannerAPI) DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Banner ID不能为空", "")
		return
	}

	if err := api.bannerService.DeleteBanner(c.Request.Context(), id); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// BatchUpdateStatus 批量更新状态
// @Summary 批量更新Banner状态
// @Description 批量启用或禁用Banner
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object true "批量更新状态请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners/batch-status [put]
func (api *BannerAPI) BatchUpdateStatus(c *gin.Context) {
	var req bookstore.BatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.bannerService.BatchUpdateStatus(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// BatchUpdateSort 批量更新排序
// @Summary 批量更新Banner排序
// @Description 批量更新Banner的排序权重
// @Tags 管理员-Banner管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object true "批量更新排序请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/admin/banners/batch-sort [put]
func (api *BannerAPI) BatchUpdateSort(c *gin.Context) {
	var req bookstore.BatchUpdateSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.bannerService.BatchUpdateSort(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// IncrementClickCount 增加点击次数
// @Summary 增加Banner点击次数
// @Description 记录Banner被点击
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path string true "Banner ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/banners/{id}/click [post]
func (api *BannerAPI) IncrementClickCount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Banner ID不能为空", "")
		return
	}

	if err := api.bannerService.IncrementClickCount(c.Request.Context(), id); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

var _ = bookstoreModel.Banner{}
