package writer

import (
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ 测试辅助函数 ============

// createTestDocument 创建测试文档
func createTestDocument(projectID, title string) *writer.Document {
	oid, _ := primitive.ObjectIDFromHex(projectID)
	return &writer.Document{
		IdentifiedEntity: base.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       base.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ProjectID:        oid,
		Title:            title,
		Type:             "chapter",
		Level:            0,
		Order:            1,
		Status:           writer.DocumentStatusCompleted,
		WordCount:        1000,
		Tags:             []string{"test", "sample"},
	}
}

// createTestProject 创建测试项目
func createTestProject(authorID, title string) *writer.Project {
	return &writer.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: primitive.NewObjectID()},
		OwnedEntity:      base.OwnedEntity{AuthorID: authorID},
		TitledEntity:     base.TitledEntity{Title: title},
		Timestamps:       base.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Summary:          "Test project summary",
		Status:           writer.StatusDraft,
	}
}
