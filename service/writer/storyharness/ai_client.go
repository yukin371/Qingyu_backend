package storyharness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"Qingyu_backend/config"
)

const (
	storyAnalysisPath     = "/api/v1/story/analyze-chapter"
	defaultAIBaseURL      = "https://api.openai.com/v1"
	defaultStoryAITimeout = 5 * time.Second
)

// ChapterAnalysisClient 抽象章节分析客户端，便于测试注入。
type ChapterAnalysisClient interface {
	AnalyzeChapter(ctx context.Context, request *ChapterAnalysisRequest) (*ChapterAnalysisResponse, error)
}

// StoryAnalysisHTTPClient 通过 HTTP 调用 Python AI 服务。
type StoryAnalysisHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewConfiguredChapterAnalysisClient 基于全局配置创建客户端。
func NewConfiguredChapterAnalysisClient() ChapterAnalysisClient {
	if config.GlobalConfig == nil || config.GlobalConfig.AI == nil {
		return nil
	}

	baseURL := strings.TrimSpace(config.GlobalConfig.AI.BaseURL)
	if baseURL == "" {
		return nil
	}
	if strings.TrimRight(baseURL, "/") == defaultAIBaseURL {
		// 默认 OpenAI BaseURL 无法提供 story analysis 路由，直接禁用并回退规则版。
		return nil
	}

	timeout := defaultStoryAITimeout
	if config.GlobalConfig.AI.AIService != nil && config.GlobalConfig.AI.AIService.Timeout > 0 {
		timeout = time.Duration(config.GlobalConfig.AI.AIService.Timeout) * time.Second
	}

	return NewStoryAnalysisHTTPClient(baseURL, &http.Client{Timeout: timeout})
}

// NewStoryAnalysisHTTPClient 创建 HTTP 客户端。
func NewStoryAnalysisHTTPClient(baseURL string, httpClient *http.Client) *StoryAnalysisHTTPClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultStoryAITimeout}
	}

	return &StoryAnalysisHTTPClient{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: httpClient,
	}
}

// AnalyzeChapter 调用 Python AI 服务的章节分析接口。
func (c *StoryAnalysisHTTPClient) AnalyzeChapter(ctx context.Context, request *ChapterAnalysisRequest) (*ChapterAnalysisResponse, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("story analysis client is not configured")
	}
	if request == nil {
		return nil, fmt.Errorf("story analysis request is nil")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal story analysis request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+storyAnalysisPath, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("build story analysis request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("call story analysis service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return nil, fmt.Errorf("story analysis service returned %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var result ChapterAnalysisResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode story analysis response: %w", err)
	}

	return &result, nil
}
