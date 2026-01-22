package types

import (
	"fmt"
	"reflect"
	"strings"
)

// NamingConvention 命名规范
type NamingConvention int

const (
	// CamelCase JSON 命名（API 对外）
	CamelCase NamingConvention = iota
	// snake_case BSON 命名（存储）
	SnakeCase
)

// FieldTags 字段标签配置
type FieldTags struct {
	JSON string
	BSON string
}

// ToTagString 生成 struct tag 字符串
func (ft FieldTags) ToTagString() string {
	var parts []string
	if ft.BSON != "" {
		parts = append(parts, fmt.Sprintf(`bson:"%s"`, ft.BSON))
	}
	if ft.JSON != "" {
		parts = append(parts, fmt.Sprintf(`json:"%s"`, ft.JSON))
	}
	return strings.Join(parts, " ")
}

// StandardFieldTags 生成标准字段标签
// JSON: camelCase, BSON: snake_case
func StandardFieldTags(fieldName string) FieldTags {
	return FieldTags{
		JSON: toCamelCase(fieldName),
		BSON: toSnakeCase(fieldName),
	}
}

// ModelTags 生成 Model 层标签（只 BSON）
func ModelTags(fieldName string) string {
	return fmt.Sprintf(`bson:"%s"`, toSnakeCase(fieldName))
}

// DTOTags 生成 DTO 层标签（只 JSON）
func DTOTags(fieldName string) string {
	return fmt.Sprintf(`json:"%s"`, toCamelCase(fieldName))
}

// TagSet 标签集合（用于生成 struct tag）
type TagSet map[string]string

// NewTagSet 创建新的标签集合
func NewTagSet() TagSet {
	return make(TagSet)
}

// String 生成 struct tag 字符串
func (t TagSet) String() string {
	if len(t) == 0 {
		return ""
	}

	var parts []string
	// 顺序：bson, json, validate, 其他
	for _, key := range []string{"bson", "json", "validate"} {
		if value, ok := t[key]; ok && value != "" {
			parts = append(parts, fmt.Sprintf(`%s:"%s"`, key, value))
			delete(t, key)
		}
	}

	// 添加其他标签
	for key, value := range t {
		if value != "" {
			parts = append(parts, fmt.Sprintf(`%s:"%s"`, key, value))
		}
	}

	return strings.Join(parts, " ")
}

// WithBSON 添加 BSON 标签
func (t TagSet) WithBSON(name string) TagSet {
	t["bson"] = name
	return t
}

// WithJSON 添加 JSON 标签
func (t TagSet) WithJSON(name string) TagSet {
	t["json"] = name
	return t
}

// WithValidate 添加 validate 标签
func (t TagSet) WithValidate(rule string) TagSet {
	t["validate"] = rule
	return t
}

// WithTag 添加任意标签
func (t TagSet) WithTag(key, value string) TagSet {
	t[key] = value
	return t
}

// toSnakeCase 将驼峰命名转换为下划线命名
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		if r >= 'A' && r <= 'Z' {
			result = append(result, r+('a'-'A'))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// toCamelCase 将下划线命名转换为驼峰命名
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i == 0 {
			// 第一个单词首字母小写
			if len(part) > 0 {
				parts[i] = strings.ToLower(string(part[0])) + part[1:]
			}
		} else {
			// 后续单词首字母大写
			if len(part) > 0 {
				parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			}
		}
	}
	return strings.Join(parts, "")
}

// StructTagBuilder struct 标签构建器
type StructTagBuilder struct {
	fields map[string]TagSet
}

// NewStructTagBuilder 创建 struct 标签构建器
func NewStructTagBuilder() *StructTagBuilder {
	return &StructTagBuilder{
		fields: make(map[string]TagSet),
	}
}

// AddField 添加字段
func (b *StructTagBuilder) AddField(fieldName string, tags TagSet) *StructTagBuilder {
	b.fields[fieldName] = tags
	return b
}

// AddStandardField 添加标准字段（自动转换命名）
func (b *StructTagBuilder) AddStandardField(fieldName string, validateRule string) *StructTagBuilder {
	tags := NewTagSet().
		WithBSON(toSnakeCase(fieldName)).
		WithJSON(toCamelCase(fieldName))
	if validateRule != "" {
		tags = tags.WithValidate(validateRule)
	}
	b.fields[fieldName] = tags
	return b
}

// Build 生成完整的 struct 定义（带标签）
func (b *StructTagBuilder) Build() string {
	var lines []string
	lines = append(lines, "type StructName struct {")
	for fieldName, tags := range b.fields {
		lines = append(lines, fmt.Sprintf("    %s %s `%s`",
			toCamelCase(fieldName), // 字段名
			"string",              // 类型（示例）
			tags.String(),         // 标签
		))
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// GetStructTags 从结构体提取所有字段标签
func GetStructTags(v interface{}) map[string]TagSet {
	tags := make(map[string]TagSet)
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return tags
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldTags := NewTagSet()

		// 提取 bson 标签
		if bsonTag := field.Tag.Get("bson"); bsonTag != "" {
			fieldTags = fieldTags.WithBSON(bsonTag)
		}

		// 提取 json 标签
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			fieldTags = fieldTags.WithJSON(jsonTag)
		}

		// 提取 validate 标签
		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			fieldTags = fieldTags.WithValidate(validateTag)
		}

		tags[field.Name] = fieldTags
	}

	return tags
}

// CommonTags 常用标签组合
var (
	// RequiredID 必填的 ID 字段
	RequiredID = NewTagSet().
		WithBSON("_id").
		WithJSON("id").
		WithValidate("required")

	// RequiredUserID 必填的用户 ID
	RequiredUserID = NewTagSet().
		WithBSON("user_id").
		WithJSON("userId").
		WithValidate("required,objectId")

	// RequiredAuthorID 必填的作者 ID
	RequiredAuthorID = NewTagSet().
		WithBSON("author_id").
		WithJSON("authorId").
		WithValidate("required,objectId")

	// OptionalString 可选字符串
	OptionalString = NewTagSet().
		WithBSON("name").
		WithJSON("name")

	// RequiredString 必填字符串
	RequiredString = NewTagSet().
		WithBSON("name").
		WithJSON("name").
		WithValidate("required")

	// Timestamp 时间戳
	Timestamp = NewTagSet().
		WithBSON("created_at").
		WithJSON("createdAt")
)
