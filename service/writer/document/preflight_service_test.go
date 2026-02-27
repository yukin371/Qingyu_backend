package document

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/interfaces/infrastructure"
)

// MockDocumentRepository 模拟DocumentRepository
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(ctx context.Context, doc *writer.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetByIDUnscoped(ctx context.Context, id string) (*writer.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Document, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error) {
	args := m.Called(ctx, projectID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writer.Document, error) {
	args := m.Called(ctx, projectID, documentType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

func (m *MockDocumentRepository) UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, projectID, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) RestoreByProject(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error) {
	args := m.Called(ctx, documentID, projectID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDocumentRepository) SoftDelete(ctx context.Context, documentID, projectID string) error {
	args := m.Called(ctx, documentID, projectID)
	return args.Error(0)
}

func (m *MockDocumentRepository) HardDelete(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockDocumentRepository) CountByProject(ctx context.Context, projectID string) (int64, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) CreateWithTransaction(ctx context.Context, doc *writer.Document, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, doc, callback)
	return args.Error(0)
}

func (m *MockDocumentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByIDs(ctx context.Context, ids []string) ([]*writer.Document, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Document), args.Error(1)
}

// MockProjectRepository 模拟ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	args := m.Called(ctx, project)
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

func (m *MockProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
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

func (m *MockProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
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


// TestPreflightService_ValidateBatchOperation 测试批量操作验证
func TestPreflightService_ValidateBatchOperation(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef:  generateStableRef(),
		OrderKey:  "a0", // DefaultOrderKey
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.IdentifiedEntity.ID = primitive.NewObjectID()
	doc1.TouchForCreate()

	// 设置mock期望
	mockRepo.On("GetByID", ctx, doc1.ID.Hex()).Return(doc1, nil).Once()

	// 测试有效的批量操作
	summary, result, err := service.ValidateBatchOperation(
		ctx,
		projectID,
		writer.BatchOpTypeDelete,
		[]string{doc1.ID.Hex()},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicyAbort,
		},
	)

	assert.NoError(t, err, "ValidateBatchOperation should not fail")
	assert.NotNil(t, summary, "Summary should not be nil")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, summary.ValidCount, "Expected 1 valid document")
	assert.Equal(t, 1, len(result.ValidIDs), "Expected 1 valid ID")
	assert.Equal(t, 0, len(result.InvalidIDs), "Expected 0 invalid IDs")
	assert.Equal(t, doc1.ID.Hex(), result.ValidIDs[0], "Valid ID should match")
	assert.NotNil(t, result.DocumentMap[doc1.ID.Hex()], "Document should be in map")

	mockRepo.AssertExpectations(t)
}

// TestPreflightService_InvalidDocument 测试无效文档
func TestPreflightService_InvalidDocument(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 使用有效的ObjectID格式但文档不存在
	nonexistentID := primitive.NewObjectID().Hex()

	// 设置mock期望 - 文档不存在
	mockRepo.On("GetByID", ctx, nonexistentID).Return(nil, nil).Once()

	_, result, err := service.ValidateBatchOperation(
		ctx,
		projectID,
		writer.BatchOpTypeDelete,
		[]string{nonexistentID},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicyAbort,
		},
	)

	assert.Error(t, err, "Expected error for invalid document")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, len(result.InvalidIDs), "Expected 1 invalid ID")
	assert.Equal(t, "document_not_found", result.InvalidIDs[0].Code, "Expected document_not_found error code")
	assert.Equal(t, nonexistentID, result.InvalidIDs[0].ID, "Invalid ID should match")

	mockRepo.AssertExpectations(t)
}

// TestPreflightService_InvalidIDFormat 测试无效ID格式
func TestPreflightService_InvalidIDFormat(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 测试无效的ID格式（不需要调用repository）
	_, result, err := service.ValidateBatchOperation(
		ctx,
		projectID,
		writer.BatchOpTypeDelete,
		[]string{"invalid-id-format"},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicyAbort,
		},
	)

	assert.Error(t, err, "Expected error for invalid ID format")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, len(result.InvalidIDs), "Expected 1 invalid ID")
	assert.Equal(t, "invalid_id_format", result.InvalidIDs[0].Code, "Expected invalid_id_format error code")
}

// TestPreflightService_WrongProject 测试错误的项目归属
func TestPreflightService_WrongProject(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	otherProjectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建属于其他项目的文档
	doc1 := &writer.Document{
		ProjectID: otherProjectID, // 属于其他项目
		Title:     "Test Doc 1",
		StableRef: generateStableRef(),
		OrderKey:  "a0",
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.IdentifiedEntity.ID = primitive.NewObjectID()
	doc1.TouchForCreate()

	// 设置mock期望
	mockRepo.On("GetByID", ctx, doc1.ID.Hex()).Return(doc1, nil).Once()

	_, result, err := service.ValidateBatchOperation(
		ctx,
		projectID, // 尝试从当前项目操作
		writer.BatchOpTypeDelete,
		[]string{doc1.ID.Hex()},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicyAbort,
		},
	)

	assert.Error(t, err, "Expected error for wrong project")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, len(result.InvalidIDs), "Expected 1 invalid ID")
	assert.Equal(t, "wrong_project", result.InvalidIDs[0].Code, "Expected wrong_project error code")

	mockRepo.AssertExpectations(t)
}

// TestPreflightService_NormalizeTargetIDs_Deduplication 测试去重
func TestPreflightService_NormalizeTargetIDs_Deduplication(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()

	// 测试去重（不包含后代）
	result, err := service.NormalizeTargetIDs(
		ctx,
		projectID,
		[]string{"id1", "id2", "id1", "id3", "id2"}, // 有重复
		false, // 不包含后代
	)

	assert.NoError(t, err, "NormalizeTargetIDs should not fail")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 3, len(result), "Expected 3 unique IDs")
}

// TestPreflightService_NormalizeTargetIDs_Descendants 测试后代节点移除
func TestPreflightService_NormalizeTargetIDs_Descendants(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()

	// 创建父文档和子文档
	parentDoc := &writer.Document{
		ProjectID: projectID,
		Title:     "Parent Doc",
		StableRef: generateStableRef(),
		OrderKey:  "a0",
		Type:      writer.TypeVolume,
		Level:     0,
		ParentID:  primitive.NilObjectID, // 根节点
	}
	parentDoc.IdentifiedEntity.ID = primitive.NewObjectID()
	parentDoc.TouchForCreate()

	childDoc := &writer.Document{
		ProjectID: projectID,
		Title:     "Child Doc",
		StableRef: generateStableRef(),
		OrderKey:  "a0",
		Type:      writer.TypeChapter,
		Level:     1,
		ParentID:  parentDoc.ID, // 父节点是parentDoc
	}
	childDoc.IdentifiedEntity.ID = primitive.NewObjectID()
	childDoc.TouchForCreate()

	// 设置mock期望
	mockRepo.On("GetByID", ctx, parentDoc.ID.Hex()).Return(parentDoc, nil).Once()
	mockRepo.On("GetByID", ctx, childDoc.ID.Hex()).Return(childDoc, nil).Once()

	// 测试后代节点移除（包含父子节点）
	result, err := service.NormalizeTargetIDs(
		ctx,
		projectID,
		[]string{parentDoc.ID.Hex(), childDoc.ID.Hex()}, // 包含父子节点
		true, // 包含后代
	)

	assert.NoError(t, err, "NormalizeTargetIDs should not fail")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, len(result), "Expected only 1 ID (parent), child should be removed")
	assert.Equal(t, parentDoc.ID.Hex(), result[0], "Expected only parent ID")

	mockRepo.AssertExpectations(t)
}

// TestPreflightService_MixedValidInvalid 测试混合有效和无效ID
func TestPreflightService_MixedValidInvalid(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建一个有效文档
	validDoc := &writer.Document{
		ProjectID: projectID,
		Title:     "Valid Doc",
		StableRef: generateStableRef(),
		OrderKey:  "a0",
		Type:      writer.TypeChapter,
		Level:     0,
	}
	validDoc.IdentifiedEntity.ID = primitive.NewObjectID()
	validDoc.TouchForCreate()

	// 使用有效的ObjectID格式但文档不存在
	nonexistentID := primitive.NewObjectID().Hex()

	// 设置mock期望
	mockRepo.On("GetByID", ctx, validDoc.ID.Hex()).Return(validDoc, nil).Once()
	mockRepo.On("GetByID", ctx, nonexistentID).Return(nil, nil).Once()

	// 测试混合有效和无效ID
	summary, result, err := service.ValidateBatchOperation(
		ctx,
		projectID,
		writer.BatchOpTypeDelete,
		[]string{validDoc.ID.Hex(), nonexistentID, "invalid-format"},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicyAbort,
		},
	)

	assert.Error(t, err, "Expected error due to invalid IDs with abort policy")
	assert.NotNil(t, summary, "Summary should not be nil")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 3, summary.TotalCount, "Total count should be 3")
	assert.Equal(t, 1, summary.ValidCount, "Expected 1 valid ID")
	assert.Equal(t, 2, summary.InvalidCount, "Expected 2 invalid IDs")
	assert.Equal(t, 1, len(result.ValidIDs), "Expected 1 valid ID in result")
	assert.Equal(t, 2, len(result.InvalidIDs), "Expected 2 invalid IDs in result")

	mockRepo.AssertExpectations(t)
}

// TestPreflightService_ContinueOnInvalid 测试遇到无效ID继续执行
func TestPreflightService_ContinueOnInvalid(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDocumentRepository)
	service := NewPreflightService(mockRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建一个有效文档
	validDoc := &writer.Document{
		ProjectID: projectID,
		Title:     "Valid Doc",
		StableRef: generateStableRef(),
		OrderKey:  "a0",
		Type:      writer.TypeChapter,
		Level:     0,
	}
	validDoc.IdentifiedEntity.ID = primitive.NewObjectID()
	validDoc.TouchForCreate()

	// 使用有效的ObjectID格式但文档不存在
	nonexistentID := primitive.NewObjectID().Hex()

	// 设置mock期望
	mockRepo.On("GetByID", ctx, validDoc.ID.Hex()).Return(validDoc, nil).Once()
	mockRepo.On("GetByID", ctx, nonexistentID).Return(nil, nil).Once()

	// 测试使用Skip策略遇到无效ID继续执行
	summary, result, err := service.ValidateBatchOperation(
		ctx,
		projectID,
		writer.BatchOpTypeDelete,
		[]string{validDoc.ID.Hex(), nonexistentID},
		&PreflightOptions{
			UserID:         userID,
			ConflictPolicy: writer.ConflictPolicySkip, // 跳过策略
		},
	)

	assert.NoError(t, err, "Should not error with skip policy")
	assert.NotNil(t, summary, "Summary should not be nil")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 1, summary.ValidCount, "Expected 1 valid ID")
	assert.Equal(t, 1, summary.InvalidCount, "Expected 1 invalid ID")

	mockRepo.AssertExpectations(t)
}

// generateStableRef 生成稳定的引用标识
func generateStableRef() string {
	// TODO: 使用ULID库生成
	return primitive.NewObjectID().Hex()
}
