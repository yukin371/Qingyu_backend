package ai

import (
	"context"
	"fmt"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"github.com/sirupsen/logrus"
)

// QuotaServiceInterface 配额服务接口
// 定义配额服务需要实现的方法，使UnifiedClient不直接依赖具体实现
type QuotaServiceInterface interface {
	ConsumeQuota(ctx context.Context, userID string, amount int, service, model, requestID string) error
}

// UnifiedClient 统一AI服务gRPC客户端
// 整合了GRPCClient和Phase3Client的所有功能，提供统一的接口
type UnifiedClient struct {
	aiServiceClient pb.AIServiceClient
	conn            *grpc.ClientConn
	endpoint        string
	timeout         time.Duration
	// 监控与追踪（可选）
	metrics      *GRPCMetrics
	tracer       *Tracer
	enableMonitor bool // 是否启用监控
	// 配额管理（可选）
	quotaService QuotaServiceInterface
	enableQuota  bool // 是否启用配额扣除
}

// NewUnifiedClient 创建统一客户端
// 通过已有的gRPC连接创建客户端实例
func NewUnifiedClient(conn *grpc.ClientConn, config *AIServiceConfig) *UnifiedClient {
	client := &UnifiedClient{
		aiServiceClient: pb.NewAIServiceClient(conn),
		conn:            conn,
		endpoint:        config.Endpoint,
		timeout:         config.Timeout,
		enableMonitor:   false, // 默认不启用监控
		enableQuota:     false, // 默认不启用配额扣除
	}

	// 如果配置了监控，初始化监控组件
	if config.EnableMonitor {
		client.enableMonitor = true
		client.metrics = NewGRPCMetrics()
		client.tracer = NewTracer(1000)
	}

	return client
}

// NewUnifiedClientWithAddress 通过地址创建统一客户端
// 创建新的gRPC连接并初始化客户端
func NewUnifiedClientWithAddress(address string) (*UnifiedClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("连接AI服务失败: %w", err)
	}

	return &UnifiedClient{
		aiServiceClient: pb.NewAIServiceClient(conn),
		conn:            conn,
		endpoint:        address,
		timeout:         30 * time.Second,
		enableMonitor:   false, // 默认不启用监控
		enableQuota:     false, // 默认不启用配额扣除
	}, nil
}

// ============ GRPCClient 原有方法 ============

// ExecuteAgent 执行 AI Agent
func (c *UnifiedClient) ExecuteAgent(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	serviceName := ServiceExecuteAgent
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 构建请求
	grpcReq := &pb.AgentExecutionRequest{
		WorkflowType: req.WorkflowType,
		ProjectId:    req.UserID, // 使用 UserID 作为 ProjectId
		Parameters:   convertInterfaceMapToStringMap(req.Parameters),
	}

	// 设置超时
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 调用 gRPC
	resp, err := c.aiServiceClient.ExecuteAgent(ctx, grpcReq)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, fmt.Errorf("gRPC ExecuteAgent failed: %w", err)
	}

	// 检查执行状态
	if resp.Status != "completed" {
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, fmt.Errorf("execution failed with status %s", resp.Status))

		return nil, fmt.Errorf("AI agent execution failed with status %s: %v", resp.Status, resp.Errors)
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	// 消费配额（如果启用）
	c.consumeQuotaIfNeeded(ctx, req.UserID, resp.TokensUsed, serviceName, "default", requestID)

	return &AgentResponse{
		Content:      resp.Result, // Result 是 JSON 字符串
		TokensUsed:   int64(resp.TokensUsed),
		WorkflowType: req.WorkflowType,
	}, nil
}

// ExecuteAgentWithRetry 执行 AI Agent（带重试）
func (c *UnifiedClient) ExecuteAgentWithRetry(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	serviceName := ServiceExecuteAgent
	var lastErr error

	for i := 0; i <= 3; i++ { // 最多重试3次
		// 如果不是第一次尝试，记录重试
		if i > 0 {
			c.recordRetry(serviceName)
		}

		resp, err := c.ExecuteAgent(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// 检查错误是否可重试
		if !isRetryableError(err) {
			return nil, err
		}

		// 等待后重试
		if i < 3 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return nil, fmt.Errorf("after retries: %w", lastErr)
}

// ============ Phase3Client 原有方法 ============

// GenerateOutline 生成故事大纲
func (c *UnifiedClient) GenerateOutline(
	ctx context.Context,
	task, userID, projectID string,
	workspaceContext map[string]string,
) (*pb.OutlineResponse, error) {
	serviceName := ServiceGenerateOutline
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	request := &pb.OutlineRequest{
		Task:             task,
		UserId:           userID,
		ProjectId:        projectID,
		WorkspaceContext: workspaceContext,
	}

	resp, err := c.aiServiceClient.GenerateOutline(ctx, request)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, err
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	// 消费配额（如果启用）
	// 注意: OutlineResponse没有TokensUsed字段，使用估算值或固定配额
	// 这里使用固定值1000作为估算，后续可根据实际情况调整
	c.consumeQuotaIfNeeded(ctx, userID, 1000, serviceName, "creative-outline", requestID)

	return resp, nil
}

// GenerateCharacters 生成角色设定
func (c *UnifiedClient) GenerateCharacters(
	ctx context.Context,
	task, userID, projectID string,
	outline *pb.OutlineData,
	workspaceContext map[string]string,
) (*pb.CharactersResponse, error) {
	serviceName := ServiceGenerateCharacters
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	request := &pb.CharactersRequest{
		Task:             task,
		UserId:           userID,
		ProjectId:        projectID,
		Outline:          outline,
		WorkspaceContext: workspaceContext,
	}

	resp, err := c.aiServiceClient.GenerateCharacters(ctx, request)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, err
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	// 消费配额（如果启用）
	// 注意: CharactersResponse没有TokensUsed字段，使用估算值
	c.consumeQuotaIfNeeded(ctx, userID, 1500, serviceName, "creative-characters", requestID)

	return resp, nil
}

// GeneratePlot 生成情节设定
func (c *UnifiedClient) GeneratePlot(
	ctx context.Context,
	task, userID, projectID string,
	outline *pb.OutlineData,
	characters *pb.CharactersData,
	workspaceContext map[string]string,
) (*pb.PlotResponse, error) {
	serviceName := ServiceGeneratePlot
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	request := &pb.PlotRequest{
		Task:             task,
		UserId:           userID,
		ProjectId:        projectID,
		Outline:          outline,
		Characters:       characters,
		WorkspaceContext: workspaceContext,
	}

	resp, err := c.aiServiceClient.GeneratePlot(ctx, request)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, err
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	// 消费配额（如果启用）
	// 注意: PlotResponse没有TokensUsed字段，使用估算值
	c.consumeQuotaIfNeeded(ctx, userID, 2000, serviceName, "creative-plot", requestID)

	return resp, nil
}

// ExecuteCreativeWorkflow 执行完整创作工作流
func (c *UnifiedClient) ExecuteCreativeWorkflow(
	ctx context.Context,
	task, userID, projectID string,
	maxReflections int32,
	enableHumanReview bool,
	workspaceContext map[string]string,
) (*pb.CreativeWorkflowResponse, error) {
	serviceName := ServiceExecuteCreativeWorkflow
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置更长的超时时间（完整工作流需要更多时间）
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	request := &pb.CreativeWorkflowRequest{
		Task:              task,
		UserId:            userID,
		ProjectId:         projectID,
		MaxReflections:    maxReflections,
		EnableHumanReview: enableHumanReview,
		WorkspaceContext:  workspaceContext,
	}

	resp, err := c.aiServiceClient.ExecuteCreativeWorkflow(ctx, request)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, err
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	// 消费配额（如果启用）
	c.consumeQuotaIfNeeded(ctx, userID, resp.TokensUsed, serviceName, "creative-workflow", requestID)

	return resp, nil
}

// ============ 通用方法 ============

// HealthCheck 健康检查
// 返回完整的健康检查响应（Phase3风格）
func (c *UnifiedClient) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	serviceName := ServiceHealthCheck
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := &pb.HealthCheckRequest{}
	resp, err := c.aiServiceClient.HealthCheck(ctx, request)
	if err != nil {
		// 检查是否超时
		if ctx.Err() == context.DeadlineExceeded {
			c.recordTimeout(serviceName)
		}

		// 记录失败
		c.recordCall(serviceName, false)
		c.recordLatency(serviceName, time.Since(startTime))
		c.endTrace(requestID, TraceStatusFailed, err)

		return nil, err
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	return resp, nil
}

// HealthCheckSimple 简单健康检查
// 返回错误或nil（GRPCClient风格）
func (c *UnifiedClient) HealthCheckSimple(ctx context.Context) error {
	serviceName := ServiceHealthCheck
	requestID := generateRequestID()
	startTime := time.Now()

	// 开始追踪
	c.startTrace(serviceName, requestID)

	// 设置超时
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 调用专门的健康检查方法
	_, err := c.aiServiceClient.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			// 检查是否超时
			if ctx.Err() == context.DeadlineExceeded {
				c.recordTimeout(serviceName)
			}

			// 记录失败
			c.recordCall(serviceName, false)
			c.recordLatency(serviceName, time.Since(startTime))
			c.endTrace(requestID, TraceStatusFailed, err)

			return fmt.Errorf("AI service unavailable: %w", err)
		}
		// 其他错误可能表示服务是可用的，只是业务逻辑错误
	}

	// 记录成功
	c.recordCall(serviceName, true)
	c.recordLatency(serviceName, time.Since(startTime))
	c.endTrace(requestID, TraceStatusSuccess, nil)

	return nil
}

// Close 关闭连接
func (c *UnifiedClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetEndpoint 获取服务端点
func (c *UnifiedClient) GetEndpoint() string {
	return c.endpoint
}

// GetTimeout 获取超时时间
func (c *UnifiedClient) GetTimeout() time.Duration {
	return c.timeout
}

// SetTimeout 设置超时时间
func (c *UnifiedClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// ============ 向后兼容性方法 ============

// AsGRPCClient 将UnifiedClient转换为GRPCClient接口
// 用于需要GRPCClient的遗留代码
func (c *UnifiedClient) AsGRPCClient() *GRPCClient {
	return &GRPCClient{
		client:   c.aiServiceClient,
		conn:     c.conn,
		endpoint: c.endpoint,
		timeout:  c.timeout,
	}
}

// AsPhase3Client 将UnifiedClient转换为Phase3Client接口
// 用于需要Phase3Client的遗留代码
func (c *UnifiedClient) AsPhase3Client() *Phase3Client {
	return &Phase3Client{
		client: c.aiServiceClient,
		conn:   c.conn,
	}
}

// ============ 辅助函数 ============

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	// 可重试的错误码
	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Aborted:
		return true
	default:
		return false
	}
}

// convertInterfaceMapToStringMap 转换 map[string]interface{} 到 map[string]string
func convertInterfaceMapToStringMap(m map[string]interface{}) map[string]string {
	if m == nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result
}

// ============ 监控与追踪方法 ============

// EnableMonitoring 启用监控
func (c *UnifiedClient) EnableMonitoring() {
	c.enableMonitor = true
	if c.metrics == nil {
		c.metrics = NewGRPCMetrics()
	}
	if c.tracer == nil {
		c.tracer = NewTracer(1000)
	}
}

// DisableMonitoring 禁用监控
func (c *UnifiedClient) DisableMonitoring() {
	c.enableMonitor = false
}

// GetMetrics 获取监控指标
func (c *UnifiedClient) GetMetrics() *GRPCMetrics {
	return c.metrics
}

// GetTracer 获取追踪器
func (c *UnifiedClient) GetTracer() *Tracer {
	return c.tracer
}

// IsMonitoringEnabled 检查是否启用监控
func (c *UnifiedClient) IsMonitoringEnabled() bool {
	return c.enableMonitor
}

// ============ 配额管理方法 ============

// SetQuotaService 设置配额服务
func (c *UnifiedClient) SetQuotaService(quotaService QuotaServiceInterface) {
	c.quotaService = quotaService
}

// EnableQuota 启用配额扣除
func (c *UnifiedClient) EnableQuota() {
	c.enableQuota = true
}

// DisableQuota 禁用配额扣除
func (c *UnifiedClient) DisableQuota() {
	c.enableQuota = false
}

// IsQuotaEnabled 检查是否启用配额扣除
func (c *UnifiedClient) IsQuotaEnabled() bool {
	return c.enableQuota && c.quotaService != nil
}

// GetQuotaReport 获取配额使用报告
// 只有在启用监控时才能获取报告
func (c *UnifiedClient) GetQuotaReport() *QuotaReport {
	if c.enableMonitor && c.metrics != nil {
		return c.metrics.GetQuotaReport()
	}
	return &QuotaReport{}
}

// FormatQuotaReport 格式化配额报告
func (c *UnifiedClient) FormatQuotaReport() string {
	if c.enableMonitor && c.metrics != nil {
		return c.metrics.FormatQuotaReport()
	}
	return "配额监控未启用"
}

// consumeQuotaIfNeeded 消费配额（如果启用）
// 在gRPC调用成功后调用，用于扣除用户配额
// 配额扣除失败不会影响主流程，只记录日志
func (c *UnifiedClient) consumeQuotaIfNeeded(ctx context.Context, userID string, tokens int32, service, model, requestID string) {
	if !c.IsQuotaEnabled() {
		return
	}

	if tokens <= 0 {
		return
	}

	// 记录配额消费到监控
	if c.enableMonitor && c.metrics != nil {
		c.metrics.RecordQuotaConsumed(userID, service, model, int64(tokens))
	}

	// 异步扣除配额，避免影响主流程
	go func() {
		err := c.quotaService.ConsumeQuota(ctx, userID, int(tokens), service, model, requestID)
		if err != nil {
			// 配额扣除失败，记录错误但不影响响应
			// 记录配额不足
			if c.enableMonitor && c.metrics != nil {
				c.metrics.RecordQuotaShortage(userID)
			}

			logrus.WithFields(logrus.Fields{
				"user_id":     userID,
				"tokens":      tokens,
				"service":     service,
				"model":       model,
				"request_id":  requestID,
				"error":       err.Error(),
			}).Error("配额扣除失败")
		} else {
			logrus.WithFields(logrus.Fields{
				"user_id":    userID,
				"tokens":     tokens,
				"service":    service,
				"model":      model,
				"request_id": requestID,
			}).Info("配额扣除成功")
		}
	}()
}

// recordCall 记录调用（内部方法）
func (c *UnifiedClient) recordCall(serviceName string, success bool) {
	if c.enableMonitor && c.metrics != nil {
		c.metrics.RecordCall(serviceName, success)
	}
}

// recordLatency 记录延迟（内部方法）
func (c *UnifiedClient) recordLatency(serviceName string, duration time.Duration) {
	if c.enableMonitor && c.metrics != nil {
		c.metrics.RecordLatency(serviceName, duration)
	}
}

// recordTimeout 记录超时（内部方法）
func (c *UnifiedClient) recordTimeout(serviceName string) {
	if c.enableMonitor && c.metrics != nil {
		c.metrics.RecordTimeout(serviceName)
	}
}

// recordRetry 记录重试（内部方法）
func (c *UnifiedClient) recordRetry(serviceName string) {
	if c.enableMonitor && c.metrics != nil {
		c.metrics.RecordRetry(serviceName)
	}
}

// startTrace 开始追踪（内部方法）
func (c *UnifiedClient) startTrace(serviceName, requestID string) {
	if c.enableMonitor && c.tracer != nil {
		c.tracer.StartTrace(serviceName, requestID)
	}
}

// endTrace 结束追踪（内部方法）
func (c *UnifiedClient) endTrace(requestID string, status string, err error) {
	if c.enableMonitor && c.tracer != nil {
		c.tracer.EndTrace(requestID, status, err)
	}
}

