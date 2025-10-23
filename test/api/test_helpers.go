package api

import (
	"Qingyu_backend/models/writer"
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/base"
)

// === Mock ProjectRepository ===

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
	if args.Error(0) == nil {
		if project.ID == "" {
			project.ID = "mock_project_id"
		}
		project.CreatedAt = time.Now()
		project.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) SoftDelete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Project, error) {
	return nil, nil
}

func (m *MockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	return 0, nil
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	return nil
}

func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	args := m.Called(ctx, projectID, ownerID, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateStatistics(ctx context.Context, id string, wordCount, chapterCount int64) error {
	args := m.Called(ctx, id, wordCount, chapterCount)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByCategory(ctx context.Context, category string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, category, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetPublicProjects(ctx context.Context, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) SearchProjects(ctx context.Context, keyword string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateStatus(ctx context.Context, id string, status writer.ProjectStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateVisibility(ctx context.Context, id string, visibility writer.Visibility) error {
	args := m.Called(ctx, id, visibility)
	return args.Error(0)
}

// === Mock EventBus ===

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return nil
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	return nil
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	return nil
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	return nil
}

// === 测试辅助函数 ===

func createTestProject(authorID string) *writer.Project {
	now := time.Now()
	return &writer.Project{
		ID:         "project123",
		AuthorID:   authorID,
		Title:      "测试项目",
		Summary:    "测试项目描述",
		Visibility: writer.VisibilityPrivate,
		Status:     writer.StatusDraft,
		Statistics: writer.ProjectStats{
			TotalWords:    10000,
			ChapterCount:  10,
			DocumentCount: 20,
			LastUpdateAt:  now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
