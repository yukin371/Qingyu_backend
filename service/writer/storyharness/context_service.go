package storyharness

import (
	"context"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// ContextService 章节上下文服务
// 负责组装 Context Lens 所需的角色快照、关系数据
type ContextService struct {
	characterRepo     writerRepo.CharacterRepository
	changeRequestRepo writerRepo.ChangeRequestRepository
}

// NewContextService 创建 ContextService 实例
func NewContextService(
	characterRepo writerRepo.CharacterRepository,
	changeRequestRepo writerRepo.ChangeRequestRepository,
) *ContextService {
	return &ContextService{
		characterRepo:     characterRepo,
		changeRequestRepo: changeRequestRepo,
	}
}

// ChapterContextData 章节上下文数据
type ChapterContextData struct {
	Characters []*writer.Character
	Relations  []*writer.CharacterRelation
	PendingCRs int64
}

// GetChapterContext 获取章节上下文
// 第一版使用现有角色/关系数据作为基线
func (s *ContextService) GetChapterContext(ctx context.Context, projectID, chapterID string) (*ChapterContextData, error) {
	// 1. 获取项目角色列表
	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("ContextService", errors.ServiceErrorInternal,
			"获取角色列表失败", "", err)
	}

	// 2. 获取项目关系列表
	relations, err := s.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return nil, errors.NewServiceError("ContextService", errors.ServiceErrorInternal,
			"获取关系列表失败", "", err)
	}

	// 3. 获取待处理建议计数
	pendingCount, _ := s.changeRequestRepo.CountPendingByChapter(ctx, projectID, chapterID)

	return &ChapterContextData{
		Characters: characters,
		Relations:  relations,
		PendingCRs: pendingCount,
	}, nil
}
