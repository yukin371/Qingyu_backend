package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// LocationService 地点服务实现
type LocationService struct {
	locationRepo writing.LocationRepository
	eventBus     base.EventBus
}

// NewLocationService 创建LocationService实例
func NewLocationService(
	locationRepo writing.LocationRepository,
	eventBus base.EventBus,
) serviceInterfaces.LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		eventBus:     eventBus,
	}
}

// Create 创建地点
func (s *LocationService) Create(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.CreateLocationRequest,
) (*writer.Location, error) {
	// 构建地点对象（使用base mixins）
	location := &writer.Location{}
	location.ProjectID = projectID
	location.Name = req.Name
	location.Description = req.Description
	location.Climate = req.Climate
	location.Culture = req.Culture
	location.Geography = req.Geography
	location.Atmosphere = req.Atmosphere
	location.ParentID = req.ParentID
	location.ImageURL = req.ImageURL

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorInternal, "create location failed", "", err)
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "location.created",
			EventData: map[string]interface{}{
				"location_id": location.ID,
				"project_id":  projectID,
				"user_id":     userID,
			},
			Timestamp: time.Now(),
			Source:    "LocationService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return location, nil
}

// GetByID 根据ID获取地点
func (s *LocationService) GetByID(
	ctx context.Context,
	locationID, projectID string,
) (*writer.Location, error) {
	location, err := s.locationRepo.FindByID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	if location.ProjectID != projectID {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorForbidden, "no permission to access this location", "", nil)
	}

	return location, nil
}

// List 获取项目下的所有地点
func (s *LocationService) List(
	ctx context.Context,
	projectID string,
) ([]*writer.Location, error) {
	return s.locationRepo.FindByProjectID(ctx, projectID)
}

// Update 更新地点
func (s *LocationService) Update(
	ctx context.Context,
	locationID, projectID string,
	req *serviceInterfaces.UpdateLocationRequest,
) (*writer.Location, error) {
	location, err := s.GetByID(ctx, locationID, projectID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Description != nil {
		location.Description = *req.Description
	}
	if req.Climate != nil {
		location.Climate = *req.Climate
	}
	if req.Culture != nil {
		location.Culture = *req.Culture
	}
	if req.Geography != nil {
		location.Geography = *req.Geography
	}
	if req.Atmosphere != nil {
		location.Atmosphere = *req.Atmosphere
	}
	if req.ParentID != nil {
		location.ParentID = *req.ParentID
	}
	if req.ImageURL != nil {
		location.ImageURL = *req.ImageURL
	}

	if err := s.locationRepo.Update(ctx, location); err != nil {
		return nil, err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "location.updated",
			EventData: map[string]interface{}{
				"location_id": locationID,
				"project_id":  projectID,
			},
			Timestamp: time.Now(),
			Source:    "LocationService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return location, nil
}

// Delete 删除地点
func (s *LocationService) Delete(
	ctx context.Context,
	locationID, projectID string,
) error {
	if _, err := s.GetByID(ctx, locationID, projectID); err != nil {
		return err
	}

	if err := s.locationRepo.Delete(ctx, locationID); err != nil {
		return err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "location.deleted",
			EventData: map[string]interface{}{
				"location_id": locationID,
				"project_id":  projectID,
			},
			Timestamp: time.Now(),
			Source:    "LocationService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return nil
}

// GetLocationTree 获取地点层级树
func (s *LocationService) GetLocationTree(
	ctx context.Context,
	projectID string,
) ([]*serviceInterfaces.LocationNode, error) {
	locations, err := s.locationRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return buildLocationTree(locations), nil
}

// GetLocationPath 获取地点的完整路径
func (s *LocationService) GetLocationPath(
	ctx context.Context,
	locationID string,
) ([]string, error) {
	path := []string{}
	currentID := locationID

	for currentID != "" {
		location, err := s.locationRepo.FindByID(ctx, currentID)
		if err != nil || location == nil {
			break
		}
		path = append([]string{location.Name}, path...)
		currentID = location.ParentID
	}

	return path, nil
}

// CreateRelation 创建地点关系
func (s *LocationService) CreateRelation(
	ctx context.Context,
	projectID string,
	req *serviceInterfaces.CreateLocationRelationRequest,
) (*writer.LocationRelation, error) {
	// 验证地点存在
	if _, err := s.GetByID(ctx, req.FromID, projectID); err != nil {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorNotFound, "source location not found", "", err)
	}
	if _, err := s.GetByID(ctx, req.ToID, projectID); err != nil {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorNotFound, "target location not found", "", err)
	}

	// 验证关系类型
	if !writer.IsValidLocationRelationType(req.Type) {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorValidation, "invalid relation type", "", nil)
	}

	// 创建关系（使用base mixins）
	relation := &writer.LocationRelation{}
	relation.ProjectID = projectID
	relation.FromID = req.FromID
	relation.ToID = req.ToID
	relation.Type = writer.LocationRelationType(req.Type)
	relation.Distance = req.Distance
	relation.Notes = req.Notes

	if err := s.locationRepo.CreateRelation(ctx, relation); err != nil {
		return nil, errors.NewServiceError("LocationService", errors.ServiceErrorInternal, "create location relation failed", "", err)
	}

	return relation, nil
}

// ListRelations 获取地点关系列表
func (s *LocationService) ListRelations(
	ctx context.Context,
	projectID string,
	locationID *string,
) ([]*writer.LocationRelation, error) {
	return s.locationRepo.FindRelations(ctx, projectID, locationID)
}

// DeleteRelation 删除地点关系
func (s *LocationService) DeleteRelation(
	ctx context.Context,
	relationID, projectID string,
) error {
	relations, err := s.locationRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return err
	}

	found := false
	for _, rel := range relations {
		if rel.ID == relationID {
			found = true
			break
		}
	}

	if !found {
		return errors.NewServiceError("LocationService", errors.ServiceErrorNotFound, "relation not found", "", nil)
	}

	return s.locationRepo.DeleteRelation(ctx, relationID)
}

// buildLocationTree 构建地点树
func buildLocationTree(locations []*writer.Location) []*serviceInterfaces.LocationNode {
	nodeMap := make(map[string]*serviceInterfaces.LocationNode)
	var roots []*serviceInterfaces.LocationNode

	for _, loc := range locations {
		nodeMap[loc.ID] = &serviceInterfaces.LocationNode{
			Location: loc,
			Children: []*serviceInterfaces.LocationNode{},
		}
	}

	for _, loc := range locations {
		node := nodeMap[loc.ID]
		if loc.ParentID == "" {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[loc.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}
