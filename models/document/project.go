package document

import "time"

// Project 表示一本小说工程
type Project struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	OwnerID     string    `bson:"owner_id" json:"ownerId"`
	Name        string    `bson:"name" json:"name"`
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
	Status    string    `bson:"status" json:"status"` // draft/review/published 等
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
