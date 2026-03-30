package writer

import (
	"context"

	writerModel "Qingyu_backend/models/writer"
	writerInterface "Qingyu_backend/repository/interfaces/writer"
)

// ============================================================================
// ExportService 适配器 - 将 Repository 接口适配到 ExportService 内部接口
// ============================================================================

// documentRepoAdapter 将 DocumentRepository 适配到 ExportService 的 DocumentRepository 接口
type documentRepoAdapter struct {
	repo writerInterface.DocumentRepository
}

// FindByID 适配 GetByID -> FindByID
func (a *documentRepoAdapter) FindByID(ctx context.Context, id string) (*writerModel.Document, error) {
	return a.repo.GetByID(ctx, id)
}

// FindByProjectID 适配 GetByProjectID -> FindByProjectID
func (a *documentRepoAdapter) FindByProjectID(ctx context.Context, projectID string) ([]*writerModel.Document, error) {
	docs, err := a.repo.GetByProjectID(ctx, projectID, 0, 0) // 0, 0 表示不分页
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// documentContentRepoAdapter 将 DocumentContentRepository 适配到 ExportService 的 DocumentContentRepository 接口
type documentContentRepoAdapter struct {
	repo writerInterface.DocumentContentRepository
}

// FindByID 直接委托给 GetByID
func (a *documentContentRepoAdapter) FindByID(ctx context.Context, id string) (*writerModel.DocumentContent, error) {
	return a.repo.GetByID(ctx, id)
}

// projectRepoAdapter 将 ProjectRepository 适配到 ExportService 的 ProjectRepository 接口
type projectRepoAdapter struct {
	repo writerInterface.ProjectRepository
}

// FindByID 适配 GetByID -> FindByID
func (a *projectRepoAdapter) FindByID(ctx context.Context, id string) (*writerModel.Project, error) {
	return a.repo.GetByID(ctx, id)
}

// NewDocumentRepoAdapter 创建 DocumentRepository 适配器
func NewDocumentRepoAdapter(repo writerInterface.DocumentRepository) *documentRepoAdapter {
	return &documentRepoAdapter{repo: repo}
}

// NewDocumentContentRepoAdapter 创建 DocumentContentRepository 适配器
func NewDocumentContentRepoAdapter(repo writerInterface.DocumentContentRepository) *documentContentRepoAdapter {
	return &documentContentRepoAdapter{repo: repo}
}

// NewProjectRepoAdapter 创建 ProjectRepository 适配器
func NewProjectRepoAdapter(repo writerInterface.ProjectRepository) *projectRepoAdapter {
	return &projectRepoAdapter{repo: repo}
}
