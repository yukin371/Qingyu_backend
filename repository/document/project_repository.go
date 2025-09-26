package document

import (
	"Qingyu_backend/models/document"
	"Qingyu_backend/repository/interfaces"
	"context"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// 继承CRUDRepository接口
	interfaces.CRUDRepository[*document.Project, string]
	// 继承 HealthRepository 接口
	interfaces.HealthRepository
	// 项目特定的查询方法
	GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*document.Project, error)
	GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*document.Project, error)
	UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error
	IsOwner(ctx context.Context, projectID, ownerID string) (bool, error)

	// 删除恢复相关
	SoftDelete(ctx context.Context, projectID, ownerID string) error
	HardDelete(ctx context.Context, projectID string) error
	Restore(ctx context.Context, projectID, ownerID string) error

	// 统计方法
	CountByOwner(ctx context.Context, ownerID string) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)

	// 事务支持
	CreateWithTransaction(ctx context.Context, project *document.Project, callback func(ctx context.Context) error) error
}

// ProjectIndexManager 项目索引管理接口
// 将索引管理从Service层分离到Repository层
type ProjectIndexManager interface {
	EnsureIndexes(ctx context.Context) error
	DropIndexes(ctx context.Context) error
	GetIndexInfo(ctx context.Context) ([]IndexInfo, error)
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name   string                 `json:"name"`
	Keys   map[string]interface{} `json:"keys"`
	Unique bool                   `json:"unique"`
	Sparse bool                   `json:"sparse"`
}

// NodeRepository 节点仓储接口
type NodeRepository interface {
	// 继承基础Repository接口
	interfaces.CRUDRepository[*document.Node, string]

	// 节点特定的查询方法
	GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*document.Node, error)
	GetByProjectAndType(ctx context.Context, projectID, nodeType string, limit, offset int64) ([]*document.Node, error)
	UpdateByProject(ctx context.Context, nodeID, projectID string, updates map[string]interface{}) error
	DeleteByProject(ctx context.Context, nodeID, projectID string) error
	RestoreByProject(ctx context.Context, nodeID, projectID string) error
	IsProjectMember(ctx context.Context, nodeID, projectID string) (bool, error)

	// 软删除相关
	SoftDelete(ctx context.Context, nodeID, projectID string) error
	HardDelete(ctx context.Context, nodeID string) error

	// 统计方法
	CountByProject(ctx context.Context, projectID string) (int64, error)

	// 事务支持
	CreateWithTransaction(ctx context.Context, node *document.Node, callback func(ctx context.Context) error) error
}

// DocumentRepository 文档仓储接口
type DocumentRepository interface {
	// 继承基础Repository接口
	interfaces.CRUDRepository[*document.Document, string]

	// 文档特定的查询方法
	GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*document.Document, error)
	GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*document.Document, error)
	UpdateByProject(ctx context.Context, documentID, projectID string, updates map[string]interface{}) error
	DeleteByProject(ctx context.Context, documentID, projectID string) error
	RestoreByProject(ctx context.Context, documentID, projectID string) error
	IsProjectMember(ctx context.Context, documentID, projectID string) (bool, error)

	// 软删除相关
	SoftDelete(ctx context.Context, documentID, projectID string) error
	HardDelete(ctx context.Context, documentID string) error

	// 统计方法
	CountByProject(ctx context.Context, projectID string) (int64, error)

	// 事务支持
	CreateWithTransaction(ctx context.Context, document *document.Document, callback func(ctx context.Context) error) error
}
