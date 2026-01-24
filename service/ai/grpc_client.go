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
		UserId:       req.UserID,
		WorkflowType: req.WorkflowType,
		Parameters:   convertMapToStruct(req.Parameters),
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

	if !resp.Success {
		return nil, fmt.Errorf("AI agent execution failed: %s", resp.ErrorMessage)
	}

	return &AgentResponse{
		Content:      resp.Content,
		TokensUsed:   resp.TokensUsed,
		Usage:        convertStructToMap(resp.Usage),
		Model:        resp.Model,
		AgentType:    resp.AgentType,
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

	// 尝试调用一个简单的 gRPC 方法（这里假设有健康检查方法）
	// 如果没有专门的健康检查方法，可以调用 ExecuteAgent 并检查响应
	_, err := c.client.ExecuteAgent(ctx, &pb.AgentExecutionRequest{
		UserId:       "health-check",
		WorkflowType: "health",
	})

	// 对于健康检查，我们主要关心连接是否可用，而不是具体的业务错误
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

// convertMapToStruct 转换 map[string]interface{} 到 protobuf Struct
func convertMapToStruct(m map[string]interface{}) *pb.Struct {
	if m == nil {
		return nil
	}

	structValue := &pb.Struct{
		Fields: make(map[string]*pb.Value),
	}

	for k, v := range m {
		structValue.Fields[k] = convertValue(v)
	}

	return structValue
}

// convertStructToMap 转换 protobuf Struct 到 map[string]interface{}
func convertStructToMap(s *pb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range s.Fields {
		result[k] = convertValueToInterface(v)
	}

	return result
}

// convertValue 转换 Go 类型到 protobuf Value
func convertValue(v interface{}) *pb.Value {
	if v == nil {
		return &pb.Value{Kind: &pb.Value_NullValue{}}
	}

	switch val := v.(type) {
	case bool:
		return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: val}}
	case float64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: val}}
	case int:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(val)}}
	case int64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(val)}}
	case string:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: val}}
	case []interface{}:
		list := &pb.ListValue{Values: make([]*pb.Value, 0)}
		for _, item := range val {
			list.Values = append(list.Values, convertValue(item))
		}
		return &pb.Value{Kind: &pb.Value_ListValue{ListValue: list}}
	case map[string]interface{}:
		return &pb.Value{Kind: &pb.Value_StructValue{StructValue: convertMapToStruct(val)}}
	default:
		return &pb.Value{Kind: &pb.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
	}
}

// convertValueToInterface 转换 protobuf Value 到 Go 类型
func convertValueToInterface(v *pb.Value) interface{} {
	if v == nil {
		return nil
	}

	switch kind := v.Kind.(type) {
	case *pb.Value_NullValue:
		return nil
	case *pb.Value_BoolValue:
		return kind.BoolValue
	case *pb.Value_NumberValue:
		return kind.NumberValue
	case *pb.Value_StringValue:
		return kind.StringValue
	case *pb.Value_ListValue:
		result := make([]interface{}, 0)
		for _, item := range kind.ListValue.Values {
			result = append(result, convertValueToInterface(item))
		}
		return result
	case *pb.Value_StructValue:
		return convertStructToMap(kind.StructValue)
	default:
		return nil
	}
}
