package storyharness

import (
	"context"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// ChangeRequestService 变更建议服务
// 负责建议的 CRUD 和状态流转
type ChangeRequestService struct {
	crRepo writerRepo.ChangeRequestRepository
}

// NewChangeRequestService 创建 ChangeRequestService 实例
func NewChangeRequestService(crRepo writerRepo.ChangeRequestRepository) *ChangeRequestService {
	return &ChangeRequestService{crRepo: crRepo}
}

// ListByChapter 获取章节下的建议列表
func (s *ChangeRequestService) ListByChapter(ctx context.Context, projectID, chapterID string) ([]*writer.ChangeRequest, error) {
	return s.crRepo.FindPendingByChapter(ctx, projectID, chapterID)
}

// GetByID 获取单条建议
func (s *ChangeRequestService) GetByID(ctx context.Context, id string) (*writer.ChangeRequest, error) {
	return s.crRepo.FindRequestByID(ctx, id)
}

// Process 处理建议（接受/忽略/延后）
func (s *ChangeRequestService) Process(ctx context.Context, requestID string, newStatus writer.ChangeRequestStatus, processedBy string) error {
	// 验证状态合法性
	switch newStatus {
	case writer.CRStatusAccepted, writer.CRStatusIgnored, writer.CRStatusDeferred:
		// valid
	default:
		return errors.NewServiceError("ChangeRequestService", errors.ServiceErrorValidation,
			"无效的处理状态", string(newStatus), nil)
	}

	return s.crRepo.UpdateRequestStatus(ctx, requestID, newStatus, processedBy)
}

// CountPending 获取章节待处理建议数
func (s *ChangeRequestService) CountPending(ctx context.Context, projectID, chapterID string) (int64, error) {
	return s.crRepo.CountPendingByChapter(ctx, projectID, chapterID)
}
