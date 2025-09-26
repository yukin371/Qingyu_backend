package document

import (
	"context"
	"fmt"
	"path/filepath"
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
	ProjectStatusDeleted ProjectStatus = "deleted" // 软删除
)

// IsValidProjectStatus 验证工程状态是否合法
func IsValidProjectStatus(status string) bool {
	switch ProjectStatus(status) {
	case ProjectStatusPublic, ProjectStatusPrivate:
		return true
	default:
		return false
	}
}

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

// Document 存放文件的内容与文件级元信息（仅当 Node.Type==file 时有对应记录）
type Document struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	ProjectID string `bson:"project_id,omitempty" json:"projectId,omitempty"` // 省略 project_id 时，视为独立文档
	NodeID    string `bson:"node_id,omitempty" json:"nodeId,omitempty"`
	Title     string `bson:"title" json:"title"`
	Content   string `bson:"content" json:"content"`
	Format    string `bson:"format" json:"format"` // markdown / txt / json 等
	Words     int    `bson:"words" json:"words"`
	Version   int    `bson:"version" json:"version"`

	// AI相关字段
	AIContext    string   `bson:"ai_context,omitempty" json:"aiContext,omitempty"`       // 文档的AI上下文，用于生成内容
	PlotThreads  []string `bson:"plot_threads,omitempty" json:"plotThreads,omitempty"`   // 文档的剧情线程，用于生成内容
	KeyPoints    []string `bson:"key_points,omitempty" json:"keyPoints,omitempty"`       // 文档的关键词点，用于生成内容
	WritingHints string   `bson:"writing_hints,omitempty" json:"writingHints,omitempty"` // 文档的写作提示，用于生成内容
	CharacterIDs []string `bson:"character_ids,omitempty" json:"characterIds,omitempty"` // 文档中出现的角色ID列表
	LocationIDs  []string `bson:"location_ids,omitempty" json:"locationIds,omitempty"`   // 文档中出现的地点ID列表
	TimelineIDs  []string `bson:"timeline_ids,omitempty" json:"timelineIds,omitempty"`   // 文档中出现的时间线ID列表

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// TouchForCreate 在创建前设置时间戳
func (f *Document) TouchForCreate() {
	now := time.Now()
	f.CreatedAt = now
	f.UpdatedAt = now
	if f.Version == 0 {
		f.Version = 1
	}
}

// TouchForUpdate 在更新前刷新更新时间戳
func (f *Document) TouchForUpdate() {
	f.UpdatedAt = time.Now()
	f.Version++
}

func (n *Node) TouchForCreate() {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	if n.Order == 0 {
		n.Order = 1
	}
}

func (n *Node) TouchForUpdate() {
	n.UpdatedAt = time.Now()
}

// BuildRelativePath 根据父路径生成相对路径
func (n *Node) BuildRelativePath(parentPath string) (string, error) {
	if n.Name == "" {
		return "", fmt.Errorf("节点名不能为空")
	}

	// 将所有的反斜杠转换为正斜杠，并清理路径中的 .. 等多余片段
	name := filepath.ToSlash(filepath.Clean(n.Name))
	if parentPath == "" || parentPath == "/" {
		return name, nil
	}

	// 确保使用正斜杠作为路径分隔符
	parentPath = filepath.ToSlash(filepath.Clean(parentPath))
	return parentPath + "/" + name, nil
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
		return fmt.Errorf("节点不存在: %v", err)
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
			return fmt.Errorf("父节点不存在: %v", err)
		}

		// 检查是否是移动到自己的后代节点（会造成循环引用）
		if newParentID == nodeID {
			return fmt.Errorf("不能将节点移动到其自身的后代")
		}

		// 检查新父节点是否是当前节点的后代
		currentID := newParentID
		for {
			var currentNode Node
			err = nodeCol().FindOne(context.Background(), bson.M{"_id": currentID, "project_id": projectID}).Decode(&currentNode)
			if err != nil {
				break
			}

			if currentNode.ParentID == nodeID {
				return fmt.Errorf("无法将节点移动到其自身的后代")
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
		bson.M{"_id": nodeID, "project_id": projectID},
		bson.M{"$set": bson.M{
			"parent_id":     node.ParentID,
			"relative_path": node.RelativePath,
			"updated_at":    time.Now(),
		}},
	)
	if err != nil {
		return fmt.Errorf("更新节点失败: %v", err)
	}

	// 递归更新所有子节点的路径
	return s.updateChildPaths(projectID, nodeID, node.RelativePath)
}

// updateChildPaths 递归更新所有子节点的路径
func (s *NodeService) updateChildPaths(projectID, parentID, parentPath string) error {
	// 查找所有直接子节点
	cursor, err := nodeCol().Find(context.Background(), bson.M{"project_id": projectID, "parent_id": parentID})
	if err != nil {
		return fmt.Errorf("查找子节点失败: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var child Node
		if decodeErr := cursor.Decode(&child); decodeErr != nil {
			return fmt.Errorf("解码子节点失败: %v", decodeErr)
		}

		// 更新子节点的路径
		newPath := parentPath + "/" + child.Name
		_, err = nodeCol().UpdateOne(
			context.Background(),
			bson.M{"_id": child.ID, "project_id": projectID},
			bson.M{"$set": bson.M{
				"relative_path": newPath,
				"updated_at":    time.Now(),
			}},
		)
		if err != nil {
			return fmt.Errorf("更新子节点路径失败: %v", err)
		}

		// 递归更新子节点的子节点
		if err := s.updateChildPaths(projectID, child.ID, newPath); err != nil {
			return fmt.Errorf("递归更新子节点失败: %v", err)
		}
	}

	return nil
}
