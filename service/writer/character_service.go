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

// CharacterService 角色服务实现
type CharacterService struct {
	characterRepo writing.CharacterRepository
	eventBus      base.EventBus
}

// NewCharacterService 创建CharacterService实例
func NewCharacterService(
	characterRepo writing.CharacterRepository,
	eventBus base.EventBus,
) serviceInterfaces.CharacterService {
	return &CharacterService{
		characterRepo: characterRepo,
		eventBus:      eventBus,
	}
}

// Create 创建角色
func (s *CharacterService) Create(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.CreateCharacterRequest,
) (*writer.Character, error) {
	// 验证关系类型（如果有）
	// TODO: 验证项目权限

	// 构建角色对象（使用base mixins）
	character := &writer.Character{}
	character.ProjectID = projectID
	character.Name = req.Name
	character.Alias = req.Alias
	character.Summary = req.Summary
	character.Traits = req.Traits
	character.Background = req.Background
	character.AvatarURL = req.AvatarURL
	character.PersonalityPrompt = req.PersonalityPrompt
	character.SpeechPattern = req.SpeechPattern
	character.CurrentState = req.CurrentState

	// 保存到数据库
	if err := s.characterRepo.Create(ctx, character); err != nil {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorInternal, "create character failed", "", err)
	}

	// 发布事件（触发向量化）
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "character.created",
			EventData: map[string]interface{}{
				"character_id": character.ID,
				"project_id":   projectID,
				"user_id":      userID,
				"name":         character.Name,
			},
			Timestamp: time.Now(),
			Source:    "CharacterService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return character, nil
}

// GetByID 根据ID获取角色
func (s *CharacterService) GetByID(
	ctx context.Context,
	characterID, projectID string,
) (*writer.Character, error) {
	character, err := s.characterRepo.FindByID(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// 验证项目权限
	if character.ProjectID != projectID {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorForbidden, "no permission to access this character", "", nil)
	}

	return character, nil
}

// List 获取项目下的所有角色
func (s *CharacterService) List(
	ctx context.Context,
	projectID string,
) ([]*writer.Character, error) {
	return s.characterRepo.FindByProjectID(ctx, projectID)
}

// Update 更新角色
func (s *CharacterService) Update(
	ctx context.Context,
	characterID, projectID string,
	req *serviceInterfaces.UpdateCharacterRequest,
) (*writer.Character, error) {
	// 获取现有角色
	character, err := s.GetByID(ctx, characterID, projectID)
	if err != nil {
		return nil, err
	}

	// 更新字段（仅更新非nil字段）
	if req.Name != nil {
		character.Name = *req.Name
	}
	if req.Alias != nil {
		character.Alias = *req.Alias
	}
	if req.Summary != nil {
		character.Summary = *req.Summary
	}
	if req.Traits != nil {
		character.Traits = *req.Traits
	}
	if req.Background != nil {
		character.Background = *req.Background
	}
	if req.AvatarURL != nil {
		character.AvatarURL = *req.AvatarURL
	}
	if req.PersonalityPrompt != nil {
		character.PersonalityPrompt = *req.PersonalityPrompt
	}
	if req.SpeechPattern != nil {
		character.SpeechPattern = *req.SpeechPattern
	}
	if req.CurrentState != nil {
		character.CurrentState = *req.CurrentState
	}

	// 保存更新
	if err := s.characterRepo.Update(ctx, character); err != nil {
		return nil, err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "character.updated",
			EventData: map[string]interface{}{
				"character_id": character.ID,
				"project_id":   projectID,
			},
			Timestamp: time.Now(),
			Source:    "CharacterService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return character, nil
}

// Delete 删除角色
func (s *CharacterService) Delete(
	ctx context.Context,
	characterID, projectID string,
) error {
	// 验证权限
	character, err := s.GetByID(ctx, characterID, projectID)
	if err != nil {
		return err
	}

	// 删除角色
	if err := s.characterRepo.Delete(ctx, characterID); err != nil {
		return err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "character.deleted",
			EventData: map[string]interface{}{
				"character_id": characterID,
				"project_id":   projectID,
			},
			Timestamp: time.Now(),
			Source:    "CharacterService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	// TODO: 删除相关的角色关系

	_ = character // 避免未使用变量警告

	return nil
}

// CreateRelation 创建角色关系
func (s *CharacterService) CreateRelation(
	ctx context.Context,
	projectID string,
	req *serviceInterfaces.CreateRelationRequest,
) (*writer.CharacterRelation, error) {
	// 验证两个角色是否存在且属于同一项目
	fromChar, err := s.GetByID(ctx, req.FromID, projectID)
	if err != nil {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorNotFound, "source character not found", "", err)
	}

	toChar, err := s.GetByID(ctx, req.ToID, projectID)
	if err != nil {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorNotFound, "target character not found", "", err)
	}

	// 验证关系类型
	if !writer.IsValidRelationType(req.Type) {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorValidation, "invalid relation type", "", nil)
	}

	// 创建关系（使用base mixins）
	relation := &writer.CharacterRelation{}
	relation.ProjectID = projectID
	relation.FromID = req.FromID
	relation.ToID = req.ToID
	relation.Type = writer.RelationType(req.Type)
	relation.Strength = req.Strength
	relation.Notes = req.Notes

	if err := s.characterRepo.CreateRelation(ctx, relation); err != nil {
		return nil, errors.NewServiceError("CharacterService", errors.ServiceErrorInternal, "create character relation failed", "", err)
	}

	_ = fromChar // 避免未使用变量警告
	_ = toChar

	return relation, nil
}

// ListRelations 获取角色关系列表
func (s *CharacterService) ListRelations(
	ctx context.Context,
	projectID string,
	characterID *string,
) ([]*writer.CharacterRelation, error) {
	return s.characterRepo.FindRelations(ctx, projectID, characterID)
}

// DeleteRelation 删除角色关系
func (s *CharacterService) DeleteRelation(
	ctx context.Context,
	relationID, projectID string,
) error {
	// 验证关系是否属于该项目
	relation, err := s.characterRepo.FindRelationByID(ctx, relationID)
	if err != nil {
		return err
	}

	if relation.ProjectID != projectID {
		return errors.NewServiceError("CharacterService", errors.ServiceErrorForbidden, "no permission to delete this relation", "", nil)
	}

	return s.characterRepo.DeleteRelation(ctx, relationID)
}

// GetCharacterGraph 获取角色关系图
func (s *CharacterService) GetCharacterGraph(
	ctx context.Context,
	projectID string,
) (*serviceInterfaces.CharacterGraph, error) {
	// 获取所有角色
	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 获取所有关系
	relations, err := s.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return nil, err
	}

	// 构建关系图
	graph := &serviceInterfaces.CharacterGraph{
		Nodes: characters,
		Edges: relations,
	}

	return graph, nil
}
