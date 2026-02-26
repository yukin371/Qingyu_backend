package writer

import (
	"context"
	"errors"
	"testing"

	writerModels "Qingyu_backend/models/writer"
	writerRepoInterfaces "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockOutlineRepository Mock大纲Repository
type MockOutlineRepository struct {
	mock.Mock
}

func (m *MockOutlineRepository) Create(ctx context.Context, outline *writerModels.OutlineNode) error {
	args := m.Called(ctx, outline)
	return args.Error(0)
}

func (m *MockOutlineRepository) FindByID(ctx context.Context, outlineID string) (*writerModels.OutlineNode, error) {
	args := m.Called(ctx, outlineID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writerModels.OutlineNode), args.Error(1)
}

func (m *MockOutlineRepository) FindByProjectID(ctx context.Context, projectID string) ([]*writerModels.OutlineNode, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.OutlineNode), args.Error(1)
}

func (m *MockOutlineRepository) Update(ctx context.Context, outline *writerModels.OutlineNode) error {
	args := m.Called(ctx, outline)
	return args.Error(0)
}

func (m *MockOutlineRepository) Delete(ctx context.Context, outlineID string) error {
	args := m.Called(ctx, outlineID)
	return args.Error(0)
}

func (m *MockOutlineRepository) FindByParentID(ctx context.Context, projectID, parentID string) ([]*writerModels.OutlineNode, error) {
	args := m.Called(ctx, projectID, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.OutlineNode), args.Error(1)
}

func (m *MockOutlineRepository) FindRoots(ctx context.Context, projectID string) ([]*writerModels.OutlineNode, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writerModels.OutlineNode), args.Error(1)
}

func (m *MockOutlineRepository) ExistsByID(ctx context.Context, outlineID string) (bool, error) {
	args := m.Called(ctx, outlineID)
	return args.Bool(0), args.Error(1)
}

func (m *MockOutlineRepository) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOutlineRepository) CountByParentID(ctx context.Context, projectID, parentID string) (int64, error) {
	args := m.Called(ctx, projectID, parentID)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

// Ensure MockOutlineRepository implements the interface
var _ writerRepoInterfaces.OutlineRepository = (*MockOutlineRepository)(nil)

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

// 辅助函数：创建测试大纲节点
func createTestOutlineNode(projectID, title, parentID string, order int) *writerModels.OutlineNode {
	id := primitive.NewObjectID()
	outline := &writerModels.OutlineNode{}
	outline.ID = id
	outline.ProjectID = projectID
	outline.Title = title
	outline.ParentID = parentID
	outline.Order = order
	outline.Summary = "测试摘要"
	outline.Type = "chapter"
	outline.Tension = 5
	outline.Characters = []string{}
	outline.Items = []string{}
	return outline
}

// TestOutlineService_Create_RootNode 测试创建根节点
func TestOutlineService_Create_RootNode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	req := &serviceInterfaces.CreateOutlineRequest{
		Title:    "第一卷",
		Summary:  "这是第一卷的内容",
		Type:     "volume",
		Tension:  7,
		Order:    0,
	}

	mockRepo.On("CountByParentID", ctx, projectID, "").Return(int64(0), nil).Maybe()
	mockRepo.On("Create", ctx, mock.AnythingOfType("*writer.OutlineNode")).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	outline, err := service.Create(ctx, projectID, userID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, outline)
	assert.Equal(t, req.Title, outline.Title)
	assert.Equal(t, req.Summary, outline.Summary)
	assert.Equal(t, req.Type, outline.Type)
	assert.Equal(t, req.Tension, outline.Tension)
	assert.Equal(t, "", outline.ParentID) // 根节点

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestOutlineService_Create_ChildNode 测试创建子节点
func TestOutlineService_Create_ChildNode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	parentNode := createTestOutlineNode(projectID, "第一卷", "", 0)

	req := &serviceInterfaces.CreateOutlineRequest{
		Title:    "第一章",
		ParentID: parentID,
		Summary:  "这是第一章的内容",
		Type:     "chapter",
		Tension:  5,
	}

	mockRepo.On("FindByID", ctx, parentID).Return(parentNode, nil).Once()
	mockRepo.On("CountByParentID", ctx, projectID, parentID).Return(int64(0), nil).Once()
	mockRepo.On("Create", ctx, mock.AnythingOfType("*writer.OutlineNode")).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	outline, err := service.Create(ctx, projectID, userID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, outline)
	assert.Equal(t, req.Title, outline.Title)
	assert.Equal(t, parentID, outline.ParentID)

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestOutlineService_Create_ParentNotFound 测试父节点不存在
func TestOutlineService_Create_ParentNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	req := &serviceInterfaces.CreateOutlineRequest{
		Title:    "第一章",
		ParentID: parentID,
	}

	mockRepo.On("FindByID", ctx, parentID).Return(nil, errors.New("not found")).Once()

	// Act
	outline, err := service.Create(ctx, projectID, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, outline)
	assert.Contains(t, err.Error(), "parent outline not found")

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_Create_CrossProjectParent 测试跨项目父节点失败
func TestOutlineService_Create_CrossProjectParent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	parentNode := createTestOutlineNode(otherProjectID, "其他项目的卷", "", 0)

	req := &serviceInterfaces.CreateOutlineRequest{
		Title:    "第一章",
		ParentID: parentID,
	}

	mockRepo.On("FindByID", ctx, parentID).Return(parentNode, nil).Once()

	// Act
	outline, err := service.Create(ctx, projectID, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, outline)
	assert.Contains(t, err.Error(), "does not belong to this project")

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_Create_InvalidData 测试无效数据
func TestOutlineService_Create_InvalidData(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	req := &serviceInterfaces.CreateOutlineRequest{
		Title: "", // 空标题
	}

	// Act
	outline, err := service.Create(ctx, projectID, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, outline)
	assert.Contains(t, err.Error(), "invalid outline data")
}

// TestOutlineService_GetByID_Success 测试成功获取大纲
func TestOutlineService_GetByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	expectedOutline := createTestOutlineNode(projectID, "第一章", "", 0)

	mockRepo.On("FindByID", ctx, outlineID).Return(expectedOutline, nil).Once()

	// Act
	outline, err := service.GetByID(ctx, outlineID, projectID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, outline)
	assert.Equal(t, expectedOutline.ID, outline.ID)
	assert.Equal(t, expectedOutline.Title, outline.Title)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetByID_NotFound 测试获取不存在的大纲
func TestOutlineService_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	mockRepo.On("FindByID", ctx, outlineID).Return(nil, errors.New("not found")).Once()

	// Act
	outline, err := service.GetByID(ctx, outlineID, projectID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, outline)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetByID_NoPermission 测试无权限访问
func TestOutlineService_GetByID_NoPermission(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	outline := createTestOutlineNode(otherProjectID, "其他项目的大纲", "", 0)

	mockRepo.On("FindByID", ctx, outlineID).Return(outline, nil).Once()

	// Act
	result, err := service.GetByID(ctx, outlineID, projectID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no permission")

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_List 测试获取项目大纲列表
func TestOutlineService_List(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()

	expectedOutlines := []*writerModels.OutlineNode{
		createTestOutlineNode(projectID, "第一章", "", 0),
		createTestOutlineNode(projectID, "第二章", "", 1),
	}

	mockRepo.On("FindByProjectID", ctx, projectID).Return(expectedOutlines, nil).Once()

	// Act
	outlines, err := service.List(ctx, projectID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, outlines, 2)
	assert.Equal(t, expectedOutlines[0].Title, outlines[0].Title)
	assert.Equal(t, expectedOutlines[1].Title, outlines[1].Title)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_Update_Success 测试成功更新
func TestOutlineService_Update_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	existingOutline := createTestOutlineNode(projectID, "原标题", "", 0)

	newTitle := "新标题"
	newSummary := "新摘要"
	newTension := 8

	req := &serviceInterfaces.UpdateOutlineRequest{
		Title:     &newTitle,
		Summary:   &newSummary,
		Tension:   &newTension,
	}

	mockRepo.On("FindByID", ctx, outlineID).Return(existingOutline, nil).Once()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*writer.OutlineNode")).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	outline, err := service.Update(ctx, outlineID, projectID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, outline)
	assert.Equal(t, newTitle, outline.Title)
	assert.Equal(t, newSummary, outline.Summary)
	assert.Equal(t, newTension, outline.Tension)

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestOutlineService_Update_ParentNode 测试更新父节点
func TestOutlineService_Update_ParentNode(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	newParentID := primitive.NewObjectID().Hex()

	existingOutline := createTestOutlineNode(projectID, "章节", "", 0)
	newParent := createTestOutlineNode(projectID, "第一卷", "", 0)

	req := &serviceInterfaces.UpdateOutlineRequest{
		ParentID: &newParentID,
	}

	mockRepo.On("FindByID", ctx, outlineID).Return(existingOutline, nil).Once()
	mockRepo.On("FindByID", ctx, newParentID).Return(newParent, nil).Once()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*writer.OutlineNode")).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	outline, err := service.Update(ctx, outlineID, projectID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, outline)
	assert.Equal(t, newParentID, outline.ParentID)

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestOutlineService_Update_CrossProjectParent 测试更新为跨项目父节点失败
func TestOutlineService_Update_CrossProjectParent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()
	newParentID := primitive.NewObjectID().Hex()

	existingOutline := createTestOutlineNode(projectID, "章节", "", 0)
	newParent := createTestOutlineNode(otherProjectID, "其他项目的卷", "", 0)

	req := &serviceInterfaces.UpdateOutlineRequest{
		ParentID: &newParentID,
	}

	mockRepo.On("FindByID", ctx, outlineID).Return(existingOutline, nil).Once()
	mockRepo.On("FindByID", ctx, newParentID).Return(newParent, nil).Once()

	// Act
	outline, err := service.Update(ctx, outlineID, projectID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, outline)
	assert.Contains(t, err.Error(), "does not belong to this project")

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_Delete_Success 测试成功删除
func TestOutlineService_Delete_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	existingOutline := createTestOutlineNode(projectID, "待删除", "", 0)

	mockRepo.On("FindByID", ctx, outlineID).Return(existingOutline, nil).Once()
	mockRepo.On("CountByParentID", ctx, projectID, outlineID).Return(int64(0), nil).Once()
	mockRepo.On("Delete", ctx, outlineID).Return(nil).Once()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Once()

	// Act
	err := service.Delete(ctx, outlineID, projectID)

	// Assert
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestOutlineService_Delete_WithChildren 测试删除有子节点的节点失败
func TestOutlineService_Delete_WithChildren(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	existingOutline := createTestOutlineNode(projectID, "有子节点的父节点", "", 0)

	mockRepo.On("FindByID", ctx, outlineID).Return(existingOutline, nil).Once()
	mockRepo.On("CountByParentID", ctx, projectID, outlineID).Return(int64(2), nil).Once()

	// Act
	err := service.Delete(ctx, outlineID, projectID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete outline with children")

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_Delete_NotFound 测试删除不存在的节点
func TestOutlineService_Delete_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	outlineID := primitive.NewObjectID().Hex()

	mockRepo.On("FindByID", ctx, outlineID).Return(nil, errors.New("not found")).Once()

	// Act
	err := service.Delete(ctx, outlineID, projectID)

	// Assert
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetTree 测试获取大纲树
func TestOutlineService_GetTree(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()

	volume := createTestOutlineNode(projectID, "第一卷", "", 0)
	chapter1 := createTestOutlineNode(projectID, "第一章", volume.ID.Hex(), 0)
	chapter2 := createTestOutlineNode(projectID, "第二章", volume.ID.Hex(), 1)

	mockRepo.On("FindRoots", ctx, projectID).Return([]*writerModels.OutlineNode{volume}, nil).Once()
	mockRepo.On("FindByParentID", ctx, projectID, volume.ID.Hex()).Return([]*writerModels.OutlineNode{chapter1, chapter2}, nil).Once()
	mockRepo.On("FindByParentID", ctx, projectID, chapter1.ID.Hex()).Return([]*writerModels.OutlineNode{}, nil).Once()
	mockRepo.On("FindByParentID", ctx, projectID, chapter2.ID.Hex()).Return([]*writerModels.OutlineNode{}, nil).Once()

	// Act
	tree, err := service.GetTree(ctx, projectID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, tree, 1)
	assert.Equal(t, volume.Title, tree[0].Title)
	assert.Len(t, tree[0].Children, 2)
	assert.Equal(t, chapter1.Title, tree[0].Children[0].Title)
	assert.Equal(t, chapter2.Title, tree[0].Children[1].Title)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetChildren 测试获取子节点列表
func TestOutlineService_GetChildren(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	parentNode := createTestOutlineNode(projectID, "第一卷", "", 0)

	expectedChildren := []*writerModels.OutlineNode{
		createTestOutlineNode(projectID, "第一章", parentID, 0),
		createTestOutlineNode(projectID, "第二章", parentID, 1),
	}

	mockRepo.On("FindByID", ctx, parentID).Return(parentNode, nil).Once()
	mockRepo.On("FindByParentID", ctx, projectID, parentID).Return(expectedChildren, nil).Once()

	// Act
	children, err := service.GetChildren(ctx, projectID, parentID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, children, 2)
	assert.Equal(t, expectedChildren[0].Title, children[0].Title)
	assert.Equal(t, expectedChildren[1].Title, children[1].Title)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetChildren_RootNodes 测试获取根节点
func TestOutlineService_GetChildren_RootNodes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()

	expectedRoots := []*writerModels.OutlineNode{
		createTestOutlineNode(projectID, "第一卷", "", 0),
		createTestOutlineNode(projectID, "第二卷", "", 1),
	}

	mockRepo.On("FindRoots", ctx, projectID).Return(expectedRoots, nil).Once()

	// Act
	children, err := service.GetChildren(ctx, projectID, "")

	// Assert
	require.NoError(t, err)
	assert.Len(t, children, 2)
	assert.Equal(t, expectedRoots[0].Title, children[0].Title)
	assert.Equal(t, expectedRoots[1].Title, children[1].Title)

	mockRepo.AssertExpectations(t)
}

// TestOutlineService_GetChildren_NoPermission 测试无权限获取子节点
func TestOutlineService_GetChildren_NoPermission(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockOutlineRepository)
	mockEventBus := new(MockEventBus)

	service := NewOutlineService(mockRepo, mockEventBus)

	projectID := primitive.NewObjectID().Hex()
	otherProjectID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	parentNode := createTestOutlineNode(otherProjectID, "其他项目的卷", "", 0)

	mockRepo.On("FindByID", ctx, parentID).Return(parentNode, nil).Once()

	// Act
	children, err := service.GetChildren(ctx, projectID, parentID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, children)
	assert.Contains(t, err.Error(), "no permission")

	mockRepo.AssertExpectations(t)
}
