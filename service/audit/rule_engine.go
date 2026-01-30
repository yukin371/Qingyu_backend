package audit

import (
	"regexp"
	"strings"

	"Qingyu_backend/models/audit"
)

// RuleEngine 规则引擎
type RuleEngine struct {
	rules []Rule
}

// Rule 审核规则接口
type Rule interface {
	Check(content string) []audit.ViolationDetail
	GetName() string
	GetPriority() int
	IsEnabled() bool
}

// NewRuleEngine 创建规则引擎
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: make([]Rule, 0),
	}
}

// AddRule 添加规则
func (e *RuleEngine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

// Check 检查内容
func (e *RuleEngine) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	for _, rule := range e.rules {
		if !rule.IsEnabled() {
			continue
		}

		ruleViolations := rule.Check(content)
		violations = append(violations, ruleViolations...)
	}

	return violations
}

// GetRules 获取所有规则
func (e *RuleEngine) GetRules() []Rule {
	return e.rules
}

// RemoveRule 移除规则
func (e *RuleEngine) RemoveRule(name string) {
	newRules := make([]Rule, 0)
	for _, rule := range e.rules {
		if rule.GetName() != name {
			newRules = append(newRules, rule)
		}
	}
	e.rules = newRules
}

// 具体规则实现

// RegexRule 正则表达式规则
type RegexRule struct {
	Name        string
	Description string
	Pattern     *regexp.Regexp
	Category    string
	Level       int
	Enabled     bool
	Priority    int
}

func (r *RegexRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	matches := r.Pattern.FindAllStringIndex(content, -1)
	for _, match := range matches {
		violation := audit.ViolationDetail{
			Type:        "regex_match",
			Category:    r.Category,
			Level:       r.Level,
			Description: r.Description,
			Position:    match[0],
			Context:     extractContext([]rune(content), match[0], match[1], 20),
			Keywords:    []string{content[match[0]:match[1]]},
		}
		violations = append(violations, violation)
	}

	return violations
}

func (r *RegexRule) GetName() string  { return r.Name }
func (r *RegexRule) GetPriority() int { return r.Priority }
func (r *RegexRule) IsEnabled() bool  { return r.Enabled }

// PhoneNumberRule 手机号检测规则
type PhoneNumberRule struct {
	Enabled  bool
	Priority int
}

func NewPhoneNumberRule() *PhoneNumberRule {
	return &PhoneNumberRule{
		Enabled:  true,
		Priority: 1,
	}
}

func (r *PhoneNumberRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	// 匹配11位手机号
	pattern := regexp.MustCompile(`1[3-9]\d{9}`)
	matches := pattern.FindAllStringIndex(content, -1)

	for _, match := range matches {
		phone := content[match[0]:match[1]]
		violation := audit.ViolationDetail{
			Type:        "phone_number",
			Category:    audit.CategoryAd,
			Level:       audit.LevelMedium,
			Description: "检测到手机号码，可能涉及广告推广",
			Position:    match[0],
			Context:     extractContext([]rune(content), match[0], match[1], 20),
			Keywords:    []string{phone},
		}
		violations = append(violations, violation)
	}

	return violations
}

func (r *PhoneNumberRule) GetName() string  { return "PhoneNumberRule" }
func (r *PhoneNumberRule) GetPriority() int { return r.Priority }
func (r *PhoneNumberRule) IsEnabled() bool  { return r.Enabled }

// URLRule URL检测规则
type URLRule struct {
	Enabled  bool
	Priority int
}

func NewURLRule() *URLRule {
	return &URLRule{
		Enabled:  true,
		Priority: 1,
	}
}

func (r *URLRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	// 匹配HTTP/HTTPS链接
	pattern := regexp.MustCompile(`https?://[^\s]+`)
	matches := pattern.FindAllStringIndex(content, -1)

	for _, match := range matches {
		url := content[match[0]:match[1]]
		violation := audit.ViolationDetail{
			Type:        "url",
			Category:    audit.CategoryAd,
			Level:       audit.LevelMedium,
			Description: "检测到外部链接，可能涉及广告推广",
			Position:    match[0],
			Context:     extractContext([]rune(content), match[0], match[1], 20),
			Keywords:    []string{url},
		}
		violations = append(violations, violation)
	}

	return violations
}

func (r *URLRule) GetName() string  { return "URLRule" }
func (r *URLRule) GetPriority() int { return r.Priority }
func (r *URLRule) IsEnabled() bool  { return r.Enabled }

// WeChatRule 微信号检测规则
type WeChatRule struct {
	Enabled  bool
	Priority int
}

func NewWeChatRule() *WeChatRule {
	return &WeChatRule{
		Enabled:  true,
		Priority: 1,
	}
}

func (r *WeChatRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	// 检测微信相关关键词
	keywords := []string{"微信", "微信号", "加微信", "wx", "weixin", "vx", "V信"}
	lowerContent := strings.ToLower(content)

	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		index := strings.Index(lowerContent, lowerKeyword)
		if index >= 0 {
			violation := audit.ViolationDetail{
				Type:        "wechat",
				Category:    audit.CategoryAd,
				Level:       audit.LevelLow,
				Description: "检测到微信相关信息，可能涉及广告推广",
				Position:    index,
				Context:     extractContext([]rune(content), index, index+len(keyword), 20),
				Keywords:    []string{keyword},
			}
			violations = append(violations, violation)
			break // 只记录第一个
		}
	}

	return violations
}

func (r *WeChatRule) GetName() string  { return "WeChatRule" }
func (r *WeChatRule) GetPriority() int { return r.Priority }
func (r *WeChatRule) IsEnabled() bool  { return r.Enabled }

// QQRule QQ号检测规则
type QQRule struct {
	Enabled  bool
	Priority int
}

func NewQQRule() *QQRule {
	return &QQRule{
		Enabled:  true,
		Priority: 1,
	}
}

func (r *QQRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	// 检测QQ相关关键词
	keywords := []string{"qq", "QQ", "扣扣", "加q", "加Q"}
	lowerContent := strings.ToLower(content)

	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		index := strings.Index(lowerContent, lowerKeyword)
		if index >= 0 {
			violation := audit.ViolationDetail{
				Type:        "qq",
				Category:    audit.CategoryAd,
				Level:       audit.LevelLow,
				Description: "检测到QQ相关信息，可能涉及广告推广",
				Position:    index,
				Context:     extractContext([]rune(content), index, index+len(keyword), 20),
				Keywords:    []string{keyword},
			}
			violations = append(violations, violation)
			break // 只记录第一个
		}
	}

	return violations
}

func (r *QQRule) GetName() string  { return "QQRule" }
func (r *QQRule) GetPriority() int { return r.Priority }
func (r *QQRule) IsEnabled() bool  { return r.Enabled }

// ExcessiveRepetitionRule 过度重复检测规则
type ExcessiveRepetitionRule struct {
	Enabled       bool
	Priority      int
	MaxRepetition int // 最大重复次数
}

func NewExcessiveRepetitionRule() *ExcessiveRepetitionRule {
	return &ExcessiveRepetitionRule{
		Enabled:       true,
		Priority:      2,
		MaxRepetition: 10, // 超过10次重复视为异常
	}
}

func (r *ExcessiveRepetitionRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	runes := []rune(content)
	if len(runes) < 2 {
		return violations
	}

	// 检测连续重复字符
	count := 1
	lastRune := runes[0]
	startPos := 0

	for i := 1; i < len(runes); i++ {
		if runes[i] == lastRune {
			count++
			if count > r.MaxRepetition {
				violation := audit.ViolationDetail{
					Type:        "excessive_repetition",
					Category:    audit.CategoryOther,
					Level:       audit.LevelLow,
					Description: "检测到过度重复字符，可能为灌水内容",
					Position:    startPos,
					Context:     extractContext([]rune(content), startPos, i, 20),
					Keywords:    []string{string(lastRune)},
				}
				violations = append(violations, violation)
				break // 只记录一次
			}
		} else {
			count = 1
			lastRune = runes[i]
			startPos = i
		}
	}

	return violations
}

func (r *ExcessiveRepetitionRule) GetName() string  { return "ExcessiveRepetitionRule" }
func (r *ExcessiveRepetitionRule) GetPriority() int { return r.Priority }
func (r *ExcessiveRepetitionRule) IsEnabled() bool  { return r.Enabled }

// ContentLengthRule 内容长度规则
type ContentLengthRule struct {
	Enabled   bool
	Priority  int
	MinLength int
	MaxLength int
}

func NewContentLengthRule() *ContentLengthRule {
	return &ContentLengthRule{
		Enabled:   true,
		Priority:  3,
		MinLength: 10,     // 最少10个字
		MaxLength: 100000, // 最多10万字
	}
}

func (r *ContentLengthRule) Check(content string) []audit.ViolationDetail {
	violations := make([]audit.ViolationDetail, 0)

	length := len([]rune(content))

	if length < r.MinLength {
		violation := audit.ViolationDetail{
			Type:        "content_too_short",
			Category:    audit.CategoryOther,
			Level:       audit.LevelLow,
			Description: "内容过短，可能为无效内容",
			Position:    0,
			Context:     content,
			Keywords:    []string{},
		}
		violations = append(violations, violation)
	}

	if length > r.MaxLength {
		violation := audit.ViolationDetail{
			Type:        "content_too_long",
			Category:    audit.CategoryOther,
			Level:       audit.LevelLow,
			Description: "内容过长，建议分段",
			Position:    0,
			Context:     content[:100] + "...",
			Keywords:    []string{},
		}
		violations = append(violations, violation)
	}

	return violations
}

func (r *ContentLengthRule) GetName() string  { return "ContentLengthRule" }
func (r *ContentLengthRule) GetPriority() int { return r.Priority }
func (r *ContentLengthRule) IsEnabled() bool  { return r.Enabled }

// LoadDefaultRules 加载默认规则
func LoadDefaultRules(engine *RuleEngine) {
	// 加载所有默认规则
	engine.AddRule(NewPhoneNumberRule())
	engine.AddRule(NewURLRule())
	engine.AddRule(NewWeChatRule())
	engine.AddRule(NewQQRule())
	engine.AddRule(NewExcessiveRepetitionRule())
	// 内容长度规则可选
	// engine.AddRule(NewContentLengthRule())
}

// 辅助函数

// extractContext函数已在dfa.go中定义，此处不再重复声明
