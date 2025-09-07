package document

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/global"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProjectStatus 工程状态
type ProjectStatus string

const (
	ProjectStatusPublic  ProjectStatus = "public"
	ProjectStatusPrivate ProjectStatus = "private"
)

// Project 表示一本小说工程
type Project struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	OwnerID     string    `bson:"owner_id" json:"ownerId"`
	Name        string    `bson:"name" json:"name"`
	Status      string    `bson:"status" json:"status"` // public | private
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// NodeType 统一描述“目录/文件”两种类型
type NodeType string

const (
	NodeTypeFolder NodeType = "folder"
	NodeTypeFile   NodeType = "file"
)

// Node 表示项目中的一个节点（目录或文件）
// 目录与文件统一为树结构，便于路径与排序管理
type Node struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	ProjectID    string    `bson:"project_id" json:"projectId"`
	ParentID     string    `bson:"parent_id,omitempty" json:"parentId,omitempty"`
	Type         NodeType  `bson:"type" json:"type"`                  // folder | file
	Name         string    `bson:"name" json:"name"`                  // 节点名（不含路径）
	Slug         string    `bson:"slug" json:"slug"`                  // 可选：URL/路径友好名
	RelativePath string    `bson:"relative_path" json:"relativePath"` // 相对项目根的路径，如 "第一卷/第一章.md"
	Order        int       `bson:"order" json:"order"`                // 同父下排序
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`

	// 非持久化：用于一次性返回树
	Children []*Node `bson:"-" json:"children,omitempty"`
}

// NovelFile 存放文件的内容与文件级元信息（仅当 Node.Type==file 时有对应记录）
type NovelFile struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	ProjectID string    `bson:"project_id" json:"projectId"`
	NodeID    string    `bson:"node_id" json:"nodeId"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`
	Format    string    `bson:"format" json:"format"` // markdown / txt / json 等
	Words     int       `bson:"words" json:"words"`
	Version   int       `bson:"version" json:"version"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// TouchForCreate 在创建前设置时间戳
func (f *NovelFile) TouchForCreate() {
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now
	if f.Version == 0 {
		f.Version = 1
	}
}

// TouchForUpdate 在更新前刷新更新时间戳
func (f *NovelFile) TouchForUpdate() {
	f.UpdatedAt = time.Now()
}

func (n *Node) TouchForCreate() {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
}

func (n *Node) TouchForUpdate() {
	n.UpdatedAt = time.Now()
}

// BuildRelativePath 根据父路径生成相对路径
func (n *Node) BuildRelativePath(parentPath string) (string, error) {
	// 检查节点名是否为空
	if n.Name == "" {
		return "", fmt.Errorf("node name cannot be empty")
	}

	if parentPath == "" || parentPath == "/" {
		return n.Name, nil
	}
	return parentPath + "/" + n.Name, nil
}

// nodeCol 返回节点的 MongoDB 集合
func nodeCol() *mongo.Collection { return global.DB.Collection("nodes") }

// NodeService 提供节点操作的服务
type NodeService struct{}

// MoveNode 移动节点到新的父节点下，并更新所有子节点的路径
func (s *NodeService) MoveNode(projectID, nodeID, newParentID string) error {
	// 检查节点是否存在
	var node Node
	err := nodeCol().FindOne(context.Background(), bson.M{"_id": nodeID, "project_id": projectID}).Decode(&node)
	if err != nil {
		return fmt.Errorf("node not found: %v", err)
	}

	// 如果是移动到根目录
	if newParentID == "" {
		node.ParentID = ""
		node.RelativePath = node.Name
	} else {
		// 检查新的父节点是否存在
		var parent Node
		err = nodeCol().FindOne(context.Background(), bson.M{"_id": newParentID, "project_id": projectID}).Decode(&parent)
		if err != nil {
			return fmt.Errorf("parent node not found: %v", err)
		}

		// 检查是否是移动到自己的后代节点（会造成循环引用）
		if newParentID == nodeID {
			return fmt.Errorf("cannot move node to itself")
		}

		// 检查新父节点是否是当前节点的后代
		currentID := newParentID
		for {
			var currentNode Node
			err = nodeCol().FindOne(context.Background(), bson.M{"_id": currentID}).Decode(&currentNode)
			if err != nil {
				break
			}

			if currentNode.ParentID == nodeID {
				return fmt.Errorf("cannot move node to its own descendant")
			}

			if currentNode.ParentID == "" {
				break
			}

			currentID = currentNode.ParentID
		}

		// 更新父节点和路径
		node.ParentID = newParentID
		node.RelativePath = parent.RelativePath + "/" + node.Name
	}

	// 更新节点
	_, err = nodeCol().UpdateOne(
		context.Background(),
		bson.M{"_id": nodeID},
		bson.M{"$set": bson.M{
			"parent_id":      node.ParentID,
			"relative_path":  node.RelativePath,
			"updated_at":     time.Now(),
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update node: %v", err)
	}

	// 递归更新所有子节点的路径
	return s.updateChildPaths(projectID, nodeID, node.RelativePath)
}

// updateChildPaths 递归更新所有子节点的路径
func (s *NodeService) updateChildPaths(projectID, parentID, parentPath string) error {
	// 查找所有直接子节点
	cursor, err := nodeCol().Find(context.Background(), bson.M{"project_id": projectID, "parent_id": parentID})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var child Node
		if err := cursor.Decode(&child); err != nil {
			return err
		}

		// 更新子节点的路径
		newPath := parentPath + "/" + child.Name
		_, err = nodeCol().UpdateOne(
			context.Background(),
			bson.M{"_id": child.ID},
			bson.M{"$set": bson.M{
				"relative_path": newPath,
				"updated_at":    time.Now(),
			}},
		)
		if err != nil {
			return err
		}

		// 递归更新子节点的子节点
		if err := s.updateChildPaths(projectID, child.ID, newPath); err != nil {
			return err
		}
	}

	return nil
}
