package finance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	sharedApi "Qingyu_backend/api/v1/shared"
	financeModel "Qingyu_backend/models/finance"
	financeService "Qingyu_backend/service/finance"
)

// AuthorRevenueAPI 作者收入API处理器
type AuthorRevenueAPI struct {
	revenueService financeService.AuthorRevenueService
}

// NewAuthorRevenueAPI 创建作者收入API实例
func NewAuthorRevenueAPI(revenueService financeService.AuthorRevenueService) *AuthorRevenueAPI {
	return &AuthorRevenueAPI{
		revenueService: revenueService,
	}
}

// GetEarnings 获取作者收入列表
//
//	@Summary		获取作者收入列表
//	@Description	获取作者的收入记录列表
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/earnings [get]
func (api *AuthorRevenueAPI) GetEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	earnings, total, err := api.revenueService.GetEarnings(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取收入列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		earnings,
		total,
		page,
		pageSize,
		"获取收入列表成功",
	))
}

// GetBookEarnings 获取某本书的收入
//
//	@Summary		获取某本书的收入
//	@Description	获取指定书籍的收入记录
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			bookId	path		string	true	"书籍ID"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/earnings/{bookId} [get]
func (api *AuthorRevenueAPI) GetBookEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	earnings, total, err := api.revenueService.GetBookEarnings(c.Request.Context(), userID.(string), bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取书籍收入失败: " + err.Error(),
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		earnings,
		total,
		page,
		pageSize,
		"获取书籍收入成功",
	))
}

// GetWithdrawals 获取提现记录
//
//	@Summary		获取提现记录
//	@Description	获取用户的提现申请记录
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/withdrawals [get]
func (api *AuthorRevenueAPI) GetWithdrawals(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	withdrawals, total, err := api.revenueService.GetWithdrawals(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取提现记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		withdrawals,
		total,
		page,
		pageSize,
		"获取提现记录成功",
	))
}

// WithdrawRequest 提现申请请求
type WithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Method      string  `json:"method" binding:"required,oneof=alipay wechat bank"`
	AccountType string  `json:"account_type" binding:"required"`
	AccountName string  `json:"account_name" binding:"required"`
	AccountNo   string  `json:"account_no" binding:"required"`
	BankName    string  `json:"bank_name,omitempty"`
	BranchName  string  `json:"branch_name,omitempty"`
}

// Withdraw 申请提现
//
//	@Summary		申请提现
//	@Description	作者申请提现
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		WithdrawRequest	true	"提现信息"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/withdraw [post]
func (api *AuthorRevenueAPI) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	accountInfo := financeModel.WithdrawAccount{
		AccountType: req.AccountType,
		AccountName: req.AccountName,
		AccountNo:   req.AccountNo,
		BankName:    req.BankName,
		BranchName:  req.BranchName,
	}
	withdrawal, err := api.revenueService.CreateWithdrawalRequest(
		c.Request.Context(),
		userID.(string),
		req.Amount,
		req.Method,
		accountInfo,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "申请提现失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "申请提现成功",
		Data:    withdrawal,
	})
}

// GetRevenueDetails 获取收入明细
//
//	@Summary		获取收入明细
//	@Description	获取作者的收入明细
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/revenue-details [get]
func (api *AuthorRevenueAPI) GetRevenueDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	details, total, err := api.revenueService.GetRevenueDetails(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取收入明细失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		details,
		total,
		page,
		pageSize,
		"获取收入明细成功",
	))
}

// GetRevenueStatistics 获取收入统计
//
//	@Summary		获取收入统计
//	@Description	获取作者的收入统计数据
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			period	query		string	false	"统计周期"	enums(daily,monthly,yearly)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/revenue-statistics [get]
func (api *AuthorRevenueAPI) GetRevenueStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	period := c.DefaultQuery("period", "monthly")

	statistics, err := api.revenueService.GetRevenueStatistics(c.Request.Context(), userID.(string), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取收入统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取收入统计成功",
		Data:    statistics,
	})
}

// GetSettlements 获取结算记录
//
//	@Summary		获取结算记录
//	@Description	获取作者的结算记录
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/settlements [get]
func (api *AuthorRevenueAPI) GetSettlements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	settlements, total, err := api.revenueService.GetSettlements(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取结算记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		settlements,
		total,
		page,
		pageSize,
		"获取结算记录成功",
	))
}

// GetSettlement 获取结算详情
//
//	@Summary		获取结算详情
//	@Description	获取指定结算记录的详细信息
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"结算ID"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/settlements/{id} [get]
func (api *AuthorRevenueAPI) GetSettlement(c *gin.Context) {
	settlementID := c.Param("id")
	if settlementID == "" {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "结算ID不能为空",
		})
		return
	}

	settlement, err := api.revenueService.GetSettlement(c.Request.Context(), settlementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取结算详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取结算详情成功",
		Data:    settlement,
	})
}

// GetTaxInfo 获取税务信息
//
//	@Summary		获取税务信息
//	@Description	获取作者的税务信息
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/tax-info [get]
func (api *AuthorRevenueAPI) GetTaxInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	taxInfo, err := api.revenueService.GetTaxInfo(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取税务信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取税务信息成功",
		Data:    taxInfo,
	})
}

// UpdateTaxInfoRequest 更新税务信息请求
type UpdateTaxInfoRequest struct {
	IDType   string `json:"id_type" binding:"required,oneof=id_card passport other"`
	IDNumber string `json:"id_number" binding:"required"`
	Name     string `json:"name" binding:"required"`
	TaxType  string `json:"tax_type" binding:"required,oneof=individual company"`
}

// UpdateTaxInfo 更新税务信息
//
//	@Summary		更新税务信息
//	@Description	更新作者的税务信息
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		UpdateTaxInfoRequest	true	"税务信息"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/author/tax-info [put]
func (api *AuthorRevenueAPI) UpdateTaxInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req UpdateTaxInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	taxInfo := &financeModel.TaxInfo{
		IDType:   req.IDType,
		IDNumber: req.IDNumber,
		Name:     req.Name,
		TaxType:  req.TaxType,
		TaxRate:  0.00, // 根据实际情况设置税率
	}

	err := api.revenueService.UpdateTaxInfo(c.Request.Context(), userID.(string), taxInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "更新税务信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "更新税务信息成功",
	})
}
