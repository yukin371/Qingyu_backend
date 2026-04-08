package writer

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/interfaces"
)

// entityService 统一实体服务实现
type entityService struct {
	characterService interfaces.CharacterService
	locationService  interfaces.LocationService
	db               *mongo.Database
}

// NewEntityService 创建统一实体服务实例
func NewEntityService(
	characterService interfaces.CharacterService,
	locationService interfaces.LocationService,
	db *mongo.Database,
) interfaces.EntityService {
	return &entityService{
		characterService: characterService,
		locationService:  locationService,
		db:               db,
	}
}

// ListEntities 查询项目下所有实体（支持按 type 筛选）
func (s *entityService) ListEntities(ctx context.Context, projectID string, entityType *string) ([]interfaces.EntitySummary, error) {
	summaries := make([]interfaces.EntitySummary, 0)

	// 确定需要查询的实体类型
	shouldQueryCharacter := entityType == nil || *entityType == string(writer.EntityTypeCharacter)
	shouldQueryItem := entityType == nil || *entityType == string(writer.EntityTypeItem)
	shouldQueryLocation := entityType == nil || *entityType == string(writer.EntityTypeLocation)

	// 查询 characters
	if shouldQueryCharacter {
		characterSummaries, err := s.listCharacters(ctx, projectID)
		if err != nil {
			log.Printf("[EntityService] 查询角色失败: %v", err)
			// 不中断，继续查询其他类型
		} else {
			summaries = append(summaries, characterSummaries...)
		}
	}

	// 查询 items
	if shouldQueryItem {
		itemSummaries, err := s.listItems(ctx, projectID)
		if err != nil {
			log.Printf("[EntityService] 查询物品失败: %v", err)
		} else {
			summaries = append(summaries, itemSummaries...)
		}
	}

	// 查询 locations
	if shouldQueryLocation && s.locationService != nil {
		locationSummaries, err := s.listLocations(ctx, projectID)
		if err != nil {
			log.Printf("[EntityService] 查询地点失败: %v", err)
		} else {
			summaries = append(summaries, locationSummaries...)
		}
	}

	return summaries, nil
}

// GetEntityGraph 统一实体图谱
func (s *entityService) GetEntityGraph(ctx context.Context, projectID string) (*interfaces.EntityGraph, error) {
	graph := &interfaces.EntityGraph{}

	// 获取所有实体节点
	nodes, err := s.ListEntities(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取实体节点失败: %w", err)
	}
	graph.Nodes = nodes

	// 获取角色关系作为边
	if s.characterService != nil {
		relations, err := s.characterService.ListRelations(ctx, projectID, nil)
		if err != nil {
			log.Printf("[EntityService] 获取角色关系失败: %v", err)
		} else {
			for _, rel := range relations {
				graph.Edges = append(graph.Edges, interfaces.RelationEdge{
					FromID:   rel.FromID,
					ToID:     rel.ToID,
					FromType: rel.FromType,
					ToType:   rel.ToType,
					Type:     string(rel.Type),
					Strength: rel.Strength,
					Notes:    rel.Notes,
				})
			}
		}
	}

	return graph, nil
}

// UpdateEntityStateFields 更新实体状态字段
// 当前仅支持 character 类型，后续可扩展到其他实体类型
func (s *entityService) UpdateEntityStateFields(ctx context.Context, entityID string, stateFields map[string]writer.StateValue) error {
	if s.characterService == nil {
		return fmt.Errorf("角色服务不可用")
	}

	// 尝试通过 character service 更新
	// 由于不知道确切的 projectID，需要先查询 character 获取 projectID
	// 这里采用遍历方式：先尝试获取所有项目，然后查找对应 character
	// 但更简洁的做法是让调用方传入 projectID
	// 暂时使用直接 MongoDB 更新方式
	collection := s.db.Collection("characters")
	objectID, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return NewWriterError(ErrInvalidInput, "无效的实体ID格式")
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"state_fields": stateFields,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新实体状态字段失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return NewWriterError(ErrCharacterNotFound, "实体不存在")
	}

	return nil
}

// listCharacters 查询角色并转为 EntitySummary
func (s *entityService) listCharacters(ctx context.Context, projectID string) ([]interfaces.EntitySummary, error) {
	if s.characterService == nil {
		return nil, nil
	}

	characters, err := s.characterService.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summaries := make([]interfaces.EntitySummary, 0, len(characters))
	for _, c := range characters {
		summary := interfaces.EntitySummary{
			ID:          c.ID.Hex(),
			Name:        c.Name,
			EntityType:  writer.EntityTypeCharacter,
			Summary:     c.Summary,
			StateFields: c.StateFields,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// listItems 查询物品并转为 EntitySummary（直接 MongoDB 查询）
func (s *entityService) listItems(ctx context.Context, projectID string) ([]interfaces.EntitySummary, error) {
	if s.db == nil {
		return nil, nil
	}
	collection := s.db.Collection("items")
	filter := bson.M{"project_id": projectID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询物品失败: %w", err)
	}
	defer cursor.Close(ctx)

	var items []writer.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("解码物品数据失败: %w", err)
	}

	summaries := make([]interfaces.EntitySummary, 0, len(items))
	for _, item := range items {
		summary := interfaces.EntitySummary{
			ID:         item.ID,
			Name:       item.Name,
			EntityType: writer.EntityTypeItem,
			Summary:    item.Description,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// listLocations 查询地点并转为 EntitySummary
func (s *entityService) listLocations(ctx context.Context, projectID string) ([]interfaces.EntitySummary, error) {
	if s.locationService == nil {
		return nil, nil
	}

	locations, err := s.locationService.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summaries := make([]interfaces.EntitySummary, 0, len(locations))
	for _, loc := range locations {
		summary := interfaces.EntitySummary{
			ID:         loc.ID.Hex(),
			Name:       loc.Name,
			EntityType: writer.EntityTypeLocation,
			Summary:    loc.Description,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
