package writer

import (
	"Qingyu_backend/models/writer"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// 继承CRUDRepository接口
	base.CRUDRepository[*writer.Project, string]
	// 继承 HealthRepository 接口
	base.HealthRepository
	// 项目特定的查询方法
	GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error)
	GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error)
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
	CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error
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
	base.CRUDRepository[*writer.Node, string]

	// 节点特定的查询方法
	GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Node, error)
	GetByProjectAndType(ctx context.Context, projectID, nodeType string, limit, offset int64) ([]*writer.Node, error)
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
	CreateWithTransaction(ctx context.Context, node *writer.Node, callback func(ctx context.Context) error) error
}

// DocumentRepository 文档仓储接口
type DocumentRepository interface {
	// 继承基础Repository接口
	base.CRUDRepository[*writer.Document, string]

	// 继承 HealthRepository 接口
	base.HealthRepository

	// 文档特定的查询方法
	GetByProjectID(ctx context.Context, projectID string, limit, offset int64) ([]*writer.Document, error)
	GetByProjectAndType(ctx context.Context, projectID, documentType string, limit, offset int64) ([]*writer.Document, error)
	GetByIDs(ctx context.Context, ids []string) ([]*writer.Document, error)
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
	CreateWithTransaction(ctx context.Context, document *writer.Document, callback func(ctx context.Context) error) error
}

// DocumentContentRepository 文档内容仓储接口
// 用于管理文档的实际内容，与Document元数据分离
type DocumentContentRepository interface {
	// 继承基础Repository接口
	base.CRUDRepository[*writer.DocumentContent, string]

	// 继承 HealthRepository 接口
	base.HealthRepository

	// 通过DocumentID查询内容
	GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error)

	// 带版本号的更新（乐观锁）
	UpdateWithVersion(ctx context.Context, documentID string, content string, expectedVersion int) error

	// 批量更新
	BatchUpdateContent(ctx context.Context, updates map[string]string) error

	// 内容统计
	GetContentStats(ctx context.Context, documentID string) (wordCount, charCount int, err error)

	// 大文件支持
	StoreToGridFS(ctx context.Context, documentID string, content []byte) (gridFSID string, err error)
	LoadFromGridFS(ctx context.Context, gridFSID string) (content []byte, err error)

	// 事务支持
	CreateWithTransaction(ctx context.Context, content *writer.DocumentContent, callback func(ctx context.Context) error) error
}
