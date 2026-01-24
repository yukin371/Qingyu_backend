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

// Deprecated: Qwen calls should go through Qingyu-Ai-Service.
// This adapter is kept for emergency fallback only.
type QwenAdapter struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Deprecated: Use gRPC client instead
func NewQwenAdapter(apiKey, baseURL string) *QwenAdapter {
	logrus.Warn("QwenAdapter is deprecated. Use Qingyu-Ai-Service gRPC API.")
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com"
	}

	return &QwenAdapter{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName 获取适配器名称
func (a *QwenAdapter) GetName() string {
	return "qwen"
}

// GetSupportedModels 获取支持的模型列表
func (a *QwenAdapter) GetSupportedModels() []string {
	return []string{
		"qwen-turbo",
		"qwen-plus",
		"qwen-max",
		"qwen-max-1201",
		"qwen-max-longcontext",
		"qwen-7b-chat",
		"qwen-14b-chat",
		"qwen-72b-chat",
		"qwen1.5-7b-chat",
		"qwen1.5-14b-chat",
		"qwen1.5-72b-chat",
		"qwen2-7b-instruct",
		"qwen2-72b-instruct",
	}
}

// QwenMessage 通义千问消息格式
type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// QwenParameters 通义千问参数
type QwenParameters struct {
	ResultFormat      string   `json:"result_format,omitempty"`
	Seed              int      `json:"seed,omitempty"`
	MaxTokens         int      `json:"max_tokens,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	RepetitionPenalty float64  `json:"repetition_penalty,omitempty"`
	Temperature       float64  `json:"temperature,omitempty"`
	Stop              []string `json:"stop,omitempty"`
	EnableSearch      bool     `json:"enable_search,omitempty"`
	IncrementalOutput bool     `json:"incremental_output,omitempty"`
}

// QwenInput 通义千问输入
type QwenInput struct {
	Messages []QwenMessage `json:"messages"`
}

// QwenRequest 通义千问API请求结构
type QwenRequest struct {
	Model      string         `json:"model"`
	Input      QwenInput      `json:"input"`
	Parameters QwenParameters `json:"parameters,omitempty"`
}

// QwenResponse 通义千问API响应结构
type QwenResponse struct {
	Output    QwenOutput `json:"output"`
	Usage     QwenUsage  `json:"usage"`
	RequestID string     `json:"request_id"`
}

// QwenOutput 通义千问输出
type QwenOutput struct {
	Text         string       `json:"text"`
	FinishReason string       `json:"finish_reason"`
	Choices      []QwenChoice `json:"choices,omitempty"`
}

// QwenChoice 通义千问选择
type QwenChoice struct {
	FinishReason string      `json:"finish_reason"`
	Message      QwenMessage `json:"message"`
}

// QwenUsage 通义千问使用统计
type QwenUsage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// TextGeneration 实现文本生成方法
func (a *QwenAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 构建通义千问请求
	qwenReq := &QwenRequest{
		Model: req.Model,
		Input: QwenInput{
			Messages: []QwenMessage{
				{
					Role:    "user",
					Content: req.Prompt,
				},
			},
		},
		Parameters: QwenParameters{
			ResultFormat:      "text",
			MaxTokens:         req.MaxTokens,
			Temperature:       req.Temperature,
			TopP:              req.TopP,
			Stop:              req.Stop,
			IncrementalOutput: false,
		},
	}

	// 发送请求
	qwenResp, err := a.sendRequest(ctx, "/api/v1/services/aigc/text-generation/generation", qwenReq)
	if err != nil {
		return nil, err
	}

	return &TextGenerationResponse{
		ID:   qwenResp.RequestID,
		Text: qwenResp.Output.Text,
		Usage: Usage{
			PromptTokens:     qwenResp.Usage.InputTokens,
			CompletionTokens: qwenResp.Usage.OutputTokens,
			TotalTokens:      qwenResp.Usage.TotalTokens,
		},
		Model:        req.Model,
		FinishReason: qwenResp.Output.FinishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 实现对话完成方法
func (a *QwenAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 转换消息格式
	var qwenMessages []QwenMessage

	for _, msg := range req.Messages {
		qwenMessages = append(qwenMessages, QwenMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 构建通义千问请求
	qwenReq := &QwenRequest{
		Model: req.Model,
		Input: QwenInput{
			Messages: qwenMessages,
		},
		Parameters: QwenParameters{
			ResultFormat:      "message",
			MaxTokens:         req.MaxTokens,
			Temperature:       req.Temperature,
			TopP:              req.TopP,
			Stop:              req.Stop,
			IncrementalOutput: false,
		},
	}

	// 发送请求
	qwenResp, err := a.sendRequest(ctx, "/api/v1/services/aigc/text-generation/generation", qwenReq)
	if err != nil {
		return nil, err
	}

	// 构建响应
	var message Message
	var finishReason string
	if len(qwenResp.Output.Choices) > 0 {
		choice := qwenResp.Output.Choices[0]
		message = Message{
			Role:    choice.Message.Role,
			Content: choice.Message.Content,
		}
		finishReason = choice.FinishReason
	} else {
		// 如果没有choices，使用text字段
		message = Message{
			Role:    "assistant",
			Content: qwenResp.Output.Text,
		}
		finishReason = qwenResp.Output.FinishReason
	}

	return &ChatCompletionResponse{
		ID:      qwenResp.RequestID,
		Message: message,
		Usage: Usage{
			PromptTokens:     qwenResp.Usage.InputTokens,
			CompletionTokens: qwenResp.Usage.OutputTokens,
			TotalTokens:      qwenResp.Usage.TotalTokens,
		},
		Model:        req.Model,
		FinishReason: finishReason,
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 实现流式文本生成
func (a *QwenAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 通义千问流式API实现较复杂，这里先返回错误
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "通义千问流式生成暂未实现",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// ImageGeneration 实现图像生成方法
func (a *QwenAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// 通义千问有图像生成功能，但使用不同的API端点
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "通义千问图像生成暂未实现",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// HealthCheck 实现健康检查
func (a *QwenAdapter) HealthCheck(ctx context.Context) error {
	// 发送一个简单的请求来检查服务状态
	req := &QwenRequest{
		Model: "qwen-turbo",
		Input: QwenInput{
			Messages: []QwenMessage{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
		},
		Parameters: QwenParameters{
			MaxTokens: 10,
		},
	}

	_, err := a.sendRequest(ctx, "/api/v1/services/aigc/text-generation/generation", req)
	return err
}

// sendRequest 发送HTTP请求到通义千问API
func (a *QwenAdapter) sendRequest(ctx context.Context, endpoint string, request interface{}) (*QwenResponse, error) {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, NewAdapterError("qwen", ErrorTypeInvalidRequest,
			fmt.Sprintf("序列化请求失败: %v", err), "REQUEST_MARSHAL_ERROR", 0, false)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewAdapterError("qwen", ErrorTypeInvalidRequest,
			fmt.Sprintf("创建HTTP请求失败: %v", err), "REQUEST_CREATE_ERROR", 0, false)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

	// 发送请求
	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, NewAdapterError("qwen", ErrorTypeTimeout,
			fmt.Sprintf("发送HTTP请求失败: %v", err), "REQUEST_SEND_ERROR", 0, true)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewAdapterError("qwen", ErrorTypeUnknown,
			fmt.Sprintf("读取响应失败: %v", err), "RESPONSE_READ_ERROR", resp.StatusCode, false)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, a.handleErrorResponse(resp.StatusCode, responseBody)
	}

	// 解析响应
	var qwenResp QwenResponse
	if err := json.Unmarshal(responseBody, &qwenResp); err != nil {
		return nil, NewAdapterError("qwen", ErrorTypeInvalidResponse,
			fmt.Sprintf("解析响应失败: %v", err), "RESPONSE_PARSE_ERROR", resp.StatusCode, false)
	}

	return &qwenResp, nil
}

// handleErrorResponse 处理错误响应
func (a *QwenAdapter) handleErrorResponse(statusCode int, body []byte) error {
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
		message = "权限不足或配额不足"
		isRetryable = false
	case 429:
		errorType = ErrorTypeRateLimit
		message = "请求频率超限"
		isRetryable = true
	case 500, 502, 503:
		errorType = ErrorTypeServiceUnavailable
		message = "通义千问服务暂时不可用"
		isRetryable = true
	default:
		errorType = ErrorTypeUnknown
		message = "未知错误"
		isRetryable = false
	}

	// 尝试解析错误详情
	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if errorMsg, ok := errorResp["message"].(string); ok {
			message = errorMsg
		} else if errorCode, ok := errorResp["code"].(string); ok {
			message = fmt.Sprintf("错误代码: %s", errorCode)
		}
	}

	return NewAdapterError("qwen", errorType, message,
		fmt.Sprintf("QWEN_%d", statusCode), statusCode, isRetryable)
}
