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

// ClaudeAdapter Claude (Anthropic) 适配器实现
type ClaudeAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewClaudeAdapter 创建Claude适配器实例
func NewClaudeAdapter(apiKey, baseURL string) *ClaudeAdapter {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}

	return &ClaudeAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName 获取适配器名称
func (a *ClaudeAdapter) GetName() string {
	return "claude"
}

// GetSupportedModels 获取支持的模型列表
func (a *ClaudeAdapter) GetSupportedModels() []string {
	return []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-5-haiku-20241022",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}
}

// ClaudeMessage Claude消息格式
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeRequest Claude API请求结构
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	Stop        []string        `json:"stop_sequences,omitempty"`
	System      string          `json:"system,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

// ClaudeResponse Claude API响应结构
type ClaudeResponse struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Role         string          `json:"role"`
	Content      []ClaudeContent `json:"content"`
	Model        string          `json:"model"`
	StopReason   string          `json:"stop_reason"`
	StopSequence string          `json:"stop_sequence"`
	Usage        ClaudeUsage     `json:"usage"`
}

// ClaudeContent Claude内容结构
type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeUsage Claude使用统计
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Claude流式响应结构体
type ClaudeStreamResponse struct {
	Type  string                 `json:"type"`
	Index int                    `json:"index,omitempty"`
	Delta *ClaudeStreamDelta     `json:"delta,omitempty"`
	Usage *ClaudeUsage          `json:"usage,omitempty"`
}

type ClaudeStreamDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// TextGeneration 实现文本生成方法
func (a *ClaudeAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 构建Claude请求
	claudeReq := &ClaudeRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
	}

	// 发送请求
	claudeResp, err := a.sendRequest(ctx, "/v1/messages", claudeReq)
	if err != nil {
		return nil, err
	}

	// 提取文本内容
	var text string
	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == "text" {
		text = claudeResp.Content[0].Text
	}

	return &TextGenerationResponse{
		ID:   claudeResp.ID,
		Text: text,
		Usage: Usage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
		Model:        claudeResp.Model,
		FinishReason: claudeResp.StopReason,
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 实现对话完成方法
func (a *ClaudeAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 转换消息格式
	var claudeMessages []ClaudeMessage
	var systemMessage string

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			claudeMessages = append(claudeMessages, ClaudeMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// 构建Claude请求
	claudeReq := &ClaudeRequest{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Messages:    claudeMessages,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		System:      systemMessage,
	}

	// 发送请求
	claudeResp, err := a.sendRequest(ctx, "/v1/messages", claudeReq)
	if err != nil {
		return nil, err
	}

	// 提取文本内容
	var text string
	if len(claudeResp.Content) > 0 && claudeResp.Content[0].Type == "text" {
		text = claudeResp.Content[0].Text
	}

	return &ChatCompletionResponse{
		ID: claudeResp.ID,
		Message: Message{
			Role:    "assistant",
			Content: text,
		},
		Usage: Usage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
		Model:        claudeResp.Model,
		FinishReason: claudeResp.StopReason,
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 实现流式文本生成
func (a *ClaudeAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
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
func (a *ClaudeAdapter) doTextGenerationStream(ctx context.Context, req *TextGenerationRequest, responseChan chan<- *TextGenerationResponse) error {
	// 构建Claude请求
	claudeReq := &ClaudeRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		Stream:      true,
	}

	// 发送流式请求
	return a.sendStreamRequest(ctx, "/v1/messages", claudeReq, responseChan)
}

// sendStreamRequest 发送流式HTTP请求
func (a *ClaudeAdapter) sendStreamRequest(ctx context.Context, endpoint string, request interface{}, responseChan chan<- *TextGenerationResponse) error {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
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
func (a *ClaudeAdapter) processStreamResponse(body io.Reader, responseChan chan<- *TextGenerationResponse) error {
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
			var streamResp ClaudeStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue // 跳过无法解析的数据
			}

			// 处理内容增量
			if streamResp.Delta != nil && streamResp.Delta.Text != "" {
				fullContent.WriteString(streamResp.Delta.Text)
				
				// 发送增量响应
				responseChan <- &TextGenerationResponse{
					ID:           primitive.NewObjectID().Hex(),
					Text:         streamResp.Delta.Text,
					Model:        "claude",
					Usage:        Usage{TotalTokens: 0}, // Claude流式响应中token信息在最后提供
					FinishReason: "",
					CreatedAt:    time.Now(),
				}
			}
		}
	}

	// 发送最终完整响应
	if fullContent.Len() > 0 {
		responseChan <- &TextGenerationResponse{
			ID:           primitive.NewObjectID().Hex(),
			Text:         fullContent.String(),
			Model:        "claude",
			Usage:        Usage{TotalTokens: 0},
			FinishReason: "stop",
			CreatedAt:    time.Now(),
		}
	}

	return scanner.Err()
}

// ImageGeneration 实现图像生成方法
func (a *ClaudeAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// Claude目前不支持图像生成
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "Claude不支持图像生成",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// HealthCheck 实现健康检查
func (a *ClaudeAdapter) HealthCheck(ctx context.Context) error {
	// 发送一个简单的请求来检查服务状态
	req := &ClaudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 10,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	_, err := a.sendRequest(ctx, "/v1/messages", req)
	return err
}

// sendRequest 发送HTTP请求到Claude API
func (a *ClaudeAdapter) sendRequest(ctx context.Context, endpoint string, request interface{}) (*ClaudeResponse, error) {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, NewAdapterError("claude", ErrorTypeInvalidRequest, 
			fmt.Sprintf("序列化请求失败: %v", err), "REQUEST_MARSHAL_ERROR", 0, false)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewAdapterError("claude", ErrorTypeInvalidRequest, 
			fmt.Sprintf("创建HTTP请求失败: %v", err), "REQUEST_CREATE_ERROR", 0, false)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, NewAdapterError("claude", ErrorTypeTimeout, 
			fmt.Sprintf("发送HTTP请求失败: %v", err), "REQUEST_SEND_ERROR", 0, true)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewAdapterError("claude", ErrorTypeUnknown, 
			fmt.Sprintf("读取响应失败: %v", err), "RESPONSE_READ_ERROR", resp.StatusCode, false)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, a.handleErrorResponse(resp.StatusCode, responseBody)
	}

	// 解析响应
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(responseBody, &claudeResp); err != nil {
		return nil, NewAdapterError("claude", ErrorTypeInvalidResponse, 
			fmt.Sprintf("解析响应失败: %v", err), "RESPONSE_PARSE_ERROR", resp.StatusCode, false)
	}

	return &claudeResp, nil
}

// handleErrorResponse 处理错误响应
func (a *ClaudeAdapter) handleErrorResponse(statusCode int, body []byte) error {
	var errorType string
	var message string
	var isRetryable bool

	switch statusCode {
	case 400:
		errorType = ErrorTypeInvalidRequest
		message = "请求参数错误"
		isRetryable = false
	case 401:
		errorType = ErrorTypeAuthentication
		message = "API密钥无效"
		isRetryable = false
	case 403:
		errorType = ErrorTypePermission
		message = "权限不足"
		isRetryable = false
	case 429:
		errorType = ErrorTypeRateLimit
		message = "请求频率超限"
		isRetryable = true
	case 500, 502, 503:
		errorType = ErrorTypeServiceUnavailable
		message = "Claude服务暂时不可用"
		isRetryable = true
	default:
		errorType = ErrorTypeUnknown
		message = "未知错误"
		isRetryable = false
	}

	// 尝试解析错误详情
	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if errorMsg, ok := errorResp["error"].(map[string]interface{}); ok {
			if msg, ok := errorMsg["message"].(string); ok {
				message = msg
			}
		}
	}

	return NewAdapterError("claude", errorType, message, 
		fmt.Sprintf("CLAUDE_%d", statusCode), statusCode, isRetryable)
}
