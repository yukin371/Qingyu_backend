package search

import (
	"context"
	"fmt"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"

	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// EmbeddingClient 封装 EmbedText gRPC 调用，用于将文本转为向量。
// 搜索时：将用户查询文本转为向量，然后做 Milvus 向量搜索。
// 索引时：将书籍文本转为向量，写入 Milvus。
type EmbeddingClient struct {
	client  pb.AIServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
	logger  *logger.Logger
}

// NewEmbeddingClient 创建 EmbedText gRPC 客户端。
// addr 格式: "host:port"（如 "localhost:50051"）。
func NewEmbeddingClient(addr string, timeout time.Duration) (*EmbeddingClient, error) {
	if addr == "" {
		return nil, fmt.Errorf("embedding client: addr is required")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("embedding client: dial %s failed: %w", addr, err)
	}

	return &EmbeddingClient{
		client:  pb.NewAIServiceClient(conn),
		conn:    conn,
		timeout: timeout,
		logger:  logger.Get().WithModule("embedding-client"),
	}, nil
}

// NewEmbeddingClientFromConn 从已有 gRPC 连接创建客户端（共享连接场景）。
func NewEmbeddingClientFromConn(conn *grpc.ClientConn, timeout time.Duration) (*EmbeddingClient, error) {
	if conn == nil {
		return nil, fmt.Errorf("embedding client: conn is required")
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &EmbeddingClient{
		client:  pb.NewAIServiceClient(conn),
		conn:    nil, // 不持有连接所有权，Close 不关它
		timeout: timeout,
		logger:  logger.Get().WithModule("embedding-client"),
	}, nil
}

// GetEmbedding 获取单条文本的向量。
func (c *EmbeddingClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("embedding client: text is required")
	}

	vecs, err := c.GetEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vecs) == 0 {
		return nil, fmt.Errorf("embedding client: no embedding returned")
	}
	return vecs[0], nil
}

// GetEmbeddings 批量获取向量。
func (c *EmbeddingClient) GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("embedding client: texts is empty")
	}

	startTime := time.Now()

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.EmbedText(ctx, &pb.EmbedRequest{
		Texts: texts,
	})
	if err != nil {
		c.logger.Error("EmbedText RPC failed",
			zap.Int("text_count", len(texts)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("embedding client: EmbedText RPC failed: %w", err)
	}

	if resp == nil || len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("embedding client: empty response from AI service")
	}

	// 转换 proto 向量为 Go [][]float32
	result := make([][]float32, 0, len(resp.Embeddings))
	for i, emb := range resp.Embeddings {
		if emb == nil {
			c.logger.Warn("nil embedding at index, skipping",
				zap.Int("index", i),
			)
			result = append(result, nil)
			continue
		}
		result = append(result, emb.Vector)
	}

	took := time.Since(startTime)
	c.logger.Info("EmbedText completed",
		zap.Int("input_count", len(texts)),
		zap.Int("output_count", len(result)),
		zap.Duration("took", took),
	)

	return result, nil
}

// Close 关闭连接（仅当通过 NewEmbeddingClient 创建时有效）。
func (c *EmbeddingClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Health 通过 AI 服务的 HealthCheck RPC 检查连接健康状态。
func (c *EmbeddingClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return fmt.Errorf("embedding client: health check failed: %w", err)
	}
	if resp != nil && resp.Status == "unhealthy" {
		return fmt.Errorf("embedding client: AI service reports unhealthy")
	}
	return nil
}
