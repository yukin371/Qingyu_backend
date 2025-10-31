package ai

import (
	"context"
	"fmt"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Phase3Client Phase3 AI服务gRPC客户端
type Phase3Client struct {
	client pb.AIServiceClient
	conn   *grpc.ClientConn
}

// NewPhase3Client 创建Phase3客户端
func NewPhase3Client(address string) (*Phase3Client, error) {
	// 创建gRPC连接
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("连接AI服务失败: %w", err)
	}

	return &Phase3Client{
		client: pb.NewAIServiceClient(conn),
		conn:   conn,
	}, nil
}

// GenerateOutline 生成故事大纲
func (c *Phase3Client) GenerateOutline(
	ctx context.Context,
	task, userID, projectID string,
	workspaceContext map[string]string,
) (*pb.OutlineResponse, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	request := &pb.OutlineRequest{
		Task:             task,
		UserId:           userID,
		ProjectId:        projectID,
		WorkspaceContext: workspaceContext,
	}

	return c.client.GenerateOutline(ctx, request)
}

// GenerateCharacters 生成角色设定
func (c *Phase3Client) GenerateCharacters(
	ctx context.Context,
	task, userID, projectID string,
	outline *pb.OutlineData,
	workspaceContext map[string]string,
) (*pb.CharactersResponse, error) {
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

	return c.client.GenerateCharacters(ctx, request)
}

// GeneratePlot 生成情节设定
func (c *Phase3Client) GeneratePlot(
	ctx context.Context,
	task, userID, projectID string,
	outline *pb.OutlineData,
	characters *pb.CharactersData,
	workspaceContext map[string]string,
) (*pb.PlotResponse, error) {
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

	return c.client.GeneratePlot(ctx, request)
}

// ExecuteCreativeWorkflow 执行完整创作工作流
func (c *Phase3Client) ExecuteCreativeWorkflow(
	ctx context.Context,
	task, userID, projectID string,
	maxReflections int32,
	enableHumanReview bool,
	workspaceContext map[string]string,
) (*pb.CreativeWorkflowResponse, error) {
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

	return c.client.ExecuteCreativeWorkflow(ctx, request)
}

// HealthCheck 健康检查
func (c *Phase3Client) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	request := &pb.HealthCheckRequest{}
	return c.client.HealthCheck(ctx, request)
}

// Close 关闭连接
func (c *Phase3Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
