package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	adminsvc "Qingyu_backend/service/admin"
)

// CategoryAdminAPI 分类管理API
type CategoryAdminAPI struct {
	categoryService adminsvc.CategoryAdminService
}

// NewCategoryAdminAPI 创建分类管理API实例
func NewCategoryAdminAPI(categoryService adminsvc.CategoryAdminService) *CategoryAdminAPI {
	return &CategoryAdminAPI{
		categoryService: categoryService,
	}
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 管理员创建新的书籍分类
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/admin/categories [post]
func (api *CategoryAdminAPI) CreateCategory(c *gin.Context) {
	var req adminsvc.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, "参数错误")
		return
	}

	category, err := api.categoryService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "创建分类失败")
		return
	}

	response.SuccessJSON(c, "创建成功", category)
}

// GetCategories 获取分类列表
// @Summary 获取分类列表
// @Description 获取书籍分类列表（支持筛选）
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param parent_id query string false "父分类ID"
// @Param level query int false "层级"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/categories [get]
func (api *CategoryAdminAPI) GetCategories(c *gin.Context) {
	var filter adminsvc.CategoryFilter

	if parentID := c.Query("parent_id"); parentID != "" {
		filter.ParentID = &parentID
	}
	if levelStr := c.Query("level"); levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			response.ErrorJSON(c, http.StatusBadRequest, "level参数格式错误")
			return
		}
		filter.Level = &level
	}

	categories, err := api.categoryService.GetCategories(c.Request.Context(), &filter)
	if err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	response.SuccessJSON(c, "获取成功", categories)
}

// GetCategoryTree 获取分类树
// @Summary 获取分类树
// @Description 获取完整的分类树形结构
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/admin/categories/tree [get]
func (api *CategoryAdminAPI) GetCategoryTree(c *gin.Context) {
	tree, err := api.categoryService.GetCategoryTree(c.Request.Context())
	if err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "获取分类树失败")
		return
	}

	response.SuccessJSON(c, "获取成功", tree)
}

// GetCategoryByID 获取分类详情
// @Summary 获取分类详情
// @Description 根据ID获取分类详细信息
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/admin/categories/{id} [get]
func (api *CategoryAdminAPI) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorJSON(c, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	category, err := api.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorJSON(c, http.StatusNotFound, "分类不存在")
		return
	}

	response.SuccessJSON(c, "获取成功", category)
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 更新分类信息
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param request body UpdateCategoryRequest true "更新内容"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/admin/categories/{id} [put]
func (api *CategoryAdminAPI) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorJSON(c, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	var req adminsvc.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, "参数错误")
		return
	}

	category, err := api.categoryService.UpdateCategory(c.Request.Context(), id, &req)
	if err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "更新分类失败")
		return
	}

	response.SuccessJSON(c, "更新成功", category)
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 删除分类（安全删除，需无子分类和关联作品）
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/admin/categories/{id} [delete]
func (api *CategoryAdminAPI) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorJSON(c, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	if err := api.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "删除分类失败")
		return
	}

	response.SuccessJSON(c, "删除成功", nil)
}

// MoveCategory 移动分类
// @Summary 移动分类
// @Description 移动分类到新的父分类下
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param request body MoveCategoryRequest true "移动信息"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/admin/categories/{id}/move [put]
func (api *CategoryAdminAPI) MoveCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorJSON(c, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	var req adminsvc.MoveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := api.categoryService.MoveCategory(c.Request.Context(), id, &req); err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "移动分类失败")
		return
	}

	response.SuccessJSON(c, "移动成功", nil)
}

// SortCategory 调整分类排序
// @Summary 调整分类排序
// @Description 调整分类的排序序号
// @Tags Admin-Category
// @Accept json
// @Produce json
// @Param id path string true "分类ID"
// @Param request body map[string]int true "排序序号 {\"sort_order\": 1}"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/admin/categories/{id}/sort [put]
func (api *CategoryAdminAPI) SortCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorJSON(c, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorJSON(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := api.categoryService.SortCategory(c.Request.Context(), id, req.SortOrder); err != nil {
		response.ErrorJSON(c, http.StatusInternalServerError, "调整排序失败")
		return
	}

	response.SuccessJSON(c, "调整成功", nil)
}

// ========== Swagger 请求结构体定义 ==========

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required" example:"玄幻"`
	Description string `json:"description" example:"玄幻小说分类"`
	ParentID    string `json:"parent_id" example:""`
	SortOrder   int    `json:"sort_order" example:"1"`
	Icon        string `json:"icon" example:"/icons/xuanhuan.png"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" example:"玄幻"`
	Description string `json:"description" example:"玄幻小说分类"`
	ParentID    string `json:"parent_id" example:""`
	SortOrder   int    `json:"sort_order" example:"1"`
	Icon        string `json:"icon" example:"/icons/xuanhuan.png"`
	IsActive    bool   `json:"is_active" example:"true"`
}

// MoveCategoryRequest 移动分类请求
type MoveCategoryRequest struct {
	NewParentID string `json:"new_parent_id" binding:"required" example:"parent_category_id"`
}
