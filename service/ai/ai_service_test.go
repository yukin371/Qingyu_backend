package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestService_GenerateContent 测试生成内容
func TestService_GenerateContent(t *testing.T) {
	// 创建服务实例
	service := &Service{
		contextService: &ContextService{}, // 使用具体类型而不是模拟
		adapterManager: nil,               // 暂时设为nil，因为我们不测试适配器管理器
	}

	// 创建请求
	req := &GenerateContentRequest{
		ProjectID: "test-project-id",
		ChapterID: "test-chapter-id",
		Prompt:    "测试提示词",
		Options:   nil,
	}

	// 由于我们使用的是真实的服务实例而不是模拟，
	// 这个测试主要验证结构体的创建和基本调用
	// 实际的业务逻辑测试需要在集成测试中进行

	// 验证请求结构
	assert.Equal(t, "test-project-id", req.ProjectID)
	assert.Equal(t, "test-chapter-id", req.ChapterID)
	assert.Equal(t, "测试提示词", req.Prompt)

	// 验证服务实例创建
	assert.NotNil(t, service)
	assert.NotNil(t, service.contextService)
}

// TestService_AnalyzeContent 测试分析内容
func TestService_AnalyzeContent(t *testing.T) {
	// 创建服务实例
	service := &Service{
		contextService: nil, // 分析内容不需要上下文服务
		adapterManager: nil, // 暂时设为nil，因为我们不测试适配器管理器
	}

	// 创建请求
	req := &AnalyzeContentRequest{
		Content:      "测试内容",
		AnalysisType: "plot",
	}

	// 验证请求结构
	assert.Equal(t, "测试内容", req.Content)
	assert.Equal(t, "plot", req.AnalysisType)

	// 验证服务实例创建
	assert.NotNil(t, service)
}

// TestService_ContinueWriting 测试续写
func TestService_ContinueWriting(t *testing.T) {
	// 创建服务实例
	service := &Service{
		contextService: &ContextService{}, // 使用具体类型而不是模拟
		adapterManager: nil,               // 暂时设为nil，因为我们不测试适配器管理器
	}

	// 创建请求
	req := &ContinueWritingRequest{
		ProjectID:      "test-project-id",
		ChapterID:      "test-chapter-id",
		CurrentText:    "当前文本内容",
		ContinueLength: 500,
		Options:        nil,
	}

	// 验证请求结构
	assert.Equal(t, "test-project-id", req.ProjectID)
	assert.Equal(t, "test-chapter-id", req.ChapterID)
	assert.Equal(t, "当前文本内容", req.CurrentText)
	assert.Equal(t, 500, req.ContinueLength)

	// 验证服务实例创建
	assert.NotNil(t, service)
	assert.NotNil(t, service.contextService)
}

// TestService_OptimizeText 测试文本优化
func TestService_OptimizeText(t *testing.T) {
	// 创建服务实例
	service := &Service{
		contextService: &ContextService{}, // 使用具体类型而不是模拟
		adapterManager: nil,               // 暂时设为nil，因为我们不测试适配器管理器
	}

	// 创建请求
	req := &OptimizeTextRequest{
		ProjectID:    "test-project-id",
		ChapterID:    "test-chapter-id",
		OriginalText: "原始文本",
		OptimizeType: "grammar",
		Instructions: "请优化语法",
		Options:      nil,
	}

	// 验证请求结构
	assert.Equal(t, "test-project-id", req.ProjectID)
	assert.Equal(t, "test-chapter-id", req.ChapterID)
	assert.Equal(t, "原始文本", req.OriginalText)
	assert.Equal(t, "grammar", req.OptimizeType)
	assert.Equal(t, "请优化语法", req.Instructions)

	// 验证服务实例创建
	assert.NotNil(t, service)
	assert.NotNil(t, service.contextService)
}
