package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Deprecated: GLM calls should go through Qingyu-Ai-Service.
// This adapter is kept for emergency fallback only.
type GLMAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Deprecated: Use gRPC client instead
func NewGLMAdapter(apiKey, baseURL string) *GLMAdapter {
	logrus.Warn("GLMAdapter is deprecated. Use Qingyu-Ai-Service gRPC API.")
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn/api/paas/v4"
	}

	return &GLMAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GetName 获取适配器名称
func (a *GLMAdapter) GetName() string {
	return "glm"
}

// GetSupportedModels 获取支持的模型列表
func (a *GLMAdapter) GetSupportedModels() []string {
	return []string{
		"glm-4-flash", // 闪现版，速度快，价格低
		"glm-4-plus",  // 增强版
		"glm-4-air",   // 轻量版
		"glm-4",       // 基础版
		"glm-3-turbo", // Turbo版
	}
}

// HealthCheck 健康检查
func (a *GLMAdapter) HealthCheck(ctx context.Context) error {
	// 简单的健康检查，发送一个最小请求
	req := &ChatCompletionRequest{
		Model:     "glm-4-flash",
		Messages:  []Message{{Role: "user", Content: "test"}},
		MaxTokens: 5,
	}

	_, err := a.ChatCompletion(ctx, req)
	return err
}

// GLMRequest 智谱AI请求结构
type GLMRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

// GLMResponse 智谱AI响应结构
type GLMResponse struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int      `json:"index"`
		Message Message  `json:"message"`
		Delta   *Message `json:"delta,omitempty"` // 用于流式响应
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// ChatCompletion 实现对话接口
func (a *GLMAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 构建 GLM API 请求
	glmReq := GLMRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   false,
	}

	// 使用默认模型
	if glmReq.Model == "" {
		glmReq.Model = "glm-4-flash"
	}

	reqBody, err := json.Marshal(glmReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/chat/completions", a.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var glmResp GLMResponse
	if err := json.Unmarshal(body, &glmResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查错误
	if glmResp.Error != nil {
		return nil, fmt.Errorf("智谱AI错误: %s", glmResp.Error.Message)
	}

	// 检查是否有响应
	if len(glmResp.Choices) == 0 {
		return nil, fmt.Errorf("智谱AI返回空响应")
	}

	// 转换为标准响应
	result := &ChatCompletionResponse{
		ID:      glmResp.ID,
		Message: glmResp.Choices[0].Message,
		Usage: Usage{
			PromptTokens:     glmResp.Usage.PromptTokens,
			CompletionTokens: glmResp.Usage.CompletionTokens,
			TotalTokens:      glmResp.Usage.TotalTokens,
		},
		Model:        glmResp.Model,
		FinishReason: "stop",
		CreatedAt:    time.Now(),
	}

	return result, nil
}

// TextGeneration 实现文本生成接口
func (a *GLMAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 将文本生成请求转换为对话请求
	chatReq := &ChatCompletionRequest{
		Model: req.Model,
		Messages: []Message{
			{Role: "user", Content: req.Prompt},
		},
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	// 调用对话接口
	chatResp, err := a.ChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, err
	}

	// 转换响应格式
	return &TextGenerationResponse{
		ID:           chatResp.ID,
		Text:         chatResp.Message.Content,
		Usage:        chatResp.Usage,
		Model:        chatResp.Model,
		FinishReason: chatResp.FinishReason,
		CreatedAt:    chatResp.CreatedAt,
	}, nil
}

// TextGenerationStream 流式文本生成（暂未实现）
func (a *GLMAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	return nil, fmt.Errorf("流式生成功能暂未实现")
}

// ImageGeneration 图像生成（智谱AI不支持）
func (a *GLMAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	return nil, fmt.Errorf("智谱AI不支持图像生成")
}
