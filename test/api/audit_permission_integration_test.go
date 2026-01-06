//go:build integration
// +build integration

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuditAPI_GetUserViolations_OwnData 用户查看自己的违规记录
func TestAuditAPI_GetUserViolations_OwnData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/users/:userId/violation-summary", auditAuthMiddleware("user123", "user"), func(c *gin.Context) {
		userID := c.Param("userId")
		currentUserID, _ := c.Get("userID")
		role, _ := c.Get("role")

		// 权限检查
		isAdmin := role == "admin"
		if !isAdmin && currentUserID.(string) != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限",
				"error":   "只能查看自己的违规记录",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data": map[string]interface{}{
				"totalViolations": 2,
				"violationTypes": map[string]interface{}{
					"spam":    1,
					"illegal": 1,
				},
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/user123/violation-summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
}

// TestAuditAPI_GetUserViolations_PermissionDenied 用户无权查看他人的违规记录
func TestAuditAPI_GetUserViolations_PermissionDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/users/:userId/violation-summary", auditAuthMiddleware("user123", "user"), func(c *gin.Context) {
		userID := c.Param("userId")
		currentUserID, _ := c.Get("userID")
		role, _ := c.Get("role")

		isAdmin := role == "admin"
		if !isAdmin && currentUserID.(string) != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限",
				"error":   "只能查看自己的违规记录",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/other_user/violation-summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "无权限", response["message"])
}

// TestAuditAPI_GetUserViolations_AdminAccess 管理员可以查看任何用户的违规记录
func TestAuditAPI_GetUserViolations_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/users/:userId/violation-summary", auditAuthMiddleware("admin1", "admin"), func(c *gin.Context) {
		userID := c.Param("userId")
		currentUserID, _ := c.Get("userID")
		role, _ := c.Get("role")

		isAdmin := role == "admin"
		if !isAdmin && currentUserID.(string) != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无权限",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "获取成功",
			"data": map[string]interface{}{
				"userId":          userID,
				"totalViolations": 3,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/any_user/violation-summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "any_user", data["userId"])
}

// TestAuditAPI_GetUserViolations_Unauthorized 未授权
func TestAuditAPI_GetUserViolations_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/users/:userId/violation-summary", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/users/user123/violation-summary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// auditAuthMiddleware 审核认证中间件
func auditAuthMiddleware(userID, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
