package adapter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GeminiAdapter Google Gemini 适配器实现
type GeminiAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewGeminiAdapter 创建Gemini适配器实例
func NewGeminiAdapter(apiKey, baseURL string) *GeminiAdapter {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}
	
	return &GeminiAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName 获取适配器名称
func (a *GeminiAdapter) GetName() string {
	return "gemini"
}

// GetSupportedModels 获取支持的模型列表
func (a *GeminiAdapter) GetSupportedModels() []string {
	return []string{
		"gemini-1.5-pro",
		"gemini-1.5-flash",
		"gemini-1.0-pro",
		"gemini-pro-vision",
	}
}

// GeminiContent Gemini内容结构
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

// GeminiPart Gemini内容部分
type GeminiPart struct {
	Text string `json:"text,omitempty"`
}

// GeminiRequest Gemini API请求结构
type GeminiRequest struct {
	Contents         []GeminiContent           `json:"contents"`
	GenerationConfig *GeminiGenerationConfig   `json:"generationConfig,omitempty"`
	SafetySettings   []GeminiSafetySetting     `json:"safetySettings,omitempty"`
}

// GeminiGenerationConfig Gemini生成配置
type GeminiGenerationConfig struct {
	Temperature     float64  `json:"temperature,omitempty"`
	TopP            float64  `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}

// GeminiSafetySetting Gemini安全设置
type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GeminiResponse Gemini API响应结构
type GeminiResponse struct {
	Candidates     []GeminiCandidate `json:"candidates"`
	UsageMetadata  GeminiUsage       `json:"usageMetadata"`
	PromptFeedback GeminiPromptFeedback `json:"promptFeedback,omitempty"`
}

// GeminiCandidate Gemini候选结果
type GeminiCandidate struct {
	Content       GeminiContent `json:"content"`
	FinishReason  string        `json:"finishReason"`
	Index         int           `json:"index"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings"`
}

// GeminiSafetyRating Gemini安全评级
type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// GeminiPromptFeedback Gemini提示反馈
type GeminiPromptFeedback struct {
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings"`
}

// GeminiUsage Gemini使用统计
type GeminiUsage struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// Gemini流式响应结构体
type GeminiStreamResponse struct {
	Candidates []GeminiStreamCandidate `json:"candidates"`
	UsageMetadata *GeminiUsage `json:"usageMetadata,omitempty"`
}

type GeminiStreamCandidate struct {
	Content GeminiContent `json:"content"`
	FinishReason string `json:"finishReason,omitempty"`
	Index int `json:"index"`
}

// TextGeneration 实现文本生成方法
func (a *GeminiAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 构建Gemini请求
	geminiReq := &GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: req.Prompt},
				},
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     req.Temperature,
			TopP:            req.TopP,
			MaxOutputTokens: req.MaxTokens,
			StopSequences:   req.Stop,
		},
		SafetySettings: a.getDefaultSafetySettings(),
	}

	// 发送请求
	geminiResp, err := a.sendRequest(ctx, fmt.Sprintf("/v1beta/models/%s:generateContent", req.Model), geminiReq)
	if err != nil {
		return nil, err
	}

	// 提取文本内容
	var text string
	var finishReason string
	if len(geminiResp.Candidates) > 0 {
		candidate := geminiResp.Candidates[0]
		if len(candidate.Content.Parts) > 0 {
			text = candidate.Content.Parts[0].Text
		}
		finishReason = candidate.FinishReason
	}

	return &TextGenerationResponse{
		ID:   primitive.NewObjectID().Hex(),
		Text: text,
		Usage: Usage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
		Model:        req.Model,
		FinishReason: finishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 实现对话完成方法
func (a *GeminiAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 转换消息格式
	var contents []GeminiContent
	
	for _, msg := range req.Messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		
		contents = append(contents, GeminiContent{
			Parts: []GeminiPart{
				{Text: msg.Content},
			},
			Role: role,
		})
	}

	// 构建Gemini请求
	geminiReq := &GeminiRequest{
		Contents: contents,
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     req.Temperature,
			TopP:            req.TopP,
			MaxOutputTokens: req.MaxTokens,
			StopSequences:   req.Stop,
		},
		SafetySettings: a.getDefaultSafetySettings(),
	}

	// 发送请求
	geminiResp, err := a.sendRequest(ctx, fmt.Sprintf("/v1beta/models/%s:generateContent", req.Model), geminiReq)
	if err != nil {
		return nil, err
	}

	// 构建响应
	var message Message
	var finishReason string
	if len(geminiResp.Candidates) > 0 {
		candidate := geminiResp.Candidates[0]
		var text string
		if len(candidate.Content.Parts) > 0 {
			text = candidate.Content.Parts[0].Text
		}
		
		message = Message{
			Role:    "assistant",
			Content: text,
		}
		finishReason = candidate.FinishReason
	}

	return &ChatCompletionResponse{
		ID:      primitive.NewObjectID().Hex(),
		Message: message,
		Usage: Usage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
		Model:        req.Model,
		FinishReason: finishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 实现流式文本生成
func (a *GeminiAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 创建响应通道
	responseChan := make(chan *TextGenerationResponse, 10)
	
	// 启动协程处理流式响应
	go func() {
		defer close(responseChan)
		
		if err := a.doTextGenerationStream(ctx, req, responseChan); err != nil {
			// 发送错误响应
			responseChan <- &TextGenerationResponse{
				ID:           primitive.NewObjectID().Hex(),
				Text:         fmt.Sprintf("流式生成失败: %v", err),
				Model:        req.Model,
				FinishReason: "error",
				CreatedAt:    time.Now(),
			}
		}
	}()
	
	return responseChan, nil
}

// doTextGenerationStream 执行流式文本生成
func (a *GeminiAdapter) doTextGenerationStream(ctx context.Context, req *TextGenerationRequest, responseChan chan<- *TextGenerationResponse) error {
	// 构建Gemini请求
	geminiReq := &GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: req.Prompt},
				},
				Role: "user",
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     req.Temperature,
			TopP:            req.TopP,
			MaxOutputTokens: req.MaxTokens,
			StopSequences:   req.Stop,
		},
		SafetySettings: a.getDefaultSafetySettings(),
	}

	// 发送流式请求
	return a.sendStreamRequest(ctx, fmt.Sprintf("/v1/models/%s:streamGenerateContent", req.Model), geminiReq, responseChan)
}

// sendStreamRequest 发送流式HTTP请求
func (a *GeminiAdapter) sendStreamRequest(ctx context.Context, endpoint string, request interface{}, responseChan chan<- *TextGenerationResponse) error {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint+"?key="+a.apiKey+"&alt=sse", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 处理流式响应
	return a.processStreamResponse(resp.Body, responseChan)
}

// processStreamResponse 处理流式响应
func (a *GeminiAdapter) processStreamResponse(body io.Reader, responseChan chan<- *TextGenerationResponse) error {
	scanner := bufio.NewScanner(body)
	var fullContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// 解析SSE数据
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// 检查是否为结束标记
			if data == "[DONE]" {
				break
			}

			// 解析JSON数据
			var streamResp GeminiStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue // 跳过无法解析的数据
			}

			// 处理候选响应
			if len(streamResp.Candidates) > 0 {
				candidate := streamResp.Candidates[0]
				if len(candidate.Content.Parts) > 0 {
					text := candidate.Content.Parts[0].Text
					if text != "" {
						fullContent.WriteString(text)
						
						// 发送增量响应
						responseChan <- &TextGenerationResponse{
							ID:           primitive.NewObjectID().Hex(),
							Text:         text,
							Model:        "gemini",
							Usage:        Usage{TotalTokens: 0}, // Gemini流式响应中token信息在最后提供
							FinishReason: "",
							CreatedAt:    time.Now(),
						}
					}
				}
			}
		}
	}

	// 发送最终完整响应
	if fullContent.Len() > 0 {
		responseChan <- &TextGenerationResponse{
			ID:           primitive.NewObjectID().Hex(),
			Text:         fullContent.String(),
			Model:        "gemini",
			Usage:        Usage{TotalTokens: 0},
			FinishReason: "stop",
			CreatedAt:    time.Now(),
		}
	}

	return scanner.Err()
}

// ImageGeneration 实现图像生成方法
func (a *GeminiAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// Gemini目前主要用于文本生成，图像生成功能有限
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "Gemini图像生成功能暂未实现",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// HealthCheck 实现健康检查
func (a *GeminiAdapter) HealthCheck(ctx context.Context) error {
	// 发送一个简单的请求来检查服务状态
	req := &GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: "Hello"},
				},
			},
		},
		GenerationConfig: &GeminiGenerationConfig{
			MaxOutputTokens: 10,
		},
		SafetySettings: a.getDefaultSafetySettings(),
	}

	_, err := a.sendRequest(ctx, "/v1beta/models/gemini-1.0-pro:generateContent", req)
	return err
}

// getDefaultSafetySettings 获取默认安全设置
func (a *GeminiAdapter) getDefaultSafetySettings() []GeminiSafetySetting {
	return []GeminiSafetySetting{
		{
			Category:  "HARM_CATEGORY_HARASSMENT",
			Threshold: "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			Category:  "HARM_CATEGORY_HATE_SPEECH",
			Threshold: "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
			Threshold: "BLOCK_MEDIUM_AND_ABOVE",
		},
		{
			Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
			Threshold: "BLOCK_MEDIUM_AND_ABOVE",
		},
	}
}

// sendRequest 发送HTTP请求到Gemini API
func (a *GeminiAdapter) sendRequest(ctx context.Context, endpoint string, request interface{}) (*GeminiResponse, error) {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("序列化请求失败: %v", err),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s%s?key=%s", a.baseURL, endpoint, a.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("创建HTTP请求失败: %v", err),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeTimeout,
			Message: fmt.Sprintf("发送HTTP请求失败: %v", err),
			Type:    ErrorTypeTimeout,
		}
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeUnknown,
			Message: fmt.Sprintf("读取响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, a.handleErrorResponse(resp.StatusCode, responseBody)
	}

	// 解析响应
	var geminiResp GeminiResponse
	if err := json.Unmarshal(responseBody, &geminiResp); err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeUnknown,
			Message: fmt.Sprintf("解析响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	return &geminiResp, nil
}

// handleErrorResponse 处理错误响应
func (a *GeminiAdapter) handleErrorResponse(statusCode int, body []byte) error {
	var errorType string
	var message string

	switch statusCode {
	case 400:
		errorType = ErrorTypeInvalidRequest
		message = "请求参数错误"
	case 401:
		errorType = ErrorTypeAuthentication
		message = "API密钥无效"
	case 403:
		errorType = ErrorTypePermission
		message = "权限不足或配额不足"
	case 429:
		errorType = ErrorTypeRateLimit
		message = "请求频率超限"
	case 500, 502, 503:
		errorType = ErrorTypeServiceUnavailable
		message = "Gemini服务暂时不可用"
	default:
		errorType = ErrorTypeUnknown
		message = "未知错误"
	}

	// 尝试解析错误详情
	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if errorInfo, ok := errorResp["error"].(map[string]interface{}); ok {
			if msg, ok := errorInfo["message"].(string); ok {
				message = msg
			}
		}
	}

	return &AdapterError{
		Code:    fmt.Sprintf("gemini_%d", statusCode),
		Message: message,
		Type:    errorType,
	}
}