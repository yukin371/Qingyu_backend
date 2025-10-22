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

// OpenAIAdapter OpenAI适配器实现
type OpenAIAdapter struct {
	apiKey       string
	baseURL      string
	client       *http.Client
	errorHandler *ErrorHandler
}

// NewOpenAIAdapter 创建OpenAI适配器实例
func NewOpenAIAdapter(apiKey, baseURL string) *OpenAIAdapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// 创建错误处理器
	retryConfig := DefaultRetryConfig()
	errorHandler := NewErrorHandler(
		retryConfig,
		5,              // 最大失败次数
		30*time.Second, // 熔断器重置时间
		100,            // 限流容量
		time.Second/10, // 限流补充速率 (10 QPS)
	)

	return &OpenAIAdapter{
		apiKey:       apiKey,
		baseURL:      baseURL,
		client:       &http.Client{Timeout: 30 * time.Second},
		errorHandler: errorHandler,
	}
}

// GetName 获取适配器名称
func (a *OpenAIAdapter) GetName() string {
	return "openai"
}

// GetSupportedModels 获取支持的模型列表
func (a *OpenAIAdapter) GetSupportedModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
		"text-davinci-003",
		"text-curie-001",
		"text-babbage-001",
		"text-ada-001",
		"dall-e-3",
		"dall-e-2",
	}
}

// OpenAIRequest OpenAI请求结构
type OpenAIRequest struct {
	Model       string                 `json:"model"`
	Messages    []OpenAIMessage        `json:"messages,omitempty"`
	Prompt      string                 `json:"prompt,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"top_p,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Stop        []string               `json:"stop,omitempty"`
	User        string                 `json:"user,omitempty"`
	Extra       map[string]interface{} `json:"-"`
}

// OpenAIMessage OpenAI消息结构
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse OpenAI响应结构
type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   OpenAIUsage    `json:"usage"`
	Error   *OpenAIError   `json:"error,omitempty"`
}

// OpenAIChoice OpenAI选择结构
type OpenAIChoice struct {
	Index        int            `json:"index"`
	Message      *OpenAIMessage `json:"message,omitempty"`
	Text         string         `json:"text,omitempty"`
	FinishReason string         `json:"finish_reason"`
}

// OpenAIUsage OpenAI使用量结构
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAIError OpenAI错误结构
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// OpenAIStreamResponse OpenAI流式响应结构
type OpenAIStreamResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []OpenAIStreamChoice `json:"choices"`
	Usage   *OpenAIUsage         `json:"usage,omitempty"`
	Error   *OpenAIError         `json:"error,omitempty"`
}

// OpenAIStreamChoice OpenAI流式选择结构
type OpenAIStreamChoice struct {
	Index        int                `json:"index"`
	Delta        *OpenAIStreamDelta `json:"delta,omitempty"`
	FinishReason *string            `json:"finish_reason"`
}

// OpenAIStreamDelta OpenAI流式增量结构
type OpenAIStreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// TextGeneration 实现文本生成方法
func (a *OpenAIAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	result, err := ExecuteWithResult(ctx, a.errorHandler.retryer, func(ctx context.Context) (interface{}, error) {
		return a.doTextGeneration(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return result.(*TextGenerationResponse), nil
}

// doTextGeneration 执行文本生成
func (a *OpenAIAdapter) doTextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	openaiReq := &OpenAIRequest{
		Model:       req.Model,
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		User:        req.User,
	}

	resp, err := a.sendRequest(ctx, "/completions", openaiReq)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidResponse,
			Message: "OpenAI API 返回空响应",
			Type:    ErrorTypeInvalidResponse,
		}
	}

	return &TextGenerationResponse{
		ID:   resp.ID,
		Text: resp.Choices[0].Text,
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Model:        resp.Model,
		FinishReason: resp.Choices[0].FinishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 实现对话完成方法
func (a *OpenAIAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	result, err := ExecuteWithResult(ctx, a.errorHandler.retryer, func(ctx context.Context) (interface{}, error) {
		return a.doChatCompletion(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ChatCompletionResponse), nil
}

// doChatCompletion 执行对话完成
func (a *OpenAIAdapter) doChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	messages := make([]OpenAIMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = OpenAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	openaiReq := &OpenAIRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		User:        req.User,
	}

	resp, err := a.sendRequest(ctx, "/chat/completions", openaiReq)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, &AdapterError{
			Code:    ErrorTypeInvalidResponse,
			Message: "OpenAI API 返回空响应",
			Type:    ErrorTypeInvalidResponse,
		}
	}

	choice := resp.Choices[0]
	var message Message
	if choice.Message != nil {
		message = Message{
			Role:    choice.Message.Role,
			Content: choice.Message.Content,
		}
	}

	return &ChatCompletionResponse{
		ID:      resp.ID,
		Message: message,
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Model:        resp.Model,
		FinishReason: choice.FinishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 实现流式文本生成方法
func (a *OpenAIAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	result, err := ExecuteWithResult(ctx, a.errorHandler.retryer, func(ctx context.Context) (interface{}, error) {
		return a.doTextGenerationStream(ctx, req)
	})

	if err != nil {
		return nil, err
	}

	return result.(<-chan *TextGenerationResponse), nil
}

// doTextGenerationStream 执行流式文本生成
func (a *OpenAIAdapter) doTextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 创建响应通道
	responseChan := make(chan *TextGenerationResponse, 10)

	// 构建OpenAI请求
	openaiReq := &OpenAIRequest{
		Model:       req.Model,
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      true, // 启用流式响应
		Stop:        req.Stop,
		User:        req.User,
	}

	// 发起流式请求
	go func() {
		defer close(responseChan)

		if err := a.sendStreamRequest(ctx, "completions", openaiReq, responseChan); err != nil {
			// 发送错误到通道
			select {
			case responseChan <- &TextGenerationResponse{
				ID:           primitive.NewObjectID().Hex(),
				Text:         "",
				Usage:        Usage{},
				Model:        req.Model,
				FinishReason: "error",
				CreatedAt:    time.Now(),
			}:
			case <-ctx.Done():
			}
		}
	}()

	return responseChan, nil
}

// sendStreamRequest 发送流式请求
func (a *OpenAIAdapter) sendStreamRequest(ctx context.Context, endpoint string, req *OpenAIRequest, responseChan chan<- *TextGenerationResponse) error {
	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/%s", a.baseURL, endpoint)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return a.handleHTTPError(resp.StatusCode, string(body))
	}

	// 处理流式响应
	return a.processStreamResponse(ctx, resp.Body, responseChan, req.Model)
}

// processStreamResponse 处理流式响应
func (a *OpenAIAdapter) processStreamResponse(ctx context.Context, body io.Reader, responseChan chan<- *TextGenerationResponse, model string) error {
	scanner := bufio.NewScanner(body)
	var fullText strings.Builder
	var totalTokens int

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// 处理SSE数据行
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// 检查结束标记
			if data == "[DONE]" {
				// 发送最终响应
				select {
				case responseChan <- &TextGenerationResponse{
					ID:           primitive.NewObjectID().Hex(),
					Text:         fullText.String(),
					Usage:        Usage{TotalTokens: totalTokens},
					Model:        model,
					FinishReason: "stop",
					CreatedAt:    time.Now(),
				}:
				case <-ctx.Done():
				}
				break
			}

			// 解析JSON数据
			var streamResp OpenAIStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue // 跳过无法解析的数据
			}

			// 处理错误
			if streamResp.Error != nil {
				return a.handleAPIError(400, streamResp.Error)
			}

			// 处理选择
			if len(streamResp.Choices) > 0 {
				choice := streamResp.Choices[0]

				// 处理增量内容
				if choice.Delta != nil && choice.Delta.Content != "" {
					fullText.WriteString(choice.Delta.Content)

					// 发送增量响应
					select {
					case responseChan <- &TextGenerationResponse{
						ID:           streamResp.ID,
						Text:         choice.Delta.Content,
						Usage:        Usage{},
						Model:        streamResp.Model,
						FinishReason: "",
						CreatedAt:    time.Now(),
					}:
					case <-ctx.Done():
						return ctx.Err()
					}
				}

				// 检查完成原因
				if choice.FinishReason != nil && *choice.FinishReason != "" {
					// 更新token使用量
					if streamResp.Usage != nil {
						totalTokens = streamResp.Usage.TotalTokens
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流式响应失败: %w", err)
	}

	return nil
}

// ImageGeneration 实现图像生成方法
func (a *OpenAIAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	result, err := ExecuteWithResult(ctx, a.errorHandler.retryer, func(ctx context.Context) (interface{}, error) {
		return a.doImageGeneration(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	return result.(*ImageGenerationResponse), nil
}

// doImageGeneration 执行图像生成
func (a *OpenAIAdapter) doImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// TODO: 实现图像生成逻辑
	return &ImageGenerationResponse{
		ID:        "openai_img_" + primitive.NewObjectID().Hex(),
		Images:    []ImageData{{URL: "https://example.com/generated-image.png"}},
		Model:     req.Model,
		CreatedAt: time.Now(),
	}, nil
}

// HealthCheck 实现健康检查方法
func (a *OpenAIAdapter) HealthCheck(ctx context.Context) error {
	return a.errorHandler.Execute(ctx, func(ctx context.Context) error {
		return a.doHealthCheck(ctx)
	})
}

// doHealthCheck 执行健康检查
func (a *OpenAIAdapter) doHealthCheck(ctx context.Context) error {
	req := &OpenAIRequest{
		Model:     "gpt-3.5-turbo",
		Messages:  []OpenAIMessage{{Role: "user", Content: "test"}},
		MaxTokens: 1,
	}

	_, err := a.sendRequest(ctx, "/chat/completions", req)
	return err
}

// sendRequest 发送HTTP请求
func (a *OpenAIAdapter) sendRequest(ctx context.Context, endpoint string, req *OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, NewAdapterError("openai", ErrorTypeInvalidRequest,
			fmt.Sprintf("序列化请求失败: %v", err), "REQUEST_MARSHAL_ERROR", 0, false)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, NewAdapterError("openai", ErrorTypeNetworkError,
			"创建请求失败", "REQUEST_CREATE_ERROR", 0, true)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, NewAdapterError("openai", ErrorTypeNetworkError,
			fmt.Sprintf("发送请求失败: %v", err), "REQUEST_SEND_ERROR", 0, true)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, NewAdapterError("openai", ErrorTypeInvalidResponse,
			fmt.Sprintf("读取响应失败: %v", err), "RESPONSE_READ_ERROR", httpResp.StatusCode, false)
	}

	var resp OpenAIResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, NewAdapterError("openai", ErrorTypeInvalidResponse,
			fmt.Sprintf("解析响应失败: %v", err), "RESPONSE_PARSE_ERROR", httpResp.StatusCode, false)
	}

	// 检查API错误
	if resp.Error != nil {
		return nil, a.handleAPIError(httpResp.StatusCode, resp.Error)
	}

	// 检查HTTP状态码
	if httpResp.StatusCode >= 400 {
		return nil, a.handleHTTPError(httpResp.StatusCode, string(body))
	}

	return &resp, nil
}

// handleAPIError 处理API错误
func (a *OpenAIAdapter) handleAPIError(statusCode int, apiErr *OpenAIError) error {
	var errorType string
	var isRetryable bool

	switch apiErr.Type {
	case "insufficient_quota":
		errorType = ErrorTypeRateLimit
		isRetryable = false
	case "invalid_request_error":
		errorType = ErrorTypeInvalidRequest
		isRetryable = false
	case "authentication_error":
		errorType = ErrorTypeAuthentication
		isRetryable = false
	case "rate_limit_exceeded":
		errorType = ErrorTypeRateLimit
		isRetryable = true
	case "server_error":
		errorType = ErrorTypeServiceUnavailable
		isRetryable = true
	default:
		errorType = ErrorTypeUnknown
		isRetryable = false
	}

	return NewAdapterError("openai", errorType, apiErr.Message,
		apiErr.Code, statusCode, isRetryable)
}

// handleHTTPError 处理HTTP错误
func (a *OpenAIAdapter) handleHTTPError(statusCode int, body string) error {
	var errorType string
	var message string
	var isRetryable bool

	switch {
	case statusCode == 429:
		errorType = ErrorTypeRateLimit
		message = "请求频率过高，请稍后重试"
		isRetryable = true
	case statusCode >= 500:
		errorType = ErrorTypeServiceUnavailable
		message = "OpenAI 服务暂时不可用"
		isRetryable = true
	case statusCode == 401:
		errorType = ErrorTypeAuthentication
		message = "OpenAI API 密钥无效"
		isRetryable = false
	case statusCode == 400:
		errorType = ErrorTypeInvalidRequest
		message = "请求参数无效"
		isRetryable = false
	default:
		errorType = ErrorTypeUnknown
		message = fmt.Sprintf("未知错误，状态码: %d", statusCode)
		isRetryable = false
	}

	if strings.Contains(body, "timeout") {
		errorType = ErrorTypeTimeout
		message = "请求超时"
		isRetryable = true
	}

	return NewAdapterError("openai", errorType, message,
		fmt.Sprintf("HTTP_%d", statusCode), statusCode, isRetryable)
}
