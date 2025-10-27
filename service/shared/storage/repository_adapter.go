package storage

import (
	"context"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// RepositoryAdapter 将 StorageRepository 适配到 FileRepository 接口
type RepositoryAdapter struct {
	repo sharedRepo.StorageRepository
}

// NewRepositoryAdapter 创建 Repository 适配器
func NewRepositoryAdapter(repo sharedRepo.StorageRepository) FileRepository {
	return &RepositoryAdapter{repo: repo}
}

// Create 创建文件记录
func (a *RepositoryAdapter) Create(ctx context.Context, file *FileInfo) error {
	return a.repo.CreateFile(ctx, file)
}

// Get 获取文件信息
func (a *RepositoryAdapter) Get(ctx context.Context, fileID string) (*FileInfo, error) {
	return a.repo.GetFile(ctx, fileID)
}

// Update 更新文件信息
func (a *RepositoryAdapter) Update(ctx context.Context, fileID string, updates map[string]interface{}) error {
	return a.repo.UpdateFile(ctx, fileID, updates)
}

// Delete 删除文件记录
func (a *RepositoryAdapter) Delete(ctx context.Context, fileID string) error {
	return a.repo.DeleteFile(ctx, fileID)
}

// List 查询文件列表
func (a *RepositoryAdapter) List(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error) {
	// 构建过滤器
	filter := &sharedRepo.FileFilter{
		UserID:   userID,
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}
	
	files, _, err := a.repo.ListFiles(ctx, filter)
	return files, err
}

// GrantAccess 授予访问权限
func (a *RepositoryAdapter) GrantAccess(ctx context.Context, fileID, userID string) error {
	return a.repo.GrantAccess(ctx, fileID, userID)
}

// RevokeAccess 撤销访问权限
func (a *RepositoryAdapter) RevokeAccess(ctx context.Context, fileID, userID string) error {
	return a.repo.RevokeAccess(ctx, fileID, userID)
}

// CheckAccess 检查访问权限
func (a *RepositoryAdapter) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	return a.repo.CheckAccess(ctx, fileID, userID)
}

