package document

import (
	"context"
	"errors"
	"time"

	"Qingyu_backend/global"
	model "Qingyu_backend/models/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VersionService struct{}

func fileCol() *mongo.Collection  { return global.DB.Collection("novel_files") }
func revCol() *mongo.Collection   { return global.DB.Collection("file_revisions") }
func patchCol() *mongo.Collection { return global.DB.Collection("file_patches") }

// BumpVersionAndCreateRevision 推进文件版本并创建快照版本记录
func (s *VersionService) BumpVersionAndCreateRevision(projectID, nodeID, authorID, message string) (*model.FileRevision, error) {
	if projectID == "" || nodeID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var f model.NovelFile
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
	id := ""
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		id = oid.Hex()
	}
	rev.ID = id
	return rev, nil
}

// CreatePatch 创建候选补丁
func (s *VersionService) CreatePatch(patch *model.FilePatch) (*model.FilePatch, error) {
	if patch == nil || patch.ProjectID == "" || patch.NodeID == "" {
		return nil, errors.New("invalid patch")
	}
	patch.Status = "pending"
	patch.CreatedAt = time.Now()
	patch.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := patchCol().InsertOne(ctx, patch)
	if err != nil {
		return nil, err
	}
	id := ""
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		id = oid.Hex()
	}
	patch.ID = id
	return patch, nil
}

// PreviewPatch 生成预览：当前简单返回原内容 + "\n" + DiffPayload（占位实现）
func (s *VersionService) PreviewPatch(projectID, nodeID, patchID string) (string, error) {
	if projectID == "" || nodeID == "" || patchID == "" {
		return "", errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var f model.NovelFile
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return "", err
	}
	var p model.FilePatch
	if err := patchCol().FindOne(ctx, bson.M{"_id": patchID, "project_id": projectID, "node_id": nodeID}).Decode(&p); err != nil {
		return "", err
	}
	preview := f.Content + "\n" + p.DiffPayload
	_, _ = patchCol().UpdateOne(ctx, bson.M{"_id": patchID}, bson.M{"$set": bson.M{"preview": preview, "updated_at": time.Now()}})
	return preview, nil
}

// ApplyPatch 应用补丁（占位实现：用 preview 或 DiffPayload 直接覆盖内容），创建版本
func (s *VersionService) ApplyPatch(projectID, nodeID, patchID, authorID, message string) (*model.FileRevision, error) {
	if projectID == "" || nodeID == "" || patchID == "" {
		return nil, errors.New("invalid arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var f model.NovelFile
	if err := fileCol().FindOne(ctx, bson.M{"project_id": projectID, "node_id": nodeID}).Decode(&f); err != nil {
		return nil, err
	}
	var p model.FilePatch
	if err := patchCol().FindOne(ctx, bson.M{"_id": patchID, "project_id": projectID, "node_id": nodeID}).Decode(&p); err != nil {
		return nil, err
	}

	// 占位：将变更结果作为新内容（实际应做三方合并或统一diff应用）
	newContent := p.Preview
	if newContent == "" {
		newContent = f.Content + "\n" + p.DiffPayload
	}

	// 更新文件内容并推进版本
	now := time.Now()
	if _, err := fileCol().UpdateOne(ctx, bson.M{"_id": f.ID}, bson.M{"$set": bson.M{"content": newContent, "updated_at": now}}); err != nil {
		return nil, err
	}
	rev, err := s.BumpVersionAndCreateRevision(projectID, nodeID, authorID, message)
	if err != nil {
		return nil, err
	}

	// 标记补丁为已应用
	_, _ = patchCol().UpdateOne(ctx, bson.M{"_id": patchID}, bson.M{"$set": bson.M{"status": "applied", "updated_at": now}})
	return rev, nil
}
