package finance

import (
	"github.com/gin-gonic/gin"

	financeApi "Qingyu_backend/api/v1/finance"
	"Qingyu_backend/middleware"
)

// RegisterFinanceRoutes 注册所有财务相关路由
func RegisterFinanceRoutes(r *gin.RouterGroup, walletAPI *financeApi.WalletAPI, membershipAPI *financeApi.MembershipAPI, authorRevenueAPI *financeApi.AuthorRevenueAPI) {
	// 财务路由需要认证
	financeGroup := r.Group("/finance")
	financeGroup.Use(middleware.JWTAuth())
	financeGroup.Use(middleware.RateLimitMiddleware(50, 60))
	{
		// ========== 钱包相关 ==========
		if walletAPI != nil {
			walletGroup := financeGroup.Group("/wallet")
			{
				// 获取钱包余额
				walletGroup.GET("/balance", walletAPI.GetBalance)

				// 获取钱包详情
				walletGroup.GET("/detail", walletAPI.GetWallet)

				// 充值
				walletGroup.POST("/recharge", walletAPI.Recharge)

				// 消费
				walletGroup.POST("/consume", walletAPI.Consume)

				// 转账
				walletGroup.POST("/transfer", walletAPI.Transfer)

				// 获取交易记录
				walletGroup.GET("/transactions", walletAPI.GetTransactions)

				// 申请提现
				walletGroup.POST("/withdraw", walletAPI.RequestWithdraw)

				// 获取提现申请列表
				walletGroup.GET("/withdraws", walletAPI.GetWithdrawRequests)
			}
		}

		// ========== 会员系统 ==========
		if membershipAPI != nil {
			membershipGroup := financeGroup.Group("/membership")
			{
				// 公开路由 - 获取套餐列表
				membershipGroup.GET("/plans", membershipAPI.GetPlans)

				// 需要认证的路由
				membershipGroup.GET("/status", membershipAPI.GetStatus)
				membershipGroup.POST("/subscribe", membershipAPI.Subscribe)
				membershipGroup.POST("/cancel", membershipAPI.Cancel)
				membershipGroup.PUT("/renew", membershipAPI.Renew)
				membershipGroup.GET("/benefits", membershipAPI.GetBenefits)
				membershipGroup.GET("/usage", membershipAPI.GetUsage)
				membershipGroup.GET("/cards", membershipAPI.ListCards)
				membershipGroup.POST("/cards/activate", membershipAPI.ActivateCard)
			}
		}

		// ========== 作者收入 ==========
		if authorRevenueAPI != nil {
			authorGroup := financeGroup.Group("/author")
			{
				// 收入查询
				authorGroup.GET("/earnings", authorRevenueAPI.GetEarnings)
				authorGroup.GET("/earnings/:bookId", authorRevenueAPI.GetBookEarnings)
				authorGroup.GET("/revenue-details", authorRevenueAPI.GetRevenueDetails)
				authorGroup.GET("/revenue-statistics", authorRevenueAPI.GetRevenueStatistics)

				// 提现管理
				authorGroup.GET("/withdrawals", authorRevenueAPI.GetWithdrawals)
				authorGroup.POST("/withdraw", authorRevenueAPI.Withdraw)

				// 结算管理
				authorGroup.GET("/settlements", authorRevenueAPI.GetSettlements)
				authorGroup.GET("/settlements/:id", authorRevenueAPI.GetSettlement)

				// 税务信息
				authorGroup.GET("/tax-info", authorRevenueAPI.GetTaxInfo)
				authorGroup.PUT("/tax-info", authorRevenueAPI.UpdateTaxInfo)
			}
		}
	}
}
