package document

import (
	"testing"
)

// TestSimpleTemplateRenderer_Render 测试Render方法
func TestSimpleTemplateRenderer_Render(t *testing.T) {
	renderer := NewSimpleTemplateRenderer()

	tests := []struct {
		name        string
		content     string
		vars        map[string]string
		expected    string
		expectError bool
	}{
		{
			name:        "空内容",
			content:     "",
			vars:        map[string]string{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "无变量",
			content:     "这是一段纯文本",
			vars:        map[string]string{},
			expected:    "这是一段纯文本",
			expectError: false,
		},
		{
			name:     "单个变量替换",
			content:  "你好,{{var:name}}",
			vars:     map[string]string{"name": "张三"},
			expected: "你好,张三",
			expectError: false,
		},
		{
			name:     "多个变量替换",
			content:  "{{var:title}} - {{var:author}}",
			vars:     map[string]string{"title": "测试标题", "author": "测试作者"},
			expected: "测试标题 - 测试作者",
			expectError: false,
		},
		{
			name:     "变量不存在时保持原样",
			content:  "你好,{{var:name}}",
			vars:     map[string]string{},
			expected: "你好,{{var:name}}",
			expectError: false,
		},
		{
			name:     "部分变量存在",
			content:  "{{var:first}}和{{var:second}}",
			vars:     map[string]string{"first": "第一"},
			expected: "第一和{{var:second}}",
			expectError: false,
		},
		{
			name:     "变量名以下划线开头",
			content:  "{{var:_private}}",
			vars:     map[string]string{"_private": "私有值"},
			expected: "私有值",
			expectError: false,
		},
		{
			name:     "变量名包含数字",
			content:  "{{var:user_name1}}",
			vars:     map[string]string{"user_name1": "用户1"},
			expected: "用户1",
			expectError: false,
		},
		{
			name:     "多个相同变量",
			content:  "{{var:name}}和{{var:name}}",
			vars:     map[string]string{"name": "重复"},
			expected: "重复和重复",
			expectError: false,
		},
		{
			name:     "混合文本和变量",
			content:  "章节:{{var:chapter}}\n字数:{{var:wordcount}}",
			vars:     map[string]string{"chapter": "第一章", "wordcount": "1000"},
			expected: "章节:第一章\n字数:1000",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderer.Render(tt.content, tt.vars)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望返回错误,但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外的错误: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("结果不匹配\n期望: %q\n实际: %q", tt.expected, result)
			}
		})
	}
}

// TestSimpleTemplateRenderer_Validate 测试Validate方法
func TestSimpleTemplateRenderer_Validate(t *testing.T) {
	renderer := NewSimpleTemplateRenderer()

	tests := []struct {
		name        string
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "空内容",
			content:     "",
			expectError: true,
			errorMsg:    "模板内容不能为空",
		},
		{
			name:        "无变量",
			content:     "这是一段纯文本",
			expectError: false,
		},
		{
			name:        "有效的变量引用",
			content:     "你好,{{var:name}}",
			expectError: false,
		},
		{
			name:        "多个有效变量",
			content:     "{{var:title}} - {{var:author}}",
			expectError: false,
		},
		{
			name:        "变量名以下划线开头",
			content:     "{{var:_private}}",
			expectError: false,
		},
		{
			name:        "变量名包含数字",
			content:     "{{var:user_name1}}",
			expectError: false,
		},
		{
			name:        "未闭合的标签-缺少}}",
			content:     "你好,{{var:name}",
			expectError: true,
			errorMsg:    "未闭合的标签",
		},
		{
			name:        "未闭合的标签-缺少{{",
			content:     "你好,var:name}}",
			expectError: true,
			errorMsg:    "未闭合的标签",
		},
		{
			name:        "无效的变量格式-缺少var:",
			content:     "你好,{{name}}",
			expectError: true,
			errorMsg:    "必须使用 {{var:variableName}} 格式",
		},
		{
			name:        "无效的变量格式-使用点号",
			content:     "你好,{{var.name}}",
			expectError: true,
			errorMsg:    "必须使用 {{var:variableName}} 格式",
		},
		{
			name:        "变量名为空",
			content:     "你好,{{var:}}",
			expectError: true,
			errorMsg:    "变量名为空",
		},
		{
			name:        "变量名以数字开头",
			content:     "你好,{{var:1name}}",
			expectError: true,
			errorMsg:    "不符合规范",
		},
		{
			name:        "变量名包含特殊字符",
			content:     "你好,{{var:user-name}}",
			expectError: true,
			errorMsg:    "不符合规范",
		},
		{
			name:        "配对正确的多个标签",
			content:     "{{var:first}}和{{var:second}}",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := renderer.Validate(tt.content)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望返回错误,但没有错误")
					return
				}

				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("错误信息不匹配\n期望包含: %q\n实际: %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("意外的错误: %v", err)
				}
			}
		})
	}
}

// TestIsValidVariableNameForRenderer 测试变量名验证函数
func TestIsValidVariableNameForRenderer(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		expected bool
	}{
		{"空字符串", "", false},
		{"字母开头", "name", true},
		{"大写字母开头", "Name", true},
		{"下划线开头", "_name", true},
		{"数字开头", "1name", false},
		{"包含数字", "name1", true},
		{"包含下划线", "user_name", true},
		{"包含多个下划线", "user_name_1", true},
		{"包含特殊字符-短横线", "user-name", false},
		{"包含特殊字符-点号", "user.name", false},
		{"包含特殊字符-空格", "user name", false},
		{"纯数字", "123", false},
		{"纯下划线", "_", true},
		{"驼峰命名", "userName", true},
		{"大写下划线", "USER_NAME", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidVariableNameForRenderer(tt.varName)
			if result != tt.expected {
				t.Errorf("变量名 %q 验证失败\n期望: %v\n实际: %v", tt.varName, tt.expected, result)
			}
		})
	}
}

// containsString 检查字符串是否包含子字符串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOf(s, substr) >= 0)
}

// indexOf 查找子字符串的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
