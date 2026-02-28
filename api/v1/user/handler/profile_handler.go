package handler

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/api/v1/user/dto"
	sharedStorage "Qingyu_backend/service/shared/storage"
)

// ProfileHandler 个人信息处理器
type ProfileHandler struct {
	userService    userServiceInterface.UserService
	storageService sharedStorage.StorageService
}

// NewProfileHandler 创建个人信息处理器实例
func NewProfileHandler(userService userServiceInterface.UserService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
	}
}

// SetStorageService 设置存储服务（可选依赖）
func (h *ProfileHandler) SetStorageService(storageSvc sharedStorage.StorageService) {
	h.storageService = storageSvc
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
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "请求参数错误", err.Error())
		return
	}

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
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	// 处理生日字段（RFC3339格式转换为time.Time）
	if req.Birthday != nil && *req.Birthday != "" {
		birthday, err := time.Parse(time.RFC3339, *req.Birthday)
		if err != nil {
			shared.BadRequest(c, "生日格式错误，请使用RFC3339格式（如1990-01-01T00:00:00Z）", err.Error())
			return
		}
		updates["birthday"] = birthday
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

// UpdatePassword 修改密码
//
//	@Summary		修改密码
//	@Description	修改当前用户的登录密码
//	@Tags			用户管理-个人信息
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		dto.UpdatePasswordRequest	true	"密码信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/password [put]
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "请求参数错误", err.Error())
		return
	}

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

// UploadAvatar 上传头像
//
//	@Summary		上传头像
//	@Description	上传用户头像图片，支持JPG、PNG、JPEG格式，最大5MB
//	@Tags			用户管理-个人信息
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"头像文件"
//	@Success		200		{object}	shared.APIResponse{data=dto.UploadAvatarResponse}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		413		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/avatar [post]
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	// 检查存储服务是否可用
	if h.storageService == nil {
		shared.InternalError(c, "存储服务不可用", nil)
		return
	}

	// 获取上传的文件
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		shared.BadRequest(c, "请选择文件", err.Error())
		return
	}

	// 验证文件大小（最大5MB）
	const maxSize = 5 * 1024 * 1024
	if fileHeader.Size > maxSize {
		shared.BadRequest(c, "文件大小不能超过5MB", "")
		return
	}

	// 验证文件类型
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif"}
	contentType := fileHeader.Header.Get("Content-Type")
	if !isAllowedType(contentType, allowedTypes) {
		shared.BadRequest(c, "只支持JPG、PNG、GIF格式", "")
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		shared.InternalError(c, "打开文件失败", err)
		return
	}
	defer file.Close()

	// 上传到存储服务
	uploadReq := &sharedStorage.UploadRequest{
		File:        file,
		Filename:    fileHeader.Filename,
		ContentType: contentType,
		Size:        fileHeader.Size,
		UserID:      userID.(string),
		IsPublic:    true,
		Category:    "avatar",
	}

	fileInfo, err := h.storageService.Upload(c.Request.Context(), uploadReq)
	if err != nil {
		shared.InternalError(c, "上传失败", err)
		return
	}

	// 更新用户头像URL（使用Path字段构建完整URL）
	avatarURL := fileInfo.Path
	// 如果不是完整URL，添加CDN前缀（根据实际配置调整）
	if !strings.HasPrefix(avatarURL, "http://") && !strings.HasPrefix(avatarURL, "https://") {
		avatarURL = "/cdn/" + avatarURL
	}

	updates := map[string]interface{}{
		"avatar": avatarURL,
	}
	serviceReq := &userServiceInterface.UpdateUserRequest{
		ID:      userID.(string),
		Updates: updates,
	}

	_, err = h.userService.UpdateUser(c.Request.Context(), serviceReq)
	if err != nil {
		// 上传成功但更新失败，记录警告
		// 注意：文件已经上传，但没有更新用户记录
		shared.InternalError(c, "上传成功但更新用户信息失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "上传成功", dto.UploadAvatarResponse{
		AvatarURL: avatarURL,
		Message:   "头像上传成功",
	})
}

// isAllowedType 检查文件类型是否允许
func isAllowedType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
	}
	return false
}

// DowngradeRoleRequest 降级角色请求
type DowngradeRoleRequest struct {
	TargetRole string `json:"target_role" binding:"required,oneof=reader author"` // 只能降级到reader或author
	Confirm    bool   `json:"confirm" binding:"required"`                         // 必须确认
}

// DowngradeRole 降级用户角色
//
//	@Summary		降级用户角色
//	@Description	将当前用户角色降级到指定角色（author只能降级到reader，admin可以降级到author或reader）
//	@Tags			用户管理-个人信息
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		DowngradeRoleRequest	true	"降级请求"
//	@Success		200		{object}	shared.APIResponse{data=map[string]interface{}}
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/user/role/downgrade [post]
func (h *ProfileHandler) DowngradeRole(c *gin.Context) {
	// 从Context中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "未认证")
		return
	}

	// 解析请求体
	var req DowngradeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 验证确认标志
	if !req.Confirm {
		shared.BadRequest(c, "请确认降级操作", "")
		return
	}

	// 验证目标角色
	if req.TargetRole != "reader" && req.TargetRole != "author" {
		shared.BadRequest(c, "目标角色无效，只能降级为reader或author", "")
		return
	}

	// 调用Service层
	serviceReq := &userServiceInterface.DowngradeRoleRequest{
		UserID:     userID.(string),
		TargetRole: req.TargetRole,
		Confirm:    req.Confirm,
	}

	resp, err := h.userService.DowngradeRole(c.Request.Context(), serviceReq)
	if err != nil {
		if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
			switch serviceErr.Type {
			case serviceInterfaces.ErrorTypeValidation:
				shared.BadRequest(c, "参数错误", serviceErr.Message)
			case serviceInterfaces.ErrorTypeBusiness:
				// 业务错误：如"已经是读者，无法降级"应该返回403
				c.JSON(403, shared.ErrorResponse{
					Code:      1003,
					Message:   serviceErr.Message,
					Timestamp: time.Now().UnixMilli(),
				})
			case serviceInterfaces.ErrorTypeNotFound:
				shared.NotFound(c, "用户不存在")
			default:
				shared.InternalError(c, "降级失败", err)
			}
			return
		}
		shared.InternalError(c, "降级失败", err)
		return
	}

	// 返回成功响应
	shared.Success(c, 200, "降级成功", map[string]interface{}{
		"current_roles": resp.CurrentRoles,
	})
}
