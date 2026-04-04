package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aiModels "Qingyu_backend/models/ai"
	writerModels "Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// StoryContextEngine 故事上下文引擎
// 负责从多个数据源组装 AI 写作所需的三层漏斗上下文：
// Layer 1: 结构化舞台（场景、角色、关系）
// Layer 2: 大纲即摘要（当前章节 + 近期章节 + 卷级结构）
// Layer 3: 最近已写文本（尾部截取）
type StoryContextEngine struct {
	documentRepo       writerRepo.DocumentRepository
	documentContentRepo writerRepo.DocumentContentRepository
	characterRepo      writerRepo.CharacterRepository
	locationRepo       writerRepo.LocationRepository
	outlineRepo        writerRepo.OutlineRepository
}

// NewStoryContextEngine 创建故事上下文引擎
func NewStoryContextEngine(
	documentRepo writerRepo.DocumentRepository,
	documentContentRepo writerRepo.DocumentContentRepository,
	characterRepo writerRepo.CharacterRepository,
	locationRepo writerRepo.LocationRepository,
	outlineRepo writerRepo.OutlineRepository,
) *StoryContextEngine {
	return &StoryContextEngine{
		documentRepo:       documentRepo,
		documentContentRepo: documentContentRepo,
		characterRepo:      characterRepo,
		locationRepo:       locationRepo,
		outlineRepo:        outlineRepo,
	}
}

// BuildStoryContext 构建完整的故事上下文
// projectID: 项目ID（string 格式）
// documentID: 当前文档ID
// mode: 写作模式 continue | rewrite | suggest
// instruction: 用户自定义写作指令
// selectedText: 选中文本（改写模式用）
func (e *StoryContextEngine) BuildStoryContext(
	ctx context.Context,
	projectID, documentID, mode, instruction, selectedText string,
) (*aiModels.StoryContext, error) {
	// 获取当前文档
	doc, err := e.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("获取文档失败: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("文档不存在: %s", documentID)
	}

	// Layer 1: 组装结构化舞台
	stage, err := e.assembleStage(ctx, projectID, doc)
	if err != nil {
		return nil, fmt.Errorf("组装舞台上下文失败: %w", err)
	}

	// Layer 2: 组装大纲上下文
	outlineCtx, err := e.assembleOutline(ctx, projectID, doc)
	if err != nil {
		return nil, fmt.Errorf("组装大纲上下文失败: %w", err)
	}

	// Layer 3: 获取最近已写文本
	recentText, err := e.fetchRecentText(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("获取最近文本失败: %w", err)
	}

	// 组装最终上下文
	sc := &aiModels.StoryContext{
		Stage:          stage,
		OutlineContext: outlineCtx,
		RAGExcerpts:    "", // RAG 检索结果由外部 RAG 服务填充
		RecentText:     recentText,
		Instruction:    instruction,
		Mode:           mode,
		SelectedText:   selectedText,
	}

	// 计算 token 估算
	stageJSON, _ := json.Marshal(stage)
	sc.StageTokens = estimateTokens(string(stageJSON))
	sc.OutlineTokens = estimateTokens(outlineCtx)
	sc.RAGTokens = estimateTokens(sc.RAGExcerpts)
	sc.TotalTokens = sc.StageTokens + sc.OutlineTokens + sc.RAGTokens + estimateTokens(recentText)

	return sc, nil
}

// assembleStage 组装结构化舞台（Layer 1）
// 从文档中提取场景状态、在场角色、地点、关系
func (e *StoryContextEngine) assembleStage(
	ctx context.Context,
	projectID string,
	doc *writerModels.Document,
) (*aiModels.StageContext, error) {
	stage := &aiModels.StageContext{
		SceneGoal:      doc.SceneGoal,
		ActiveConflict: doc.ActiveConflict,
		PlotThreads:    doc.PlotThreads,
		Characters:     make([]*aiModels.StageCharacter, 0),
		Relations:      make([]*aiModels.StageRelation, 0),
	}

	// 收集在场角色的 hex ID，用于后续关系过滤
	characterHexIDs := make(map[string]bool, len(doc.CharacterIDs))

	// 遍历在场角色，获取详细信息
	for _, charID := range doc.CharacterIDs {
		hexID := charID.Hex()
		characterHexIDs[hexID] = true

		char, err := e.characterRepo.FindByID(ctx, hexID)
		if err != nil {
			// 单个角色查询失败不阻断整个流程，跳过即可
			continue
		}
		if char == nil {
			continue
		}

		stage.Characters = append(stage.Characters, &aiModels.StageCharacter{
			Name:          char.Name,
			ShortDesc:     char.ShortDescription,
			CurrentState:  char.CurrentState,
			Personality:   char.PersonalityPrompt,
			SpeechPattern: char.SpeechPattern,
		})
	}

	// 遍历在场地点，获取名称和氛围
	var locationParts []string
	for _, locID := range doc.LocationIDs {
		loc, err := e.locationRepo.FindByID(ctx, locID.Hex())
		if err != nil {
			continue
		}
		if loc == nil {
			continue
		}

		locDesc := loc.Name
		if loc.Atmosphere != "" {
			locDesc += "（" + loc.Atmosphere + "）"
		}
		locationParts = append(locationParts, locDesc)
	}
	if len(locationParts) > 0 {
		stage.Location = strings.Join(locationParts, "、")
	}

	// 获取项目全部角色关系，过滤出在场角色之间的关系
	relations, err := e.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		// 关系查询失败不阻断，返回已有数据
		return stage, nil
	}

	for _, rel := range relations {
		// 只保留两端角色都在场的關係
		if !characterHexIDs[rel.FromID] || !characterHexIDs[rel.ToID] {
			continue
		}

		stage.Relations = append(stage.Relations, &aiModels.StageRelation{
			From:    rel.FromID,
			To:      rel.ToID,
			Type:    string(rel.Type),
			Tension: strengthToTension(rel.Strength),
		})
	}

	return stage, nil
}

// assembleOutline 组装大纲上下文（Layer 2）
// 返回 markdown 格式的大纲摘要，包含当前章节、近期章节、卷级结构
func (e *StoryContextEngine) assembleOutline(
	ctx context.Context,
	projectID string,
	doc *writerModels.Document,
) (string, error) {
	var sb strings.Builder

	// 获取当前章节的大纲节点
	if doc.OutlineNodeID != "" {
		currentNode, err := e.outlineRepo.FindByDocumentID(ctx, doc.ID.Hex())
		if err == nil && currentNode != nil {
			sb.WriteString("## 当前章节\n")
			sb.WriteString(fmt.Sprintf("- 标题: %s\n", currentNode.Title))
			if currentNode.Summary != "" {
				sb.WriteString(fmt.Sprintf("- 摘要: %s\n", currentNode.Summary))
			}
			if currentNode.Tension > 0 {
				sb.WriteString(fmt.Sprintf("- 紧张度: %d/10\n", currentNode.Tension))
			}
			sb.WriteString("\n")
		}
	}

	// 获取所有大纲节点，用于查找近期章节和卷级结构
	allNodes, err := e.outlineRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		// 大纲查询失败，返回已有内容
		return sb.String(), nil
	}

	// 查找当前节点的同级节点（相同 ParentID），取前3个作为近期章节
	if doc.OutlineNodeID != "" {
		currentNode, _ := e.outlineRepo.FindByDocumentID(ctx, doc.ID.Hex())
		if currentNode != nil {
			var siblings []*writerModels.OutlineNode
			for _, node := range allNodes {
				if node.ParentID == currentNode.ParentID && node.ID.Hex() != currentNode.ID.Hex() {
					siblings = append(siblings, node)
				}
			}

			// 限制最多3个近期章节
			recentCount := 0
			sb.WriteString("## 近期章节\n")
			for _, sib := range siblings {
				if recentCount >= 3 {
					break
				}
				sb.WriteString(fmt.Sprintf("- %s", sib.Title))
				if sib.Summary != "" {
					sb.WriteString(fmt.Sprintf(": %s", sib.Summary))
				}
				sb.WriteString("\n")
				recentCount++
			}
			sb.WriteString("\n")
		}
	}

	// 卷级结构：只保留 ParentID 为空的顶层节点
	sb.WriteString("## 卷级结构\n")
	for _, node := range allNodes {
		if node.ParentID == "" {
			sb.WriteString(fmt.Sprintf("- %s\n", node.Title))
		}
	}

	return sb.String(), nil
}

// fetchRecentText 获取最近已写文本（Layer 3）
// 取文档内容的最后 500 字符
func (e *StoryContextEngine) fetchRecentText(ctx context.Context, documentID string) (string, error) {
	// 防御性检查：documentContentRepo 可能为 nil
	if e.documentContentRepo == nil {
		return "", nil
	}

	docContent, err := e.documentContentRepo.GetByDocumentID(ctx, documentID)
	if err != nil {
		return "", nil
	}
	if docContent == nil || docContent.Content == "" {
		return "", nil
	}

	content := docContent.Content
	// 取最后 500 字符
	const maxTailLen = 500
	runes := []rune(content)
	if len(runes) > maxTailLen {
		return string(runes[len(runes)-maxTailLen:]), nil
	}
	return content, nil
}

// BuildPrompt 将 StoryContext 转为最终 prompt string
// 用于传递给 AI 生成服务
func BuildPrompt(sc *aiModels.StoryContext) string {
	var sb strings.Builder

	// 系统指令行
	sb.WriteString("你是一位专业的小说写作助手。请根据以下上下文信息进行创作。\n\n")

	// 大纲上下文
	if sc.OutlineContext != "" {
		sb.WriteString("=== 大纲上下文 ===\n")
		sb.WriteString(sc.OutlineContext)
		sb.WriteString("\n")
	}

	// 场景状态 JSON
	if sc.Stage != nil {
		stageJSON, err := json.Marshal(sc.Stage)
		if err == nil && len(stageJSON) > 2 { // 不是空对象 "{}"
			sb.WriteString("=== 场景状态 ===\n")
			sb.WriteString(string(stageJSON))
			sb.WriteString("\n\n")
		}
	}

	// 参考内容（RAG 检索结果）
	if sc.RAGExcerpts != "" {
		sb.WriteString("=== 参考资料 ===\n")
		sb.WriteString(sc.RAGExcerpts)
		sb.WriteString("\n\n")
	}

	// 最近文本
	if sc.RecentText != "" {
		sb.WriteString("=== 最近文本 ===\n")
		sb.WriteString(sc.RecentText)
		sb.WriteString("\n\n")
	}

	// 根据模式构建不同的写作指令
	switch sc.Mode {
	case "continue":
		sb.WriteString("=== 写作指令 ===\n")
		sb.WriteString("请继续续写以下内容，保持风格和语调一致。")
		if sc.Instruction != "" {
			sb.WriteString(fmt.Sprintf("\n额外要求: %s", sc.Instruction))
		}
		sb.WriteString("\n")

	case "rewrite":
		sb.WriteString("=== 写作指令 ===\n")
		sb.WriteString("请改写以下选中的文本内容。")
		if sc.SelectedText != "" {
			sb.WriteString(fmt.Sprintf("\n选中内容:\n%s", sc.SelectedText))
		}
		if sc.Instruction != "" {
			sb.WriteString(fmt.Sprintf("\n改写要求: %s", sc.Instruction))
		}
		sb.WriteString("\n")

	case "suggest":
		sb.WriteString("=== 写作指令 ===\n")
		sb.WriteString("请基于当前上下文，给出 3 条写作建议。")
		if sc.Instruction != "" {
			sb.WriteString(fmt.Sprintf("\n关注方向: %s", sc.Instruction))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// strengthToTension 将关系强度 (0-100) 转为紧张度等级
func strengthToTension(strength int) string {
	if strength >= 75 {
		return "high"
	}
	if strength >= 40 {
		return "medium"
	}
	return "low"
}

// estimateTokens 估算中文文本的 token 数
// 中文约 1.5 字/token
func estimateTokens(text string) int {
	if text == "" {
		return 0
	}
	runeCount := len([]rune(text))
	// 每 2 个字符约 3 个 token
	return runeCount * 3 / 2
}
