package finance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	sharedApi "Qingyu_backend/api/v1/shared"
	financeService "Qingyu_backend/service/finance"
)

// MembershipAPI 会员API处理器
type MembershipAPI struct {
	membershipService financeService.MembershipService
}

// NewMembershipAPI 创建会员API实例
func NewMembershipAPI(membershipService financeService.MembershipService) *MembershipAPI {
	return &MembershipAPI{
		membershipService: membershipService,
	}
}

// GetPlans 获取会员套餐列表
//
//	@Summary		获取会员套餐列表
//	@Description	获取所有可用的会员套餐
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/plans [get]
func (api *MembershipAPI) GetPlans(c *gin.Context) {
	plans, err := api.membershipService.GetPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取套餐列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取套餐列表成功",
		Data:    plans,
	})
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	PlanID        string `json:"plan_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat bank wallet"`
}

// Subscribe 订阅会员
//
//	@Summary		订阅会员
//	@Description	用户订阅会员套餐
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		SubscribeRequest	true	"订阅信息"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/subscribe [post]
func (api *MembershipAPI) Subscribe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	membership, err := api.membershipService.Subscribe(c.Request.Context(), userID.(string), req.PlanID, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "订阅失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "订阅成功",
		Data:    membership,
	})
}

// GetStatus 获取会员状态
//
//	@Summary		获取会员状态
//	@Description	获取当前用户的会员状态信息
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/status [get]
func (api *MembershipAPI) GetStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	membership, err := api.membershipService.GetMembership(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取会员状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取会员状态成功",
		Data:    membership,
	})
}

// Cancel 取消自动续费
//
//	@Summary		取消自动续费
//	@Description	取消会员自动续费
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/cancel [post]
func (api *MembershipAPI) Cancel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	err := api.membershipService.CancelMembership(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "取消自动续费失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "取消自动续费成功",
	})
}

// Renew 手动续费
//
//	@Summary		手动续费
//	@Description	手动续费会员
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/renew [put]
func (api *MembershipAPI) Renew(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	membership, err := api.membershipService.RenewMembership(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "续费失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "续费成功",
		Data:    membership,
	})
}

// GetBenefits 获取会员权益列表
//
//	@Summary		获取会员权益列表
//	@Description	获取会员权益列表
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			level	query		string	false	"会员等级"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/benefits [get]
func (api *MembershipAPI) GetBenefits(c *gin.Context) {
	level := c.Query("level")

	benefits, err := api.membershipService.GetBenefits(c.Request.Context(), level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取权益列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取权益列表成功",
		Data:    benefits,
	})
}

// GetUsage 获取会员权益使用情况
//
//	@Summary		获取会员权益使用情况
//	@Description	获取当前用户的权益使用情况
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/usage [get]
func (api *MembershipAPI) GetUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	usage, err := api.membershipService.GetUsage(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取权益使用情况失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "获取权益使用情况成功",
		Data:    usage,
	})
}

// ListCards 获取会员卡列表
//
//	@Summary		获取会员卡列表
//	@Description	获取会员卡列表（管理员）
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			status		query		string	false	"状态"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/cards [get]
func (api *MembershipAPI) ListCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	filter := map[string]interface{}{}
	if status != "" {
		filter["status"] = status
	}

	cards, total, err := api.membershipService.ListCards(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "获取会员卡列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.PaginatedResponseHelper(
		cards,
		total,
		page,
		pageSize,
		"获取会员卡列表成功",
	))
}

// ActivateCardRequest 激活会员卡请求
type ActivateCardRequest struct {
	Code string `json:"code" binding:"required"`
}

// ActivateCard 激活会员卡
//
//	@Summary		激活会员卡
//	@Description	使用卡密激活会员
//	@Tags			会员系统
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		ActivateCardRequest	true	"激活信息"
//	@Success 200 {object} APIResponse
//	@Router			/api/v1/finance/membership/cards/activate [post]
func (api *MembershipAPI) ActivateCard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharedApi.APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	var req ActivateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharedApi.APIResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	membership, err := api.membershipService.ActivateCard(c.Request.Context(), userID.(string), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharedApi.APIResponse{
			Code:    500,
			Message: "激活会员卡失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sharedApi.APIResponse{
		Code:    200,
		Message: "激活会员卡成功",
		Data:    membership,
	})
}
