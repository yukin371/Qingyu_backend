package writer

import (
	"fmt"
	"path/filepath"
	"time"
)

// Node 节点模型（用于大纲等树形结构）
type Node struct {
	ID           string                 `bson:"_id,omitempty" json:"id"`
	ProjectID    string                 `bson:"project_id" json:"projectId" validate:"required"`
	ParentID     string                 `bson:"parent_id,omitempty" json:"parentId,omitempty"`
	Name         string                 `bson:"name" json:"name" validate:"required,min=1,max=100"`
	Type         string                 `bson:"type" json:"type"`                                      // outline | timeline | location
	RelativePath string                 `bson:"relative_path,omitempty" json:"relativePath,omitempty"` // 相对路径
	Order        int                    `bson:"order" json:"order"`                                    // 同级排序
	Level        int                    `bson:"level" json:"level"`                                    // 层级深度
	Description  string                 `bson:"description,omitempty" json:"description,omitempty"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`

	// 时间戳
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// BuildRelativePath 构建相对路径
func (n *Node) BuildRelativePath(parentPath string) (string, error) {
	if n.Name == "" {
		return "", fmt.Errorf("节点名称不能为空")
	}

	if parentPath == "" {
		return n.Name, nil
	}

	return filepath.Join(parentPath, n.Name), nil
}

// TouchForCreate 设置创建时的默认值
func (n *Node) TouchForCreate() {
	now := time.Now()
	if n.CreatedAt.IsZero() {
		n.CreatedAt = now
	}
	if n.UpdatedAt.IsZero() {
		n.UpdatedAt = now
	}
}

// TouchForUpdate 设置更新时的默认值
func (n *Node) TouchForUpdate() {
	n.UpdatedAt = time.Now()
}

// IsRoot 判断是否为根节点
func (n *Node) IsRoot() bool {
	return n.ParentID == ""
}

// Validate 验证节点数据
func (n *Node) Validate() error {
	if n.ProjectID == "" {
		return fmt.Errorf("项目ID不能为空")
	}
	if n.Name == "" {
		return fmt.Errorf("节点名称不能为空")
	}
	if len(n.Name) > 100 {
		return fmt.Errorf("节点名称不能超过100字符")
	}
	return nil
}
