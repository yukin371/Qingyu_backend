package document

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/global"
	model "Qingyu_backend/models/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VersionService 版本管理服务
type VersionService struct{}

func fileCol() *mongo.Collection   { return global.DB.Collection("novel_files") }    // 文件集合
func revCol() *mongo.Collection    { return global.DB.Collection("file_revisions") } // 版本集合
func patchCol() *mongo.Collection  { return global.DB.Collection("file_patches") }   // 补丁集合
func commitCol() *mongo.Collection { return global.DB.Collection("commits") }        // 提交集合

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
	if _, err := revCol().Indexes().CreateMany(ctx, revIdxes); err != nil {
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
	if _, err := patchCol().Indexes().CreateMany(ctx, patchIdxes); err != nil {
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
	if _, err := commitCol().Indexes().CreateMany(ctx, commitIdxes); err != nil {
		return err
	}

	return nil
}

// BumpVersionAndCreateRevision 创建新版本并记录修订
func (s *VersionService) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*model.FileRevision, error) {
	if s == nil {
		return nil, errors.New("VersionService is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var f model.Document
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 版本推进
	next := f.Version + 1
	if _, err := fileCol().UpdateOne(ctx, bson.M{"_id": f.ID}, bson.M{"$set": bson.M{"version": next, "updated_at": time.Now()}}); err != nil {
		return nil, err
	}

	// 使用快照存储策略
	snapshot, storageRef, err := s.StoreSnapshot(f.Content, projectID, nodeID, next)
	if err != nil {
		return nil, err
	}

	rev := &model.FileRevision{
		ProjectID:   projectID,
		NodeID:      nodeID,
		Version:     next,
		AuthorID:    authorID,
		Message:     message,
		Snapshot:    snapshot,
		StorageRef:  storageRef,
		CreatedAt:   time.Now(),
	}
	res, err := revCol().InsertOne(ctx, rev)
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
func (s *VersionService) UpdateContentWithVersion(projectID, nodeID, authorID, message, newContent string, expectedVersion int) (*model.FileRevision, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 只在版本匹配时更新
	res, err := fileCol().UpdateOne(ctx,
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
func (s *VersionService) RollbackToVersion(projectID, nodeID string, targetVersion int, authorID, message string) (*model.FileRevision, error) {
	if projectID == "" || nodeID == "" || targetVersion <= 0 {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 找到目标修订
	var rev model.FileRevision
	if err := revCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID, "version": targetVersion}).Decode(&rev); err != nil {
		return nil, err
	}

	// 获取快照内容
	content, err := s.RetrieveSnapshot(rev.Snapshot, rev.StorageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve snapshot: %w", err)
	}

	// 读取当前文档版本
	var f model.Document
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 使用乐观锁更新内容（期望为当前版本）
	return s.UpdateContentWithVersion(projectID, nodeID, authorID, message, content, f.Version)
}

// CreatePatch 提交一个候选补丁（状态为 pending）
func (s *VersionService) CreatePatch(projectID, nodeID string, baseVersion int, diffFormat, diffPayload, createdBy, message string) (*model.FilePatch, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 使用字符串 id（Hex）来避免类型不一致
	id := primitive.NewObjectID().Hex()
	p := &model.FilePatch{
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
	_, err := patchCol().InsertOne(ctx, bson.M{"_id": id, "project_id": p.ProjectID, "node_id": p.NodeID, "base_version": p.BaseVersion, "diff_format": p.DiffFormat, "diff_payload": p.DiffPayload, "created_by": p.CreatedBy, "status": p.Status, "preview": p.Preview, "created_at": p.CreatedAt, "updated_at": p.UpdatedAt})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ApplyPatch 审核并应用补丁（仅在 baseVersion 匹配时直接应用）
func (s *VersionService) ApplyPatch(projectID, patchID, applierID string) (*model.FileRevision, error) {
	if projectID == "" || patchID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 查找补丁
	var p model.FilePatch
	if err := patchCol().FindOne(ctx, bson.M{"_id": patchID, "project_id": projectID}).Decode(&p); err != nil {
		return nil, err
	}
	if p.Status != "pending" {
		return nil, errors.New("patch not in pending state")
	}

	// 获取当前文档
	var f model.Document
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": p.NodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 简化：只支持完整替换的 diffFormat 为 "full"
	if p.DiffFormat != "full" {
		return nil, errors.New("only full diffFormat supported currently")
	}

	// 要求 baseVersion 匹配当前版本以直接应用
	if p.BaseVersion != f.Version {
		return nil, errors.New("version_conflict")
	}

	// 使用乐观锁更新内容
	rev, err := s.UpdateContentWithVersion(projectID, p.NodeID, applierID, p.Preview, p.DiffPayload, f.Version)
	if err != nil {
		return nil, err
	}

	// 标记补丁为 applied
	if _, err := patchCol().UpdateOne(ctx, bson.M{"_id": patchID}, bson.M{"$set": bson.M{"status": "applied", "updated_at": time.Now()}}); err != nil {
		// 不致命，仍返回 rev
	}

	return rev, nil
}

// ListRevisions 列表修订（按版本倒序）
func (s *VersionService) ListRevisions(ctx context.Context, projectID, nodeID string, limit, offset int64) ([]*model.FileRevision, error) {
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
	cur, err := revCol().Find(ctx, bson.M{"project_id": projectID, "node_id": nodeID}, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []*model.FileRevision
	for cur.Next(ctx) {
		var r model.FileRevision
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
func (s *VersionService) DetectConflicts(ctx context.Context, projectID, nodeID string, expectedVersion int) (*model.ConflictInfo, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}

	// 获取当前文件状态
	var currentFile struct {
		Version   int    `bson:"version"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
	err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&currentFile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file_not_found")
		}
		return nil, err
	}

	// 如果版本匹配，没有冲突
	if currentFile.Version == expectedVersion {
		return &model.ConflictInfo{
			HasConflict:     false,
			CurrentVersion:  currentFile.Version,
			ExpectedVersion: expectedVersion,
		}, nil
	}

	// 获取冲突的修订记录
	cursor, err := revCol().Find(ctx, 
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

	var conflictingRevisions []model.FileRevision
	if err := cursor.All(ctx, &conflictingRevisions); err != nil {
		return nil, err
	}

	return &model.ConflictInfo{
		HasConflict:           true,
		CurrentVersion:        currentFile.Version,
		ExpectedVersion:       expectedVersion,
		ConflictingRevisions:  conflictingRevisions,
		LastModified:          currentFile.UpdatedAt,
	}, nil
}

// BatchDetectConflicts 批量检测多个文件的版本冲突
func (s *VersionService) BatchDetectConflicts(ctx context.Context, projectID string, files []struct {
	NodeID          string `json:"node_id"`
	ExpectedVersion int    `json:"expected_version"`
}) (*model.BatchConflictResult, error) {
	if projectID == "" || len(files) == 0 {
		return nil, errors.New("invalid arguments")
	}

	result := &model.BatchConflictResult{
		ProjectID:    projectID,
		HasConflicts: false,
		Conflicts:    make(map[string]*model.ConflictInfo),
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
func (s *VersionService) CreateCommit(ctx context.Context, projectID, authorID, message string, files []model.CommitFile) (*model.Commit, error) {
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
	commit := &model.Commit{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		AuthorID:  authorID,
		Message:   message,
		FileCount: len(files),
		CreatedAt: time.Now(),
	}

	// 使用MongoDB事务确保原子性
	session, err := global.DB.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 插入提交记录
		_, err := commitCol().InsertOne(sessCtx, commit)
		if err != nil {
			return nil, fmt.Errorf("failed to create commit: %w", err)
		}

		// 为每个文件创建修订记录
		var revisions []interface{}
		for _, file := range files {
			// 更新文件内容和版本
			updateResult, err := fileCol().UpdateOne(sessCtx,
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
		revision := &model.FileRevision{
			ID:          primitive.NewObjectID().Hex(),
			ProjectID:   projectID,
			NodeID:      file.NodeID,
			Version:     file.ExpectedVersion + 1,
			AuthorID:    authorID,
			Message:     message,
			Snapshot:    snapshot,
			StorageRef:  storageRef,
			ParentVers:  file.ExpectedVersion,
			CommitID:    commit.ID,
			CreatedAt:   time.Now(),
		}
			revisions = append(revisions, revision)
		}

		// 批量插入修订记录
		if len(revisions) > 0 {
			_, err = revCol().InsertMany(sessCtx, revisions)
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

	return result.(*model.Commit), nil
}

// StoreSnapshot 根据策略存储快照内容
func (s *VersionService) StoreSnapshot(content string, projectID, nodeID string, version int) (string, string, error) {
	contentSize := len([]byte(content))
	strategy := model.GetSnapshotStrategy(contentSize)
	
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
func (s *VersionService) ListCommits(ctx context.Context, projectID string, authorID string, limit, offset int64) ([]*model.Commit, error) {
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

	cursor, err := commitCol().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var commits []*model.Commit
	if err := cursor.All(ctx, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}

// GetCommitDetails 获取提交详情，包括相关的文件修订
func (s *VersionService) GetCommitDetails(ctx context.Context, projectID, commitID string) (*model.Commit, []*model.FileRevision, error) {
	if projectID == "" || commitID == "" {
		return nil, nil, errors.New("invalid arguments")
	}

	// 获取提交信息
	var commit model.Commit
	err := commitCol().FindOne(ctx, bson.M{"_id": commitID, "project_id": projectID}).Decode(&commit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, errors.New("commit_not_found")
		}
		return nil, nil, err
	}

	// 获取相关的文件修订
	cursor, err := revCol().Find(ctx, 
		bson.M{"commit_id": commitID, "project_id": projectID},
		options.Find().SetSort(bson.D{{Key: "node_id", Value: 1}}),
	)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var revisions []*model.FileRevision
	if err := cursor.All(ctx, &revisions); err != nil {
		return nil, nil, err
	}

	return commit, revisions, nil
}

// ResolveBatchConflicts 批量解决冲突
func (s *VersionService) ResolveBatchConflicts(ctx context.Context, req *model.BatchConflictResolution) (*model.Commit, error) {
	// 首先检测所有文件的冲突状态
	var files []struct {
		NodeID          string
		ExpectedVersion int
		Content         string
	}

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

		files = append(files, struct {
			NodeID          string
			ExpectedVersion int
			Content         string
		}{
			NodeID:          nodeID,
			ExpectedVersion: currentVersion,
			Content:         resolvedContent,
		})
	}

	if len(files) == 0 {
		return nil, errors.New("no conflicts to resolve")
	}

	// 创建批量提交来应用解决方案
	message := fmt.Sprintf("Resolve conflicts: %s", req.Message)
	return s.CreateCommit(ctx, req.ProjectID, req.AuthorID, message, files)
}

// AutoResolveConflicts 自动解决简单冲突
func (s *VersionService) AutoResolveConflicts(ctx context.Context, projectID, nodeID string, conflictingRevisions []*model.FileRevision) (string, error) {
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

	return content, nil
}
