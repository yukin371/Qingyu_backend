package storyharness

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TriggerIndexResult 本次章节索引触发结果
type TriggerIndexResult struct {
	BatchID      string `json:"batchId"`
	Generated    int    `json:"generated"`
	Pending      int    `json:"pending"`
	Deduplicated int    `json:"deduplicated"`
	Source       string `json:"source"`
}

// IndexerService 章节索引与规则引擎服务
// 第一版仅做低成本规则扫描，并生成 ChangeRequestBatch + ChangeRequest。
type IndexerService struct {
	documentRepo        writerRepo.DocumentRepository
	documentContentRepo writerRepo.DocumentContentRepository
	characterRepo       writerRepo.CharacterRepository
	changeRequestRepo   writerRepo.ChangeRequestRepository
}

// NewIndexerService 创建 IndexerService 实例
func NewIndexerService(
	documentRepo writerRepo.DocumentRepository,
	documentContentRepo writerRepo.DocumentContentRepository,
	characterRepo writerRepo.CharacterRepository,
	changeRequestRepo writerRepo.ChangeRequestRepository,
) *IndexerService {
	return &IndexerService{
		documentRepo:        documentRepo,
		documentContentRepo: documentContentRepo,
		characterRepo:       characterRepo,
		changeRequestRepo:   changeRequestRepo,
	}
}

// TriggerChapterIndex 触发章节索引，生成最小规则建议。
func (s *IndexerService) TriggerChapterIndex(ctx context.Context, projectID, chapterID string) (*TriggerIndexResult, error) {
	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorValidation, "无效的项目ID", projectID, err)
	}

	chapterOID, err := primitive.ObjectIDFromHex(chapterID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorValidation, "无效的章节ID", chapterID, err)
	}

	chapter, err := s.documentRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "获取章节失败", chapterID, err)
	}
	if chapter == nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorNotFound, "章节不存在", chapterID, nil)
	}
	if chapter.ProjectID != projectOID {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorForbidden, "章节不属于当前项目", chapterID, nil)
	}

	text, err := s.loadDocumentText(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "获取角色列表失败", projectID, err)
	}

	existingPending, err := s.changeRequestRepo.FindPendingByChapter(ctx, projectID, chapterID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "获取待处理建议失败", chapterID, err)
	}

	batch := &writer.ChangeRequestBatch{
		ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectOID},
		ChapterID:           chapterOID,
		Source:              "rule",
		SaveTriggerID:       stringPtr(chapterID),
	}
	batch.TouchForCreate()

	if err := s.changeRequestRepo.CreateBatch(ctx, batch); err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "创建建议批次失败", chapterID, err)
	}

	suggestions := s.buildRuleSuggestions(text, projectOID, chapterOID, characters)
	pendingKeys := make(map[string]struct{}, len(existingPending))
	for _, item := range existingPending {
		pendingKeys[s.changeRequestKey(item)] = struct{}{}
	}

	created := 0
	deduplicated := 0
	currentKeys := make(map[string]struct{}, len(suggestions))

	for _, item := range suggestions {
		item.BatchID = batch.ID
		key := s.changeRequestKey(item)
		if _, exists := pendingKeys[key]; exists {
			deduplicated++
			continue
		}
		if _, exists := currentKeys[key]; exists {
			deduplicated++
			continue
		}
		currentKeys[key] = struct{}{}

		item.TouchForCreate()
		if err := item.Validate(); err != nil {
			deduplicated++
			continue
		}
		if err := s.changeRequestRepo.CreateRequest(ctx, item); err != nil {
			return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "创建变更建议失败", item.Title, err)
		}
		created++
	}

	if err := s.changeRequestRepo.UpdateBatchCounts(ctx, batch.ID.Hex(), created, created); err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "更新批次计数失败", batch.ID.Hex(), err)
	}

	return &TriggerIndexResult{
		BatchID:      batch.ID.Hex(),
		Generated:    created,
		Pending:      created,
		Deduplicated: deduplicated,
		Source:       batch.Source,
	}, nil
}

// TriggerChapterIndexByDocument 根据文档ID触发章节索引。
func (s *IndexerService) TriggerChapterIndexByDocument(ctx context.Context, documentID string) (*TriggerIndexResult, error) {
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "获取文档失败", documentID, err)
	}
	if document == nil {
		return nil, errors.NewServiceError("IndexerService", errors.ServiceErrorNotFound, "文档不存在", documentID, nil)
	}

	return s.TriggerChapterIndex(ctx, document.ProjectID.Hex(), documentID)
}

func (s *IndexerService) loadDocumentText(ctx context.Context, chapterID string) (string, error) {
	content, err := s.documentContentRepo.GetByDocumentID(ctx, chapterID)
	if err != nil {
		return "", errors.NewServiceError("IndexerService", errors.ServiceErrorInternal, "获取章节内容失败", chapterID, err)
	}
	if content == nil || strings.TrimSpace(content.Content) == "" {
		return "", nil
	}

	if content.ContentType == "tiptap_json" {
		if text := extractPlainTextFromTipTap(content.Content); strings.TrimSpace(text) != "" {
			return text, nil
		}
	}

	return content.Content, nil
}

func (s *IndexerService) buildRuleSuggestions(
	text string,
	projectID primitive.ObjectID,
	chapterID primitive.ObjectID,
	characters []*writer.Character,
) []*writer.ChangeRequest {
	if strings.TrimSpace(text) == "" || len(characters) == 0 {
		return nil
	}

	sentences := splitSentences(text)
	results := make([]*writer.ChangeRequest, 0, len(sentences))

	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)
		if trimmed == "" {
			continue
		}

		mentioned := findMentionedCharacters(trimmed, characters)
		if len(mentioned) == 0 {
			continue
		}

		if summary, ok := detectCharacterState(trimmed); ok {
			for _, character := range mentioned {
				results = append(results, &writer.ChangeRequest{
					ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
					ChapterID:           chapterID,
					Category:            writer.CRCategoryCharacterState,
					Priority:            writer.CRPriorityMedium,
					Status:              writer.CRStatusPending,
					Title:               fmt.Sprintf("建议更新角色状态：%s", character.Name),
					Description:         fmt.Sprintf("检测到 %s 在本章可能出现状态变化：%s", character.Name, summary),
					SuggestedChange: map[string]interface{}{
						"characterId":   character.ID.Hex(),
						"characterName": character.Name,
						"stateSummary":  summary,
					},
					Evidence: []writer.EvidenceRef{buildEvidenceRef(chapterID.Hex(), trimmed)},
					Source:   "rule",
				})
			}
		}

		if len(mentioned) >= 2 {
			if relationLabel, strength, ok := detectRelationChange(trimmed); ok {
				left := mentioned[0]
				right := mentioned[1]
				results = append(results, &writer.ChangeRequest{
					ProjectScopedEntity: writerBase.ProjectScopedEntity{ProjectID: projectID},
					ChapterID:           chapterID,
					Category:            writer.CRCategoryRelationChange,
					Priority:            writer.CRPriorityHigh,
					Status:              writer.CRStatusPending,
					Title:               fmt.Sprintf("建议更新关系：%s / %s", left.Name, right.Name),
					Description:         fmt.Sprintf("检测到 %s 与 %s 的关系可能发生变化：%s", left.Name, right.Name, relationLabel),
					SuggestedChange: map[string]interface{}{
						"fromId":   left.ID.Hex(),
						"toId":     right.ID.Hex(),
						"fromName": left.Name,
						"toName":   right.Name,
						"relation": relationLabel,
						"strength": strength,
					},
					Evidence: []writer.EvidenceRef{buildEvidenceRef(chapterID.Hex(), trimmed)},
					Source:   "rule",
				})
			}
		}
	}

	return results
}

func (s *IndexerService) changeRequestKey(item *writer.ChangeRequest) string {
	return strings.Join([]string{string(item.Category), item.Title}, "|")
}

func buildEvidenceRef(documentID, quote string) writer.EvidenceRef {
	return writer.EvidenceRef{
		DocumentID:   documentID,
		ParagraphIdx: 0,
		QuoteText:    truncateText(quote, 240),
	}
}

func splitSentences(text string) []string {
	replacer := strings.NewReplacer("。", "\n", "！", "\n", "？", "\n", "!", "\n", "?", "\n", ";", "\n", "；", "\n")
	normalized := replacer.Replace(text)
	lines := strings.Split(normalized, "\n")
	results := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}

func findMentionedCharacters(sentence string, characters []*writer.Character) []*writer.Character {
	matches := make([]*writer.Character, 0, 2)
	for _, character := range characters {
		if character == nil || strings.TrimSpace(character.Name) == "" {
			continue
		}
		if strings.Contains(sentence, character.Name) {
			matches = append(matches, character)
			continue
		}
		for _, alias := range character.Alias {
			if alias != "" && strings.Contains(sentence, alias) {
				matches = append(matches, character)
				break
			}
		}
	}
	return matches
}

func detectCharacterState(sentence string) (string, bool) {
	rules := []struct {
		keywords []string
		summary  string
	}{
		{keywords: []string{"断臂", "重伤", "受伤", "流血", "伤口", "昏迷"}, summary: "身体状态受损"},
		{keywords: []string{"疲惫", "虚弱", "力竭", "脱力", "喘息"}, summary: "体力明显下降"},
		{keywords: []string{"愤怒", "暴怒", "崩溃", "绝望", "恐惧", "害怕"}, summary: "情绪状态发生明显波动"},
	}

	for _, rule := range rules {
		for _, keyword := range rule.keywords {
			if strings.Contains(sentence, keyword) {
				return rule.summary, true
			}
		}
	}

	return "", false
}

func detectRelationChange(sentence string) (string, int, bool) {
	rules := []struct {
		keywords []string
		label    string
		strength int
	}{
		{keywords: []string{"不再信任", "信任崩塌", "怀疑", "猜忌"}, label: "信任下降", strength: 35},
		{keywords: []string{"决裂", "敌视", "敌意", "反目"}, label: "关系恶化", strength: 20},
		{keywords: []string{"和解", "并肩", "保护", "依赖", "信任"}, label: "关系改善", strength: 70},
	}

	for _, rule := range rules {
		for _, keyword := range rule.keywords {
			if strings.Contains(sentence, keyword) {
				return rule.label, rule.strength, true
			}
		}
	}

	return "", 0, false
}

func truncateText(value string, max int) string {
	if len([]rune(value)) <= max {
		return value
	}

	runes := []rune(value)
	return string(runes[:max])
}

func extractPlainTextFromTipTap(raw string) string {
	var payload interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return raw
	}

	var builder strings.Builder
	walkTipTapText(payload, &builder)
	return builder.String()
}

func walkTipTapText(node interface{}, builder *strings.Builder) {
	switch value := node.(type) {
	case map[string]interface{}:
		nodeType, _ := value["type"].(string)
		if nodeType == "text" {
			if text, ok := value["text"].(string); ok {
				builder.WriteString(text)
				builder.WriteString(" ")
			}
		}
		if content, ok := value["content"].([]interface{}); ok {
			for _, child := range content {
				walkTipTapText(child, builder)
			}
		}
	case []interface{}:
		for _, child := range value {
			walkTipTapText(child, builder)
		}
	}
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
