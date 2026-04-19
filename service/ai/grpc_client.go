package ai

import (
	"context"
	"fmt"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCClient AI服务 gRPC 客户端
type GRPCClient struct {
	client   pb.AIServiceClient
	conn     *grpc.ClientConn
	endpoint string
	timeout  time.Duration
}

// NewGRPCClient 创建 gRPC 客户端
func NewGRPCClient(conn *grpc.ClientConn, config *AIServiceConfig) *GRPCClient {
	return &GRPCClient{
		client:   pb.NewAIServiceClient(conn),
		conn:     conn,
		endpoint: config.Endpoint,
		timeout:  config.Timeout,
	}
}

// AgentRequest AI Agent 请求
type AgentRequest struct {
	UserID       string
	WorkflowType string
	Parameters   map[string]interface{}
}

// AgentResponse AI Agent 响应
type AgentResponse struct {
	Content      string
	TokensUsed   int64
	Usage        map[string]interface{}
	Model        string
	AgentType    string
	WorkflowType string
}

// ExecuteAgent 执行 AI Agent
func (c *GRPCClient) ExecuteAgent(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
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
	resp, err := c.client.ExecuteAgent(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC ExecuteAgent failed: %w", err)
	}

	// 检查执行状态
	if resp.Status != "completed" {
		return nil, fmt.Errorf("AI agent execution failed with status %s: %v", resp.Status, resp.Errors)
	}

	return &AgentResponse{
		Content:      resp.Result, // Result 是 JSON 字符串
		TokensUsed:   int64(resp.TokensUsed),
		WorkflowType: req.WorkflowType,
	}, nil
}

// ExecuteAgentWithRetry 执行 AI Agent（带重试）
func (c *GRPCClient) ExecuteAgentWithRetry(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	var lastErr error

	for i := 0; i <= 3; i++ { // 最多重试3次
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

// HealthCheck 健康检查
func (c *GRPCClient) HealthCheck(ctx context.Context) error {
	// 设置超时
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 调用专门的健康检查方法
	_, err := c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unavailable {
			return fmt.Errorf("AI service unavailable: %w", err)
		}
		// 其他错误可能表示服务是可用的，只是业务逻辑错误
	}

	return nil
}

// StoryWriteRequest 故事写作请求
type StoryWriteRequest struct {
	ProjectID     string
	DocumentID    string
	Mode          string
	Instruction   string
	SelectedText  string
	AssembledPrompt string
	Options       *GenerateOptions
}

// StoryWriteResponse 故事写作响应
type StoryWriteResponse struct {
	Content       string
	TokensUsed    int32
	Model         string
	GeneratedAt   int64
	ContextStats  *StoryContextStats
}

// StoryContextStats 故事上下文统计
type StoryContextStats struct {
	StageTokens   int32
	OutlineTokens int32
	RAGTokens     int32
	TotalTokens   int32
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	Model       string
	MaxTokens   int32
	Temperature float32
	Stop        []string
	Stream      bool
}

// StoryWrite 故事上下文写作
func (c *GRPCClient) StoryWrite(ctx context.Context, req *StoryWriteRequest) (*StoryWriteResponse, error) {
	// 构建请求
	grpcReq := &pb.StoryContextRequest{
		ProjectId:      req.ProjectID,
		DocumentId:     req.DocumentID,
		Mode:           req.Mode,
		Instruction:    req.Instruction,
		SelectedText:   req.SelectedText,
		AssembledPrompt: req.AssembledPrompt,
	}

	// 设置生成选项
	if req.Options != nil {
		grpcReq.Options = &pb.GenerateOptions{
			Model:       req.Options.Model,
			MaxTokens:   req.Options.MaxTokens,
			Temperature: req.Options.Temperature,
			Stop:        req.Options.Stop,
			Stream:      req.Options.Stream,
		}
	}

	// 设置超时
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// 调用 gRPC
	resp, err := c.client.StoryWrite(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC StoryWrite failed: %w", err)
	}

	// 转换响应
	response := &StoryWriteResponse{
		Content:     resp.Content,
		TokensUsed:  resp.TokensUsed,
		Model:       resp.Model,
		GeneratedAt: resp.GeneratedAt,
	}

	if resp.ContextStats != nil {
		response.ContextStats = &StoryContextStats{
			StageTokens:   resp.ContextStats.StageTokens,
			OutlineTokens: resp.ContextStats.OutlineTokens,
			RAGTokens:     resp.ContextStats.RagTokens,
			TotalTokens:   resp.ContextStats.TotalTokens,
		}
	}

	return response, nil
}

// Close 关闭连接
func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
