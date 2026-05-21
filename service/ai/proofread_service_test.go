package ai

import (
	"context"
	"strings"
	"testing"

	"Qingyu_backend/service/ai/dto"
)

type mockProofreadTextGenerator struct {
	responseText string
	lastPrompt   string
	callCount    int
}

func (m *mockProofreadTextGenerator) GenerateText(ctx context.Context, req *TextGenerateRequest) (*TextGenerateResponse, error) {
	m.lastPrompt = req.Prompt
	m.callCount++
	return &TextGenerateResponse{
		Text:       m.responseText,
		TokensUsed: 42,
		Model:      "test-model",
	}, nil
}

func TestProofreadContentReturnsStructuredReviewMetadata(t *testing.T) {
	generator := &mockProofreadTextGenerator{
		responseText: `{
			"score": 82,
			"issues": [
				{
					"type": "typo",
					"severity": "error",
					"message": "疑似错别字",
					"start": 2,
					"end": 4,
					"originalText": "在见",
					"suggestions": [{"text": "再见", "reason": "告别语应使用再见"}]
				}
			]
		}`,
	}
	service := NewProofreadService(generator)

	result, err := service.ProofreadContent(context.Background(), &dto.ProofreadRequest{
		Content:       "张三在见李四。\n\n" + strings.Repeat("这是一段很长的移动端阅读测试文本", 18),
		CheckTypes:    []string{"typo", "grammar", "readability"},
		ReaderProfile: "移动端网文读者",
	})
	if err != nil {
		t.Fatalf("ProofreadContent returned error: %v", err)
	}

	if result.ReviewID == "" {
		t.Fatal("expected review id")
	}
	if !strings.HasPrefix(result.ContentHash, "sha256:") {
		t.Fatalf("expected sha256 content hash, got %q", result.ContentHash)
	}
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Type != "typo" {
		t.Fatalf("expected typo issue, got %q", result.Issues[0].Type)
	}
	if result.Issues[0].Suggestions[0] != "再见" {
		t.Fatalf("expected suggestion to be normalized")
	}
	if len(result.Issues[0].SuggestionDetails) != 1 {
		t.Fatalf("expected structured suggestion details")
	}
	if result.Issues[0].SuggestionDetails[0].Reason != "告别语应使用再见" {
		t.Fatalf("expected suggestion reason to be preserved")
	}
	if len(result.PreviewWarnings) == 0 {
		t.Fatal("expected mobile preview warnings for long paragraph")
	}
	if !strings.Contains(generator.lastPrompt, "错别字") {
		t.Fatal("expected typo check to map into Chinese proofreading prompt")
	}
	if !strings.Contains(generator.lastPrompt, "移动端阅读体验") {
		t.Fatal("expected readability check in prompt")
	}
}

func TestProofreadContentKeepsLegacyStringSuggestions(t *testing.T) {
	generator := &mockProofreadTextGenerator{
		responseText: `{
			"issues": [
				{
					"type": "spelling",
					"severity": "medium",
					"message": "疑似错别字",
					"line": 1,
					"column": 3,
					"original": "在见",
					"suggestions": ["再见"]
				}
			]
		}`,
	}
	service := NewProofreadService(generator)

	result, err := service.ProofreadContent(context.Background(), &dto.ProofreadRequest{
		Content: "张三在见李四。",
	})
	if err != nil {
		t.Fatalf("ProofreadContent returned error: %v", err)
	}

	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Type != "typo" {
		t.Fatalf("expected spelling to normalize to typo, got %q", result.Issues[0].Type)
	}
	if result.Issues[0].Severity != "warning" {
		t.Fatalf("expected medium to normalize to warning, got %q", result.Issues[0].Severity)
	}
	if result.Issues[0].Suggestions[0] != "再见" {
		t.Fatalf("expected legacy string suggestion")
	}
}

func TestProofreadContentCachesResultsByContentHashAndOptions(t *testing.T) {
	generator := &mockProofreadTextGenerator{
		responseText: `{
			"issues": [
				{
					"type": "typo",
					"severity": "error",
					"message": "疑似错别字",
					"start": 2,
					"end": 4,
					"originalText": "在见",
					"suggestions": ["再见"]
				}
			]
		}`,
	}
	service := NewProofreadService(generator)
	request := &dto.ProofreadRequest{
		Content:    "张三在见李四。",
		CheckTypes: []string{"typo", "grammar"},
	}

	first, err := service.ProofreadContent(context.Background(), request)
	if err != nil {
		t.Fatalf("first ProofreadContent returned error: %v", err)
	}
	first.Issues[0].Message = "调用方本地修改"
	first.Issues[0].Suggestions[0] = "调用方修改建议"
	first.Issues[0].SuggestionDetails[0].Text = "调用方修改详情"

	second, err := service.ProofreadContent(context.Background(), &dto.ProofreadRequest{
		Content:    "张三在见李四。",
		CheckTypes: []string{"typo", "grammar"},
	})
	if err != nil {
		t.Fatalf("second ProofreadContent returned error: %v", err)
	}

	if generator.callCount != 1 {
		t.Fatalf("expected cached second call, got generator call count %d", generator.callCount)
	}
	if second.ReviewID != first.ReviewID {
		t.Fatalf("expected cached review id to be reused")
	}
	if second.Issues[0].Message == "调用方本地修改" {
		t.Fatalf("expected cached response to be cloned")
	}
	if second.Issues[0].Suggestions[0] == "调用方修改建议" {
		t.Fatalf("expected cached suggestion slice to be cloned")
	}
	if second.Issues[0].SuggestionDetails[0].Text == "调用方修改详情" {
		t.Fatalf("expected cached suggestion details to be cloned")
	}

	_, err = service.ProofreadContent(context.Background(), &dto.ProofreadRequest{
		Content:    "张三在见李四。",
		CheckTypes: []string{"readability"},
	})
	if err != nil {
		t.Fatalf("third ProofreadContent returned error: %v", err)
	}
	if generator.callCount != 2 {
		t.Fatalf("expected different check types to bypass cache, got call count %d", generator.callCount)
	}
}
