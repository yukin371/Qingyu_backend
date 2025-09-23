package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WenxinAdapter 百度文心一言适配器实现
type WenxinAdapter struct {
	apiKey    string
	secretKey string
	baseURL   string
	client    *http.Client
	token     string
	tokenExp  time.Time
}

// NewWenxinAdapter 创建文心一言适配器实例
func NewWenxinAdapter(apiKey, secretKey, baseURL string) *WenxinAdapter {
	if baseURL == "" {
		baseURL = "https://aip.baidubce.com"
	}

	return &WenxinAdapter{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetName 获取适配器名称
func (a *WenxinAdapter) GetName() string {
	return "wenxin"
}

// GetSupportedModels 获取支持的模型列表
func (a *WenxinAdapter) GetSupportedModels() []string {
	return []string{
		"ernie-4.0-8k",
		"ernie-4.0-8k-preview",
		"ernie-3.5-8k",
		"ernie-3.5-8k-0205",
		"ernie-turbo-8k",
		"ernie-speed-8k",
		"ernie-lite-8k",
		"ernie-tiny-8k",
	}
}

// WenxinMessage 文心一言消息格式
type WenxinMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// WenxinRequest 文心一言API请求结构
type WenxinRequest struct {
	Messages        []WenxinMessage `json:"messages"`
	Temperature     float64         `json:"temperature,omitempty"`
	TopP            float64         `json:"top_p,omitempty"`
	PenaltyScore    float64         `json:"penalty_score,omitempty"`
	Stream          bool            `json:"stream,omitempty"`
	System          string          `json:"system,omitempty"`
	Stop            []string        `json:"stop,omitempty"`
	DisableSearch   bool            `json:"disable_search,omitempty"`
	EnableCitation  bool            `json:"enable_citation,omitempty"`
	MaxOutputTokens int             `json:"max_output_tokens,omitempty"`
}

// WenxinResponse 文心一言API响应结构
type WenxinResponse struct {
	ID               string      `json:"id"`
	Object           string      `json:"object"`
	Created          int64       `json:"created"`
	SentenceID       int         `json:"sentence_id"`
	IsEnd            bool        `json:"is_end"`
	IsTruncated      bool        `json:"is_truncated"`
	Result           string      `json:"result"`
	NeedClearHistory bool        `json:"need_clear_history"`
	BanRound         int         `json:"ban_round"`
	Usage            WenxinUsage `json:"usage"`
}

// WenxinUsage 文心一言使用统计
type WenxinUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// WenxinTokenResponse 获取访问令牌的响应
type WenxinTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// TextGeneration 实现文本生成方法
func (a *WenxinAdapter) TextGeneration(ctx context.Context, req *TextGenerationRequest) (*TextGenerationResponse, error) {
	// 确保有有效的访问令牌
	if err := a.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	// 构建文心一言请求
	wenxinReq := &WenxinRequest{
		Messages: []WenxinMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		Stop:            req.Stop,
		MaxOutputTokens: req.MaxTokens,
		Stream:          false,
	}

	// 发送请求
	wenxinResp, err := a.sendRequest(ctx, a.getModelEndpoint(req.Model), wenxinReq)
	if err != nil {
		return nil, err
	}

	return &TextGenerationResponse{
		ID:   wenxinResp.ID,
		Text: wenxinResp.Result,
		Usage: Usage{
			PromptTokens:     wenxinResp.Usage.PromptTokens,
			CompletionTokens: wenxinResp.Usage.CompletionTokens,
			TotalTokens:      wenxinResp.Usage.TotalTokens,
		},
		Model:        req.Model,
		FinishReason: a.getFinishReason(wenxinResp),
		CreatedAt:    time.Now(),
	}, nil
}

// ChatCompletion 实现对话完成方法
func (a *WenxinAdapter) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 确保有有效的访问令牌
	if err := a.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	// 转换消息格式
	var wenxinMessages []WenxinMessage
	var systemMessage string

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			wenxinMessages = append(wenxinMessages, WenxinMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// 构建文心一言请求
	wenxinReq := &WenxinRequest{
		Messages:        wenxinMessages,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		Stop:            req.Stop,
		MaxOutputTokens: req.MaxTokens,
		System:          systemMessage,
		Stream:          false,
	}

	// 发送请求
	wenxinResp, err := a.sendRequest(ctx, a.getModelEndpoint(req.Model), wenxinReq)
	if err != nil {
		return nil, err
	}

	return &ChatCompletionResponse{
		ID: wenxinResp.ID,
		Message: Message{
			Role:    "assistant",
			Content: wenxinResp.Result,
		},
		Usage: Usage{
			PromptTokens:     wenxinResp.Usage.PromptTokens,
			CompletionTokens: wenxinResp.Usage.CompletionTokens,
			TotalTokens:      wenxinResp.Usage.TotalTokens,
		},
		Model:        req.Model,
		FinishReason: a.getFinishReason(wenxinResp),
		CreatedAt:    time.Now(),
	}, nil
}

// TextGenerationStream 实现流式文本生成
func (a *WenxinAdapter) TextGenerationStream(ctx context.Context, req *TextGenerationRequest) (<-chan *TextGenerationResponse, error) {
	// 文心一言流式API实现较复杂，这里先返回错误
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "文心一言流式生成暂未实现",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// ImageGeneration 实现图像生成方法
func (a *WenxinAdapter) ImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	// 文心一言有图像生成功能，但API结构不同，这里先返回错误
	return nil, &AdapterError{
		Code:    ErrorTypeServiceUnavailable,
		Message: "文心一言图像生成暂未实现",
		Type:    ErrorTypeServiceUnavailable,
	}
}

// HealthCheck 实现健康检查
func (a *WenxinAdapter) HealthCheck(ctx context.Context) error {
	// 确保有有效的访问令牌
	if err := a.ensureValidToken(ctx); err != nil {
		return err
	}

	// 发送一个简单的请求来检查服务状态
	req := &WenxinRequest{
		Messages: []WenxinMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		MaxOutputTokens: 10,
	}

	_, err := a.sendRequest(ctx, a.getModelEndpoint("ernie-lite-8k"), req)
	return err
}

// ensureValidToken 确保有有效的访问令牌
func (a *WenxinAdapter) ensureValidToken(ctx context.Context) error {
	// 检查令牌是否过期
	if a.token != "" && time.Now().Before(a.tokenExp) {
		return nil
	}

	// 获取新的访问令牌
	return a.refreshToken(ctx)
}

// refreshToken 刷新访问令牌
func (a *WenxinAdapter) refreshToken(ctx context.Context) error {
	url := fmt.Sprintf("%s/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		a.baseURL, a.apiKey, a.secretKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return &AdapterError{
			Code:    ErrorTypeInvalidRequest,
			Message: fmt.Sprintf("创建令牌请求失败: %v", err),
			Type:    ErrorTypeInvalidRequest,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return &AdapterError{
			Code:    ErrorTypeTimeout,
			Message: fmt.Sprintf("获取访问令牌失败: %v", err),
			Type:    ErrorTypeTimeout,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &AdapterError{
			Code:    ErrorTypeUnknown,
			Message: fmt.Sprintf("读取令牌响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	var tokenResp WenxinTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return &AdapterError{
			Code:    ErrorTypeUnknown,
			Message: fmt.Sprintf("解析令牌响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	if tokenResp.Error != "" {
		return &AdapterError{
			Code:    ErrorTypeAuthentication,
			Message: fmt.Sprintf("获取访问令牌失败: %s", tokenResp.ErrorDescription),
			Type:    ErrorTypeAuthentication,
		}
	}

	a.token = tokenResp.AccessToken
	a.tokenExp = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// getModelEndpoint 获取模型对应的API端点
func (a *WenxinAdapter) getModelEndpoint(model string) string {
	endpoints := map[string]string{
		"ernie-4.0-8k":         "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro",
		"ernie-4.0-8k-preview": "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-4.0-8k-preview",
		"ernie-3.5-8k":         "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions",
		"ernie-3.5-8k-0205":    "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-3.5-8k-0205",
		"ernie-turbo-8k":       "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant",
		"ernie-speed-8k":       "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie_speed",
		"ernie-lite-8k":        "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-lite-8k",
		"ernie-tiny-8k":        "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-tiny-8k",
	}

	if endpoint, ok := endpoints[model]; ok {
		return endpoint
	}

	// 默认使用ernie-3.5-8k的端点
	return "/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions"
}

// getFinishReason 获取完成原因
func (a *WenxinAdapter) getFinishReason(resp *WenxinResponse) string {
	if resp.IsEnd {
		if resp.IsTruncated {
			return "length"
		}
		return "stop"
	}
	return "unknown"
}

// sendRequest 发送HTTP请求到文心一言API
func (a *WenxinAdapter) sendRequest(ctx context.Context, endpoint string, request interface{}) (*WenxinResponse, error) {
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
	url := fmt.Sprintf("%s%s?access_token=%s", a.baseURL, endpoint, a.token)
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
	var wenxinResp WenxinResponse
	if err := json.Unmarshal(responseBody, &wenxinResp); err != nil {
		return nil, &AdapterError{
			Code:    ErrorTypeUnknown,
			Message: fmt.Sprintf("解析响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	return &wenxinResp, nil
}

// handleErrorResponse 处理错误响应
func (a *WenxinAdapter) handleErrorResponse(statusCode int, body []byte) error {
	var errorType string
	var message string

	switch statusCode {
	case 400:
		errorType = ErrorTypeInvalidRequest
		message = "请求参数错误"
	case 401:
		errorType = ErrorTypeAuthentication
		message = "访问令牌无效"
	case 403:
		errorType = ErrorTypePermission
		message = "权限不足或配额不足"
	case 429:
		errorType = ErrorTypeRateLimit
		message = "请求频率超限"
	case 500, 502, 503:
		errorType = ErrorTypeServiceUnavailable
		message = "文心一言服务暂时不可用"
	default:
		errorType = ErrorTypeUnknown
		message = "未知错误"
	}

	// 尝试解析错误详情
	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if errorMsg, ok := errorResp["error_msg"].(string); ok {
			message = errorMsg
		} else if errorCode, ok := errorResp["error_code"].(float64); ok {
			message = fmt.Sprintf("错误代码: %.0f", errorCode)
		}
	}

	return &AdapterError{
		Code:    fmt.Sprintf("wenxin_%d", statusCode),
		Message: message,
		Type:    errorType,
	}
}
