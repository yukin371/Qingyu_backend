package writer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/interfaces"
)

// mockCharacterServiceForEntity 用于 entity_service 测试的 mock
type mockCharacterServiceForEntity struct {
	mock.Mock
}

func (m *mockCharacterServiceForEntity) Create(ctx context.Context, projectID, userID string, req *interfaces.CreateCharacterRequest) (*writer.Character, error) {
	args := m.Called(ctx, projectID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Character), args.Error(1)
}

func (m *mockCharacterServiceForEntity) GetByID(ctx context.Context, characterID, projectID string) (*writer.Character, error) {
	args := m.Called(ctx, characterID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Character), args.Error(1)
}

func (m *mockCharacterServiceForEntity) List(ctx context.Context, projectID string) ([]*writer.Character, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Character), args.Error(1)
}

func (m *mockCharacterServiceForEntity) Update(ctx context.Context, characterID, projectID string, req *interfaces.UpdateCharacterRequest) (*writer.Character, error) {
	args := m.Called(ctx, characterID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Character), args.Error(1)
}

func (m *mockCharacterServiceForEntity) Delete(ctx context.Context, characterID, projectID string) error {
	args := m.Called(ctx, characterID, projectID)
	return args.Error(0)
}

func (m *mockCharacterServiceForEntity) CreateRelation(ctx context.Context, projectID string, req *interfaces.CreateRelationRequest) (*writer.CharacterRelation, error) {
	args := m.Called(ctx, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.CharacterRelation), args.Error(1)
}

func (m *mockCharacterServiceForEntity) ListRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	args := m.Called(ctx, projectID, characterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.CharacterRelation), args.Error(1)
}

func (m *mockCharacterServiceForEntity) DeleteRelation(ctx context.Context, relationID, projectID string) error {
	args := m.Called(ctx, relationID, projectID)
	return args.Error(0)
}

func (m *mockCharacterServiceForEntity) GetCharacterGraph(ctx context.Context, projectID string) (*interfaces.CharacterGraph, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.CharacterGraph), args.Error(1)
}

func (m *mockCharacterServiceForEntity) CreateRelationTimelineEvent(ctx context.Context, projectID string, req *interfaces.CreateRelationTimelineEventRequest) (*writer.RelationTimelineEvent, error) {
	args := m.Called(ctx, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.RelationTimelineEvent), args.Error(1)
}

func (m *mockCharacterServiceForEntity) GetRelationTimeline(ctx context.Context, relationID, projectID string) ([]*writer.RelationTimelineEvent, error) {
	args := m.Called(ctx, relationID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.RelationTimelineEvent), args.Error(1)
}

func (m *mockCharacterServiceForEntity) UpdateRelationTimelineEvent(ctx context.Context, eventID, projectID string, req *interfaces.UpdateRelationTimelineEventRequest) (*writer.RelationTimelineEvent, error) {
	args := m.Called(ctx, eventID, projectID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.RelationTimelineEvent), args.Error(1)
}

func (m *mockCharacterServiceForEntity) DeleteRelationTimelineEvent(ctx context.Context, eventID, projectID string) error {
	args := m.Called(ctx, eventID, projectID)
	return args.Error(0)
}

// TestListEntities_CharactersOnly 测试仅查询角色类型实体
func TestListEntities_CharactersOnly(t *testing.T) {
	mockCharSvc := new(mockCharacterServiceForEntity)
	svc := NewEntityService(mockCharSvc, nil, nil)

	ctx := context.Background()
	projectID := "test-project-id"
	entityType := string(writer.EntityTypeCharacter)

	mockCharSvc.On("List", ctx, projectID).Return([]*writer.Character{}, nil)

	result, err := svc.ListEntities(ctx, projectID, &entityType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockCharSvc.AssertExpectations(t)
}

// TestListEntities_AllTypes 测试查询所有类型实体（无筛选）
func TestListEntities_AllTypes(t *testing.T) {
	mockCharSvc := new(mockCharacterServiceForEntity)
	svc := NewEntityService(mockCharSvc, nil, nil)

	ctx := context.Background()
	projectID := "test-project-id"

	mockCharSvc.On("List", ctx, projectID).Return([]*writer.Character{}, nil)

	result, err := svc.ListEntities(ctx, projectID, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockCharSvc.AssertExpectations(t)
}

// TestListEntities_InvalidFilter 测试无效的筛选类型不返回角色
func TestListEntities_InvalidFilter(t *testing.T) {
	mockCharSvc := new(mockCharacterServiceForEntity)
	svc := NewEntityService(mockCharSvc, nil, nil)

	ctx := context.Background()
	projectID := "test-project-id"
	entityType := "organization"

	result, err := svc.ListEntities(ctx, projectID, &entityType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockCharSvc.AssertNotCalled(t, "List")
}

// TestGetEntityGraph 测试获取实体图谱
func TestGetEntityGraph(t *testing.T) {
	mockCharSvc := new(mockCharacterServiceForEntity)
	svc := NewEntityService(mockCharSvc, nil, nil)

	ctx := context.Background()
	projectID := "test-project-id"

	mockCharSvc.On("List", ctx, projectID).Return([]*writer.Character{}, nil)
	mockCharSvc.On("ListRelations", ctx, projectID, (*string)(nil)).Return([]*writer.CharacterRelation{}, nil)

	graph, err := svc.GetEntityGraph(ctx, projectID)

	assert.NoError(t, err)
	assert.NotNil(t, graph)
	assert.Empty(t, graph.Nodes)
	assert.Empty(t, graph.Edges)
	mockCharSvc.AssertExpectations(t)
}

// TestListEntities_WithItems 测试仅查询物品类型
func TestListEntities_WithItems(t *testing.T) {
	mockCharSvc := new(mockCharacterServiceForEntity)
	svc := NewEntityService(mockCharSvc, nil, nil)

	ctx := context.Background()
	projectID := "test-project-id"
	entityType := string(writer.EntityTypeItem)

	result, err := svc.ListEntities(ctx, projectID, &entityType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockCharSvc.AssertNotCalled(t, "List")
}
