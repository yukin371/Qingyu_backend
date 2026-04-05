package handler

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/user/dto"
	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/pkg/response"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
)

// PublicUserHandler 公开用户信息处理器
type PublicUserHandler struct {
	userService      userServiceInterface.UserService
	bookstoreService BookstoreService // 可选依赖
}

// BookstoreService 公开用户作品查询端口（最小依赖）
type BookstoreService interface {
	GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
}

// NewPublicUserHandler 创建公开用户信息处理器实例
func NewPublicUserHandler(userService userServiceInterface.UserService) *PublicUserHandler {
	return &PublicUserHandler{
		userService: userService,
	}
}

// SetBookstoreService 设置BookstoreService（可选依赖注入）
func (h *PublicUserHandler) SetBookstoreService(bookstoreSvc BookstoreService) {
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
//	@Success		200		{object}	response.APIResponse{data=dto.UserProfileResponse}
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/users/{id} [get]
func (h *PublicUserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
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
				response.NotFound(c, "用户不存在")
			default:
				response.InternalError(c, err)
			}
			return
		}
		response.InternalError(c, err)
		return
	}

	// 构建公开信息响应（不包含敏感信息）
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	// 解析 CreatedAt 字符串为 time.Time
	createdAt, _ := time.Parse(time.RFC3339, resp.User.CreatedAt)
	publicProfile := dto.PublicUserProfileResponse{
		UserID:    resp.User.ID,
		Username:  resp.User.Username,
		Avatar:    resp.User.Avatar,
		Nickname:  resp.User.Nickname,
		Bio:       resp.User.Bio,
		Role:      role,
		CreatedAt: createdAt,
	}

	response.Success(c, publicProfile)
}

// GetUserProfile 获取用户详细资料（公开访问）
//
//	@Summary		获取用户详细资料
//	@Description	获取指定用户的详细资料，用于展示用户主页
//	@Tags			用户管理-公开信息
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"用户ID"
//	@Success		200		{object}	response.APIResponse{data=dto.PublicUserProfileResponse}
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/users/{id}/profile [get]
func (h *PublicUserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
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
				response.NotFound(c, "用户不存在")
			default:
				response.InternalError(c, err)
			}
			return
		}
		response.InternalError(c, err)
		return
	}

	// 构建公开信息响应（不包含敏感信息）
	role := ""
	if len(resp.User.Roles) > 0 {
		role = resp.User.Roles[0]
	}
	// 解析 CreatedAt 字符串为 time.Time
	createdAt, _ := time.Parse(time.RFC3339, resp.User.CreatedAt)
	publicProfile := dto.PublicUserProfileResponse{
		UserID:    resp.User.ID,
		Username:  resp.User.Username,
		Avatar:    resp.User.Avatar,
		Nickname:  resp.User.Nickname,
		Bio:       resp.User.Bio,
		Role:      role,
		CreatedAt: createdAt,
	}

	response.Success(c, publicProfile)
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
//	@Success		200		{object}	response.APIResponse{data=dto.GetUserBooksResponse}
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/users/{id}/books [get]
func (h *PublicUserHandler) GetUserBooks(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	_ = c.Query("status") // status 参数预留，可用于进一步筛选

	// 如果没有设置BookstoreService，返回空列表
	if h.bookstoreService == nil {
		response.Success(c, dto.GetUserBooksResponse{
			Books: []map[string]interface{}{},
			Total: 0,
			Page:  page,
			Size:  size,
		})
		return
	}

	// 调用BookstoreService查询用户的已发布作品
	booksRaw, total, err := h.bookstoreService.GetBooksByAuthorID(c.Request.Context(), userID, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 将返回的书籍转换为响应格式
	resp := dto.GetUserBooksResponse{
		Books: booksRaw,
		Total: int(total),
		Page:  page,
		Size:  size,
	}

	response.Success(c, resp)
}

// GetBatchUsers 批量获取用户信息（公开访问）
//
//	@Summary		批量获取用户信息
//	@Description	根据用户ID列表批量获取用户的公开信息
//	@Tags			用户管理-公开信息
//	@Accept			json
//	@Produce		json
//	@Param			ids		query		string	true	"用户ID列表，用逗号分隔"
//	@Success		200		{object}	response.APIResponse{data=[]dto.PublicUserProfileResponse}
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/user/users/batch [get]
func (h *PublicUserHandler) GetBatchUsers(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		response.BadRequest(c, "参数错误", "ids参数不能为空")
		return
	}

	// 解析ID列表
	idList := parseIDs(idsParam)
	if len(idList) == 0 {
		response.BadRequest(c, "参数错误", "未提供有效的用户ID")
		return
	}

	// 限制批量大小，防止滥用
	maxBatchSize := 50
	if len(idList) > maxBatchSize {
		response.BadRequest(c, "参数错误", "批量大小不能超过50")
		return
	}

	// 批量获取用户信息
	users := make([]*dto.PublicUserProfileResponse, 0, len(idList))
	for _, userID := range idList {
		serviceReq := &userServiceInterface.GetUserRequest{
			ID: userID,
		}

		resp, err := h.userService.GetUser(c.Request.Context(), serviceReq)
		if err != nil {
			// 单个用户获取失败不影响其他用户
			continue
		}

		role := ""
		if len(resp.User.Roles) > 0 {
			role = resp.User.Roles[0]
		}
		createdAt, _ := time.Parse(time.RFC3339, resp.User.CreatedAt)
		publicProfile := dto.PublicUserProfileResponse{
			UserID:    resp.User.ID,
			Username:  resp.User.Username,
			Avatar:    resp.User.Avatar,
			Nickname:  resp.User.Nickname,
			Bio:       resp.User.Bio,
			Role:      role,
			CreatedAt: createdAt,
		}
		users = append(users, &publicProfile)
	}

	response.Success(c, gin.H{
		"users": users,
		"count": len(users),
	})
}

// parseIDs 解析逗号分隔的ID列表
func parseIDs(idsParam string) []string {
	if idsParam == "" {
		return nil
	}
	parts := strings.Split(idsParam, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
