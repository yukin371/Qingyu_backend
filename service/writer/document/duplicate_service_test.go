package document

import (
	"Qingyu_backend/models/writer"
	"Qingyu_backend/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestDuplicateService_validateRequest 测试参数验证
func TestDuplicateService_validateRequest(t *testing.T) {
	service := &DuplicateService{
		serviceName: "DuplicateService",
	}

	tests := []struct {
		name        string
		documentID  string
		req         *DuplicateRequest
		expectError bool
		errorMsg    string
	}{
		{
			name:        "空文档ID",
			documentID:  "",
			req:         &DuplicateRequest{Position: "inner"},
			expectError: true,
			errorMsg:    "文档ID不能为空",
		},
		{
			name:        "空请求",
			documentID:  "doc-123",
			req:         nil,
			expectError: true,
			errorMsg:    "请求参数不能为空",
		},
		{
			name:        "无效的Position",
			documentID:  "doc-123",
			req:         &DuplicateRequest{Position: "invalid"},
			expectError: true,
			errorMsg:    "无效的position参数",
		},
		{
			name:        "有效的请求 - inner",
			documentID:  "doc-123",
			req:         &DuplicateRequest{Position: "inner"},
			expectError: false,
		},
		{
			name:        "有效的请求 - before",
			documentID:  "doc-123",
			req:         &DuplicateRequest{Position: "before"},
			expectError: false,
		},
		{
			name:        "有效的请求 - after",
			documentID:  "doc-123",
			req:         &DuplicateRequest{Position: "after"},
			expectError: false,
		},
		{
			name:        "有效的请求 - 空Position",
			documentID:  "doc-123",
			req:         &DuplicateRequest{Position: ""},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateRequest(tt.documentID, tt.req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDuplicateService_createDuplicateDocument 测试创建复制的文档
func TestDuplicateService_createDuplicateDocument(t *testing.T) {
	service := &DuplicateService{
		serviceName: "DuplicateService",
	}

	projectID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()

	sourceDoc := &writer.Document{
		ProjectID:    projectID,
		Title:        "第一章",
		Type:         "chapter",
		Level:        1,
		Order:        0,
		Status:       writer.DocumentStatusPlanned,
		StableRef:    "chapter-1",
		OrderKey:     "a0",
		WordCount:    1000,
		ParentID:     parentID,
		CharacterIDs: []primitive.ObjectID{primitive.NewObjectID()},
		LocationIDs:  []primitive.ObjectID{primitive.NewObjectID()},
		TimelineIDs:  []primitive.ObjectID{primitive.NewObjectID()},
		PlotThreads:  []string{"plot1", "plot2"},
		KeyPoints:    []string{"key1", "key2"},
		WritingHints: []string{"hint1"},
		Tags:         []string{"tag1"},
		Notes:        "这是备注",
	}

	t.Run("基本复制", func(t *testing.T) {
		req := &DuplicateRequest{
			Position:    "inner",
			CopyContent: false,
		}

		newDoc := service.createDuplicateDocument(sourceDoc, req)

		// 验证基本属性
		assert.NotEqual(t, sourceDoc.ID, newDoc.ID)
		assert.Equal(t, projectID, newDoc.ProjectID)
		assert.Equal(t, "Copy - 第一章", newDoc.Title)
		assert.Equal(t, "chapter", newDoc.Type)
		assert.Equal(t, 1, newDoc.Level)
		assert.Equal(t, 0, newDoc.Order)
		assert.Equal(t, writer.DocumentStatusPlanned, newDoc.Status)
		assert.Equal(t, 0, newDoc.WordCount) // 字数重置为0

		// 验证StableRef和OrderKey
		assert.Equal(t, "chapter-1-copy", newDoc.StableRef)
		assert.Equal(t, "a00", newDoc.OrderKey) // GenerateSiblingOrderKey("a0") = "a00"

		// 验证关联信息被复制
		assert.Equal(t, len(sourceDoc.CharacterIDs), len(newDoc.CharacterIDs))
		assert.Equal(t, len(sourceDoc.LocationIDs), len(newDoc.LocationIDs))
		assert.Equal(t, len(sourceDoc.TimelineIDs), len(newDoc.TimelineIDs))
		assert.Equal(t, sourceDoc.PlotThreads, newDoc.PlotThreads)
		assert.Equal(t, sourceDoc.KeyPoints, newDoc.KeyPoints)
		assert.Equal(t, sourceDoc.WritingHints, newDoc.WritingHints)
		assert.Equal(t, sourceDoc.Tags, newDoc.Tags)
		assert.Equal(t, sourceDoc.Notes, newDoc.Notes)
	})

	t.Run("指定目标父节点", func(t *testing.T) {
		newParentID := primitive.NewObjectID()
		newParentIDStr := newParentID.Hex()

		req := &DuplicateRequest{
			TargetParentID: &newParentIDStr,
			Position:       "inner",
			CopyContent:    false,
		}

		newDoc := service.createDuplicateDocument(sourceDoc, req)

		// 验证父节点已更改
		assert.Equal(t, newParentID, newDoc.ParentID)
	})

	t.Run("不指定目标父节点", func(t *testing.T) {
		req := &DuplicateRequest{
			Position:    "inner",
			CopyContent: false,
		}

		newDoc := service.createDuplicateDocument(sourceDoc, req)

		// 验证使用原文档的父节点
		assert.Equal(t, sourceDoc.ParentID, newDoc.ParentID)
	})
}

// TestGenerateSiblingOrderKey 测试GenerateSiblingOrderKey函数
func TestGenerateSiblingOrderKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空输入",
			input:    "",
			expected: "a0",
		},
		{
			name:     "在a0后生成",
			input:    "a0",
			expected: "a00",
		},
		{
			name:     "在a00后生成",
			input:    "a00",
			expected: "a000",
		},
		{
			name:     "在a000后生成",
			input:    "a000",
			expected: "a0000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GenerateSiblingOrderKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDuplicateService_ToResponse 测试ToResponse方法
func TestDuplicateService_ToResponse(t *testing.T) {
	service := &DuplicateService{
		serviceName: "DuplicateService",
	}

	docID := primitive.NewObjectID()
	doc := &writer.Document{
		Title:     "测试文档",
		StableRef: "test-doc",
	}
	doc.ID = docID

	result := service.ToResponse(doc)

	assert.NotNil(t, result)
	assert.Equal(t, docID.Hex(), result.DocumentID)
	assert.Equal(t, "测试文档", result.Title)
	assert.Equal(t, "test-doc", result.StableRef)
}

// TestDuplicateService_BaseServiceInterface 测试BaseService接口实现
func TestDuplicateService_BaseServiceInterface(t *testing.T) {
	service := &DuplicateService{
		serviceName: "DuplicateService",
		version:     "1.0.0",
	}

	// 测试接口方法
	assert.Equal(t, "DuplicateService", service.GetServiceName())
	assert.Equal(t, "1.0.0", service.GetVersion())

	ctx := context.Background()
	assert.NoError(t, service.Initialize(ctx))
	assert.NoError(t, service.Close(ctx))
}

// TestDuplicateRequest 测试DuplicateRequest结构
func TestDuplicateRequest(t *testing.T) {
	req := &DuplicateRequest{
		Position:    "inner",
		CopyContent: true,
	}

	assert.Equal(t, "inner", req.Position)
	assert.True(t, req.CopyContent)
	assert.Nil(t, req.TargetParentID)
}

// TestDuplicateResponse 测试DuplicateResponse结构
func TestDuplicateResponse(t *testing.T) {
	resp := &DuplicateResponse{
		DocumentID: "doc-123",
		Title:      "Copy - 测试文档",
		StableRef:  "test-doc-copy",
	}

	assert.Equal(t, "doc-123", resp.DocumentID)
	assert.Equal(t, "Copy - 测试文档", resp.Title)
	assert.Equal(t, "test-doc-copy", resp.StableRef)
}

// BenchmarkCreateDuplicateDocument 性能测试
func BenchmarkCreateDuplicateDocument(b *testing.B) {
	service := &DuplicateService{
		serviceName: "DuplicateService",
	}

	projectID := primitive.NewObjectID()
	sourceDoc := &writer.Document{
		ProjectID:    projectID,
		Title:        "第一章",
		Type:         "chapter",
		Level:        1,
		Order:        0,
		Status:       writer.DocumentStatusPlanned,
		StableRef:    "chapter-1",
		OrderKey:     "a0",
		CharacterIDs: make([]primitive.ObjectID, 10),
		LocationIDs:  make([]primitive.ObjectID, 5),
		PlotThreads:  make([]string, 20),
		KeyPoints:    make([]string, 20),
		WritingHints: make([]string, 20),
		Tags:         make([]string, 10),
	}

	req := &DuplicateRequest{
		Position:    "inner",
		CopyContent: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.createDuplicateDocument(sourceDoc, req)
	}
}
