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

func fileCol() *mongo.Collection { return global.DB.Collection("novel_files") }
func revCol() *mongo.Collection  { return global.DB.Collection("file_revisions") }

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
	id := ""
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		id = oid.Hex()
	}
	rev.ID = id
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
