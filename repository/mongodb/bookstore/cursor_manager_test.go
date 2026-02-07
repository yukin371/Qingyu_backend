package mongodb

import (
	"Qingyu_backend/models/bookstore"
	"testing"
	"time"
)

// TestCursorType 测试游标类型常量
func TestCursorType(t *testing.T) {
	tests := []struct {
		name       string
		cursorType bookstore.CursorType
		expected   string
	}{
		{"Offset游标", bookstore.CursorTypeOffset, "offset"},
		{"Timestamp游标", bookstore.CursorTypeTimestamp, "timestamp"},
		{"ID游标", bookstore.CursorTypeID, "id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.cursorType) != tt.expected {
				t.Errorf("CursorType = %v, want %v", tt.cursorType, tt.expected)
			}
		})
	}
}

// TestStreamCursor 测试游标结构
func TestStreamCursor(t *testing.T) {
	now := time.Now()
	cursor := &bookstore.StreamCursor{
		Type:      bookstore.CursorTypeTimestamp,
		Value:     "1234567890",
		Timestamp: now.Unix(),
		TTL:       3600,
	}

	if cursor.Type != bookstore.CursorTypeTimestamp {
		t.Errorf("Type = %v, want %v", cursor.Type, bookstore.CursorTypeTimestamp)
	}
	if cursor.Value != "1234567890" {
		t.Errorf("Value = %v, want %v", cursor.Value, "1234567890")
	}
	if cursor.Timestamp != now.Unix() {
		t.Errorf("Timestamp = %v, want %v", cursor.Timestamp, now.Unix())
	}
	if cursor.TTL != 3600 {
		t.Errorf("TTL = %v, want %v", cursor.TTL, 3600)
	}
}

// TestEncodeCursor 测试游标编码 - 基础功能
func TestEncodeCursor(t *testing.T) {
	cm := NewCursorManager()

	tests := []struct {
		name        string
		cursorType  bookstore.CursorType
		value       interface{}
		expectError bool
	}{
		{
			name:        "编码Offset游标",
			cursorType:  bookstore.CursorTypeOffset,
			value:       100,
			expectError: false,
		},
		{
			name:        "编码Timestamp游标",
			cursorType:  bookstore.CursorTypeTimestamp,
			value:       int64(1706140800000),
			expectError: false,
		},
		{
			name:        "编码ID游标",
			cursorType:  bookstore.CursorTypeID,
			value:       "507f1f77bcf86cd799439011",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := cm.EncodeCursor(tt.cursorType, tt.value)

			if tt.expectError && err == nil {
				t.Error("期望返回错误，但没有")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望错误，但得到: %v", err)
			}
			if !tt.expectError && encoded == "" {
				t.Error("编码结果不应为空")
			}
			if !tt.expectError && encoded != "" {
				// 验证是Base64编码
				if len(encoded) < 10 {
					t.Errorf("编码结果太短: %s", encoded)
				}
			}
		})
	}
}

// TestDecodeCursor 测试游标解码
func TestDecodeCursor(t *testing.T) {
	cm := NewCursorManager()

	// 先编码一个游标
	originalValue := int64(1706140800000)
	encoded, err := cm.EncodeCursor(bookstore.CursorTypeTimestamp, originalValue)
	if err != nil {
		t.Fatalf("编码失败: %v", err)
	}

	// 测试解码
	t.Run("解码有效游标", func(t *testing.T) {
		decoded, err := cm.DecodeCursor(encoded)
		if err != nil {
			t.Errorf("解码失败: %v", err)
		}
		if decoded == nil {
			t.Fatal("解码结果不应为nil")
		}
		if decoded.Type != bookstore.CursorTypeTimestamp {
			t.Errorf("Type = %v, want %v", decoded.Type, bookstore.CursorTypeTimestamp)
		}
		if decoded.Timestamp == 0 {
			t.Error("Timestamp不应为0")
		}
		if decoded.TTL == 0 {
			t.Error("TTL不应为0")
		}
	})

	t.Run("解码无效游标", func(t *testing.T) {
		invalidCursor := "invalid-base64!@#$%"
		_, err := cm.DecodeCursor(invalidCursor)
		if err == nil {
			t.Error("期望返回错误，但没有")
		}
	})

	t.Run("解码空游标", func(t *testing.T) {
		_, err := cm.DecodeCursor("")
		if err == nil {
			t.Error("期望返回错误，但没有")
		}
	})
}

// TestEncodeDecodeRoundTrip 测试编解码往返
func TestEncodeDecodeRoundTrip(t *testing.T) {
	cm := NewCursorManager()

	tests := []struct {
		name       string
		cursorType bookstore.CursorType
		value      interface{}
	}{
		{"Offset往返", bookstore.CursorTypeOffset, 100},
		{"Timestamp往返", bookstore.CursorTypeTimestamp, int64(1706140800000)},
		{"ID往返", bookstore.CursorTypeID, "507f1f77bcf86cd799439011"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := cm.EncodeCursor(tt.cursorType, tt.value)
			if err != nil {
				t.Fatalf("编码失败: %v", err)
			}

			decoded, err := cm.DecodeCursor(encoded)
			if err != nil {
				t.Fatalf("解码失败: %v", err)
			}

			if decoded.Type != tt.cursorType {
				t.Errorf("Type = %v, want %v", decoded.Type, tt.cursorType)
			}
		})
	}
}

// TestValidateCursor 测试游标验证
func TestValidateCursor(t *testing.T) {
	cm := NewCursorManager()

	t.Run("验证有效游标", func(t *testing.T) {
		encoded, _ := cm.EncodeCursor(bookstore.CursorTypeTimestamp, int64(1706140800000))
		valid := cm.ValidateCursor(encoded)
		if !valid {
			t.Error("有效游标验证失败")
		}
	})

	t.Run("验证无效游标", func(t *testing.T) {
		valid := cm.ValidateCursor("invalid-cursor")
		if valid {
			t.Error("无效游标应该验证失败")
		}
	})

	t.Run("验证空游标", func(t *testing.T) {
		valid := cm.ValidateCursor("")
		if valid {
			t.Error("空游标应该验证失败")
		}
	})
}

// TestIsCursorExpired 测试游标过期检查
func TestIsCursorExpired(t *testing.T) {
	cm := NewCursorManager()

	t.Run("未过期的游标", func(t *testing.T) {
		encoded, _ := cm.EncodeCursor(bookstore.CursorTypeTimestamp, int64(time.Now().UnixMilli()))
		expired := cm.IsCursorExpired(encoded)
		if expired {
			t.Error("新创建的游标不应该过期")
		}
	})
}

// TestBuildCursorFilter 测试构建游标过滤条件
func TestBuildCursorFilter(t *testing.T) {
	cm := NewCursorManager()

	t.Run("Timestamp游标过滤", func(t *testing.T) {
		encoded, _ := cm.EncodeCursor(bookstore.CursorTypeTimestamp, int64(1706140800000))
		filter, err := cm.BuildCursorFilter(encoded, "created_at", -1)
		if err != nil {
			t.Fatalf("构建过滤条件失败: %v", err)
		}
		if filter == nil {
			t.Fatal("过滤条件不应为nil")
		}
	})

	t.Run("无效游标", func(t *testing.T) {
		_, err := cm.BuildCursorFilter("invalid", "", 0)
		if err == nil {
			t.Error("期望返回错误，但没有")
		}
	})
}
