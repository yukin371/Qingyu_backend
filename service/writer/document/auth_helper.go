package document

import (
	"context"

	"Qingyu_backend/models/writer"
	pkgErrors "Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// AuthHelper 权限验证辅助结构
type AuthHelper struct {
	projectRepo  writerRepo.ProjectRepository
	documentRepo writerRepo.DocumentRepository
	serviceName  string
}

// NewAuthHelper 创建权限验证辅助实例
func NewAuthHelper(
	projectRepo writerRepo.ProjectRepository,
	documentRepo writerRepo.DocumentRepository,
	serviceName string,
) *AuthHelper {
	return &AuthHelper{
		projectRepo:  projectRepo,
		documentRepo: documentRepo,
		serviceName:  serviceName,
	}
}

// VerifyProjectEdit 验证用户是否有项目编辑权限
// 返回: userID, project, error
func (h *AuthHelper) VerifyProjectEdit(ctx context.Context, projectID string) (string, *writer.Project, error) {
	// 1. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorUnauthorized,
			"用户未登录",
			"",
			nil,
		)
	}

	// 2. 查询项目
	project, err := h.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询项目失败",
			"",
			err,
		)
	}

	if project == nil {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"项目不存在",
			"",
			nil,
		)
	}

	// 3. 验证编辑权限
	if !project.CanEdit(userID) {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorForbidden,
			"无权限编辑该项目",
			"",
			nil,
		)
	}

	return userID, project, nil
}

// VerifyProjectView 验证用户是否有项目查看权限
// 返回: userID, project, error
func (h *AuthHelper) VerifyProjectView(ctx context.Context, projectID string) (string, *writer.Project, error) {
	// 1. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorUnauthorized,
			"用户未登录",
			"",
			nil,
		)
	}

	// 2. 查询项目
	project, err := h.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询项目失败",
			"",
			err,
		)
	}

	if project == nil {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"项目不存在",
			"",
			nil,
		)
	}

	// 3. 验证查看权限
	if !project.CanView(userID) {
		return "", nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorForbidden,
			"无权限查看该项目",
			"",
			nil,
		)
	}

	return userID, project, nil
}

// VerifyDocumentEdit 验证用户是否有文档编辑权限（通过项目间接验证）
// 返回: userID, document, project, error
func (h *AuthHelper) VerifyDocumentEdit(ctx context.Context, documentID string) (string, *writer.Document, *writer.Project, error) {
	// 1. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorUnauthorized,
			"用户未登录",
			"",
			nil,
		)
	}

	// 2. 查询文档
	document, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询文档失败",
			"",
			err,
		)
	}

	if document == nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"文档不存在",
			"",
			nil,
		)
	}

	// 3. 查询项目并验证编辑权限
	project, err := h.projectRepo.GetByID(ctx, document.ProjectID.Hex())
	if err != nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询项目失败",
			"",
			err,
		)
	}

	if project == nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"文档所属项目不存在",
			"",
			nil,
		)
	}

	if !project.CanEdit(userID) {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorForbidden,
			"无权限编辑该文档",
			"",
			nil,
		)
	}

	return userID, document, project, nil
}

// VerifyDocumentView 验证用户是否有文档查看权限（通过项目间接验证）
// 返回: userID, document, project, error
func (h *AuthHelper) VerifyDocumentView(ctx context.Context, documentID string) (string, *writer.Document, *writer.Project, error) {
	// 1. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorUnauthorized,
			"用户未登录",
			"",
			nil,
		)
	}

	// 2. 查询文档
	document, err := h.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询文档失败",
			"",
			err,
		)
	}

	if document == nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"文档不存在",
			"",
			nil,
		)
	}

	// 3. 查询项目并验证查看权限
	project, err := h.projectRepo.GetByID(ctx, document.ProjectID.Hex())
	if err != nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorInternal,
			"查询项目失败",
			"",
			err,
		)
	}

	if project == nil {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorNotFound,
			"文档所属项目不存在",
			"",
			nil,
		)
	}

	if !project.CanView(userID) {
		return "", nil, nil, pkgErrors.NewServiceError(
			h.serviceName,
			pkgErrors.ServiceErrorForbidden,
			"无权限查看该文档",
			"",
			nil,
		)
	}

	return userID, document, project, nil
}

// GetUserID 从上下文中获取用户ID
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok && userID != ""
}

// MustGetUserID 从上下文中获取用户ID，如果不存在则返回错误
func MustGetUserID(ctx context.Context, serviceName string) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", pkgErrors.NewServiceError(
			serviceName,
			pkgErrors.ServiceErrorUnauthorized,
			"用户未登录",
			"",
			nil,
		)
	}
	return userID, nil
}
