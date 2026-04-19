package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/config"
	"Qingyu_backend/pkg/emailcode"
)

func TestValidateRegisterVerificationCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("code manager disabled", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ok := ValidateRegisterVerificationCode(c, nil, "test@example.com", "")

		assert.True(t, ok)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing verification code", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		previousConfig := config.GlobalConfig
		config.GlobalConfig = &config.Config{
			Email: &config.EmailConfig{
				Enabled:   true,
				FixedCode: "123456",
			},
		}
		defer func() {
			config.GlobalConfig = previousConfig
		}()
		manager := emailcode.NewManager()

		ok := ValidateRegisterVerificationCode(c, manager, "test@example.com", "")

		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "请先填写邮箱验证码")
	})
}

func TestSendRegisterVerificationCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("disabled manager returns bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/send-verification-code", nil)
		manager := emailcode.NewManager()

		ok := SendRegisterVerificationCode(c.Request.Context(), c, manager, "test@example.com")

		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "邮箱验证码功能未启用")
	})
}
