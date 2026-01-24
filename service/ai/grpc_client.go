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

// Close 关闭连接
func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

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
