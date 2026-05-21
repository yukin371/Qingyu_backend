package ai

import (
	"context"
	"fmt"
	"time"

	pb "Qingyu_backend/pkg/grpc/pb"
)

// TextGenerator is the narrow dependency used by local AI helper services.
// It keeps helper services off the deprecated provider adapter package.
type TextGenerator interface {
	GenerateText(ctx context.Context, req *TextGenerateRequest) (*TextGenerateResponse, error)
}

type TextGenerateRequest struct {
	ProjectID    string
	ChapterID    string
	Prompt       string
	Model        string
	MaxTokens    int
	Temperature  float64
	WorkflowType string
}

type TextGenerateResponse struct {
	Text       string
	TokensUsed int
	Model      string
}

type phase3TextGenerator struct {
	client *Phase3Client
}

func NewPhase3TextGenerator(client *Phase3Client) TextGenerator {
	return &phase3TextGenerator{client: client}
}

func (g *phase3TextGenerator) GenerateText(ctx context.Context, req *TextGenerateRequest) (*TextGenerateResponse, error) {
	if g == nil || g.client == nil {
		return nil, fmt.Errorf("AI text generator is not configured")
	}
	if req == nil {
		return nil, fmt.Errorf("text generation request is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	grpcReq := &pb.GenerateContentRequest{
		ProjectId: req.ProjectID,
		ChapterId: req.ChapterID,
		Prompt:    req.Prompt,
		Options: &pb.GenerateOptions{
			Model:       req.Model,
			MaxTokens:   int32(req.MaxTokens),
			Temperature: float32(req.Temperature),
		},
	}

	resp, err := g.client.client.GenerateContent(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("generate text via AI service: %w", err)
	}

	return &TextGenerateResponse{
		Text:       resp.GetContent(),
		TokensUsed: int(resp.GetTokensUsed()),
		Model:      resp.GetModel(),
	}, nil
}
