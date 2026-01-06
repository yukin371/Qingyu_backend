package interfaces

import (
	"Qingyu_backend/models/writer"
	"context"
)

// LocationService 地点服务接口
type LocationService interface {
	// 基础CRUD
	Create(ctx context.Context, projectID, userID string, req *CreateLocationRequest) (*writer.Location, error)
	GetByID(ctx context.Context, locationID, projectID string) (*writer.Location, error)
	List(ctx context.Context, projectID string) ([]*writer.Location, error)
	Update(ctx context.Context, locationID, projectID string, req *UpdateLocationRequest) (*writer.Location, error)
	Delete(ctx context.Context, locationID, projectID string) error

	// 层级管理
	GetLocationTree(ctx context.Context, projectID string) ([]*LocationNode, error)
	GetLocationPath(ctx context.Context, locationID string) ([]string, error)

	// 关系管理
	CreateRelation(ctx context.Context, projectID string, req *CreateLocationRelationRequest) (*writer.LocationRelation, error)
	ListRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error)
	DeleteRelation(ctx context.Context, relationID, projectID string) error
}

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Climate     string `json:"climate"`
	Culture     string `json:"culture"`
	Geography   string `json:"geography"`
	Atmosphere  string `json:"atmosphere"`
	ParentID    string `json:"parentId"`
	ImageURL    string `json:"imageUrl"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Climate     *string `json:"climate"`
	Culture     *string `json:"culture"`
	Geography   *string `json:"geography"`
	Atmosphere  *string `json:"atmosphere"`
	ParentID    *string `json:"parentId"`
	ImageURL    *string `json:"imageUrl"`
}

// CreateLocationRelationRequest 创建地点关系请求
type CreateLocationRelationRequest struct {
	FromID   string `json:"fromId" validate:"required"`
	ToID     string `json:"toId" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Distance string `json:"distance"`
	Notes    string `json:"notes"`
}

// LocationNode 地点树节点
type LocationNode struct {
	Location *writer.Location `json:"location"`
	Children []*LocationNode  `json:"children"`
}
