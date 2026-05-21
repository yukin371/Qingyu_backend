package finance

import (
	"strconv"

	"github.com/gin-gonic/gin"

	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/pkg/response"
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

func objectIDToString(id interface{ Hex() string }) string {
	if id == nil {
		return ""
	}
	return id.Hex()
}

func mapAuthorEarning(item *financeModel.AuthorEarning) gin.H {
	if item == nil {
		return gin.H{}
	}

	return gin.H{
		"id":                 objectIDToString(item.ID),
		"author_id":          item.AuthorID,
		"book_id":            objectIDToString(item.BookID),
		"book_title":         item.BookTitle,
		"chapter_id":         objectIDToString(item.ChapterID),
		"chapter_title":      item.ChapterTitle,
		"type":               item.Type,
		"amount":             item.AuthorIncome.ToYuan(),
		"amount_cents":       item.AuthorIncome.ToCents(),
		"gross_amount":       item.Amount.ToYuan(),
		"gross_amount_cents": item.Amount.ToCents(),
		"platform_fee":       item.PlatformFee.ToYuan(),
		"platform_fee_cents": item.PlatformFee.ToCents(),
		"reader_id":          item.ReaderID,
		"reader_nickname":    item.ReaderNickname,
		"is_settled":         item.IsSettled,
		"created_at":         item.CreatedAt,
		"updated_at":         item.UpdatedAt,
	}
}

func mapAuthorWithdrawal(item *financeModel.WithdrawalRequest) gin.H {
	if item == nil {
		return gin.H{}
	}

	return gin.H{
		"id":                  objectIDToString(item.ID),
		"user_id":             item.UserID,
		"amount":              item.Amount.ToYuan(),
		"amount_cents":        item.Amount.ToCents(),
		"fee":                 item.Fee.ToYuan(),
		"fee_cents":           item.Fee.ToCents(),
		"actual_amount":       item.ActualAmount.ToYuan(),
		"actual_amount_cents": item.ActualAmount.ToCents(),
		"method":              item.Method,
		"account_info": gin.H{
			"account_type": item.AccountInfo.AccountType,
			"account_name": item.AccountInfo.AccountName,
			"account_no":   item.AccountInfo.AccountNo,
			"bank_name":    item.AccountInfo.BankName,
			"branch_name":  item.AccountInfo.BranchName,
		},
		"status":         item.Status,
		"reject_reason":  item.RejectReason,
		"approved_by":    item.ApprovedBy,
		"approved_at":    item.ApprovedAt,
		"completed_at":   item.CompletedAt,
		"transaction_id": item.TransactionID,
		"note":           item.Note,
		"created_at":     item.CreatedAt,
		"updated_at":     item.UpdatedAt,
	}
}

func mapRevenueDetail(item *financeModel.RevenueDetail) gin.H {
	if item == nil {
		return gin.H{}
	}

	return gin.H{
		"id":                  objectIDToString(item.ID),
		"author_id":           item.AuthorID,
		"book_id":             objectIDToString(item.BookID),
		"book_title":          item.BookTitle,
		"type":                item.Type,
		"total_amount":        item.TotalAmount.ToYuan(),
		"total_amount_cents":  item.TotalAmount.ToCents(),
		"total_income":        item.TotalIncome.ToYuan(),
		"total_income_cents":  item.TotalIncome.ToCents(),
		"transaction_count":   item.TransactionCount,
		"first_earning_at":    item.FirstEarningAt,
		"last_earning_at":     item.LastEarningAt,
		"created_at":          item.CreatedAt,
		"updated_at":          item.UpdatedAt,
	}
}

func mapRevenueStatistics(item *financeModel.RevenueStatistics) gin.H {
	if item == nil {
		return gin.H{}
	}

	return gin.H{
		"author_id":               item.AuthorID,
		"period":                  item.Period,
		"period_start":            item.PeriodStart,
		"period_end":              item.PeriodEnd,
		"total_revenue":           item.TotalRevenue.ToYuan(),
		"total_revenue_cents":     item.TotalRevenue.ToCents(),
		"chapter_income":          item.ChapterIncome.ToYuan(),
		"chapter_income_cents":    item.ChapterIncome.ToCents(),
		"reward_income":           item.RewardIncome.ToYuan(),
		"reward_income_cents":     item.RewardIncome.ToCents(),
		"vip_reading_income":      item.VIPReadingIncome.ToYuan(),
		"vip_reading_income_cents": item.VIPReadingIncome.ToCents(),
		"transaction_count":       item.TransactionCount,
		"reader_count":            item.ReaderCount,
		"book_count":              item.BookCount,
		"created_at":              item.CreatedAt,
		"updated_at":              item.UpdatedAt,
	}
}

func mapSettlement(item *financeModel.Settlement) gin.H {
	if item == nil {
		return gin.H{}
	}

	return gin.H{
		"id":                  objectIDToString(item.ID),
		"author_id":           item.AuthorID,
		"author_nickname":     item.AuthorNickname,
		"period_start":        item.PeriodStart,
		"period_end":          item.PeriodEnd,
		"total_revenue":       item.TotalRevenue.ToYuan(),
		"total_revenue_cents": item.TotalRevenue.ToCents(),
		"platform_fee":        item.PlatformFee.ToYuan(),
		"platform_fee_cents":  item.PlatformFee.ToCents(),
		"actual_income":       item.ActualIncome.ToYuan(),
		"actual_income_cents": item.ActualIncome.ToCents(),
		"tax_fee":             item.TaxFee.ToYuan(),
		"tax_fee_cents":       item.TaxFee.ToCents(),
		"final_income":        item.FinalIncome.ToYuan(),
		"final_income_cents":  item.FinalIncome.ToCents(),
		"earning_count":       item.EarningCount,
		"status":              item.Status,
		"processed_at":        item.ProcessedAt,
		"transaction_id":      item.TransactionID,
		"note":                item.Note,
		"created_at":          item.CreatedAt,
		"updated_at":          item.UpdatedAt,
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/earnings [get]
func (api *AuthorRevenueAPI) GetEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	earnings, total, err := api.revenueService.GetEarnings(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	items := make([]gin.H, 0, len(earnings))
	for _, item := range earnings {
		items = append(items, mapAuthorEarning(item))
	}

	response.Paginated(c, items, total, page, pageSize, "获取收入列表成功")
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/earnings/{bookId} [get]
func (api *AuthorRevenueAPI) GetBookEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		response.BadRequest(c, "书籍ID不能为空", "")
		return
	}

	earnings, total, err := api.revenueService.GetBookEarnings(c.Request.Context(), userID.(string), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	items := make([]gin.H, 0, len(earnings))
	for _, item := range earnings {
		items = append(items, mapAuthorEarning(item))
	}

	response.Paginated(c, items, total, page, pageSize, "获取书籍收入成功")
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/withdrawals [get]
func (api *AuthorRevenueAPI) GetWithdrawals(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	withdrawals, total, err := api.revenueService.GetWithdrawals(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	items := make([]gin.H, 0, len(withdrawals))
	for _, item := range withdrawals {
		items = append(items, mapAuthorWithdrawal(item))
	}

	response.Paginated(c, items, total, page, pageSize, "获取提现记录成功")
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/withdraw [post]
func (api *AuthorRevenueAPI) Withdraw(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
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
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "申请提现成功", mapAuthorWithdrawal(withdrawal))
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/revenue-details [get]
func (api *AuthorRevenueAPI) GetRevenueDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	details, total, err := api.revenueService.GetRevenueDetails(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	items := make([]gin.H, 0, len(details))
	for _, item := range details {
		items = append(items, mapRevenueDetail(item))
	}

	response.Paginated(c, items, total, page, pageSize, "获取收入明细成功")
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/revenue-statistics [get]
func (api *AuthorRevenueAPI) GetRevenueStatistics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	period := c.DefaultQuery("period", "monthly")

	statistics, err := api.revenueService.GetRevenueStatistics(c.Request.Context(), userID.(string), period)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	items := make([]gin.H, 0, len(statistics))
	for _, item := range statistics {
		items = append(items, mapRevenueStatistics(item))
	}

	response.SuccessWithMessage(c, "获取收入统计成功", items)
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/settlements [get]
func (api *AuthorRevenueAPI) GetSettlements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	settlements, total, err := api.revenueService.GetSettlements(c.Request.Context(), userID.(string), page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	items := make([]gin.H, 0, len(settlements))
	for _, item := range settlements {
		items = append(items, mapSettlement(item))
	}

	response.Paginated(c, items, total, page, pageSize, "获取结算记录成功")
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/settlements/{id} [get]
func (api *AuthorRevenueAPI) GetSettlement(c *gin.Context) {
	settlementID := c.Param("id")
	if settlementID == "" {
		response.BadRequest(c, "结算ID不能为空", "")
		return
	}

	settlement, err := api.revenueService.GetSettlement(c.Request.Context(), settlementID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取结算详情成功", mapSettlement(settlement))
}

// GetTaxInfo 获取税务信息
//
//	@Summary		获取税务信息
//	@Description	获取作者的税务信息
//	@Tags			作者收入
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/tax-info [get]
func (api *AuthorRevenueAPI) GetTaxInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	taxInfo, err := api.revenueService.GetTaxInfo(c.Request.Context(), userID.(string))
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取税务信息成功", taxInfo)
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
//	@Success 200 {object} response.APIResponse
//	@Router			/api/v1/finance/author/tax-info [put]
func (api *AuthorRevenueAPI) UpdateTaxInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未认证")
		return
	}

	var req UpdateTaxInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误", err.Error())
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
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新税务信息成功", nil)
}
