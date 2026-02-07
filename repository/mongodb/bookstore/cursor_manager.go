package mongodb

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CursorManager 游标管理器
type CursorManager struct {
	defaultTTL int64 // 默认TTL（秒）
}

// NewCursorManager 创建游标管理器
func NewCursorManager() *CursorManager {
	return &CursorManager{
		defaultTTL: 3600, // 默认1小时
	}
}

// EncodeCursor 编码游标
func (cm *CursorManager) EncodeCursor(cursorType bookstore.CursorType, value interface{}) (string, error) {
	// 将值序列化为JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	// 构建游标数据
	cursorData := map[string]interface{}{
		"type": cursorType,
		"value": string(valueJSON),
		"ts":   time.Now().Unix(),
		"ttl":  cm.defaultTTL,
	}

	// 序列化整个游标数据
	data, err := json.Marshal(cursorData)
	if err != nil {
		return "", err
	}

	// Base64编码
	return base64.URLEncoding.EncodeToString(data), nil
}

// DecodeCursor 解码游标
func (cm *CursorManager) DecodeCursor(encoded string) (*bookstore.StreamCursor, error) {
	if encoded == "" {
		return nil, errors.New("cursor is empty")
	}

	// Base64解码
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// 反序列化
	var cursorData map[string]interface{}
	if err := json.Unmarshal(data, &cursorData); err != nil {
		return nil, err
	}

	// 提取字段
	cursorType, _ := cursorData["type"].(string)
	value, _ := cursorData["value"].(string)
	timestamp, _ := cursorData["ts"].(float64)
	ttl, _ := cursorData["ttl"].(float64)

	return &bookstore.StreamCursor{
		Type:      bookstore.CursorType(cursorType),
		Value:     value,
		Timestamp: int64(timestamp),
		TTL:       int64(ttl),
	}, nil
}

// ValidateCursor 验证游标是否有效
func (cm *CursorManager) ValidateCursor(encoded string) bool {
	if encoded == "" {
		return false
	}

	_, err := cm.DecodeCursor(encoded)
	return err == nil
}

// IsCursorExpired 检查游标是否过期
func (cm *CursorManager) IsCursorExpired(encoded string) bool {
	cursor, err := cm.DecodeCursor(encoded)
	if err != nil {
		return true // 无效游标视为已过期
	}

	now := time.Now().Unix()
	expirationTime := cursor.Timestamp + cursor.TTL
	return now > expirationTime
}

// BuildCursorFilter 根据游标构建MongoDB过滤条件
func (cm *CursorManager) BuildCursorFilter(encoded string, sortField string, sortOrder int) (bson.M, error) {
	cursor, err := cm.DecodeCursor(encoded)
	if err != nil {
		return nil, err
	}

	filter := bson.M{}

	switch cursor.Type {
	case bookstore.CursorTypeTimestamp:
		// 时间戳游标：使用时间比较
		if sortField == "" {
			sortField = "created_at"
		}
		// 解码时间戳值
		var timestamp int64
		if err := json.Unmarshal([]byte(cursor.Value), &timestamp); err == nil {
			if sortOrder < 0 {
				// 降序：查找小于当前时间戳的记录
				filter[sortField] = bson.M{"$lt": timestamp}
			} else {
				// 升序：查找大于当前时间戳的记录
				filter[sortField] = bson.M{"$gt": timestamp}
			}
		}

	case bookstore.CursorTypeID:
		// ID游标：使用ID比较
		if sortField == "" {
			sortField = "_id"
		}
		objectID, err := primitive.ObjectIDFromHex(cursor.Value)
		if err == nil {
			if sortOrder < 0 {
				// 降序：查找小于当前ID的记录
				filter[sortField] = bson.M{"$lt": objectID}
			} else {
				// 升序：查找大于当前ID的记录
				filter[sortField] = bson.M{"$gt": objectID}
			}
		}

	case bookstore.CursorTypeOffset:
		// Offset游标：直接使用偏移量（不构建filter，在应用层处理）
		// 返回空filter
	}

	return filter, nil
}

// GetCursorValue 获取游标的值
func (cm *CursorManager) GetCursorValue(encoded string) (interface{}, error) {
	cursor, err := cm.DecodeCursor(encoded)
	if err != nil {
		return nil, err
	}

	return cursor.Value, nil
}

// GenerateNextCursor 从书籍生成下一个游标
func (cm *CursorManager) GenerateNextCursor(book *bookstore.Book, cursorType bookstore.CursorType, sortField string) (string, error) {
	if book == nil {
		return "", errors.New("book is nil")
	}

	var value interface{}
	switch cursorType {
	case bookstore.CursorTypeTimestamp:
		// 使用创建时间或更新时间
		if sortField == "updated_at" {
			value = book.UpdatedAt.UnixMilli()
		} else {
			value = book.CreatedAt.UnixMilli()
		}

	case bookstore.CursorTypeID:
		// 使用书籍ID
		value = book.ID.Hex()

	case bookstore.CursorTypeOffset:
		// Offset游标需要传入具体的偏移量值
		// 这里使用0作为默认值，实际使用时需要在外部计算
		value = 0

	default:
		return "", errors.New("unsupported cursor type")
	}

	return cm.EncodeCursor(cursorType, value)
}

// GetDefaultTTL 获取默认TTL
func (cm *CursorManager) GetDefaultTTL() int64 {
	return cm.defaultTTL
}

// SetDefaultTTL 设置默认TTL
func (cm *CursorManager) SetDefaultTTL(ttl int64) {
	cm.defaultTTL = ttl
}

// base64Encode Base64编码辅助方法
func (cm *CursorManager) base64Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// base64Decode Base64解码辅助方法
func (cm *CursorManager) base64Decode(encoded string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(encoded)
}
