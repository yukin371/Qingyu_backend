package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDConverter ID 类型转换工具
// 用于在 string ID 和 ObjectID 之间进行转换，保持向后兼容

// StringToObjectId 将字符串 ID 转换为 ObjectID
func StringToObjectId(id string) (primitive.ObjectID, error) {
	if id == "" {
		return primitive.NilObjectID, nil
	}
	return primitive.ObjectIDFromHex(id)
}

// ObjectIdToString 将 ObjectID 转换为字符串 ID
func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

// StringSliceToObjectIDSlice 将字符串切片转换为 ObjectID 切片
func StringSliceToObjectIDSlice(ids []string) ([]primitive.ObjectID, error) {
	if ids == nil {
		return nil, nil
	}
	result := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid ID format '%s': %w", id, err)
		}
		result = append(result, objectID)
	}
	return result, nil
}

// ObjectIDSliceToStringSlice 将 ObjectID 切片转换为字符串切片
func ObjectIDSliceToStringSlice(ids []primitive.ObjectID) []string {
	if ids == nil {
		return nil
	}
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, id.Hex())
	}
	return result
}

// StringPtrToObjectIdPtr 将字符串指针转换为 ObjectID 指针
func StringPtrToObjectIdPtr(idPtr *string) (*primitive.ObjectID, error) {
	if idPtr == nil || *idPtr == "" {
		return nil, nil
	}
	objectID, err := primitive.ObjectIDFromHex(*idPtr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return &objectID, nil
}

// ObjectIdPtrToStringPtr 将 ObjectID 指针转换为字符串指针
func ObjectIdPtrToStringPtr(idPtr *primitive.ObjectID) *string {
	if idPtr == nil {
		return nil
	}
	hex := idPtr.Hex()
	return &hex
}

// ConvertCategoryIdsForBook 将 Book 模型的 CategoryIDs (ObjectID) 转换为 string 切片
// 用于与新的 Category 模型 (string ID) 进行交互
func ConvertCategoryIdsToString(categoryIds []primitive.ObjectID) []string {
	return ObjectIDSliceToStringSlice(categoryIds)
}

// ConvertCategoryIdsFromBook 将 string 切片转换为 Book 模型的 CategoryIDs (ObjectID)
func ConvertCategoryIdsToObjectID(categoryIds []string) ([]primitive.ObjectID, error) {
	return StringSliceToObjectIDSlice(categoryIds)
}

// ============================================
// 统一ID解析函数（推荐使用）
// ============================================

// ParseID 解析必需的ID，空字符串返回错误
// 用于必需ID的场景，如GetByID、UpdateByID等
func ParseID(id string) (primitive.ObjectID, error) {
	if id == "" {
		return primitive.NilObjectID, ErrEmptyID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, id)
	}
	return oid, nil
}

// ParseOptionalID 解析可选ID，空字符串返回nil（不报错）
// 用于过滤条件等场景，如"不限分类"时category_id为空
func ParseOptionalID(id string) (*primitive.ObjectID, error) {
	if id == "" {
		return nil, nil
	}
	oid, err := ParseID(id)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}

// ParseIDs 批量解析ID列表，空列表返回nil
// 每个ID都必须有效，空字符串会报错
func ParseIDs(ids []string) ([]primitive.ObjectID, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	result := make([]primitive.ObjectID, 0, len(ids))
	for i, id := range ids {
		oid, err := ParseID(id)
		if err != nil {
			return nil, fmt.Errorf("ids[%d]: %w", i, err)
		}
		result = append(result, oid)
	}
	return result, nil
}

// ParseOptionalIDs 批量解析可选ID列表，跳过空字符串
// 用于批量过滤场景，空字符串被视为"不限制"
func ParseOptionalIDs(ids []string) ([]primitive.ObjectID, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	result := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		oid, err := ParseID(id)
		if err != nil {
			return nil, err
		}
		result = append(result, oid)
	}
	return result, nil
}
