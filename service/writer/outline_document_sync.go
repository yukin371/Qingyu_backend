package writer

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/distlock"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OutlineDocumentSyncService 大纲-文档双向同步服务
type OutlineDocumentSyncService struct {
	outlineRepo  writerRepo.OutlineRepository
	documentRepo writerRepo.DocumentRepository
	projectRepo  writerRepo.ProjectRepository
	outlineSvc   *OutlineService
	distLock     *distlock.RedisLockService
}

// NewOutlineDocumentSyncService 创建同步服务
func NewOutlineDocumentSyncService(
	outlineRepo writerRepo.OutlineRepository,
	documentRepo writerRepo.DocumentRepository,
	projectRepo writerRepo.ProjectRepository,
	outlineSvc *OutlineService,
	distLock *distlock.RedisLockService,
) *OutlineDocumentSyncService {
	return &OutlineDocumentSyncService{
		outlineRepo:  outlineRepo,
		documentRepo: documentRepo,
		projectRepo:  projectRepo,
		outlineSvc:   outlineSvc,
		distLock:     distLock,
	}
}

// getOutlineLevel 获取大纲节点的层级深度
// Root (ParentID=="") = 0, Level 1 = 1, Level 2 = 2, etc.
func (s *OutlineDocumentSyncService) getOutlineLevel(ctx context.Context, node *writer.OutlineNode) int {
	if node.ParentID == "" {
		return 0
	}
	parent, err := s.outlineRepo.FindByID(ctx, node.ParentID)
	if err != nil || parent == nil {
		return 0
	}
	return s.getOutlineLevel(ctx, parent) + 1
}

// outlineLevelToDocType 映射大纲层级到文档类型
// Level 0: 全局总纲（不映射到文档）
// Level 1: 卷（volume）→ arc 类型
// Level 2: 章节（chapter）→ scene 类型
func outlineLevelToDocType(level int) string {
	switch level {
	case 1:
		return "volume"
	case 2:
		return "chapter"
	default:
		return ""
	}
}

// findOrCreateGlobalOutline 查找或创建项目的全局总纲节点
// 全局总纲不关联任何文档，作为所有卷的父节点
// 使用 Redis 分布式锁保证并发安全
func (s *OutlineDocumentSyncService) findOrCreateGlobalOutline(ctx context.Context, projectID string) (string, error) {
	// 如果没有分布式锁服务，降级为纯数据库原子操作
	if s.distLock == nil {
		return s.findOrCreateGlobalOutlineWithUpsert(ctx, projectID)
	}

	// 使用分布式锁保证并发安全
	lockKey := fmt.Sprintf("global_outline:%s", projectID)
	lockID, err := s.distLock.AcquireWithRetry(ctx, lockKey, 5*time.Second, 3, 500*time.Millisecond)
	if err != nil {
		return s.findOrCreateGlobalOutlineWithUpsert(ctx, projectID)
	}
	defer func() {
		_ = s.distLock.Release(ctx, lockKey, lockID)
	}()

	// 临界区内使用原子性 upsert
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	existing, err := s.outlineRepo.FindByGlobalOutline(ctx, projectOID)
	if err != nil {
		return "", err
	}

	if existing != nil {
		return existing.ID.Hex(), nil
	}

	return "", fmt.Errorf("findOrCreateGlobalOutline 失败: FindByGlobalOutline 返回 nil")
}

// findOrCreateGlobalOutlineWithUpsert 使用数据库 upsert 原子操作查找或创建全局总纲（无锁降级方案）
func (s *OutlineDocumentSyncService) findOrCreateGlobalOutlineWithUpsert(ctx context.Context, projectID string) (string, error) {
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	existing, err := s.outlineRepo.FindByGlobalOutline(ctx, projectOID)
	if err != nil {
		return "", err
	}

	if existing != nil {
		return existing.ID.Hex(), nil
	}

	// 如果返回 nil 但没有错误，说明 FindByGlobalOutline 没有正确实现 upsert
	// 这不应该发生，但如果发生了，我们尝试最后一次查找
	existing, err = s.outlineRepo.FindByGlobalOutline(ctx, projectOID)
	if err == nil && existing != nil {
		return existing.ID.Hex(), nil
	}

	return "", fmt.Errorf("findOrCreateGlobalOutline 失败: FindByGlobalOutline 返回 nil")
}

// docTypeToOutlineParentID 获取文档对应大纲的父节点ID
// volume → global outline, chapter → volume's outline
func (s *OutlineDocumentSyncService) docTypeToOutlineParentID(ctx context.Context, doc *writer.Document, projectID string) string {
	if doc.Type == "volume" {
		// 卷的父节点是全局总纲
		globalOutlineID, err := s.findOrCreateGlobalOutline(ctx, projectID)
		if err != nil {
			return ""
		}
		return globalOutlineID
	}
	if doc.Type == "chapter" && !doc.ParentID.IsZero() {
		// 找到父文档（卷）对应的 outline 节点
		parentDoc, err := s.documentRepo.GetByID(ctx, doc.ParentID.Hex())
		if err == nil && parentDoc != nil && parentDoc.OutlineNodeID != "" {
			return parentDoc.OutlineNodeID
		}
	}
	return ""
}

// SyncFromOutlineCreation 大纲创建后自动创建对应文档并建立双向引用
func (s *OutlineDocumentSyncService) SyncFromOutlineCreation(
	ctx context.Context,
	projectID string,
	outlineNode *writer.OutlineNode,
) (*writer.Document, error) {
	// 如果已经关联了文档，跳过
	if outlineNode.DocumentID != "" {
		return nil, nil
	}

	// 判断层级
	level := s.getOutlineLevel(ctx, outlineNode)
	docType := outlineLevelToDocType(level)
	if docType == "" {
		return nil, nil // 不需要映射的层级
	}

	// 计算文档层级
	docLevel := level // Level 1 outline → Level 1 doc, Level 2 → Level 2

	// 创建文档
	doc := &writer.Document{
		ProjectID:     outlineNode.ProjectID,
		Title:         outlineNode.Title,
		Type:          docType,
		Level:         docLevel,
		Order:         outlineNode.Order,
		Status:        writer.DocumentStatusPlanned,
		OutlineNodeID: outlineNode.ID.Hex(),
		StableRef:     outlineNode.ID.Hex(), // 使用大纲ID作为稳定引用
		OrderKey:      fmt.Sprintf("%04d", outlineNode.Order*1000),
	}

	// 持久化文档到数据库
	if err := s.documentRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	// 回写大纲节点的 document_id 引用
	outlineNode.DocumentID = doc.ID.Hex()
	_ = s.outlineRepo.Update(ctx, outlineNode)
	return doc, nil
}

// SyncFromDocumentCreation 文档创建后自动创建对应大纲节点并建立双向引用
func (s *OutlineDocumentSyncService) SyncFromDocumentCreation(
	ctx context.Context,
	projectID string,
	doc *writer.Document,
) (*writer.OutlineNode, error) {
	// 如果已经关联了大纲，跳过
	if doc.OutlineNodeID != "" {
		return nil, nil
	}

	// 检查是否已存在相同 document_id 的大纲节点（防止前端重复创建）
	existingOutline, err := s.outlineRepo.FindByDocumentID(ctx, doc.ID.Hex())
	if err == nil && existingOutline != nil {
		// 已存在大纲节点，只需回写文档引用
		updates := map[string]interface{}{
			"outline_node_id": existingOutline.ID.Hex(),
		}
		if updateErr := s.documentRepo.Update(ctx, doc.ID.Hex(), updates); updateErr != nil {
			return nil, updateErr
		}
		return existingOutline, nil
	}

	// 只有 volume 和 chapter 需要映射
	outlineParentID := ""
	outlineType := ""

	switch doc.Type {
	case "volume":
		outlineParentID = s.docTypeToOutlineParentID(ctx, doc, projectID)
		outlineType = "arc"
	case "chapter":
		outlineParentID = s.docTypeToOutlineParentID(ctx, doc, projectID)
		outlineType = "scene"
	default:
		return nil, nil
	}

	// 构建大纲节点
	outlineNode := &writer.OutlineNode{}
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	outlineNode.ProjectID = projectOID
	outlineNode.Title = doc.Title
	outlineNode.ParentID = outlineParentID
	outlineNode.Order = doc.Order
	outlineNode.Type = outlineType
	outlineNode.DocumentID = doc.ID.Hex()
	outlineNode.Tags = []string{}

	// 持久化大纲节点到数据库
	if err := s.outlineRepo.Create(ctx, outlineNode); err != nil {
		return nil, err
	}

	// 回写文档的 outline_node_id 引用
	updates := map[string]interface{}{
		"outline_node_id": outlineNode.ID.Hex(),
	}
	_ = s.documentRepo.Update(ctx, doc.ID.Hex(), updates)
	return outlineNode, nil
}

// SyncTitleToDocument 大纲标题变更同步到文档
func (s *OutlineDocumentSyncService) SyncTitleToDocument(ctx context.Context, outlineNodeID string, newTitle string) error {
	outlineNode, err := s.outlineRepo.FindByID(ctx, outlineNodeID)
	if err != nil || outlineNode == nil {
		return nil
	}
	if outlineNode.DocumentID == "" {
		return nil
	}

	updates := map[string]interface{}{
		"title": newTitle,
	}
	if err := s.documentRepo.Update(ctx, outlineNode.DocumentID, updates); err != nil {
		return err
	}
	return nil
}

// SyncTitleToOutline 文档标题变更同步到大纲
func (s *OutlineDocumentSyncService) SyncTitleToOutline(ctx context.Context, documentID string, newTitle string) error {
	doc, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil || doc == nil {
		return nil
	}
	if doc.OutlineNodeID == "" {
		return nil
	}

	outlineNode, err := s.outlineRepo.FindByID(ctx, doc.OutlineNodeID)
	if err != nil || outlineNode == nil {
		return nil
	}

	outlineNode.Title = newTitle
	if err := s.outlineRepo.Update(ctx, outlineNode); err != nil {
		return err
	}
	return nil
}

// HandleOutlineDeletion 大纲删除时处理关联文档
func (s *OutlineDocumentSyncService) HandleOutlineDeletion(ctx context.Context, outlineNode *writer.OutlineNode) error {
	if outlineNode.DocumentID == "" {
		return nil
	}

	// 清除文档的反向引用
	updates := map[string]interface{}{
		"outline_node_id": "",
	}
	if err := s.documentRepo.Update(ctx, outlineNode.DocumentID, updates); err != nil {
		return err
	}
	return nil
}

// HandleDocumentDeletion 文档删除时处理关联大纲
func (s *OutlineDocumentSyncService) HandleDocumentDeletion(ctx context.Context, doc *writer.Document) error {
	if doc.OutlineNodeID == "" {
		return nil
	}

	// 清除大纲的正向引用
	outlineNode, err := s.outlineRepo.FindByID(ctx, doc.OutlineNodeID)
	if err != nil || outlineNode == nil {
		return nil
	}
	outlineNode.DocumentID = ""
	if err := s.outlineRepo.Update(ctx, outlineNode); err != nil {
		return err
	}
	return nil
}

// syncOutlineToDocument 大纲创建后同步文档的 outline_node_id 引用（内部方法）
func (s *OutlineDocumentSyncService) syncOutlineToDocument(ctx context.Context, documentID, outlineNodeID string) error {
	updates := map[string]interface{}{
		"outline_node_id": outlineNodeID,
	}
	if err := s.documentRepo.Update(ctx, documentID, updates); err != nil {
		return err
	}
	return nil
}

// SyncLevelChangeToDocument 大纲层级变化时同步到文档
func (s *OutlineDocumentSyncService) SyncLevelChangeToDocument(ctx context.Context, outlineNodeID string, newLevel int) error {
	// 获取大纲节点
	outlineNode, err := s.outlineRepo.FindByID(ctx, outlineNodeID)
	if err != nil || outlineNode == nil {
		return nil
	}

	// 检查是否有关联文档
	if outlineNode.DocumentID == "" {
		return nil
	}

	// 根据新层级确定文档类型
	newDocType := outlineLevelToDocType(newLevel)
	if newDocType == "" {
		// 不需要映射的层级，删除文档的关联
		updates := map[string]interface{}{
			"outline_node_id": "",
		}
		_ = s.documentRepo.Update(ctx, outlineNode.DocumentID, updates)
		return nil
	}

	// 更新文档的 Level 和 Type
	updates := map[string]interface{}{
		"level": newLevel,
		"type":  newDocType,
	}
	if err := s.documentRepo.Update(ctx, outlineNode.DocumentID, updates); err != nil {
		return err
	}
	return nil
}
