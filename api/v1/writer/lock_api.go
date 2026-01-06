package writer

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/lock"
)

// LockAPI 文档锁定API
type LockAPI struct {
	lockService lock.DocumentLockService
}

// NewLockAPI 创建文档锁定API实例
func NewLockAPI(lockService lock.DocumentLockService) *LockAPI {
	return &LockAPI{
		lockService: lockService,
	}
}

// LockDocument 锁定文档
//
//	@Summary		锁定文档
//	@Description	锁定文档以进行编辑，防止多人同时编辑冲突
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"文档ID"
//	@Param			request	body	LockDocumentRequest	true	"锁定请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		409		{object}	shared.APIResponse	"文档已被锁定"
//	@Router			/api/v1/writer/documents/{id}/lock [post]
func (api *LockAPI) LockDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	var req LockDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userName := ""
	if name, exists := c.Get("userName"); exists {
		userName = name.(string)
	}

	// 获取设备ID
	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = "unknown"
	}

	// 锁定文档（默认30分钟）
	ttl := 30 * time.Minute
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	lock, err := api.lockService.LockDocument(c.Request.Context(), documentID, userID.(string), userName, deviceID, req.AutoExtend, ttl)
	if err != nil {
		if isLockedError(err) {
			shared.Error(c, http.StatusConflict, "文档已被锁定", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "锁定失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "锁定成功", lock)
}

// UnlockDocument 解锁文档
//
//	@Summary		解锁文档
//	@Description	解锁文档，允许其他用户编辑
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"文档ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse	"无权解锁"
//	@Router			/api/v1/writer/documents/{id}/lock [delete]
func (api *LockAPI) UnlockDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	// 解锁文档
	if err := api.lockService.UnlockDocument(c.Request.Context(), documentID, userID.(string)); err != nil {
		if isPermissionError(err) {
			shared.Error(c, http.StatusForbidden, "无权操作", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "解锁失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "解锁成功", nil)
}

// RefreshLock 刷新锁（心跳）
//
//	@Summary		刷新锁
//	@Description	发送心跳以保持文档锁定状态
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"文档ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/lock/refresh [put]
func (api *LockAPI) RefreshLock(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	// 刷新锁（延长30分钟）
	ttl := 30 * time.Minute

	if err := api.lockService.RefreshLock(c.Request.Context(), documentID, userID.(string), ttl); err != nil {
		shared.Error(c, http.StatusInternalServerError, "刷新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "刷新成功", nil)
}

// GetLockStatus 获取锁状态
//
//	@Summary		获取锁状态
//	@Description	获取文档当前的锁定状态
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"文档ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/lock/status [get]
func (api *LockAPI) GetLockStatus(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	// 获取用户信息
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	status, err := api.lockService.GetLockStatus(c.Request.Context(), documentID, userID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", status)
}

// ForceUnlock 强制解锁（管理员）
//
//	@Summary		强制解锁
//	@Description	管理员强制解锁文档
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"文档ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse	"需要管理员权限"
//	@Router			/api/v1/writer/documents/{id}/lock/force [post]
func (api *LockAPI) ForceUnlock(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	// 检查管理员权限
	isAdmin := false
	if role, exists := c.Get("userRole"); exists {
		isAdmin = (role == "admin" || role == "super_admin")
	}

	if !isAdmin {
		shared.Error(c, http.StatusForbidden, "权限不足", "需要管理员权限")
		return
	}

	// 强制解锁
	if err := api.lockService.ForceUnlock(c.Request.Context(), documentID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "强制解锁失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "强制解锁成功", nil)
}

// ExtendLock 延长锁时间
//
//	@Summary		延长锁
//	@Description	延长文档锁定时间
//	@Tags			Writer-Lock
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string				true	"文档ID"
//	@Param			request	body	ExtendLockRequest	true	"延长请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/lock/extend [post]
func (api *LockAPI) ExtendLock(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	var req ExtendLockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	// 延长锁
	ttl := time.Duration(req.TTL) * time.Second
	if req.TTL <= 0 {
		ttl = 30 * time.Minute // 默认30分钟
	}

	if err := api.lockService.ExtendLock(c.Request.Context(), documentID, userID.(string), ttl); err != nil {
		shared.Error(c, http.StatusInternalServerError, "延长失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "延长成功", nil)
}

// LockDocumentRequest 锁定请求
type LockDocumentRequest struct {
	AutoExtend bool `json:"autoExtend"` // 是否自动续期
	TTL        int  `json:"ttl"`        // 锁定时长（秒），0表示默认30分钟
}

// ExtendLockRequest 延长请求
type ExtendLockRequest struct {
	TTL int `json:"ttl" binding:"required,min=1"` // 延长时长（秒）
}

// isLockedError 检查是否是锁定错误
func isLockedError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return contains(errMsg, "locked") || contains(errMsg, "锁定")
}

// isPermissionError 检查是否是权限错误
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return contains(errMsg, "permission") || contains(errMsg, "权限") || contains(errMsg, "denied")
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (
		strings.HasPrefix(strings.ToLower(s), strings.ToLower(substr)) ||
		strings.HasSuffix(strings.ToLower(s), strings.ToLower(substr)) ||
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))))
}
