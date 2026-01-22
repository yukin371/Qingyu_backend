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
