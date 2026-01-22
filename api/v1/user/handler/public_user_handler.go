package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
)

// PublicUserHandler 公开用户信息处理器
type PublicUserHandler struct {
	userService      userServiceInterface.UserService
	bookstoreService interface{} // 可选依赖
}

// NewPublicUserHandler 创建公开用户信息处理器实例
func NewPublicUserHandler(userService userServiceInterface.UserService) *PublicUserHandler {
	return &PublicUserHandler{
		userService: userService,
	}
}

// SetBookstoreService 设置BookstoreService（可选依赖注入）
func (h *PublicUserHandler) SetBookstoreService(bookstoreSvc interface{}) {
	h.bookstoreService = bookstoreSvc
}

// GetUser 获取用户信息（公开访问）
//
//	@Summary		获取用户信息
//	@Description	获取指定用户的公开信息
//	@Tags			用户管理-公开信息
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse{data=dto.UserProfileResponse}
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/users/{id} [get]
func (h *PublicUserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 调用Service层获取用户信息
	serviceReq := &userServiceInterface.GetUserRequest{
		ID: userID,
	}

	resp, err := h.userService.GetUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "获取用户信息失败", err)
			}
			return
		}
		shared.InternalError(c, "获取用户信息失败", err)
		return
	}

	// 构建公开信息响应（不包含敏感信息）
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	publicProfile := dto.PublicUserProfileResponse{
		UserID:    resp.User.ID.Hex(),
		Username:  resp.User.Username,
		Avatar:    resp.User.Avatar,
		Nickname:  resp.User.Nickname,
		Bio:       resp.User.Bio,
		Role:      role,
		CreatedAt: resp.User.CreatedAt,
	}

	shared.Success(c, http.StatusOK, "获取成功", publicProfile)
}

// GetUserProfile 获取用户详细资料（公开访问）
//
//	@Summary		获取用户详细资料
//	@Description	获取指定用户的详细资料，用于展示用户主页
//	@Tags			用户管理-公开信息
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse{data=dto.PublicUserProfileResponse}
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/users/{id}/profile [get]
func (h *PublicUserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 调用Service层获取用户信息
	serviceReq := &userServiceInterface.GetUserRequest{
		ID: userID,
	}

	resp, err := h.userService.GetUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "获取用户信息失败", err)
			}
			return
		}
		shared.InternalError(c, "获取用户信息失败", err)
		return
	}

	// 构建公开信息响应（不包含敏感信息）
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	publicProfile := dto.PublicUserProfileResponse{
		UserID:    resp.User.ID.Hex(),
		Username:  resp.User.Username,
		Avatar:    resp.User.Avatar,
		Nickname:  resp.User.Nickname,
		Bio:       resp.User.Bio,
		Role:      role,
		CreatedAt: resp.User.CreatedAt,
	}

	shared.Success(c, http.StatusOK, "获取成功", publicProfile)
}

// GetUserBooks 获取用户的作品列表
//
//	@Summary		获取用户作品列表
//	@Description	获取指定用户的已发布作品列表
//	@Tags			用户管理-公开信息
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"用户ID"
//	@Param			page	query		int		false	"页码"		default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Param			status	query		string	false	"状态筛选"	Enums(published, completed)
//	@Success		200		{object}	shared.APIResponse{data=dto.GetUserBooksResponse}
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/users/{id}/books [get]
func (h *PublicUserHandler) GetUserBooks(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		shared.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	_ = c.Query("status") // status 参数预留，可用于进一步筛选

	// 如果没有设置BookstoreService，返回空列表
	if h.bookstoreService == nil {
		shared.Success(c, http.StatusOK, "获取成功", dto.GetUserBooksResponse{
			Books: []map[string]interface{}{},
			Total: 0,
			Page:  page,
			Size:  size,
		})
		return
	}

	// 使用类型断言调用BookstoreService的GetBooksByAuthorID方法
	type BookstoreService interface {
		GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) (interface{}, int64, error)
	}

	bookstoreSvc, ok := h.bookstoreService.(BookstoreService)
	if !ok {
		shared.InternalError(c, "服务配置错误", nil)
		return
	}

	// 调用BookstoreService查询用户的已发布作品
	booksRaw, total, err := bookstoreSvc.GetBooksByAuthorID(c.Request.Context(), userID, page, size)
	if err != nil {
		shared.InternalError(c, "获取用户作品失败", err)
		return
	}

	// 将返回的书籍转换为响应格式
	response := dto.GetUserBooksResponse{
		Books: booksRaw,
		Total: int(total),
		Page:  page,
		Size:  size,
	}

	shared.Success(c, http.StatusOK, "获取成功", response)
}
