package writer

import (
	"context"

	"Qingyu_backend/models/writer"
)

// projectReader 是 writer 资产创建场景所需的最小项目读取 port。
// 项目权限语义仍由 writer.Project.CanEdit 统一负责。
type projectReader interface {
	GetByID(ctx context.Context, id string) (*writer.Project, error)
}
