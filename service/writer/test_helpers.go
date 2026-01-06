package writer

import (
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
)

// ============ 测试辅助函数 ============

// createTestDocument 创建测试文档
func createTestDocument(id, projectID, title string) *writer.Document {
	return &writer.Document{
		IdentifiedEntity:     base.IdentifiedEntity{ID: id},
		ProjectScopedEntity:  base.ProjectScopedEntity{ProjectID: projectID},
		TitledEntity:         base.TitledEntity{Title: title},
		Timestamps:           base.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Type:                 "chapter",
		Level:                0,
		Order:                1,
		Status:               writer.DocumentStatusCompleted,
		WordCount:            1000,
		Tags:                 []string{"test", "sample"},
	}
}

// createTestProject 创建测试项目
func createTestProject(id, authorID, title string) *writer.Project {
	return &writer.Project{
		IdentifiedEntity: base.IdentifiedEntity{ID: id},
		OwnedEntity:      base.OwnedEntity{AuthorID: authorID},
		TitledEntity:     base.TitledEntity{Title: title},
		Timestamps:       base.Timestamps{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Summary:          "Test project summary",
		Status:           writer.StatusDraft,
	}
}
