package document

import (
	"context"
	"errors"
	"fmt"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository"
	writerInterface "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidTargetID        = errors.New("invalid target document ID")
	ErrDocumentNotFound       = errors.New("document not found")
	ErrVersionConflict        = errors.New("document version conflict")
	ErrInsufficientPermission = errors.New("insufficient permission")
)

// PreflightService 预检查服务接口
type PreflightService interface {
	// ValidateBatchOperation 验证批量操作
	ValidateBatchOperation(
		ctx context.Context,
		projectID primitive.ObjectID,
		opType writer.BatchOperationType,
		targetIDs []string,
		options *PreflightOptions,
	) (*writer.PreflightSummary, *PreflightResult, error)

	// NormalizeTargetIDs 规范化目标ID（去重、移除后代节点）
	NormalizeTargetIDs(
		ctx context.Context,
		projectID primitive.ObjectID,
		targetIDs []string,
		includeDescendants bool,
	) ([]string, error)
}

// PreflightOptions 预检查选项
type PreflightOptions struct {
	ExpectedVersions   map[string]int        // 版本检查
	ConflictPolicy     writer.ConflictPolicy // 冲突策略
	IncludeDescendants bool                  // 是否包含后代节点
	UserID             primitive.ObjectID    // 当前用户ID
}

// PreflightResult 预检查结果
type PreflightResult struct {
	ValidIDs    []string                    // 有效ID列表
	InvalidIDs  []InvalidTarget             // 无效ID列表
	SkippedIDs  []string                    // 跳过的ID列表
	Warnings    []string                    // 警告信息
	DocumentMap map[string]*writer.Document // ID -> Document 映射
}

// InvalidTarget 无效目标
type InvalidTarget struct {
	ID     string // 目标ID
	Reason string // 失败原因
	Code   string // 错误代码
}

// PreflightServiceImpl 预检查服务实现
type PreflightServiceImpl struct {
	docRepo     writerInterface.DocumentRepository
	contentRepo writerInterface.DocumentContentRepository
}

// NewPreflightService 创建预检查服务
func NewPreflightService(
	docRepo writerInterface.DocumentRepository,
	contentRepo writerInterface.DocumentContentRepository,
) PreflightService {
	return &PreflightServiceImpl{
		docRepo:     docRepo,
		contentRepo: contentRepo,
	}
}

// ValidateBatchOperation 验证批量操作
func (s *PreflightServiceImpl) ValidateBatchOperation(
	ctx context.Context,
	projectID primitive.ObjectID,
	opType writer.BatchOperationType,
	targetIDs []string,
	options *PreflightOptions,
) (*writer.PreflightSummary, *PreflightResult, error) {

	if options == nil {
		options = &PreflightOptions{}
	}
	hasVersionConflict := false

	result := &PreflightResult{
		ValidIDs:    make([]string, 0),
		InvalidIDs:  make([]InvalidTarget, 0),
		SkippedIDs:  make([]string, 0),
		Warnings:    make([]string, 0),
		DocumentMap: make(map[string]*writer.Document),
	}

	// 1. 规范化目标ID（去重、移除后代）
	normalizedIDs, err := s.NormalizeTargetIDs(ctx, projectID, targetIDs, options.IncludeDescendants)
	if err != nil {
		return nil, nil, err
	}

	// 2. 验证每个目标ID
	for _, id := range normalizedIDs {
		// 将字符串ID转换为ObjectID验证格式
		_, err := repository.ParseID(id)
		if err != nil {
			result.InvalidIDs = append(result.InvalidIDs, InvalidTarget{
				ID:     id,
				Reason: "invalid document ID format",
				Code:   "invalid_id_format",
			})
			continue
		}

		// 查询文档
		doc, err := s.docRepo.GetByID(ctx, id)
		if err != nil || doc == nil {
			result.InvalidIDs = append(result.InvalidIDs, InvalidTarget{
				ID:     id,
				Reason: "document not found",
				Code:   "document_not_found",
			})
			continue
		}

		// 验证项目归属
		if doc.ProjectID != projectID {
			result.InvalidIDs = append(result.InvalidIDs, InvalidTarget{
				ID:     id,
				Reason: "document does not belong to this project",
				Code:   "wrong_project",
			})
			continue
		}

		// 版本检查
		if options.ExpectedVersions != nil {
			if expectedVersion, ok := options.ExpectedVersions[id]; ok {
				if s.contentRepo == nil {
					return nil, nil, fmt.Errorf("expectedVersions 预检查依赖 DocumentContentRepository")
				}

				currentVersion := 1
				content, err := s.contentRepo.GetByDocumentID(ctx, id)
				if err != nil {
					return nil, nil, fmt.Errorf("查询文档 %s 版本失败: %w", id, err)
				}
				if content != nil && content.Version > 0 {
					currentVersion = content.Version
				}

				if currentVersion != expectedVersion {
					switch options.ConflictPolicy {
					case writer.ConflictPolicyOverwrite:
						result.Warnings = append(result.Warnings,
							fmt.Sprintf("Document %s has version conflict, will overwrite (expected=%d current=%d)",
								id, expectedVersion, currentVersion))
					case writer.ConflictPolicySkip:
						result.SkippedIDs = append(result.SkippedIDs, id)
						result.Warnings = append(result.Warnings,
							fmt.Sprintf("Document %s has version conflict, skipped (expected=%d current=%d)",
								id, expectedVersion, currentVersion))
						continue
					default:
						hasVersionConflict = true
						result.InvalidIDs = append(result.InvalidIDs, InvalidTarget{
							ID:     id,
							Reason: fmt.Sprintf("version conflict: expected=%d current=%d", expectedVersion, currentVersion),
							Code:   "version_conflict",
						})
						continue
					}
				}
			}
		}

		// 权限检查框架（TODO: 实现具体权限验证逻辑）
		// if !s.hasPermission(ctx, doc, options.UserID, opType) {
		// 	result.InvalidIDs = append(result.InvalidIDs, InvalidTarget{
		// 		ID:     id,
		// 		Reason: "insufficient permission",
		// 		Code:   "permission_denied",
		// 	})
		// 	continue
		// }

		// 所有验证通过
		result.ValidIDs = append(result.ValidIDs, id)
		result.DocumentMap[id] = doc
	}

	// 3. 构建摘要
	summary := &writer.PreflightSummary{
		TotalCount:   len(targetIDs),
		ValidCount:   len(result.ValidIDs),
		InvalidCount: len(result.InvalidIDs),
		SkippedCount: len(result.SkippedIDs),
	}

	// 4. 如果要求原子操作且有无效ID，返回错误
	if len(result.InvalidIDs) > 0 && options.ConflictPolicy == writer.ConflictPolicyAbort {
		if hasVersionConflict {
			return summary, result, ErrVersionConflict
		}
		return summary, result, ErrInvalidTargetID
	}

	return summary, result, nil
}

// NormalizeTargetIDs 规范化目标ID（去重、移除后代节点）
func (s *PreflightServiceImpl) NormalizeTargetIDs(
	ctx context.Context,
	projectID primitive.ObjectID,
	targetIDs []string,
	includeDescendants bool,
) ([]string, error) {

	// 1. 去重
	uniqueIDs := make(map[string]bool)
	for _, id := range targetIDs {
		uniqueIDs[id] = true
	}

	// 2. 如果不包含后代，直接返回
	if !includeDescendants {
		result := make([]string, 0, len(uniqueIDs))
		for id := range uniqueIDs {
			result = append(result, id)
		}
		return result, nil
	}

	// 3. 如果包含后代，需要移除父节点在列表中的后代节点
	// 这是一个简化版本，实际实现需要查询文档树结构
	ancestorSet := make(map[string]bool)

	for id := range uniqueIDs {
		// 验证ID格式
		_, err := repository.ParseID(id)
		if err != nil {
			continue
		}

		// 查询文档
		doc, err := s.docRepo.GetByID(ctx, id)
		if err != nil || doc == nil {
			continue
		}

		// 如果是根节点或没有父节点在列表中，保留
		isDescendant := false
		if !doc.ParentID.IsZero() {
			parentHex := doc.ParentID.Hex()
			if uniqueIDs[parentHex] {
				isDescendant = true
			}
		}

		if !isDescendant {
			ancestorSet[id] = true
		}
	}

	result := make([]string, 0, len(ancestorSet))
	for id := range ancestorSet {
		result = append(result, id)
	}

	return result, nil
}
