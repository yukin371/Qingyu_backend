package shared

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/finance/wallet"
)

// WalletAPI 钱包服务API处理器
type WalletAPI struct {
	walletService wallet.WalletService
}

// NewWalletAPI 创建钱包API实例
func NewWalletAPI(walletService wallet.WalletService) *WalletAPI {
	return &WalletAPI{
		walletService: walletService,
	}
}

// GetBalance 查询余额
//
//	@Summary		查询余额
//	@Description	查询用户钱包余额
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/wallet/balance [get]
func (api *WalletAPI) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	balance, err := api.walletService.GetBalance(c.Request.Context(), userID.(string))
	if err != nil {
		InternalError(c, "查询余额失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "查询余额成功",
		Data:    balance,
	})
}

// GetWallet 获取钱包信息
//
//	@Summary		获取钱包信息
//	@Description	获取用户完整钱包信息
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Failure		401	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/shared/wallet [get]
func (api *WalletAPI) GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	walletInfo, err := api.walletService.GetWallet(c.Request.Context(), userID.(string))
	if err != nil {
		InternalError(c, "获取钱包信息失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取钱包信息成功",
		Data:    walletInfo,
	})
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount float64 `json:"amount" binding:"required" validate:"positive_amount,amount_range"`
	Method string  `json:"method" binding:"required,oneof=alipay wechat bank"`
}

// Recharge 充值
//
//	@Summary		钱包充值
//	@Description	用户钱包充值
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		RechargeRequest	true	"充值信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/wallet/recharge [post]
func (api *WalletAPI) Recharge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req RechargeRequest
	if !ValidateRequest(c, &req) {
		return
	}

	transaction, err := api.walletService.Recharge(c.Request.Context(), userID.(string), req.Amount, req.Method)
	if err != nil {
		InternalError(c, "充值失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "充值成功",
		Data:    transaction,
	})
}

// ConsumeRequest 消费请求
type ConsumeRequest struct {
	Amount float64 `json:"amount" binding:"required" validate:"positive_amount,amount_range"`
	Reason string  `json:"reason" binding:"required,min=1,max=200"`
}

// Consume 消费
//
//	@Summary		钱包消费
//	@Description	用户钱包消费
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ConsumeRequest	true	"消费信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/wallet/consume [post]
func (api *WalletAPI) Consume(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req ConsumeRequest
	if !ValidateRequest(c, &req) {
		return
	}

	transaction, err := api.walletService.Consume(c.Request.Context(), userID.(string), req.Amount, req.Reason)
	if err != nil {
		InternalError(c, "消费失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "消费成功",
		Data:    transaction,
	})
}

// TransferRequest 转账请求
type TransferRequest struct {
	ToUserID string  `json:"to_user_id" binding:"required,min=1"`
	Amount   float64 `json:"amount" binding:"required" validate:"positive_amount,amount_range"`
	Reason   string  `json:"reason" validate:"omitempty,max=200"`
}

// Transfer 转账
//
//	@Summary		用户转账
//	@Description	向其他用户转账
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		TransferRequest	true	"转账信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/wallet/transfer [post]
func (api *WalletAPI) Transfer(c *gin.Context) {
	fromUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req TransferRequest
	if !ValidateRequest(c, &req) {
		return
	}

	transaction, err := api.walletService.Transfer(c.Request.Context(), fromUserID.(string), req.ToUserID, req.Amount, req.Reason)
	if err != nil {
		InternalError(c, "转账失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "转账成功",
		Data:    transaction,
	})
}

// GetTransactions 查询交易记录
//
//	@Summary		查询交易记录
//	@Description	查询用户交易记录列表
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			type		query		string	false	"交易类型"
//	@Success 200 {object} APIResponse
//	@Failure		401			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/shared/wallet/transactions [get]
func (api *WalletAPI) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	transactionType := c.Query("type")

	req := &wallet.ListTransactionsRequest{
		TransactionType: transactionType,
		Page:            page,
		PageSize:        pageSize,
	}

	transactions, err := api.walletService.ListTransactions(c.Request.Context(), userID.(string), req)
	if err != nil {
		InternalError(c, "查询交易记录失败", err)
		return
	}

	c.JSON(http.StatusOK, PaginatedResponseHelper(
		transactions,
		int64(len(transactions)),
		page,
		pageSize,
		"查询交易记录成功",
	))
}

// WithdrawRequest 提现请求
type WithdrawRequest struct {
	Amount  float64 `json:"amount" binding:"required" validate:"positive_amount,amount_range"`
	Account string  `json:"account" binding:"required" validate:"withdraw_account"`
}

// RequestWithdraw 申请提现
//
//	@Summary		申请提现
//	@Description	用户申请提现
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		WithdrawRequest	true	"提现信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/wallet/withdraw [post]
func (api *WalletAPI) RequestWithdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req WithdrawRequest
	if !ValidateRequest(c, &req) {
		return
	}

	withdrawReq, err := api.walletService.RequestWithdraw(c.Request.Context(), userID.(string), req.Amount, req.Account)
	if err != nil {
		InternalError(c, "申请提现失败", err)
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "申请提现成功",
		Data:    withdrawReq,
	})
}

// GetWithdrawRequests 查询提现申请
//
//	@Summary		查询提现申请
//	@Description	查询用户提现申请列表
//	@Tags			钱包
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			status		query		string	false	"状态"
//	@Success 200 {object} APIResponse
//	@Failure		401			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/shared/wallet/withdrawals [get]
func (api *WalletAPI) GetWithdrawRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	req := &wallet.ListWithdrawRequestsRequest{
		UserID:   userID.(string),
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	requests, err := api.walletService.ListWithdrawRequests(c.Request.Context(), req)
	if err != nil {
		InternalError(c, "查询提现申请失败", err)
		return
	}

	c.JSON(http.StatusOK, PaginatedResponseHelper(
		requests,
		int64(len(requests)),
		page,
		pageSize,
		"查询提现申请成功",
	))
}
