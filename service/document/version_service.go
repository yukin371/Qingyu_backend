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

func fileCol() *mongo.Collection  { return global.DB.Collection("novel_files") }    // 文件集合
func revCol() *mongo.Collection   { return global.DB.Collection("file_revisions") } // 版本集合
func patchCol() *mongo.Collection { return global.DB.Collection("file_patches") }   // 补丁集合

// EnsureIndexes 创建版本相关的 MongoDB 索引（幂等）
func (s *VersionService) EnsureIndexes(ctx context.Context) error {
	// file_revisions: project_id + node_id + version(desc)
	idxes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "version", Value: -1}},
			Options: nil,
		},
		{
			Keys:    bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: nil,
		},
	}
	if _, err := revCol().Indexes().CreateMany(ctx, idxes); err != nil {
		return err
	}
	// file_patches 索引
	_, _ = global.DB.Collection("file_patches").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "node_id", Value: 1}, {Key: "status", Value: 1}, {Key: "created_at", Value: -1}},
	})
	return nil
}

// BumpVersionAndCreateRevision 推进文件版本并创建快照版本记录
func (s *VersionService) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*model.FileRevision, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
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

	rev := &model.FileRevision{
		ProjectID:  projectID,
		NodeID:     nodeID,
		Version:    next,
		AuthorID:   authorID,
		Message:    message,
		Snapshot:   f.Content,
		ParentVers: f.Version,
		CreatedAt:  time.Now(),
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

	// Snapshot 必须存在（当前未实现外部存储 retrieval）
	if rev.Snapshot == "" {
		return nil, errors.New("snapshot not stored inline; external storage retrieval not implemented")
	}

	// 读取当前文档版本
	var f model.Document
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}

	// 使用乐观锁更新内容（期望为当前版本）
	return s.UpdateContentWithVersion(projectID, nodeID, authorID, message, rev.Snapshot, f.Version)
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
