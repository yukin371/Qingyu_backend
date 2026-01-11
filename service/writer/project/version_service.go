package project

import (
	"Qingyu_backend/models/writer"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VersionService 版本管理服务
type VersionService struct {
	db *mongo.Database
}

// NewVersionService 创建版本服务
func NewVersionService(db *mongo.Database) *VersionService {
	return &VersionService{db: db}
}

func (s *VersionService) fileCol() *mongo.Collection    { return s.db.Collection("novel_files") }       // 文件集合（Document元数据）
func (s *VersionService) contentCol() *mongo.Collection { return s.db.Collection("document_contents") } // 文档内容集合
func (s *VersionService) revCol() *mongo.Collection     { return s.db.Collection("file_revisions") }    // 版本集合
func (s *VersionService) patchCol() *mongo.Collection   { return s.db.Collection("file_patches") }      // 补丁集合
func (s *VersionService) commitCol() *mongo.Collection  { return s.db.Collection("commits") }           // 提交集合

// getDocumentContent 获取文档内容（辅助函数）
func (s *VersionService) getDocumentContent(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	var content writer.DocumentContent
	err := s.contentCol().FindOne(ctx, bson.M{"document_id": documentID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文档内容失败: %w", err)
	}
	return &content, nil
}

// EnsureIndexes 创建版本相关的 MongoDB 索引（幂等）
func (s *VersionService) EnsureIndexes(ctx context.Context) error {
	// file_revisions 索引
	revIdxes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "version", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "commit_id", Value: 1}},
			Options: nil,
		},
	}
	if _, err := s.revCol().Indexes().CreateMany(ctx, revIdxes); err != nil {
		return err
	}

	// file_patches 索引
	patchIdxes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
		},
	}
	if _, err := s.patchCol().Indexes().CreateMany(ctx, patchIdxes); err != nil {
		return err
	}

	// commits 索引
	commitIdxes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "author_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
	}
	if _, err := s.commitCol().Indexes().CreateMany(ctx, commitIdxes); err != nil {
		return err
	}

	return nil
}

// BumpVersionAndCreateRevision 创建新版本并记录修订
func (s *VersionService) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*writer.FileRevision, error) {
	if s == nil {
		return nil, errors.New("VersionService is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查询Document元数据
	var f writer.Document
	if err := s.fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 查询DocumentContent获取版本号和内容
	docContent, err := s.getDocumentContent(ctx, f.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}
	if docContent == nil {
		return nil, errors.New("文档内容不存在")
	}

	// 版本推进（在DocumentContent中）
	next := docContent.Version + 1
	if _, err := s.contentCol().UpdateOne(ctx,
		bson.M{"document_id": f.ID},
		bson.M{"$set": bson.M{"version": next, "updated_at": time.Now()}}); err != nil {
		return nil, err
	}

	// 使用快照存储策略
	snapshot, storageRef, err := s.StoreSnapshot(docContent.Content, projectID, nodeID, next)
	if err != nil {
		return nil, err
	}

	rev := &writer.FileRevision{
		ProjectID:  projectID,
		NodeID:     nodeID,
		Version:    next,
		AuthorID:   authorID,
		Message:    message,
		Snapshot:   snapshot,
		StorageRef: storageRef,
		CreatedAt:  time.Now(),
	}
	res, err := s.revCol().InsertOne(ctx, rev)
	if err != nil {
		return nil, err
	}
	// 尝试从 InsertedID 中提取字符串 id（兼容 primitive.ObjectID）
	switch v := res.InsertedID.(type) {
	case string:
		rev.ID = v
	case interface{ Hex() string }:
		rev.ID = v.Hex()
	default:
		// 使用默认的格式化作为回退
		rev.ID = fmt.Sprintf("%v", res.InsertedID)
	}
	return rev, nil
}

// UpdateContentWithVersion 使用乐观并发控制更新内容，成功后创建新版本
func (s *VersionService) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*writer.FileRevision, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 只在版本匹配时更新
	res, err := s.fileCol().UpdateOne(ctx,
		bson.M{"project_id": projectID, "node_id": nodeID, "version": expectedVersion},
		bson.M{"$set": bson.M{"content": newContent, "updated_at": time.Now()}},
	)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, errors.New("version_conflict")
	}

	// 推进版本并记录修订
	return s.BumpVersionAndCreateRevision(projectID, nodeID, authorID, message)
}

// RollbackToVersion 回滚到指定的历史版本（通过创建新版本实现回滚）
func (s *VersionService) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*writer.FileRevision, error) {
	if projectID == "" || nodeID == "" || targetVersion <= 0 {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 找到目标修订
	var rev writer.FileRevision
	if err := s.revCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID, "version": targetVersion}).Decode(&rev); err != nil {
		return nil, err
	}

	// 获取快照内容
	content, err := s.RetrieveSnapshot(rev.Snapshot, rev.StorageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve snapshot: %w", err)
	}

	// 读取当前文档
	var f writer.Document
	if err := s.fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 获取当前DocumentContent版本
	docContent, err := s.getDocumentContent(ctx, f.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}
	if docContent == nil {
		return nil, errors.New("文档内容不存在")
	}

	// 使用乐观锁更新内容（期望为当前版本）
	return s.UpdateContentWithVersion(projectID, nodeID, authorID, message, content, docContent.Version)
}

// CreatePatch 提交一个候选补丁（状态为 pending）
func (s *VersionService) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*writer.FilePatch, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 使用字符串 id（Hex）来避免类型不一致
	id := primitive.NewObjectID().Hex()
	p := &writer.FilePatch{
		ID:          id,
		ProjectID:   projectID,
		NodeID:      nodeID,
		BaseVersion: baseVersion,
		DiffFormat:  diffFormat,
		DiffPayload: diffPayload,
		CreatedBy:   createdBy,
		Status:      "pending",
		Preview:     message,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// 手动指定 _id 字段为字符串 id
	_, err := s.patchCol().InsertOne(ctx, bson.M{"_id": id, "project_id": p.ProjectID, "node_id": p.NodeID, "base_version": p.BaseVersion, "diff_format": p.DiffFormat, "diff_payload": p.DiffPayload, "created_by": p.CreatedBy, "status": p.Status, "preview": p.Preview, "created_at": p.CreatedAt, "updated_at": p.UpdatedAt})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ApplyPatch 审核并应用补丁（仅在 baseVersion 匹配时直接应用）
func (s *VersionService) ApplyPatch(projectID, patchID, applierID string) (*writer.FileRevision, error) {
	if projectID == "" || patchID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 查找补丁
	var p writer.FilePatch
	if err := s.patchCol().FindOne(ctx, bson.M{"_id": patchID, "project_id": projectID}).Decode(&p); err != nil {
		return nil, err
	}
	if p.Status != "pending" {
		return nil, errors.New("patch not in pending state")
	}

	// 获取当前文档
	var f writer.Document
	if err := s.fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": p.NodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 获取当前DocumentContent版本
	docContent, err := s.getDocumentContent(ctx, f.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("获取文档内容失败: %w", err)
	}
	if docContent == nil {
		return nil, errors.New("文档内容不存在")
	}

	// 简化：只支持完整替换的 diffFormat 为 "full"
	if p.DiffFormat != "full" {
		return nil, errors.New("only full diffFormat supported currently")
	}

	// 要求 baseVersion 匹配当前版本以直接应用
	if p.BaseVersion != docContent.Version {
		return nil, errors.New("version_conflict")
	}

	// 使用乐观锁更新内容
	rev, err := s.UpdateContentWithVersion(projectID, p.NodeID, applierID, p.Preview, p.DiffPayload, docContent.Version)
	if err != nil {
		return nil, err
	}

	// 标记补丁为 applied
	if _, err := s.patchCol().UpdateOne(ctx, bson.M{"_id": patchID}, bson.M{"$set": bson.M{"status": "applied", "updated_at": time.Now()}}); err != nil {
		// 不致命，仍返回 rev
	}

	return rev, nil
}

// ListRevisions 列表修订（按版本倒序）
func (s *VersionService) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*writer.FileRevision, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	findOpts := &options.FindOptions{}
	if limit > 0 {
		findOpts.SetLimit(limit)
	}
	if offset > 0 {
		findOpts.SetSkip(offset)
	}
	cur, err := s.revCol().Find(ctx, bson.M{"project_id": projectID, "node_id": nodeID}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []*writer.FileRevision
	for cur.Next(ctx) {
		var r writer.FileRevision
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		res = append(res, &r)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// DetectConflicts 检测文件的版本冲突
func (s *VersionService) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*writer.ConflictInfo, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}

	// 获取当前文件状态
	var currentFile struct {
		Version   int       `bson:"version"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
	err := s.fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&currentFile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file_not_found")
		}
		return nil, err
	}

	// 如果版本匹配，没有冲突
	if currentFile.Version == expectedVersion {
		return &writer.ConflictInfo{
			HasConflict:     false,
			CurrentVersion:  currentFile.Version,
			ExpectedVersion: expectedVersion,
		}, nil
	}

	// 获取冲突的修订记录
	cursor, err := s.revCol().Find(ctx,
		bson.M{
			"project_id": projectID,
			"node_id":    nodeID,
			"version":    bson.M{"$gt": expectedVersion, "$lte": currentFile.Version},
		},
		options.Find().SetSort(bson.D{{Key: "version", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var conflictingRevisions []writer.FileRevision
	if err := cursor.All(ctx, &conflictingRevisions); err != nil {
		return nil, err
	}

	return &writer.ConflictInfo{
		HasConflict:          true,
		CurrentVersion:       currentFile.Version,
		ExpectedVersion:      expectedVersion,
		ConflictingRevisions: conflictingRevisions,
		LastModified:         currentFile.UpdatedAt,
	}, nil
}

// BatchDetectConflicts 批量检测多个文件的版本冲突
func (s *VersionService) BatchDetectConflicts(ctx context.Context, projectID string, files []struct {
	NodeID          string `json:"node_id"`
	ExpectedVersion int    `json:"expected_version"`
}) (*writer.BatchConflictResult, error) {
	if projectID == "" || len(files) == 0 {
		return nil, errors.New("invalid arguments")
	}

	result := &writer.BatchConflictResult{
		ProjectID:    projectID,
		HasConflicts: false,
		Conflicts:    make(map[string]*writer.ConflictInfo),
	}

	for _, file := range files {
		conflict, err := s.DetectConflicts(ctx, projectID, file.NodeID, file.ExpectedVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to detect conflicts for file %s: %w", file.NodeID, err)
		}

		result.Conflicts[file.NodeID] = conflict
		if conflict.HasConflict {
			result.HasConflicts = true
		}
	}

	return result, nil
}

// CreateCommit 创建批量提交（使用MongoDB事务）
func (s *VersionService) CreateCommit(ctx context.Context, projectID, authorID, message string, files []writer.CommitFile) (*writer.Commit, error) {
	if projectID == "" || authorID == "" || len(files) == 0 {
		return nil, errors.New("invalid arguments")
	}

	// 首先检测所有文件的冲突
	var conflictFiles []struct {
		NodeID          string `json:"node_id"`
		ExpectedVersion int    `json:"expected_version"`
	}
	for _, file := range files {
		conflictFiles = append(conflictFiles, struct {
			NodeID          string `json:"node_id"`
			ExpectedVersion int    `json:"expected_version"`
		}{
			NodeID:          file.NodeID,
			ExpectedVersion: file.ExpectedVersion,
		})
	}

	conflicts, err := s.BatchDetectConflicts(ctx, projectID, conflictFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to detect conflicts: %w", err)
	}

	if conflicts.HasConflicts {
		return nil, fmt.Errorf("commit_conflicts_detected")
	}

	// 创建提交记录
	commit := &writer.Commit{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		AuthorID:  authorID,
		Message:   message,
		FileCount: len(files),
		CreatedAt: time.Now(),
	}

	// 使用MongoDB事务确保原子性
	session, err := s.db.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 插入提交记录
		_, err := s.commitCol().InsertOne(sessCtx, commit)
		if err != nil {
			return nil, fmt.Errorf("failed to create commit: %w", err)
		}

		// 为每个文件创建修订记录
		var revisions []interface{}
		for _, file := range files {
			// 更新文件内容和版本
			updateResult, err := s.fileCol().UpdateOne(sessCtx,
				bson.M{"project_id": projectID, "node_id": file.NodeID, "version": file.ExpectedVersion},
				bson.M{
					"$set": bson.M{
						"content":    file.Content,
						"updated_at": time.Now(),
					},
					"$inc": bson.M{"version": 1},
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to update file %s: %w", file.NodeID, err)
			}
			if updateResult.MatchedCount == 0 {
				return nil, fmt.Errorf("version conflict for file %s", file.NodeID)
			}

			// 使用快照存储策略
			snapshot, storageRef, err := s.StoreSnapshot(file.Content, projectID, file.NodeID, file.ExpectedVersion+1)
			if err != nil {
				return nil, fmt.Errorf("failed to store snapshot for file %s: %w", file.NodeID, err)
			}

			// 创建修订记录
			revision := &writer.FileRevision{
				ID:         primitive.NewObjectID().Hex(),
				ProjectID:  projectID,
				NodeID:     file.NodeID,
				Version:    file.ExpectedVersion + 1,
				AuthorID:   authorID,
				Message:    message,
				Snapshot:   snapshot,
				StorageRef: storageRef,
				ParentVers: file.ExpectedVersion,
				CommitID:   commit.ID,
				CreatedAt:  time.Now(),
			}
			revisions = append(revisions, revision)
		}

		// 批量插入修订记录
		if len(revisions) > 0 {
			_, err = s.revCol().InsertMany(sessCtx, revisions)
			if err != nil {
				return nil, fmt.Errorf("failed to create revisions: %w", err)
			}
		}

		return commit, nil
	}

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return nil, err
	}

	return result.(*writer.Commit), nil
}

// StoreSnapshot 根据策略存储快照内容
func (s *VersionService) StoreSnapshot(content string, projectID, nodeID string, version int) (string, string, error) {
	contentSize := len([]byte(content))
	strategy := writer.GetSnapshotStrategy(contentSize)

	switch strategy {
	case "inline":
		// 直接存储在数据库中
		return content, "", nil
	case "external":
		// 存储到外部文件系统（简化实现）
		// 在实际项目中，这里可以集成对象存储服务如MinIO、AWS S3等
		externalPath := fmt.Sprintf("snapshots/%s/%s/v%d.txt", projectID, nodeID, version)
		// 这里只返回路径，实际存储逻辑需要根据具体需求实现
		return "", externalPath, nil
	default:
		return content, "", nil
	}
}

// RetrieveSnapshot 根据存储策略检索快照内容
func (s *VersionService) RetrieveSnapshot(snapshot, storageRef string) (string, error) {
	if storageRef != "" {
		// 从外部存储检索内容（简化实现）
		// 在实际项目中，这里需要实现从对象存储服务读取文件的逻辑
		return "", errors.New("external storage retrieval not implemented")
	}
	// 内联存储，直接返回快照内容
	return snapshot, nil
}

// ListCommits 查询提交历史
func (s *VersionService) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*writer.Commit, error) {
	if projectID == "" {
		return nil, errors.New("project_id is required")
	}

	filter := bson.M{"project_id": projectID}
	if authorID != "" {
		filter["author_id"] = authorID
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := s.commitCol().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commits []*writer.Commit
	if err := cursor.All(ctx, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// GetCommitDetails 获取提交详情，包括相关的文件修订
func (s *VersionService) GetCommitDetails(ctx context.Context, projectID, commitID string) (*writer.Commit, []*writer.FileRevision, error) {
	if projectID == "" || commitID == "" {
		return nil, nil, errors.New("invalid arguments")
	}

	// 获取提交信息
	var commit writer.Commit
	err := s.commitCol().FindOne(ctx, bson.M{"_id": commitID, "project_id": projectID}).Decode(&commit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, errors.New("commit_not_found")
		}
		return nil, nil, err
	}

	// 获取相关的文件修订
	cursor, err := s.revCol().Find(ctx,
		bson.M{"commit_id": commitID, "project_id": projectID},
		options.Find().SetSort(bson.D{{Key: "node_id", Value: 1}}),
	)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var revisions []*writer.FileRevision
	if err := cursor.All(ctx, &revisions); err != nil {
		return nil, nil, err
	}

	return &commit, revisions, nil
}

// GetCurrentVersion 获取文件的当前版本号
func (s *VersionService) GetCurrentVersion(ctx context.Context, projectID, nodeID string) (int, error) {
	if projectID == "" || nodeID == "" {
		return 0, errors.New("invalid arguments")
	}

	var file struct {
		Version int `bson:"version"`
	}
	err := s.fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, errors.New("file_not_found")
		}
		return 0, err
	}

	return file.Version, nil
}

// ResolveBatchConflicts 批量解决冲突
func (s *VersionService) ResolveBatchConflicts(ctx context.Context, req *writer.BatchConflictResolution) (*writer.Commit, error) {
	// 首先检测所有文件的冲突状态
	var commitFiles []writer.CommitFile

	for nodeID, resolution := range req.Resolutions {
		// 验证冲突是否仍然存在
		conflict, err := s.DetectConflicts(ctx, req.ProjectID, nodeID, 0) // 使用0表示检查当前版本
		if err != nil {
			return nil, fmt.Errorf("failed to detect conflicts for file %s: %w", nodeID, err)
		}

		if !conflict.HasConflict {
			// 如果没有冲突，跳过此文件
			continue
		}

		// 根据解决策略处理
		var resolvedContent string
		switch resolution.Strategy {
		case "auto":
			// 自动合并（简化实现，实际项目中需要更复杂的合并算法）
			resolvedContent = resolution.MergedContent
		case "manual":
			// 手动解决
			resolvedContent = resolution.MergedContent
		case "force":
			// 强制覆盖
			resolvedContent = resolution.MergedContent
		default:
			return nil, fmt.Errorf("unsupported resolution strategy: %s", resolution.Strategy)
		}

		// 获取当前版本作为期望版本
		currentVersion, err := s.GetCurrentVersion(ctx, req.ProjectID, nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current version for file %s: %w", nodeID, err)
		}

		commitFiles = append(commitFiles, writer.CommitFile{
			NodeID:          nodeID,
			ExpectedVersion: currentVersion,
			Content:         resolvedContent,
		})
	}

	if len(commitFiles) == 0 {
		return nil, errors.New("no conflicts to resolve")
	}

	// 创建批量提交来应用解决方案
	message := fmt.Sprintf("Resolve conflicts: %s", req.Message)
	return s.CreateCommit(ctx, req.ProjectID, req.AuthorID, message, commitFiles)
}

// AutoResolveConflicts 自动解决简单冲突
func (s *VersionService) AutoResolveConflicts(ctx context.Context, projectID, nodeID string, conflictingRevisions []writer.FileRevision) (string, error) {
	// 简化的自动合并实现
	// 在实际项目中，这里需要实现更复杂的三路合并算法

	if len(conflictingRevisions) < 2 {
		return "", errors.New("insufficient revisions for auto-merge")
	}

	// 获取最新版本的内容作为基础
	latestRevision := conflictingRevisions[0]
	for _, rev := range conflictingRevisions {
		if rev.Version > latestRevision.Version {
			latestRevision = rev
		}
	}

	// 检索快照内容
	content, err := s.RetrieveSnapshot(latestRevision.Snapshot, latestRevision.StorageRef)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve snapshot: %w", err)
	}

	// 简化的合并策略：返回最新版本的内容
	// 在实际项目中，这里需要实现更智能的合并算法
	return content, nil
}

// GetVersionHistory 获取版本历史
func (s *VersionService) GetVersionHistory(ctx context.Context, documentID string, page, pageSize int) (*VersionHistoryResponse, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询版本历史
	filter := bson.M{"node_id": documentID}
	opts := options.Find().
		SetSort(bson.D{{Key: "version", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(pageSize))

	cursor, err := s.revCol().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询版本历史失败: %w", err)
	}
	defer cursor.Close(ctx)

	var revisions []writer.FileRevision
	if err := cursor.All(ctx, &revisions); err != nil {
		return nil, fmt.Errorf("解析版本历史失败: %w", err)
	}

	// 统计总数
	total, err := s.revCol().CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计版本数量失败: %w", err)
	}

	// 转换为响应格式
	versions := make([]*VersionInfo, 0, len(revisions))
	for _, rev := range revisions {
		versions = append(versions, &VersionInfo{
			VersionID: rev.ID,
			Version:   rev.Version,
			Message:   rev.Message,
			CreatedAt: rev.CreatedAt,
			CreatedBy: rev.AuthorID,
			WordCount: 0, // TODO: 从快照中获取字数
		})
	}

	return &VersionHistoryResponse{
		Versions: versions,
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetVersion 获取特定版本
func (s *VersionService) GetVersion(ctx context.Context, documentID, versionID string) (*VersionDetail, error) {
	// 查询版本 - versionID直接是string类型
	var revision writer.FileRevision
	err := s.revCol().FindOne(ctx, bson.M{"_id": versionID, "node_id": documentID}).Decode(&revision)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("版本不存在")
		}
		return nil, fmt.Errorf("查询版本失败: %w", err)
	}

	// 获取内容
	content, err := s.RetrieveSnapshot(revision.Snapshot, revision.StorageRef)
	if err != nil {
		return nil, fmt.Errorf("获取版本内容失败: %w", err)
	}

	return &VersionDetail{
		VersionID:  revision.ID,
		DocumentID: revision.NodeID,
		Version:    revision.Version,
		Content:    content,
		Message:    revision.Message,
		CreatedAt:  revision.CreatedAt,
		CreatedBy:  revision.AuthorID,
		WordCount:  len(content), // 简单字数统计
	}, nil
}

// CompareVersions 比较两个版本
func (s *VersionService) CompareVersions(ctx context.Context, documentID, fromVersionID, toVersionID string) (*VersionDiff, error) {
	// 获取两个版本
	fromVersion, err := s.GetVersion(ctx, documentID, fromVersionID)
	if err != nil {
		return nil, fmt.Errorf("获取源版本失败: %w", err)
	}

	toVersion, err := s.GetVersion(ctx, documentID, toVersionID)
	if err != nil {
		return nil, fmt.Errorf("获取目标版本失败: %w", err)
	}

	// 简单的行差异比较
	fromLines := splitLines(fromVersion.Content)
	toLines := splitLines(toVersion.Content)

	changes := make([]ChangeItem, 0)
	addedLines := 0
	deletedLines := 0

	// 简单的差异算法（实际应该使用更好的diff算法）
	maxLen := len(fromLines)
	if len(toLines) > maxLen {
		maxLen = len(toLines)
	}

	for i := 0; i < maxLen; i++ {
		fromLine := ""
		toLine := ""
		if i < len(fromLines) {
			fromLine = fromLines[i]
		}
		if i < len(toLines) {
			toLine = toLines[i]
		}

		if fromLine != toLine {
			if fromLine == "" {
				// 新增行
				changes = append(changes, ChangeItem{
					Type:    "added",
					Line:    i + 1,
					Content: toLine,
				})
				addedLines++
			} else if toLine == "" {
				// 删除行
				changes = append(changes, ChangeItem{
					Type:    "deleted",
					Line:    i + 1,
					Content: fromLine,
				})
				deletedLines++
			} else {
				// 修改行
				changes = append(changes, ChangeItem{
					Type:    "modified",
					Line:    i + 1,
					Content: toLine,
				})
			}
		}
	}

	return &VersionDiff{
		FromVersion:  fromVersionID,
		ToVersion:    toVersionID,
		Changes:      changes,
		AddedLines:   addedLines,
		DeletedLines: deletedLines,
	}, nil
}

// RestoreVersion 恢复到特定版本
func (s *VersionService) RestoreVersion(ctx context.Context, documentID, versionID string) error {
	// 获取要恢复的版本
	version, err := s.GetVersion(ctx, documentID, versionID)
	if err != nil {
		return fmt.Errorf("获取版本失败: %w", err)
	}

	// 获取文档当前内容
	currentContent, err := s.getDocumentContent(ctx, documentID)
	if err != nil {
		return fmt.Errorf("获取当前文档失败: %w", err)
	}

	// 更新文档内容
	if currentContent == nil {
		// 创建新内容
		_, err = s.contentCol().InsertOne(ctx, bson.M{
			"document_id": documentID,
			"content":     version.Content,
			"updated_at":  time.Now(),
		})
	} else {
		// 更新现有内容
		_, err = s.contentCol().UpdateOne(ctx,
			bson.M{"document_id": documentID},
			bson.M{
				"$set": bson.M{
					"content":    version.Content,
					"updated_at": time.Now(),
				},
			},
		)
	}

	if err != nil {
		return fmt.Errorf("更新文档内容失败: %w", err)
	}

	// TODO: 创建一个新的版本记录，标记为恢复操作

	return nil
}

// splitLines 将文本按行分割
func splitLines(text string) []string {
	if text == "" {
		return []string{}
	}
	lines := []string{}
	currentLine := ""
	for _, ch := range text {
		if ch == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(ch)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}
