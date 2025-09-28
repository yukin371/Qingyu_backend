package project

import (
	"context"
	"errors"
	"time"

	"Qingyu_backend/global"
	model "Qingyu_backend/models/document"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeService struct{}

// NewNodeService 创建设置服务
func NewNodeService() *NodeService {
	return &NodeService{}
}

func nodeCol() *mongo.Collection { return global.DB.Collection("nodes") }

// RenameNode 重命名节点，并对所有子孙节点的相对路径进行级联更新
func (s *NodeService) RenameNode(projectID, nodeID, newName string) error {
	if projectID == "" || nodeID == "" || newName == "" {
		return errors.New("invalid arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var n model.Node
	if err := nodeCol().FindOne(ctx, bson.M{"_id": nodeID, "project_id": projectID}).Decode(&n); err != nil {
		return err
	}

	oldPath := n.RelativePath
	// 计算新路径：使用父路径 + 新名称
	parentPath := ""
	if n.ParentID != "" {
		var parent model.Node
		if err := nodeCol().FindOne(ctx, bson.M{"_id": n.ParentID, "project_id": projectID}).Decode(&parent); err != nil {
			return err
		}
		parentPath = parent.RelativePath
	}
	newPath, err := (&model.Node{Name: newName}).BuildRelativePath(parentPath)
	if err != nil {
		return err
	}

	// 更新当前节点
	update := bson.M{
		"$set": bson.M{
			"name":          newName,
			"relative_path": newPath,
			"updated_at":    time.Now(),
		},
	}
	if _, err := nodeCol().UpdateOne(ctx, bson.M{"_id": nodeID, "project_id": projectID}, update); err != nil {
		return err
	}

	// 级联更新子孙节点：查出以 oldPath 为前缀的所有节点
	cursor, err := nodeCol().Find(ctx, bson.M{
		"project_id":    projectID,
		"relative_path": bson.M{"$regex": "^" + oldPath + "/"},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var child model.Node
		if err := cursor.Decode(&child); err != nil {
			return err
		}
		// 将前缀 oldPath 替换为 newPath
		rel := child.RelativePath
		if len(rel) >= len(oldPath) && rel[:len(oldPath)] == oldPath {
			rel = newPath + rel[len(oldPath):]
		}
		if _, err := nodeCol().UpdateOne(ctx, bson.M{"_id": child.ID}, bson.M{"$set": bson.M{"relative_path": rel, "updated_at": time.Now()}}); err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

// MoveNode 移动节点到新的父节点，并级联更新路径
func (s *NodeService) MoveNode(projectID, nodeID, newParentID string) error {
	if projectID == "" || nodeID == "" {
		return errors.New("invalid arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var n model.Node
	if err := nodeCol().FindOne(ctx, bson.M{"_id": nodeID, "project_id": projectID}).Decode(&n); err != nil {
		return err
	}
	oldPath := n.RelativePath

	// 计算新的父路径
	parentPath := ""
	if newParentID != "" {
		var parent model.Node
		if err := nodeCol().FindOne(ctx, bson.M{"_id": newParentID, "project_id": projectID}).Decode(&parent); err != nil {
			return err
		}
		parentPath = parent.RelativePath
	}
	newPath, err := (&model.Node{Name: n.Name}).BuildRelativePath(parentPath)
	if err != nil {
		return err
	}

	// 更新当前节点的 ParentID 和 RelativePath
	if _, err := nodeCol().UpdateOne(ctx, bson.M{"_id": nodeID}, bson.M{"$set": bson.M{
		"parent_id":     newParentID,
		"relative_path": newPath,
		"updated_at":    time.Now(),
	}}); err != nil {
		return err
	}

	// 级联更新子孙节点
	cursor, err := nodeCol().Find(ctx, bson.M{
		"project_id":    projectID,
		"relative_path": bson.M{"$regex": "^" + oldPath + "/"},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var child model.Node
		if err := cursor.Decode(&child); err != nil {
			return err
		}
		rel := child.RelativePath
		if len(rel) >= len(oldPath) && rel[:len(oldPath)] == oldPath {
			rel = newPath + rel[len(oldPath):]
		}
		if _, err := nodeCol().UpdateOne(ctx, bson.M{"_id": child.ID}, bson.M{"$set": bson.M{"relative_path": rel, "updated_at": time.Now()}}); err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}
