package integration_test

import (
	"Qingyu_backend/models/writer"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"
	"Qingyu_backend/service/writer/project"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// 尝试直接连接 MongoDB（避免 import cycle），如果失败则跳过测试
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Skipping integration tests: cannot load config:", err)
		os.Exit(0)
	}

	if config.GlobalConfig.Database == nil {
		fmt.Println("Skipping integration tests: database config is nil")
		os.Exit(0)
	}

	// 检查MongoDB配置
	if config.GlobalConfig.Database.Primary.Type != config.DatabaseTypeMongoDB ||
		config.GlobalConfig.Database.Primary.MongoDB == nil {
		fmt.Println("Skipping integration tests: MongoDB config is missing")
		os.Exit(0)
	}

	mongoCfg := config.GlobalConfig.Database.Primary.MongoDB
	clientOpts := options.Client().ApplyURI(mongoCfg.URI)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println("Skipping integration tests: cannot connect to MongoDB:", err)
		os.Exit(0)
	}
	// ping
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Skipping integration tests: cannot ping MongoDB:", err)
		os.Exit(0)
	}
	global.MongoClient = client
	global.DB = client.Database(mongoCfg.Database)

	code := m.Run()

	// cleanup
	_ = global.MongoClient.Disconnect(ctx)
	os.Exit(code)
}

func fileCol() *mongo.Collection {
	return global.DB.Collection("novel_files")
}

func contentCol() *mongo.Collection {
	return global.DB.Collection("document_contents")
}

func revCol() *mongo.Collection {
	return global.DB.Collection("file_revisions")
}

func patchCol() *mongo.Collection {
	return global.DB.Collection("file_patches")
}

func cleanupCollections(t *testing.T, projectID string) {
	ctx := context.Background()
	global.DB.Collection("novel_files").DeleteMany(ctx, bson.M{"project_id": projectID})
	global.DB.Collection("document_contents").DeleteMany(ctx, bson.M{}) // 清理文档内容
	global.DB.Collection("file_revisions").DeleteMany(ctx, bson.M{"project_id": projectID})
	global.DB.Collection("file_patches").DeleteMany(ctx, bson.M{"project_id": projectID})
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
