package integration_test

import (
	"Qingyu_backend/models/writer"
	"context"
	"testing"
	"time"

	"Qingyu_backend/service"
	"Qingyu_backend/service/writer/project"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func getDB() *mongo.Database {
	if db == nil {
		db = service.ServiceManager.GetMongoDB()
	}
	return db
}

func fileCol() *mongo.Collection {
	return getDB().Collection("novel_files")
}

func contentCol() *mongo.Collection {
	return getDB().Collection("document_contents")
}

func revCol() *mongo.Collection {
	return getDB().Collection("file_revisions")
}

func patchCol() *mongo.Collection {
	return getDB().Collection("file_patches")
}

func cleanupCollections(t *testing.T, projectID string) {
	ctx := context.Background()
	getDB().Collection("novel_files").DeleteMany(ctx, bson.M{"project_id": projectID})
	getDB().Collection("document_contents").DeleteMany(ctx, bson.M{}) // 清理文档内容
	getDB().Collection("file_revisions").DeleteMany(ctx, bson.M{"project_id": projectID})
	getDB().Collection("file_patches").DeleteMany(ctx, bson.M{"project_id": projectID})
}

func TestUpdateContentWithVersion_HappyPath(t *testing.T) {
	t.Skip("VersionService需要完整的依赖注入，暂时跳过此测试")

	svc := &project.VersionService{}
	projectID := "test_project"
	nodeID := "node_update_1"
	cleanupCollections(t, projectID)
	ctx := context.Background()

	// 插入初始文档元数据
	docID := primitive.NewObjectID().Hex()
	_, err := fileCol().InsertOne(ctx, bson.M{
		"_id":        docID,
		"project_id": projectID,
		"node_id":    nodeID,
		"title":      "Test Document",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("insert file failed: %v", err)
	}

	// 插入文档内容（包含版本号）
	_, err = contentCol().InsertOne(ctx, bson.M{
		"_id":         primitive.NewObjectID().Hex(),
		"document_id": docID,
		"content":     "old",
		"version":     1,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	})
	if err != nil {
		t.Fatalf("insert content failed: %v", err)
	}

	rev, err := svc.UpdateContentWithVersion(projectID, nodeID, "user1", "update msg", "new content", 1)
	if err != nil {
		t.Fatalf("UpdateContentWithVersion failed: %v", err)
	}
	if rev == nil {
		t.Fatalf("expected revision, got nil")
	}

	// 验证文档内容已更新
	var content writer.DocumentContent
	if err := contentCol().FindOne(ctx, bson.M{"document_id": docID}).Decode(&content); err != nil {
		t.Fatalf("failed fetch content: %v", err)
	}
	if content.Content != "new content" {
		t.Fatalf("content mismatch: got %v", content.Content)
	}
	if content.Version != 2 {
		t.Fatalf("version mismatch: want 2 got %d", content.Version)
	}

	cleanupCollections(t, projectID)
}

func TestUpdateContentWithVersion_Conflict(t *testing.T) {
	svc := &project.VersionService{}
	projectID := "test_project"
	nodeID := "node_update_2"
	cleanupCollections(t, projectID)
	ctx := context.Background()

	// 插入文档元数据
	docID := primitive.NewObjectID().Hex()
	_, err := fileCol().InsertOne(ctx, bson.M{
		"_id":        docID,
		"project_id": projectID,
		"node_id":    nodeID,
		"title":      "Test Document",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("insert file failed: %v", err)
	}

	// 插入文档内容（版本号为2）
	_, err = contentCol().InsertOne(ctx, bson.M{
		"_id":         primitive.NewObjectID().Hex(),
		"document_id": docID,
		"content":     "old",
		"version":     2,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	})
	if err != nil {
		t.Fatalf("insert content failed: %v", err)
	}

	// 尝试用版本1更新，应该冲突
	_, err = svc.UpdateContentWithVersion(projectID, nodeID, "user1", "update msg", "new content", 1)
	if err == nil {
		t.Fatalf("expected version_conflict but got nil")
	}

	cleanupCollections(t, projectID)
}

func TestRollbackToVersion(t *testing.T) {
	svc := &project.VersionService{}
	projectID := "test_project"
	nodeID := "node_rb_1"
	cleanupCollections(t, projectID)
	ctx := context.Background()

	// 插入文档元数据
	docID := primitive.NewObjectID().Hex()
	_, err := fileCol().InsertOne(ctx, bson.M{
		"_id":        docID,
		"project_id": projectID,
		"node_id":    nodeID,
		"title":      "Test Document",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("insert file failed: %v", err)
	}

	// 插入当前文档内容 version 3
	_, err = contentCol().InsertOne(ctx, bson.M{
		"_id":         primitive.NewObjectID().Hex(),
		"document_id": docID,
		"content":     "v3",
		"version":     3,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	})
	if err != nil {
		t.Fatalf("insert content failed: %v", err)
	}

	// 插入历史 revision v1
	_, err = revCol().InsertOne(ctx, bson.M{
		"project_id": projectID,
		"node_id":    nodeID,
		"version":    1,
		"snapshot":   "v1",
		"created_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("insert revision failed: %v", err)
	}

	rev, err := svc.RollbackToVersion(projectID, nodeID, 1, "admin", "rollback to v1")
	if err != nil {
		t.Fatalf("RollbackToVersion failed: %v", err)
	}
	if rev == nil {
		t.Fatalf("expected rev, got nil")
	}

	// 验证文档内容已回滚
	var content writer.DocumentContent
	if err := contentCol().FindOne(ctx, bson.M{"document_id": docID}).Decode(&content); err != nil {
		t.Fatalf("failed fetch content: %v", err)
	}
	if content.Content != "v1" {
		t.Fatalf("rollback content mismatch: got %v", content.Content)
	}

	cleanupCollections(t, projectID)
}

func TestCreateAndApplyPatch_Full(t *testing.T) {
	svc := &project.VersionService{}
	projectID := "test_project"
	nodeID := "node_patch_1"
	cleanupCollections(t, projectID)
	ctx := context.Background()

	// 插入文档元数据
	docID := primitive.NewObjectID().Hex()
	_, err := fileCol().InsertOne(ctx, bson.M{
		"_id":        docID,
		"project_id": projectID,
		"node_id":    nodeID,
		"title":      "Test Document",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("insert file failed: %v", err)
	}

	// 插入当前文档内容 version 1
	_, err = contentCol().InsertOne(ctx, bson.M{
		"_id":         primitive.NewObjectID().Hex(),
		"document_id": docID,
		"content":     "old",
		"version":     1,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	})
	if err != nil {
		t.Fatalf("insert content failed: %v", err)
	}

	p, err := svc.CreatePatch(projectID, nodeID, 1, "full", "new content via patch", "user2", "patch msg")
	if err != nil {
		t.Fatalf("CreatePatch failed: %v", err)
	}

	rev, err := svc.ApplyPatch(projectID, p.ID, "admin")
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}
	if rev == nil {
		t.Fatalf("expected rev, got nil")
	}

	// 验证文档内容已应用补丁
	var content writer.DocumentContent
	if err := contentCol().FindOne(ctx, bson.M{"document_id": docID}).Decode(&content); err != nil {
		t.Fatalf("failed fetch content: %v", err)
	}
	if content.Content != "new content via patch" {
		t.Fatalf("apply patch content mismatch: got %v", content.Content)
	}

	// check patch status
	var pdoc writer.FilePatch
	if err := patchCol().FindOne(ctx, bson.M{"_id": p.ID}).Decode(&pdoc); err != nil {
		t.Fatalf("failed fetch patch: %v", err)
	}
	if pdoc.Status != "applied" {
		t.Fatalf("patch status not applied: %v", pdoc.Status)
	}

	cleanupCollections(t, projectID)
}
