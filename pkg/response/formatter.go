package response

import (
	"reflect"
	"strings"
)

// SensitiveFields 敏感字段列表（需要过滤）
var SensitiveFields = []string{
	"password",
	"secret",
	"token",
	"api_key",
	"private_key",
	"access_token",
	"refresh_token",
}

// FilterSensitiveFields 过滤响应数据中的敏感字段
func FilterSensitiveFields(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return FilterSensitiveFields(v.Elem().Interface())
	case reflect.Struct:
		return filterStruct(v)
	case reflect.Slice, reflect.Array:
		return filterSlice(v)
	case reflect.Map:
		return filterMap(v)
	default:
		return data
	}
}

// filterStruct 过滤结构体中的敏感字段
func filterStruct(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 检查字段是否可导出
		if !field.IsExported() {
			continue
		}

		// 获取JSON标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 解析JSON标签
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// 检查是否为敏感字段
		if isSensitiveField(fieldName) {
			continue // 跳过敏感字段
		}

		// 递归过滤嵌套结构
		result[fieldName] = FilterSensitiveFields(fieldValue.Interface())
	}

	return result
}

// filterSlice 过滤切片中的敏感字段
func filterSlice(v reflect.Value) []interface{} {
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = FilterSensitiveFields(v.Index(i).Interface())
	}
	return result
}

// filterMap 过滤Map中的敏感字段
func filterMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyStr := ""
		if key.Kind() == reflect.String {
			keyStr = key.String()
		} else {
			keyStr = key.String()
		}

		// 检查键是否为敏感字段
		if isSensitiveField(keyStr) {
			continue
		}

		result[keyStr] = FilterSensitiveFields(value.Interface())
	}
	return result
}

// isSensitiveField 检查字段名是否为敏感字段
func isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, sensitive := range SensitiveFields {
		if lowerField == sensitive || strings.Contains(lowerField, sensitive) {
			return true
		}
	}
	return false
}

// AddSensitiveField 添加自定义敏感字段
func AddSensitiveField(field string) {
	SensitiveFields = append(SensitiveFields, strings.ToLower(field))
}

// RemoveSensitiveField 移除敏感字段
func RemoveSensitiveField(field string) {
	lowerField := strings.ToLower(field)
	for i, f := range SensitiveFields {
		if f == lowerField {
			SensitiveFields = append(SensitiveFields[:i], SensitiveFields[i+1:]...)
			break
		}
	}
}
