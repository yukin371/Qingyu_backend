package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sanitizer 敏感信息脱敏器
type Sanitizer struct {
	// 是否启用脱敏
	Enabled bool
	// 自定义脱敏规则
	CustomRules map[string]func(string) string
}

// DefaultSanitizer 默认脱敏器
var DefaultSanitizer = &Sanitizer{
	Enabled:     true,
	CustomRules: make(map[string]func(string) string),
}

// MaskString 脱敏字符串
func MaskString(s string, visibleChars int) string {
	if s == "" {
		return ""
	}

	runeCount := utf8.RuneCountInString(s)
	if runeCount <= visibleChars {
		// 如果字符串长度小于等于可见字符数，则显示前一半
		show := runeCount / 2
		if show == 0 {
			show = 1
		}
		return SubString(s, 0, show) + strings.Repeat("*", runeCount-show)
	}

	// 显示前visibleChars个字符
	return SubString(s, 0, visibleChars) + strings.Repeat("*", runeCount-visibleChars)
}

// MaskEmail 脱敏邮箱
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	// 用户名显示前2个字符
	if len(username) <= 2 {
		username = strings.Repeat("*", len(username))
	} else {
		username = SubString(username, 0, 2) + strings.Repeat("*", len(username)-2)
	}

	return username + "@" + domain
}

// MaskPhone 脱敏手机号
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}

	// 去除非数字字符
	re := regexp.MustCompile(`[^0-9]`)
	phone = re.ReplaceAllString(phone, "")

	if len(phone) != 11 {
		return phone
	}

	// 显示前3位和后4位
	return phone[:3] + "****" + phone[7:]
}

// MaskIDCard 脱敏身份证号
func MaskIDCard(idCard string) string {
	if idCard == "" {
		return ""
	}

	length := len(idCard)
	if length < 8 {
		return strings.Repeat("*", length)
	}

	// 显示前3位和后4位
	return idCard[:3] + strings.Repeat("*", length-7) + idCard[length-4:]
}

// MaskBankCard 脱敏银行卡号
func MaskBankCard(cardNo string) string {
	if cardNo == "" {
		return ""
	}

	// 去除非数字字符
	re := regexp.MustCompile(`[^0-9]`)
	cardNo = re.ReplaceAllString(cardNo, "")

	length := len(cardNo)
	if length < 8 {
		return strings.Repeat("*", length)
	}

	// 显示前4位和后4位
	return cardNo[:4] + strings.Repeat("*", length-8) + cardNo[length-4:]
}

// MaskObjectID 脱敏MongoDB ObjectID
func MaskObjectID(id primitive.ObjectID) string {
	hex := id.Hex()
	if len(hex) <= 8 {
		return hex
	}
	// 显示前4位和后4位
	return hex[:4] + "****" + hex[len(hex)-4:]
}

// MaskToken 脱敏Token
func MaskToken(token string) string {
	if token == "" {
		return ""
	}

	length := len(token)
	if length <= 16 {
		return strings.Repeat("*", length)
	}

	// 显示前8位和后8位
	return token[:8] + strings.Repeat("*", length-16) + token[length-8:]
}

// MaskPassword 脱敏密码（完全隐藏）
func MaskPassword(password string) string {
	if password == "" {
		return ""
	}
	return strings.Repeat("*", len(password))
}

// MaskAddress 脱敏地址
func MaskAddress(address string) string {
	if address == "" {
		return ""
	}

	// 如果地址太短，返回前半部分
	words := strings.Fields(address)
	if len(words) <= 2 {
		return SubString(address, 0, len(address)/2) + "***"
	}

	// 保留前2个词和最后1个词
	result := words[0]
	if len(words) > 1 {
		result += " " + words[1]
	}
	if len(words) > 2 {
		result += " ****"
	}
	if len(words) > 3 {
		result += " " + words[len(words)-1]
	}

	return result
}

// MaskJSON 脱敏JSON中的敏感字段
func MaskJSON(data []byte, sensitiveFields ...string) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	// 脱敏敏感字段
	for _, field := range sensitiveFields {
		if value, exists := obj[field]; exists {
			if str, ok := value.(string); ok {
				obj[field] = maskFieldByType(field, str)
			}
		}
	}

	return json.Marshal(obj)
}

// maskFieldByType 根据字段类型脱敏
func maskFieldByType(field, value string) string {
	switch strings.ToLower(field) {
	case "email", "emailaddress", "mail":
		return MaskEmail(value)
	case "phone", "mobile", "tel", "telephone":
		return MaskPhone(value)
	case "password", "passwd", "pwd":
		return MaskPassword(value)
	case "idcard", "id_card", "idnumber":
		return MaskIDCard(value)
	case "bankcard", "bank_card", "cardno", "card_no":
		return MaskBankCard(value)
	case "token", "accesstoken", "refreshtoken":
		return MaskToken(value)
	case "address":
		return MaskAddress(value)
	default:
		return MaskString(value, 4)
	}
}

// SanitizeUser 脱敏用户信息
type UserSanitizer struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	IDCard   string `json:"id_card,omitempty"`
	BankCard string `json:"bank_card,omitempty"`
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

// Sanitize 脱敏
func (u *UserSanitizer) Sanitize() *UserSanitizer {
	if u.Email != "" {
		u.Email = MaskEmail(u.Email)
	}
	if u.Phone != "" {
		u.Phone = MaskPhone(u.Phone)
	}
	if u.IDCard != "" {
		u.IDCard = MaskIDCard(u.IDCard)
	}
	if u.BankCard != "" {
		u.BankCard = MaskBankCard(u.BankCard)
	}
	if u.Token != "" {
		u.Token = MaskToken(u.Token)
	}
	if u.Password != "" {
		u.Password = MaskPassword(u.Password)
	}
	return u
}

// SanitizeResponse 脱敏响应数据
func SanitizeResponse(data interface{}, sensitiveFields ...string) interface{} {
	// 将data转换为JSON再转换回来，应用脱敏规则
	jsonData, err := json.Marshal(data)
	if err != nil {
		return data
	}

	// 如果有敏感字段，进行脱敏
	if len(sensitiveFields) > 0 {
		masked, err := MaskJSON(jsonData, sensitiveFields...)
		if err != nil {
			return data
		}
		jsonData = masked
	}

	// 解析回interface{}
	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return data
	}

	return result
}

// SubString 获取子字符串（按字符而非字节）
func SubString(s string, start, length int) string {
	runes := []rune(s)
	runeCount := len(runes)

	if start >= runeCount {
		return ""
	}

	end := start + length
	if end > runeCount {
		end = runeCount
	}

	return string(runes[start:end])
}

// MaskWithPattern 使用正则表达式模式脱敏
func MaskWithPattern(s string, pattern string, replacement string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(s, replacement)
}

// MaskJSONFields 脱敏JSON对象中的敏感字段
func MaskJSONFields(obj interface{}, fields []string) interface{} {
	// 使用JSON序列化/反序列化
	data, err := json.Marshal(obj)
	if err != nil {
		return obj
	}

	var mapData map[string]interface{}
	if err := json.Unmarshal(data, &mapData); err != nil {
		// 如果不是对象，直接返回
		return obj
	}

	// 创建敏感字段的map用于快速查找
	sensitiveMap := make(map[string]bool)
	for _, field := range fields {
		sensitiveMap[field] = true
	}

	// 递归脱敏
	maskMapRecursive(mapData, sensitiveMap)

	return mapData
}

// maskMapRecursive 递归脱敏map
func maskMapRecursive(m map[string]interface{}, sensitiveFields map[string]bool) {
	for key, value := range m {
		if sensitiveFields[key] {
			// 脱敏敏感字段
			if str, ok := value.(string); ok {
				m[key] = maskFieldByType(key, str)
			}
		} else {
			// 递归处理嵌套对象
			switch v := value.(type) {
			case map[string]interface{}:
				maskMapRecursive(v, sensitiveFields)
			case []interface{}:
				// 处理数组
				for i, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						maskMapRecursive(itemMap, sensitiveFields)
						v[i] = itemMap
					}
				}
			}
		}
	}
}

// SensitiveString 敏感字符串类型
type SensitiveString struct {
	value string
}

// NewSensitiveString 创建敏感字符串
func NewSensitiveString(s string) *SensitiveString {
	return &SensitiveString{value: s}
}

// String 实现Stringer接口
func (s *SensitiveString) String() string {
	return MaskString(s.value, 4)
}

// Value 获取原始值
func (s *SensitiveString) Value() string {
	return s.value
}

// MarshalJSON 实现JSON序列化
func (s *SensitiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(MaskString(s.value, 4))
}

// UnmarshalJSON 实现JSON反序列化
func (s *SensitiveString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	s.value = str
	return nil
}

// LogString 返回用于日志的脱敏字符串
func LogString(v interface{}) string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	// 自定义编码器处理敏感字段
	if err := encodeSensitive(v, encoder); err != nil {
		return fmt.Sprintf("%v", v)
	}

	return buf.String()
}

// encodeSensitive 编码敏感数据
func encodeSensitive(v interface{}, encoder *json.Encoder) error {
	// 检查是否实现了敏感接口
	if sanitizer, ok := v.(interface{ Sanitize() interface{} }); ok {
		v = sanitizer.Sanitize()
	}

	return encoder.Encode(v)
}
