package document

import (
	"context"
	"errors"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	pkgErrors "Qingyu_backend/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ==================== Test Helper ====================

func createTestProject(ownerID string) *writer.Project {
	return &writer.Project{
		OwnedEntity: base.OwnedEntity{
			AuthorID: ownerID,
		},
		Visibility: writer.VisibilityPrivate,
	}
}

func createTestProjectWithCollab(ownerID string, collabUserIDs []string) *writer.Project {
	project := &writer.Project{
		OwnedEntity: base.OwnedEntity{
			AuthorID: ownerID,
		},
		Visibility:    writer.VisibilityPrivate,
		Collaborators: make([]writer.Collaborator, 0),
	}
	for _, uid := range collabUserIDs {
		acceptedAt := getTimePtr()
		project.Collaborators = append(project.Collaborators, writer.Collaborator{
			UserID:     uid,
			Role:       writer.RoleEditor,
			AcceptedAt: acceptedAt,
		})
	}
	return project
}

func getTimePtr() *time.Time {
	now := time.Now()
	return &now
}

// ==================== Tests ====================

func TestVerifyProjectEdit_Success(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "user123")

	testProject := createTestProject("user123")
	mockProjectRepo.On("GetByID", ctx, "project123").Return(testProject, nil)

	// Act
	userID, project, err := helper.VerifyProjectEdit(ctx, "project123")

	// Assert
	if err != nil {
		t.Errorf("期望无错误，实际得到: %v", err)
	}
	if userID != "user123" {
		t.Errorf("期望 userID=user123，实际得到: %s", userID)
	}
	if project == nil {
		t.Error("期望返回项目，实际为 nil")
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestVerifyProjectEdit_Unauthorized(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.Background() // 无 userID

	// Act
	_, _, err := helper.VerifyProjectEdit(ctx, "project123")

	// Assert
	if err == nil {
		t.Error("期望返回未授权错误，实际无错误")
	}
	var svcErr *pkgErrors.ServiceError
	if errors.As(err, &svcErr) {
		if svcErr.Type != pkgErrors.ServiceErrorUnauthorized {
			t.Errorf("期望错误类型=%s，实际=%s", pkgErrors.ServiceErrorUnauthorized, svcErr.Type)
		}
	} else {
		t.Errorf("期望 ServiceError 类型，实际: %T", err)
	}
}

func TestVerifyProjectEdit_Forbidden(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "otherUser") // 当前用户是 otherUser

	testProject := createTestProject("owner123") // 项目所有者是 owner123
	mockProjectRepo.On("GetByID", ctx, "project123").Return(testProject, nil)

	// Act
	_, _, err := helper.VerifyProjectEdit(ctx, "project123")

	// Assert
	if err == nil {
		t.Error("期望返回禁止访问错误，实际无错误")
	}
	var svcErr *pkgErrors.ServiceError
	if errors.As(err, &svcErr) {
		if svcErr.Type != pkgErrors.ServiceErrorForbidden {
			t.Errorf("期望错误类型=%s，实际=%s", pkgErrors.ServiceErrorForbidden, svcErr.Type)
		}
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestVerifyProjectEdit_ProjectNotFound(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "user123")

	mockProjectRepo.On("GetByID", ctx, "project123").Return(nil, nil) // 项目不存在

	// Act
	_, _, err := helper.VerifyProjectEdit(ctx, "project123")

	// Assert
	if err == nil {
		t.Error("期望返回项目不存在错误，实际无错误")
	}
	var svcErr *pkgErrors.ServiceError
	if errors.As(err, &svcErr) {
		if svcErr.Type != pkgErrors.ServiceErrorNotFound {
			t.Errorf("期望错误类型=%s，实际=%s", pkgErrors.ServiceErrorNotFound, svcErr.Type)
		}
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestVerifyProjectEdit_CollaboratorCanEdit(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "collabUser")

	testProject := createTestProjectWithCollab("owner123", []string{"collabUser"})
	mockProjectRepo.On("GetByID", ctx, "project123").Return(testProject, nil)

	// Act
	userID, project, err := helper.VerifyProjectEdit(ctx, "project123")

	// Assert
	if err != nil {
		t.Errorf("协作者应该有编辑权限，实际得到错误: %v", err)
	}
	if userID != "collabUser" {
		t.Errorf("期望 userID=collabUser，实际得到: %s", userID)
	}
	if project == nil {
		t.Error("期望返回项目，实际为 nil")
	}
	mockProjectRepo.AssertExpectations(t)
}

func TestVerifyDocumentEdit_Success(t *testing.T) {
	// Arrange
	projectID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "user123")

	testProject := createTestProject("user123")
	testDoc := &writer.Document{
		ProjectID: projectID,
	}

	mockDocRepo.On("GetByID", ctx, "doc123").Return(testDoc, nil)
	mockProjectRepo.On("GetByID", ctx, "507f1f77bcf86cd799439011").Return(testProject, nil)

	// Act
	userID, doc, project, err := helper.VerifyDocumentEdit(ctx, "doc123")

	// Assert
	if err != nil {
		t.Errorf("期望无错误，实际得到: %v", err)
	}
	if userID != "user123" {
		t.Errorf("期望 userID=user123，实际得到: %s", userID)
	}
	if doc == nil {
		t.Error("期望返回文档，实际为 nil")
	}
	if project == nil {
		t.Error("期望返回项目，实际为 nil")
	}
	mockDocRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestVerifyDocumentEdit_DocumentNotFound(t *testing.T) {
	// Arrange
	mockProjectRepo := new(MockProjectRepository)
	mockDocRepo := new(MockDocumentRepository)
	helper := NewAuthHelper(mockProjectRepo, mockDocRepo, "TestService")
	ctx := context.WithValue(context.Background(), "userID", "user123")

	mockDocRepo.On("GetByID", ctx, "doc123").Return(nil, nil) // 文档不存在

	// Act
	_, _, _, err := helper.VerifyDocumentEdit(ctx, "doc123")

	// Assert
	if err == nil {
		t.Error("期望返回文档不存在错误，实际无错误")
	}
	var svcErr *pkgErrors.ServiceError
	if errors.As(err, &svcErr) {
		if svcErr.Type != pkgErrors.ServiceErrorNotFound {
			t.Errorf("期望错误类型=%s，实际=%s", pkgErrors.ServiceErrorNotFound, svcErr.Type)
		}
	}
	mockDocRepo.AssertExpectations(t)
}

func TestGetUserID_Success(t *testing.T) {
	// Arrange
	ctx := context.WithValue(context.Background(), "userID", "user123")

	// Act
	userID, ok := GetUserID(ctx)

	// Assert
	if !ok {
		t.Error("期望返回 true，实际返回 false")
	}
	if userID != "user123" {
		t.Errorf("期望 userID=user123，实际得到: %s", userID)
	}
}

func TestGetUserID_Empty(t *testing.T) {
	// Arrange
	ctx := context.WithValue(context.Background(), "userID", "")

	// Act
	_, ok := GetUserID(ctx)

	// Assert
	if ok {
		t.Error("期望返回 false（空字符串），实际返回 true")
	}
}

func TestGetUserID_NotPresent(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	_, ok := GetUserID(ctx)

	// Assert
	if ok {
		t.Error("期望返回 false（无 userID），实际返回 true")
	}
}

func TestMustGetUserID_Success(t *testing.T) {
	// Arrange
	ctx := context.WithValue(context.Background(), "userID", "user123")

	// Act
	userID, err := MustGetUserID(ctx, "TestService")

	// Assert
	if err != nil {
		t.Errorf("期望无错误，实际得到: %v", err)
	}
	if userID != "user123" {
		t.Errorf("期望 userID=user123，实际得到: %s", userID)
	}
}

func TestMustGetUserID_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	_, err := MustGetUserID(ctx, "TestService")

	// Assert
	if err == nil {
		t.Error("期望返回错误，实际无错误")
	}
}
