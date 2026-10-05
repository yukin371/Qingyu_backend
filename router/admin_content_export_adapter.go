package router

import (
	"context"

	"Qingyu_backend/models/bookstore"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"
	serviceinterfaces "Qingyu_backend/service/interfaces"
)

// adminContentBookRepository 是组合根适配器需要的最小仓储能力。
type adminContentBookRepository interface {
	List(context.Context, infrastructure.Filter) ([]*bookstore.Book, error)
}

// adminContentChapterRepository 是组合根适配器需要的最小仓储能力。
type adminContentChapterRepository interface {
	List(context.Context, infrastructure.Filter) ([]*bookstore.Chapter, error)
}

// adminContentBookPortAdapter 将仓储过滤器适配为 admin service 端口过滤器。
type adminContentBookPortAdapter struct {
	repository adminContentBookRepository
}

var _ serviceinterfaces.AdminContentBookPort = (*adminContentBookPortAdapter)(nil)

func newAdminContentBookPort(repository adminContentBookRepository) serviceinterfaces.AdminContentBookPort {
	if repository == nil {
		return nil
	}
	return &adminContentBookPortAdapter{repository: repository}
}

func (a *adminContentBookPortAdapter) List(ctx context.Context, filter serviceinterfaces.AdminContentFilter) ([]*bookstore.Book, error) {
	return a.repository.List(ctx, adminContentInfrastructureFilter(filter, false))
}

// adminContentChapterPortAdapter 将仓储过滤器适配为 admin service 端口过滤器。
type adminContentChapterPortAdapter struct {
	repository adminContentChapterRepository
}

var _ serviceinterfaces.AdminContentChapterPort = (*adminContentChapterPortAdapter)(nil)

func newAdminContentChapterPort(repository adminContentChapterRepository) serviceinterfaces.AdminContentChapterPort {
	if repository == nil {
		return nil
	}
	return &adminContentChapterPortAdapter{repository: repository}
}

func (a *adminContentChapterPortAdapter) List(ctx context.Context, filter serviceinterfaces.AdminContentFilter) ([]*bookstore.Chapter, error) {
	return a.repository.List(ctx, adminContentInfrastructureFilter(filter, true))
}

// adminContentInfrastructureFilter 是唯一的 repository 过滤器映射点。
// Mongo 查询直接使用持久化字段名，避免 handler 传入驼峰字段导致日期筛选失效。
func adminContentInfrastructureFilter(filter serviceinterfaces.AdminContentFilter, chapter bool) infrastructure.Filter {
	conditions := make(map[string]interface{})
	if filter.Status != "" {
		conditions["status"] = filter.Status
	}
	if filter.Author != "" {
		conditions["author"] = filter.Author
	}
	if filter.BookID != "" {
		conditions["book_id"] = filter.BookID
	}
	if filter.CreatedAt != nil {
		if filter.CreatedAt.Start != nil {
			conditions["created_at.$gte"] = *filter.CreatedAt.Start
		}
		if filter.CreatedAt.End != nil {
			conditions["created_at.$lte"] = *filter.CreatedAt.End
		}
	}
	if chapter && filter.PublishedAt != nil {
		if filter.PublishedAt.Start != nil {
			conditions["publish_time.$gte"] = *filter.PublishedAt.Start
		}
		if filter.PublishedAt.End != nil {
			conditions["publish_time.$lte"] = *filter.PublishedAt.End
		}
	}
	return &infrastructure.BaseFilter{Conditions: conditions}
}
