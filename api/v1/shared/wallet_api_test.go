package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWalletAPI_GetBalance_InvalidUserIDReturnsUnauthorized(t *testing.T) {
	api := NewWalletAPI(nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Set("user_id", 123)

	assert.NotPanics(t, func() {
		api.GetBalance(c)
	})
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "用户信息格式错误", resp["message"])
}

func TestWalletAPI_ProtectedActions_InvalidUserIDReturnsUnauthorized(t *testing.T) {
	for _, action := range []struct {
		name   string
		handle func(*WalletAPI, *gin.Context)
	}{
		{name: "获取钱包", handle: (*WalletAPI).GetWallet},
		{name: "充值", handle: (*WalletAPI).Recharge},
		{name: "消费", handle: (*WalletAPI).Consume},
		{name: "转账", handle: (*WalletAPI).Transfer},
		{name: "交易记录", handle: (*WalletAPI).GetTransactions},
		{name: "申请提现", handle: (*WalletAPI).RequestWithdraw},
		{name: "提现申请", handle: (*WalletAPI).GetWithdrawRequests},
	} {
		t.Run(action.name, func(t *testing.T) {
			api := NewWalletAPI(nil)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Set("user_id", 123)

			assert.NotPanics(t, func() {
				action.handle(api, c)
			})
			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			resp := decodeSharedAPIResponse(t, recorder)
			assert.Equal(t, "用户信息格式错误", resp["message"])
		})
	}
}
