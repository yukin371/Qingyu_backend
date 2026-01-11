package writer

import (
	"fmt"
	"path/filepath"

	"Qingyu_backend/models/writer/base"
)

// Node 节点模型（通用树形结构节点）
// 用于大纲、时间线、地点分类等树形结构
type Node struct {
	base.IdentifiedEntity    `bson:",inline"` // ID
	base.Timestamps          `bson:",inline"` // CreatedAt, UpdatedAt, DeletedAt
	base.ProjectScopedEntity `bson:",inline"` // ProjectID
	base.NamedEntity         `bson:",inline"` // Name

	// 树形结构属性（保持BSON字段名不变，确保数据库兼容）
	ParentID     string                 `bson:"parent_id,omitempty" json:"parentId,omitempty"`         // 父节点ID
	Type         string                 `bson:"type" json:"type" validate:"required"`                   // outline | timeline | location
	RelativePath string                 `bson:"relative_path,omitempty" json:"relativePath,omitempty"`  // 相对路径
	Order        int                    `bson:"order" json:"order" validate:"min=0"`                    // 同级排序
	Level        int                    `bson:"level" json:"level" validate:"min=0,max=10"`             // 层级深度
	Description  string                 `bson:"description,omitempty" json:"description,omitempty" validate:"max=500"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// TouchForCreate 创建时设置默认值
func (n *Node) TouchForCreate() {
	n.IdentifiedEntity.GenerateID()
	n.Timestamps.TouchForCreate()
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

// TouchForUpdate 更新时设置默认值
func (n *Node) TouchForUpdate() {
	n.Timestamps.Touch()
}

// IsRoot 判断是否为根节点
func (n *Node) IsRoot() bool {
	return n.ParentID == ""
}

// Validate 验证节点数据
func (n *Node) Validate() error {
	// 验证名称
	if err := base.ValidateName(n.Name, 100); err != nil {
		return err
	}

	// 验证项目ID
	if n.ProjectID == "" {
		return base.ErrProjectIDRequired
	}

	// 验证类型
	if n.Type == "" {
		return fmt.Errorf("节点类型不能为空")
	}

	return nil
}
