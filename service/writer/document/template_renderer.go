package document

import (
	"fmt"
	"regexp"
	"strings"
)

// TemplateRenderer 模板渲染器接口
type TemplateRenderer interface {
	// Render 渲染模板内容
	// content: 模板内容,包含 {{var:variableName}} 格式的变量
	// vars: 变量键值对
	// 返回: 渲染后的内容和可能的错误
	Render(content string, vars map[string]string) (string, error)

	// Validate 验证模板内容
	// 检查 {{ }} 是否配对,变量名是否符合规范
	Validate(content string) error
}

// SimpleTemplateRenderer 简单模板渲染器
// 支持简单的变量替换,不包含复杂的模板语法(if/for/表达式)
type SimpleTemplateRenderer struct {
	// variableRegex 编译后的变量匹配正则表达式
	variableRegex *regexp.Regexp
}

// NewSimpleTemplateRenderer 创建简单模板渲染器实例
func NewSimpleTemplateRenderer() *SimpleTemplateRenderer {
	// 匹配 {{var:variableName}} 格式的变量
	// 变量名规范: 以字母或下划线开头,只能包含字母、数字和下划线
	re := regexp.MustCompile(`\{\{var:([a-zA-Z_][a-zA-Z0-9_]*)\}\}`)

	return &SimpleTemplateRenderer{
		variableRegex: re,
	}
}

// Render 渲染模板内容,将 {{var:variableName}} 替换为对应的变量值
// 如果变量不存在,保持原样 {{var:variableName}}
func (r *SimpleTemplateRenderer) Render(content string, vars map[string]string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("模板内容不能为空")
	}

	// 使用正则表达式查找并替换所有变量
	result := r.variableRegex.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名
		submatches := r.variableRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			// 正常情况下不会到这里,因为正则表达式保证了格式正确
			return match // 保持原样
		}

		varName := submatches[1]

		// 查找变量值
		value, ok := vars[varName]
		if !ok {
			// 变量不存在,保持原样
			return match
		}

		return value
	})

	return result, nil
}

// Validate 验证模板内容
// 检查:
// 1. {{ 和 }} 是否配对
// 2. 变量名是否符合规范
// 3. 是否有无效的变量引用
func (r *SimpleTemplateRenderer) Validate(content string) error {
	if content == "" {
		return fmt.Errorf("模板内容不能为空")
	}

	// 检查 {{ 和 }} 是否配对
	openBraces := strings.Count(content, "{{")
	closeBraces := strings.Count(content, "}}")

	if openBraces != closeBraces {
		return fmt.Errorf("模板语法错误: 未闭合的标签({{ }}不匹配,找到 %d 个 {{ 和 %d 个 }})", openBraces, closeBraces)
	}

	// 查找所有 {{xxx}} 格式的内容,检查是否为有效的变量引用
	allBracesRegex := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := allBracesRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// 检查是否为 var: 开头
		fullContent := match[1]
		if !strings.HasPrefix(fullContent, "var:") {
			return fmt.Errorf("无效的变量引用: {{%s}},必须使用 {{var:variableName}} 格式", fullContent)
		}

		// 提取变量名并验证
		varName := strings.TrimPrefix(fullContent, "var:")
		if varName == "" {
			return fmt.Errorf("变量名为空: {{%s}}", fullContent)
		}

		// 验证变量名是否符合规范
		if !isValidVariableNameForRenderer(varName) {
			return fmt.Errorf("变量名 '%s' 不符合规范: 必须以字母或下划线开头,只能包含字母、数字和下划线", varName)
		}
	}

	return nil
}

// isValidVariableNameForRenderer 验证变量名是否有效
// 变量名规范: 以字母或下划线开头,只能包含字母、数字和下划线
func isValidVariableNameForRenderer(name string) bool {
	if name == "" {
		return false
	}

	// 变量名必须以字母或下划线开头
	firstChar := name[0]
	if !((firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		firstChar == '_') {
		return false
	}

	// 变量名只能包含字母、数字和下划线
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_') {
			return false
		}
	}

	return true
}
