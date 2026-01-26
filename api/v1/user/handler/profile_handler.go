package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
)

// ProfileHandler 个人信息处理器
type ProfileHandler struct {
	userService userServiceInterface.UserService
}

// NewProfileHandler 创建个人信息处理器实例
func NewProfileHandler(userService userServiceInterface.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// GetProfile 获取当前用户信息
//
//	@Summary		获取当前用户信息
//	@Description	获取当前登录用户的详细信息
//	@Tags			用户管理-个人信息
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200		{object}	shared.APIResponse{data=dto.UserProfileResponse}
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		404		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// 从Context中获取当前用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.GetUserRequest{
		ID: userID.(string),
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

	// 直接返回 Service 层的 DTO
	shared.Success(c, http.StatusOK, "获取成功", resp.User)
}

// UpdateProfile 更新当前用户信息
//
//	@Summary		更新当前用户信息
//	@Description	更新当前登录用户的个人信息
//	@Tags			用户管理-个人信息
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		dto.UpdateProfileRequest	true	"更新信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	var req dto.UpdateProfileRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	// 调用Service层
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID.(string),
		Updates: updates,
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "更新失败", serviceErr.Message)
			default:
				shared.InternalError(c, "更新失败", err)
			}
			return
		}
		shared.InternalError(c, "更新失败", err)
		return
	}

	// 返回更新后的用户信息
	shared.Success(c, http.StatusOK, "更新成功", updatedUser)
}

// ChangePassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改当前用户的登录密码
//	@Tags			用户管理-个人信息
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		dto.ChangePasswordRequest	true	"密码信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/password [put]
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	var req dto.ChangePasswordRequest
	if !shared.ValidateRequest(c, &req) {
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.UpdatePasswordRequest{
		ID:          userID.(string),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	_, err := h.userService.UpdatePassword(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			case serviceInterfaces.ErrorTypeUnauthorized:
				shared.Unauthorized(c, "旧密码错误")
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "修改密码失败", serviceErr.Message)
			default:
				shared.InternalError(c, "修改密码失败", err)
			}
			return
		}
		shared.InternalError(c, "修改密码失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "密码修改成功", nil)
}
