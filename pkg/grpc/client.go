package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"Qingyu_backend/pkg/grpc/pb"
)

// AIClient AI Service gRPC 客户端
type AIClient struct {
	conn   *grpc.ClientConn
	client pb.AIServiceClient
}

// NewAIClient 创建 AI Service 客户端
func NewAIClient(address string) (*AIClient, error) {
	// 配置 keepalive 参数
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // 每 10 秒发送 ping
		Timeout:             time.Second,      // ping 超时时间
		PermitWithoutStream: true,             // 允许在没有活动流时发送 ping
	}

	// 连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 开发环境使用不安全连接
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100MB 最大接收消息大小
			grpc.MaxCallSendMsgSize(100*1024*1024), // 100MB 最大发送消息大小
		),
	}

	// 建立连接
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}

	// 创建客户端
	client := pb.NewAIServiceClient(conn)

	return &AIClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close 关闭连接
func (c *AIClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GenerateContent 生成内容
func (c *AIClient) GenerateContent(ctx context.Context, req *pb.GenerateContentRequest) (*pb.GenerateContentResponse, error) {
	return c.client.GenerateContent(ctx, req)
}

// QueryKnowledge RAG 查询
func (c *AIClient) QueryKnowledge(ctx context.Context, req *pb.RAGQueryRequest) (*pb.RAGQueryResponse, error) {
	return c.client.QueryKnowledge(ctx, req)
}

// GetContext 获取上下文
func (c *AIClient) GetContext(ctx context.Context, req *pb.ContextRequest) (*pb.ContextResponse, error) {
	return c.client.GetContext(ctx, req)
}

// ExecuteAgent 执行 Agent 工作流
func (c *AIClient) ExecuteAgent(ctx context.Context, req *pb.AgentExecutionRequest) (*pb.AgentExecutionResponse, error) {
	return c.client.ExecuteAgent(ctx, req)
}

// EmbedText 向量化文本
func (c *AIClient) EmbedText(ctx context.Context, req *pb.EmbedRequest) (*pb.EmbedResponse, error) {
	return c.client.EmbedText(ctx, req)
}

// HealthCheck 健康检查
func (c *AIClient) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
}
