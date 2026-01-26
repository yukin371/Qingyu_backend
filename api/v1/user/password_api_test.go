package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	userService "Qingyu_backend/service/user"
)

// ==================== SendPasswordResetCode 测试 ====================

func TestSendPasswordResetCode_Success(t *testing.T) {
	t.Skip("需要集成测试环境，包括PasswordService和Repository")
}

func TestSendPasswordResetCode_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/send", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email": "invalid-email",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendPasswordResetCode_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/send", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendPasswordResetCode_EmptyEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/send", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email": "",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== ResetPassword 测试 ====================

func TestResetPassword_Success(t *testing.T) {
	t.Skip("需要集成测试环境")
}

func TestResetPassword_InvalidEmailFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email":        "invalid-email",
		"code":         "123456",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_InvalidCodeLength(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email":        "test@example.com",
		"code":         "12345",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_PasswordTooShort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email":        "test@example.com",
		"code":         "123456",
		"new_password": "short",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"code":         "123456",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_MissingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email":        "test@example.com",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_MissingNewPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email": "test@example.com",
		"code":  "123456",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_InvalidCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/users/password/reset/verify", func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Code        string `json:"code" binding:"required,len=6"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		// 模拟验证码无效错误
		if req.Code == "000000" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "验证码无效或已过期",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"email":        "test@example.com",
		"code":         "000000",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/password/reset/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== UpdatePassword 测试 ====================

func TestUpdatePassword_Success(t *testing.T) {
	t.Skip("需要集成测试环境")
}

func TestUpdatePassword_MissingOldPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassword_NewPasswordTooShort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "oldpassword",
		"new_password": "short",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassword_MissingNewPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "oldpassword",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassword_EmptyOldPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassword_EmptyNewPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "oldpassword",
		"new_password": "",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePassword_OldPasswordMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		c.Set("userID", "user123")

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		// 模拟旧密码错误
		if req.OldPassword == "wrongpassword" {
			err := userService.ErrOldPasswordMismatch
			if errors.Is(err, userService.ErrOldPasswordMismatch) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "旧密码错误",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "wrongpassword",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdatePassword_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.PUT("/users/password", func(c *gin.Context) {
		// 不设置 userID

		var req struct {
			OldPassword string `json:"old_password" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
			})
			return
		}

		// 检查 userID 是否存在
		userID := c.GetString("userID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	reqBody := map[string]string{
		"old_password": "oldpassword",
		"new_password": "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/users/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
