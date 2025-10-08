package bookstore

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category 分类模型
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string            `bson:"name" json:"name" validate:"required,min=1,max=50"`        // 分类名称
	Description string            `bson:"description" json:"description" validate:"max=200"`        // 分类描述
	Icon        string            `bson:"icon" json:"icon"`                                         // 分类图标URL
	ParentID    *primitive.ObjectID `bson:"parent_id,omitempty" json:"parentId,omitempty"`         // 父分类ID
	Level       int               `bson:"level" json:"level"`                                       // 分类层级 0-顶级分类
	SortOrder   int               `bson:"sort_order" json:"sortOrder"`                             // 排序权重
	BookCount   int64             `bson:"book_count" json:"bookCount"`                             // 书籍数量
	IsActive    bool              `bson:"is_active" json:"isActive"`                               // 是否启用
	CreatedAt   time.Time         `bson:"created_at" json:"createdAt"`                             // 创建时间
	UpdatedAt   time.Time         `bson:"updated_at" json:"updatedAt"`                             // 更新时间
}

// CategoryFilter 分类查询过滤器
type CategoryFilter struct {
	ParentID  *primitive.ObjectID `json:"parentId,omitempty"`
	Level     *int               `json:"level,omitempty"`
	IsActive  *bool              `json:"isActive,omitempty"`
	Keyword   *string            `json:"keyword,omitempty"`
	SortBy    string             `json:"sortBy,omitempty"`    // sort_order, book_count, created_at
	SortOrder string             `json:"sortOrder,omitempty"` // asc, desc
	Limit     int                `json:"limit,omitempty"`
	Offset    int                `json:"offset,omitempty"`
}

// CategoryTree 分类树结构
type CategoryTree struct {
	Category
	Children []*CategoryTree `json:"children,omitempty"`
}
