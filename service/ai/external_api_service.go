package ai

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
	"Qingyu_backend/models/ai"
	"Qingyu_backend/service/ai/adapter"
)

// ExternalAPIService 外部AI API服务（已弃用，保留用于向后兼容）
// 推荐使用 adapter.AdapterManager 进行AI服务调用
type ExternalAPIService struct {
	httpClient     *http.Client
	config         *config.AIConfig
	adapterManager *adapter.AdapterManager // 新增适配器管理器
}

// NewExternalAPIService 创建外部API服务
func NewExternalAPIService(cfg *config.AIConfig) *ExternalAPIService {
	return &ExternalAPIService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 使用固定超时时间
		},
		config:         cfg,
		adapterManager: nil, // 暂时设为nil，由ai_service.go中统一管理
	}
}

// GenerateRequest AI生成请求
type GenerateRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GenerateResponse AI生成响应
type GenerateResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GenerateContent 生成内容（向后兼容方法，内部使用适配器管理器）
func (s *ExternalAPIService) GenerateContent(ctx context.Context, aiContext *ai.AIContext, prompt string, options *ai.GenerateOptions) (*ai.GenerateResult, error) {
	// 构建系统提示词
	systemPrompt := s.buildSystemPrompt(aiContext)

	// 构建用户提示词
	userPrompt := s.buildUserPrompt(aiContext, prompt)

	// 组合完整提示词
	fullPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, userPrompt)

	// 使用适配器管理器生成内容
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      fullPrompt,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Model:       options.Model,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("生成内容失败: %w", err)
	}

	// 转换为旧格式响应
	return &ai.GenerateResult{
		Content:      result.Text,
		TokensUsed:   result.Usage.TotalTokens,
		Model:        result.Model,
		FinishReason: result.FinishReason,
	}, nil
}

// buildSystemPrompt 构建系统提示词
func (s *ExternalAPIService) buildSystemPrompt(aiContext *ai.AIContext) string {
	var prompt strings.Builder

	prompt.WriteString("你是一个专业的小说写作助手。请根据以下项目信息和上下文来协助创作：\n\n")

	// 当前章节信息
	if aiContext.CurrentChapter != nil {
		prompt.WriteString(fmt.Sprintf("当前章节：%s\n", aiContext.CurrentChapter.Title))
		if aiContext.CurrentChapter.Summary != "" {
			prompt.WriteString(fmt.Sprintf("章节摘要：%s\n", aiContext.CurrentChapter.Summary))
		}
	}

	// 角色信息
	if len(aiContext.ActiveCharacters) > 0 {
		prompt.WriteString("\n活跃角色：\n")
		for _, char := range aiContext.ActiveCharacters {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", char.Name, char.Summary))
			if len(char.Traits) > 0 {
				prompt.WriteString(fmt.Sprintf("  性格特点：%s\n", strings.Join(char.Traits, "、")))
			}
		}
	}

	// 地点信息
	if len(aiContext.CurrentLocations) > 0 {
		prompt.WriteString("\n当前场景：\n")
		for _, loc := range aiContext.CurrentLocations {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", loc.Name, loc.Description))
		}
	}

	// 情节线索
	if len(aiContext.PlotThreads) > 0 {
		prompt.WriteString("\n活跃情节线索：\n")
		for _, thread := range aiContext.PlotThreads {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", thread.Name, thread.Description))
		}
	}

	// 世界观设定
	if aiContext.WorldSettings != nil {
		prompt.WriteString(fmt.Sprintf("\n世界观设定：%s\n", aiContext.WorldSettings.Description))
		if len(aiContext.WorldSettings.Rules) > 0 {
			prompt.WriteString("世界规则：\n")
			for _, rule := range aiContext.WorldSettings.Rules {
				prompt.WriteString(fmt.Sprintf("- %s\n", rule))
			}
		}
	}

	prompt.WriteString("\n请基于以上信息协助创作，保持角色一致性和情节连贯性。")

	return prompt.String()
}

// buildUserPrompt 构建用户提示词
func (s *ExternalAPIService) buildUserPrompt(aiContext *ai.AIContext, prompt string) string {
	var userPrompt strings.Builder

	// 当前章节内容
	if aiContext.CurrentChapter != nil {
		if aiContext.CurrentChapter.Content != "" {
			userPrompt.WriteString(fmt.Sprintf("当前章节内容：\n%s\n\n", aiContext.CurrentChapter.Content))
		}

		// 章节关键点
		if len(aiContext.CurrentChapter.KeyPoints) > 0 {
			userPrompt.WriteString("章节关键点：\n")
			for _, point := range aiContext.CurrentChapter.KeyPoints {
				userPrompt.WriteString(fmt.Sprintf("- %s\n", point))
			}
			userPrompt.WriteString("\n")
		}

		// 写作提示
		if aiContext.CurrentChapter.WritingHints != "" {
			userPrompt.WriteString(fmt.Sprintf("写作提示：%s\n\n", aiContext.CurrentChapter.WritingHints))
		}
	}

	// 前序章节摘要
	if len(aiContext.PreviousChapters) > 0 {
		userPrompt.WriteString("前序章节摘要：\n")
		for _, chapter := range aiContext.PreviousChapters {
			userPrompt.WriteString(fmt.Sprintf("- %s: %s\n", chapter.Title, chapter.Summary))
		}
		userPrompt.WriteString("\n")
	}

	// 用户具体请求
	userPrompt.WriteString(fmt.Sprintf("用户请求：%s", prompt))

	return userPrompt.String()
}

// sendRequest 发送HTTP请求（已弃用，保留用于向后兼容）
func (s *ExternalAPIService) sendRequest(ctx context.Context, request *GenerateRequest) (*GenerateResponse, error) {
	// 使用简化的配置
	apiURL := s.config.BaseURL + "/chat/completions"

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBody))
	}

	// 解析响应
	var response GenerateResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &response, nil
}

// AnalyzeContent 分析内容（向后兼容方法，内部使用适配器管理器）
func (s *ExternalAPIService) AnalyzeContent(ctx context.Context, content string, analysisType string) (*ai.AnalysisResult, error) {
	var prompt string
	switch analysisType {
	case "plot":
		prompt = "请分析以下文本的情节结构，包括主要情节点、冲突、转折等："
	case "character":
		prompt = "请分析以下文本中的角色表现，包括性格特点、行为动机、对话特色等："
	case "style":
		prompt = "请分析以下文本的写作风格，包括语言特色、叙述方式、文体特点等："
	default:
		prompt = "请对以下文本进行综合分析："
	}

	// 构建完整提示词
	fullPrompt := fmt.Sprintf("你是一个专业的文学分析师，请对提供的文本进行深入分析。\n\n%s\n\n文本内容：\n%s", prompt, content)

	// 使用适配器管理器生成分析
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      fullPrompt,
		Temperature: 0.3,
		MaxTokens:   1000,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("分析内容失败: %w", err)
	}

	// 转换为旧格式响应
	return &ai.AnalysisResult{
		Type:       analysisType,
		Analysis:   result.Text,
		TokensUsed: result.Usage.TotalTokens,
		Model:      result.Model,
	}, nil
}
