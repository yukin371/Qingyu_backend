package ai

import (
	"net/http"
	"time"

	"Qingyu_backend/middleware"
	pb "Qingyu_backend/pkg/grpc/pb"
	"Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// CreativeAPI Phase3创作API处理器
type CreativeAPI struct {
	phase3Client *ai.Phase3Client
}

// NewCreativeAPI 创建创作API处理器
func NewCreativeAPI(phase3Client *ai.Phase3Client) *CreativeAPI {
	return &CreativeAPI{
		phase3Client: phase3Client,
	}
}

// GenerateOutline 生成大纲
// @Summary 生成故事大纲
// @Description 根据任务描述生成完整的故事大纲，包含章节、故事结构等
// @Tags Phase3-创作
// @Accept json
// @Produce json
// @Param request body GenerateOutlineRequest true "大纲生成请求"
// @Success 200 {object} middleware.Response{data=GenerateOutlineResponse} "成功"
// @Failure 400 {object} middleware.Response "参数错误"
// @Failure 500 {object} middleware.Response "服务器错误"
// @Router /api/v1/ai/creative/outline [post]
func (a *CreativeAPI) GenerateOutline(c *gin.Context) {
	var req GenerateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationError(c, err.Error())
		return
	}

	// 获取用户ID
	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	// 调用gRPC服务
	ctx, cancel := c.Request.Context(), func() {}
	defer cancel()

	startTime := time.Now()
	grpcResp, err := a.phase3Client.GenerateOutline(
		ctx,
		req.Task,
		userIDStr,
		req.ProjectID,
		req.WorkspaceContext,
	)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, "大纲生成失败", err.Error())
		return
	}

	// 转换响应
	response := &GenerateOutlineResponse{
		Outline:       convertProtoOutlineToModel(grpcResp.Outline),
		ExecutionTime: grpcResp.ExecutionTime,
	}

	middleware.Success(c, response)
}

// GenerateCharacters 生成角色
// @Summary 生成角色设定
// @Description 根据大纲生成角色设定，包含角色关系网络
// @Tags Phase3-创作
// @Accept json
// @Produce json
// @Param request body GenerateCharactersRequest true "角色生成请求"
// @Success 200 {object} middleware.Response{data=GenerateCharactersResponse} "成功"
// @Failure 400 {object} middleware.Response "参数错误"
// @Failure 500 {object} middleware.Response "服务器错误"
// @Router /api/v1/ai/creative/characters [post]
func (a *CreativeAPI) GenerateCharacters(c *gin.Context) {
	var req GenerateCharactersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	ctx, cancel := c.Request.Context(), func() {}
	defer cancel()

	// 转换大纲数据为proto格式
	var outlineProto *pb.OutlineData
	if req.Outline != nil {
		outlineProto = convertModelOutlineToProto(req.Outline)
	}

	grpcResp, err := a.phase3Client.GenerateCharacters(
		ctx,
		req.Task,
		userIDStr,
		req.ProjectID,
		outlineProto,
		req.WorkspaceContext,
	)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, "角色生成失败", err.Error())
		return
	}

	response := &GenerateCharactersResponse{
		Characters:    convertProtoCharactersToModel(grpcResp.Characters),
		ExecutionTime: grpcResp.ExecutionTime,
	}

	middleware.Success(c, response)
}

// GeneratePlot 生成情节
// @Summary 生成情节设定
// @Description 根据大纲和角色生成情节，包含时间线事件、情节线索等
// @Tags Phase3-创作
// @Accept json
// @Produce json
// @Param request body GeneratePlotRequest true "情节生成请求"
// @Success 200 {object} middleware.Response{data=GeneratePlotResponse} "成功"
// @Failure 400 {object} middleware.Response "参数错误"
// @Failure 500 {object} middleware.Response "服务器错误"
// @Router /api/v1/ai/creative/plot [post]
func (a *CreativeAPI) GeneratePlot(c *gin.Context) {
	var req GeneratePlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	ctx, cancel := c.Request.Context(), func() {}
	defer cancel()

	// 转换数据为proto格式
	var outlineProto *pb.OutlineData
	if req.Outline != nil {
		outlineProto = convertModelOutlineToProto(req.Outline)
	}

	var charactersProto *pb.CharactersData
	if req.Characters != nil {
		charactersProto = convertModelCharactersToProto(req.Characters)
	}

	grpcResp, err := a.phase3Client.GeneratePlot(
		ctx,
		req.Task,
		userIDStr,
		req.ProjectID,
		outlineProto,
		charactersProto,
		req.WorkspaceContext,
	)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, "情节生成失败", err.Error())
		return
	}

	response := &GeneratePlotResponse{
		Plot:          convertProtoPlotToModel(grpcResp.Plot),
		ExecutionTime: grpcResp.ExecutionTime,
	}

	middleware.Success(c, response)
}

// ExecuteCreativeWorkflow 执行完整创作工作流
// @Summary 执行完整创作工作流
// @Description 一次性生成大纲、角色、情节的完整创作流程
// @Tags Phase3-创作
// @Accept json
// @Produce json
// @Param request body ExecuteCreativeWorkflowRequest true "工作流执行请求"
// @Success 200 {object} middleware.Response{data=ExecuteCreativeWorkflowResponse} "成功"
// @Failure 400 {object} middleware.Response "参数错误"
// @Failure 500 {object} middleware.Response "服务器错误"
// @Router /api/v1/ai/creative/workflow [post]
func (a *CreativeAPI) ExecuteCreativeWorkflow(c *gin.Context) {
	var req ExecuteCreativeWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ValidationError(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	// 设置默认值
	if req.MaxReflections == 0 {
		req.MaxReflections = 3
	}

	ctx, cancel := c.Request.Context(), func() {}
	defer cancel()

	grpcResp, err := a.phase3Client.ExecuteCreativeWorkflow(
		ctx,
		req.Task,
		userIDStr,
		req.ProjectID,
		req.MaxReflections,
		req.EnableHumanReview,
		req.WorkspaceContext,
	)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, "工作流执行失败", err.Error())
		return
	}

	response := &ExecuteCreativeWorkflowResponse{
		ExecutionID:      grpcResp.ExecutionId,
		ReviewPassed:     grpcResp.ReviewPassed,
		ReflectionCount:  grpcResp.ReflectionCount,
		Outline:          convertProtoOutlineToModel(grpcResp.Outline),
		Characters:       convertProtoCharactersToModel(grpcResp.Characters),
		Plot:             convertProtoPlotToModel(grpcResp.Plot),
		DiagnosticReport: convertProtoDiagnosticReportToModel(grpcResp.DiagnosticReport),
		Reasoning:        grpcResp.Reasoning,
		ExecutionTimes:   grpcResp.ExecutionTimes,
		TokensUsed:       grpcResp.TokensUsed,
	}

	middleware.Success(c, response)
}

// HealthCheck 健康检查
// @Summary Phase3服务健康检查
// @Description 检查Phase3 AI服务的健康状态
// @Tags Phase3-创作
// @Produce json
// @Success 200 {object} middleware.Response{data=map[string]interface{}} "成功"
// @Failure 500 {object} middleware.Response "服务器错误"
// @Router /api/v1/ai/creative/health [get]
func (a *CreativeAPI) HealthCheck(c *gin.Context) {
	ctx, cancel := c.Request.Context(), func() {}
	defer cancel()

	resp, err := a.phase3Client.HealthCheck(ctx)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, "健康检查失败", err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"status": resp.Status,
		"checks": resp.Checks,
	})
}
