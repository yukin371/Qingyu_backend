package impl

import (
	"context"

	"Qingyu_backend/models/writer"
	serviceInterfaces "Qingyu_backend/service/interfaces"
	serviceWriter "Qingyu_backend/service/interfaces/writer"
	writerservice "Qingyu_backend/service/writer"
)

// ContentManagementImpl 内容管理端口实现
type ContentManagementImpl struct {
	characterService writerservice.CharacterService
	locationService  writerservice.LocationService
	timelineService  writerservice.TimelineService
	serviceName      string
	version          string
}

// NewContentManagementImpl 创建内容管理端口实现
func NewContentManagementImpl(
	characterService writerservice.CharacterService,
	locationService writerservice.LocationService,
	timelineService writerservice.TimelineService,
) serviceWriter.ContentManagementPort {
	return &ContentManagementImpl{
		characterService: characterService,
		locationService:  locationService,
		timelineService:  timelineService,
		serviceName:      "ContentManagementPort",
		version:          "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (c *ContentManagementImpl) Initialize(ctx context.Context) error {
	return nil
}

func (c *ContentManagementImpl) Health(ctx context.Context) error {
	return nil
}

func (c *ContentManagementImpl) Close(ctx context.Context) error {
	return nil
}

func (c *ContentManagementImpl) GetServiceName() string {
	return c.serviceName
}

func (c *ContentManagementImpl) GetVersion() string {
	return c.version
}

// ============================================================================
// ContentManagementPort 角色管理方法
// ============================================================================

// CreateCharacter 创建角色
func (c *ContentManagementImpl) CreateCharacter(ctx context.Context, projectID, userID string, req *serviceWriter.CreateCharacterRequest) (*writer.Character, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateCharacterRequest{
		Name:              req.Name,
		Alias:             req.Alias,
		Summary:           req.Summary,
		Traits:            req.Traits,
		Background:        req.Background,
		AvatarURL:         req.AvatarURL,
		PersonalityPrompt: req.PersonalityPrompt,
		SpeechPattern:     req.SpeechPattern,
		CurrentState:      req.CurrentState,
	}
	return c.characterService.Create(ctx, projectID, userID, convertReq)
}

// GetCharacterByID 根据ID获取角色
func (c *ContentManagementImpl) GetCharacterByID(ctx context.Context, characterID, projectID string) (*writer.Character, error) {
	return c.characterService.GetByID(ctx, characterID, projectID)
}

// ListCharacters 获取项目下的所有角色
func (c *ContentManagementImpl) ListCharacters(ctx context.Context, projectID string) ([]*writer.Character, error) {
	return c.characterService.List(ctx, projectID)
}

// UpdateCharacter 更新角色
func (c *ContentManagementImpl) UpdateCharacter(ctx context.Context, characterID, projectID string, req *serviceWriter.UpdateCharacterRequest) (*writer.Character, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.UpdateCharacterRequest{
		Name:              req.Name,
		Alias:             nil, // Port DTO 使用 string，需要转换为 *[]string
		Summary:           req.Summary,
		Traits:            nil, // Port DTO 使用 []string，需要转换为 *[]string
		Background:        req.Background,
		AvatarURL:         req.AvatarURL,
		PersonalityPrompt: req.PersonalityPrompt,
		SpeechPattern:     req.SpeechPattern,
		CurrentState:      req.CurrentState,
	}
	return c.characterService.Update(ctx, characterID, projectID, convertReq)
}

// DeleteCharacter 删除角色
func (c *ContentManagementImpl) DeleteCharacter(ctx context.Context, characterID, projectID string) error {
	return c.characterService.Delete(ctx, characterID, projectID)
}

// CreateCharacterRelation 创建角色关系
func (c *ContentManagementImpl) CreateCharacterRelation(ctx context.Context, projectID string, req *serviceWriter.CreateRelationRequest) (*writer.CharacterRelation, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateRelationRequest{
		FromID:   req.FromID,
		ToID:     req.ToID,
		Type:     req.Type,
		Strength: req.Strength,
		Notes:    req.Notes,
	}
	return c.characterService.CreateRelation(ctx, projectID, convertReq)
}

// ListCharacterRelations 获取角色关系列表
func (c *ContentManagementImpl) ListCharacterRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	return c.characterService.ListRelations(ctx, projectID, characterID)
}

// DeleteCharacterRelation 删除角色关系
func (c *ContentManagementImpl) DeleteCharacterRelation(ctx context.Context, relationID, projectID string) error {
	return c.characterService.DeleteRelation(ctx, relationID, projectID)
}

// GetCharacterGraph 获取角色关系图
func (c *ContentManagementImpl) GetCharacterGraph(ctx context.Context, projectID string) (*serviceWriter.CharacterGraph, error) {
	graph, err := c.characterService.GetCharacterGraph(ctx, projectID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.CharacterGraph{
		Nodes: graph.Nodes,
		Edges: graph.Edges,
	}, nil
}

// ============================================================================
// ContentManagementPort 地点管理方法
// ============================================================================

// CreateLocation 创建地点
func (c *ContentManagementImpl) CreateLocation(ctx context.Context, projectID, userID string, req *serviceWriter.CreateLocationRequest) (*writer.Location, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateLocationRequest{
		Name:        req.Name,
		Description: req.Description,
		Climate:     req.Climate,
		Culture:     req.Culture,
		Geography:   req.Geography,
		Atmosphere:  req.Atmosphere,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
	}
	return c.locationService.Create(ctx, projectID, userID, convertReq)
}

// GetLocationByID 根据ID获取地点
func (c *ContentManagementImpl) GetLocationByID(ctx context.Context, locationID, projectID string) (*writer.Location, error) {
	return c.locationService.GetByID(ctx, locationID, projectID)
}

// ListLocations 获取项目下的所有地点
func (c *ContentManagementImpl) ListLocations(ctx context.Context, projectID string) ([]*writer.Location, error) {
	return c.locationService.List(ctx, projectID)
}

// UpdateLocation 更新地点
func (c *ContentManagementImpl) UpdateLocation(ctx context.Context, locationID, projectID string, req *serviceWriter.UpdateLocationRequest) (*writer.Location, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.UpdateLocationRequest{
		Name:        req.Name,
		Description: req.Description,
		Climate:     req.Climate,
		Culture:     req.Culture,
		Geography:   req.Geography,
		Atmosphere:  req.Atmosphere,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
	}
	return c.locationService.Update(ctx, locationID, projectID, convertReq)
}

// DeleteLocation 删除地点
func (c *ContentManagementImpl) DeleteLocation(ctx context.Context, locationID, projectID string) error {
	return c.locationService.Delete(ctx, locationID, projectID)
}

// GetLocationTree 获取地点层级树
func (c *ContentManagementImpl) GetLocationTree(ctx context.Context, projectID string) ([]*serviceWriter.LocationNode, error) {
	tree, err := c.locationService.GetLocationTree(ctx, projectID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	result := make([]*serviceWriter.LocationNode, 0, len(tree))
	for _, node := range tree {
		result = append(result, c.convertLocationNode(node))
	}
	return result, nil
}

// convertLocationNode 递归转换地点节点
func (c *ContentManagementImpl) convertLocationNode(node *serviceInterfaces.LocationNode) *serviceWriter.LocationNode {
	if node == nil {
		return nil
	}
	children := make([]*serviceWriter.LocationNode, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, c.convertLocationNode(child))
	}
	return &serviceWriter.LocationNode{
		Location: node.Location,
		Children: children,
	}
}

// GetLocationPath 获取地点的完整路径
func (c *ContentManagementImpl) GetLocationPath(ctx context.Context, locationID string) ([]string, error) {
	return c.locationService.GetLocationPath(ctx, locationID)
}

// CreateLocationRelation 创建地点关系
func (c *ContentManagementImpl) CreateLocationRelation(ctx context.Context, projectID string, req *serviceWriter.CreateLocationRelationRequest) (*writer.LocationRelation, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateLocationRelationRequest{
		FromID:   req.FromID,
		ToID:     req.ToID,
		Type:     req.Type,
		Distance: "", // Port DTO 使用 *int，Service DTO 使用 string
		Notes:    req.Notes,
	}
	if req.Distance != nil {
		// convert int to string
		convertReq.Distance = string(rune(*req.Distance))
	}
	return c.locationService.CreateRelation(ctx, projectID, convertReq)
}

// ListLocationRelations 获取地点关系列表
func (c *ContentManagementImpl) ListLocationRelations(ctx context.Context, projectID string, locationID *string) ([]*writer.LocationRelation, error) {
	return c.locationService.ListRelations(ctx, projectID, locationID)
}

// DeleteLocationRelation 删除地点关系
func (c *ContentManagementImpl) DeleteLocationRelation(ctx context.Context, relationID, projectID string) error {
	return c.locationService.DeleteRelation(ctx, relationID, projectID)
}

// ============================================================================
// ContentManagementPort 时间线管理方法
// ============================================================================

// CreateTimeline 创建时间线
func (c *ContentManagementImpl) CreateTimeline(ctx context.Context, projectID string, req *serviceWriter.CreateTimelineRequest) (*writer.Timeline, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateTimelineRequest{
		Name:        req.Name,
		Description: req.Description,
		StartTime:   nil, // Port DTO 使用 string，Service DTO 使用 *writer.StoryTime
		EndTime:     nil, // Port DTO 使用 string，Service DTO 使用 *writer.StoryTime
	}
	return c.timelineService.CreateTimeline(ctx, projectID, convertReq)
}

// GetTimeline 根据ID获取时间线
func (c *ContentManagementImpl) GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error) {
	return c.timelineService.GetTimeline(ctx, timelineID, projectID)
}

// ListTimelines 获取项目下的所有时间线
func (c *ContentManagementImpl) ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error) {
	return c.timelineService.ListTimelines(ctx, projectID)
}

// DeleteTimeline 删除时间线
func (c *ContentManagementImpl) DeleteTimeline(ctx context.Context, timelineID, projectID string) error {
	return c.timelineService.DeleteTimeline(ctx, timelineID, projectID)
}

// CreateTimelineEvent 创建时间线事件
func (c *ContentManagementImpl) CreateTimelineEvent(ctx context.Context, projectID string, req *serviceWriter.CreateTimelineEventRequest) (*writer.TimelineEvent, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.CreateTimelineEventRequest{
		TimelineID:   req.TimelineID,
		Title:        req.Title,
		Description:  req.Description,
		EventType:    req.EventType,
		Importance:   req.Importance,
		Participants: req.Participants,
		LocationIDs:  req.LocationIDs,
		ChapterIDs:   req.ChapterIDs,
		StoryTime:    nil, // Port DTO 使用 string，Service DTO 使用 *writer.StoryTime
		Duration:     "",  // Port DTO 使用 *int，Service DTO 使用 string
		Impact:       req.Impact,
	}
	return c.timelineService.CreateEvent(ctx, projectID, convertReq)
}

// GetTimelineEvent 根据ID获取事件
func (c *ContentManagementImpl) GetTimelineEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error) {
	return c.timelineService.GetEvent(ctx, eventID, projectID)
}

// ListTimelineEvents 获取时间线下的所有事件
func (c *ContentManagementImpl) ListTimelineEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error) {
	return c.timelineService.ListEvents(ctx, timelineID)
}

// UpdateTimelineEvent 更新时间线事件
func (c *ContentManagementImpl) UpdateTimelineEvent(ctx context.Context, eventID, projectID string, req *serviceWriter.UpdateTimelineEventRequest) (*writer.TimelineEvent, error) {
	// 转换请求类型
	convertReq := &serviceInterfaces.UpdateTimelineEventRequest{
		Title:        req.Title,
		Description:  req.Description,
		EventType:    req.EventType,
		Importance:   req.Importance,
		Participants: nil, // Port DTO 使用 []string，Service DTO 使用 *[]string
		LocationIDs:  nil, // Port DTO 使用 []string，Service DTO 使用 *[]string
		ChapterIDs:   nil, // Port DTO 使用 []string，Service DTO 使用 *[]string
		StoryTime:    nil, // Port DTO 使用 *string，Service DTO 使用 *writer.StoryTime
		Duration:     nil, // Port DTO 使用 *int，Service DTO 使用 *string
		Impact:       req.Impact,
	}
	return c.timelineService.UpdateEvent(ctx, eventID, projectID, convertReq)
}

// DeleteTimelineEvent 删除时间线事件
func (c *ContentManagementImpl) DeleteTimelineEvent(ctx context.Context, eventID, projectID string) error {
	return c.timelineService.DeleteEvent(ctx, eventID, projectID)
}

// GetTimelineVisualization 获取时间线可视化数据
func (c *ContentManagementImpl) GetTimelineVisualization(ctx context.Context, timelineID string) (*serviceWriter.TimelineVisualization, error) {
	vis, err := c.timelineService.GetTimelineVisualization(ctx, timelineID)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	events := make([]*serviceWriter.TimelineEventNode, 0, len(vis.Events))
	for _, event := range vis.Events {
		events = append(events, &serviceWriter.TimelineEventNode{
			ID:         event.ID,
			Title:      event.Title,
			StoryTime:  event.StoryTime,
			EventType:  event.EventType,
			Importance: event.Importance,
			Characters: event.Characters,
		})
	}
	connections := make([]*serviceWriter.EventConnection, 0, len(vis.Connections))
	for _, conn := range vis.Connections {
		connections = append(connections, &serviceWriter.EventConnection{
			FromEventID: conn.FromEventID,
			ToEventID:   conn.ToEventID,
			Type:        conn.Type,
		})
	}
	return &serviceWriter.TimelineVisualization{
		Events:      events,
		Connections: connections,
	}, nil
}
