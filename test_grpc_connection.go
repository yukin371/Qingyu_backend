package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到Python AI服务的gRPC端口
	conn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewAIServiceClient(conn)

	// 测试健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthReq := &pb.HealthCheckRequest{}
	healthResp, err := client.HealthCheck(ctx, healthReq)
	if err != nil {
		log.Fatalf("健康检查失败: %v", err)
	}
	fmt.Printf("✅ gRPC连接成功！健康状态: %s\n", healthResp.Status)
	fmt.Printf("检查项: %v\n", healthResp.Checks)

	// 测试生成内容接口
	genReq := &pb.GenerateContentRequest{
		ProjectId: "test-project-001",
		ChapterId: "test-chapter-001",
		Prompt:    "这是一个测试提示词",
		Options: &pb.GenerateOptions{
			Model:       "gpt-4",
			MaxTokens:   100,
			Temperature: 0.7,
		},
	}

	genResp, err := client.GenerateContent(ctx, genReq)
	if err != nil {
		log.Fatalf("生成内容失败: %v", err)
	}
	fmt.Printf("\n✅ 生成内容成功！\n")
	fmt.Printf("内容: %s\n", genResp.Content)
	fmt.Printf("模型: %s\n", genResp.Model)
	fmt.Printf("Token使用: %d\n", genResp.TokensUsed)

	fmt.Println("\n🎉 所有测试通过！Python AI服务与Go后端gRPC通信正常。")
}
