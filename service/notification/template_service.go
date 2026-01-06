package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/notification"
	repo "Qingyu_backend/repository/interfaces/notification"
	"Qingyu_backend/pkg/errors"
)

// TemplateService 通知模板服务接口
type TemplateService interface {
	// 模板基础操作
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*notification.NotificationTemplate, error)
	GetTemplate(ctx context.Context, id string) (*notification.NotificationTemplate, error)
	ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error)
	UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, id string) error

	// 模板使用
	RenderTemplate(ctx context.Context, templateType notification.NotificationType, action string, variables map[string]interface{}, language string) (string, string, error)

	// 预置模板管理
	InitializeDefaultTemplates(ctx context.Context) error
}

// templateServiceImpl 模板服务实现
type templateServiceImpl struct {
	templateRepo repo.NotificationTemplateRepository
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(templateRepo repo.NotificationTemplateRepository) TemplateService {
	return &templateServiceImpl{
		templateRepo: templateRepo,
	}
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Type      notification.NotificationType `json:"type" validate:"required"`
	Action    string                         `json:"action" validate:"required,min=1,max=100"`
	Title     string                         `json:"title" validate:"required,min=1,max=200"`
	Content   string                         `json:"content" validate:"required,min=1,max=5000"`
	Variables []string                       `json:"variables"`
	Data      map[string]interface{}         `json:"data"`
	Language  string                         `json:"language" validate:"required"`
	IsActive  bool                           `json:"isActive"`
}

// ListTemplatesRequest 获取模板列表请求
type ListTemplatesRequest struct {
	Type     *notification.NotificationType `json:"type"`
	Action   *string                        `json:"action"`
	Language *string                        `json:"language"`
	IsActive *bool                          `json:"isActive"`
	Limit    int                            `json:"limit" validate:"min=1,max=100"`
	Offset   int                            `json:"offset" validate:"min=0"`
}

// ListTemplatesResponse 获取模板列表响应
type ListTemplatesResponse struct {
	Templates []*notification.NotificationTemplate `json:"templates"`
	Total     int64                                `json:"total"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Title     *string                `json:"title" validate:"omitempty,min=1,max=200"`
	Content   *string                `json:"content" validate:"omitempty,min=1,max=5000"`
	Variables *[]string              `json:"variables"`
	Data      *map[string]interface{} `json:"data"`
	IsActive  *bool                  `json:"isActive"`
}

// CreateTemplate 创建模板
func (s *templateServiceImpl) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*notification.NotificationTemplate, error) {
	now := time.Now()
	template := &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      req.Type,
		Action:    req.Action,
		Title:     req.Title,
		Content:   req.Content,
		Variables: req.Variables,
		Data:      req.Data,
		Language:  req.Language,
		IsActive:  req.IsActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("TEMPLATE_CREATE_FAILED", "创建模板失败", err)
	}

	return template, nil
}

// GetTemplate 获取模板详情
func (s *templateServiceImpl) GetTemplate(ctx context.Context, id string) (*notification.NotificationTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("TEMPLATE_GET_FAILED", "获取模板失败", err)
	}

	if template == nil {
		return nil, errors.BookstoreServiceFactory.NotFoundError("NotificationTemplate", id)
	}

	return template, nil
}

// ListTemplates 获取模板列表
func (s *templateServiceImpl) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 构建筛选条件
	filter := &repo.TemplateFilter{
		Type:     req.Type,
		Action:   req.Action,
		Language: req.Language,
		IsActive: req.IsActive,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	// 获取模板列表
	templates, err := s.templateRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("TEMPLATE_LIST_FAILED", "获取模板列表失败", err)
	}

	// 计算总数
	total := int64(len(templates))

	return &ListTemplatesResponse{
		Templates: templates,
		Total:     total,
	}, nil
}

// UpdateTemplate 更新模板
func (s *templateServiceImpl) UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) error {
	// 检查模板是否存在
	exists, err := s.templateRepo.Exists(ctx, id)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("TEMPLATE_CHECK_FAILED", "检查模板是否存在失败", err)
	}
	if !exists {
		return errors.BookstoreServiceFactory.NotFoundError("NotificationTemplate", id)
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Variables != nil {
		updates["variables"] = *req.Variables
	}
	if req.Data != nil {
		updates["data"] = *req.Data
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	updates["updated_at"] = time.Now()

	if len(updates) == 0 {
		return nil // 没有更新
	}

	if err := s.templateRepo.Update(ctx, id, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("TEMPLATE_UPDATE_FAILED", "更新模板失败", err)
	}

	return nil
}

// DeleteTemplate 删除模板
func (s *templateServiceImpl) DeleteTemplate(ctx context.Context, id string) error {
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		return errors.BookstoreServiceFactory.InternalError("TEMPLATE_DELETE_FAILED", "删除模板失败", err)
	}
	return nil
}

// RenderTemplate 渲染模板
func (s *templateServiceImpl) RenderTemplate(ctx context.Context, templateType notification.NotificationType, action string, variables map[string]interface{}, language string) (string, string, error) {
	// 获取模板
	if language == "" {
		language = "zh-CN" // 默认中文
	}

	template, err := s.templateRepo.GetActiveTemplate(ctx, templateType, action, language)
	if err != nil {
		return "", "", errors.BookstoreServiceFactory.InternalError("TEMPLATE_GET_FAILED", "获取模板失败", err)
	}

	if template == nil {
		return "", "", errors.BookstoreServiceFactory.NotFoundError("NotificationTemplate", fmt.Sprintf("%s:%s:%s", templateType, action, language))
	}

	// 替换变量
	title := s.replaceVariables(template.Title, variables)
	content := s.replaceVariables(template.Content, variables)

	return title, content, nil
}

// replaceVariables 替换模板变量
func (s *templateServiceImpl) replaceVariables(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int, int32, int64:
			valueStr = fmt.Sprintf("%d", v)
		case float32, float64:
			valueStr = fmt.Sprintf("%.2f", v)
		default:
			// 尝试序列化为JSON
			if jsonBytes, err := json.Marshal(v); err == nil {
				valueStr = string(jsonBytes)
			} else {
				valueStr = fmt.Sprintf("%v", v)
			}
		}
		result = replaceAll(result, placeholder, valueStr)
	}
	return result
}

// replaceAll 替换所有匹配的字符串
func replaceAll(s, old, new string) string {
	result := s
	for {
		idx := indexOf(result, old)
		if idx == -1 {
			break
		}
		result = result[:idx] + new + result[idx+len(old):]
	}
	return result
}

// indexOf 查找子字符串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// InitializeDefaultTemplates 初始化默认模板
func (s *templateServiceImpl) InitializeDefaultTemplates(ctx context.Context) error {
	templates := s.getDefaultTemplates()

	for _, template := range templates {
		// 检查模板是否已存在
		existingTemplates, err := s.templateRepo.GetByTypeAndAction(ctx, template.Type, template.Action)
		if err != nil {
			continue
		}

		// 检查是否已有相同语言的模板
		exists := false
		for _, existing := range existingTemplates {
			if existing.Language == template.Language {
				exists = true
				break
			}
		}

		if !exists {
			// 创建新模板
			if err := s.templateRepo.Create(ctx, template); err != nil {
				continue
			}
		}
	}

	return nil
}

// getDefaultTemplates 获取默认模板列表
func (s *templateServiceImpl) getDefaultTemplates() []*notification.NotificationTemplate {
	now := time.Now()
	templates := make([]*notification.NotificationTemplate, 0)

	// 系统通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeSystem,
		Action:    "announcement",
		Title:     "平台公告",
		Content:   "{{title}}\n\n{{content}}",
		Variables: []string{"title", "content"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeSystem,
		Action:    "maintenance",
		Title:     "系统维护通知",
		Content:   "尊敬的用户，系统将于{{startTime}}至{{endTime}}进行维护，期间部分功能可能无法使用，敬请谅解。",
		Variables: []string{"startTime", "endTime"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 社交通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeSocial,
		Action:    "follow",
		Title:     "您有新的关注者",
		Content:   "{{followerName}}关注了您，点击查看TA的主页。",
		Variables: []string{"followerName", "followerId"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeSocial,
		Action:    "like",
		Title:     "作品收到点赞",
		Content:   "{{likerName}}点赞了您的作品《{{bookTitle}}》。",
		Variables: []string{"likerName", "bookTitle", "bookId"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeSocial,
		Action:    "comment",
		Title:     "作品收到新评论",
		Content:   "{{commenterName}}评论了您的作品《{{bookTitle}}》：{{commentContent}}",
		Variables: []string{"commenterName", "bookTitle", "bookId", "commentContent"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 内容通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeContent,
		Action:    "review_approved",
		Title:     "作品审核通过",
		Content:   "恭喜！您的作品《{{bookTitle}}》已通过审核，现已上架。",
		Variables: []string{"bookTitle", "bookId"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeContent,
		Action:    "review_rejected",
		Title:     "作品审核未通过",
		Content:   "很遗憾，您的作品《{{bookTitle}}》未通过审核。原因：{{reason}}",
		Variables: []string{"bookTitle", "bookId", "reason"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeContent,
		Action:    "book_offline",
		Title:     "作品下架通知",
		Content:   "您的作品《{{bookTitle}}》已被下架。原因：{{reason}}",
		Variables: []string{"bookTitle", "bookId", "reason"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 打赏通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeReward,
		Action:    "received",
		Title:     "收到打赏",
		Content:   "{{senderName}}打赏了您的作品《{{bookTitle}}》，金额：{{amount}}书币。",
		Variables: []string{"senderName", "bookTitle", "bookId", "amount"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 私信通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeMessage,
		Action:    "received",
		Title:     "收到新私信",
		Content:   "{{senderName}}给您发送了一条私信，点击查看。",
		Variables: []string{"senderName", "senderId"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 更新通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeUpdate,
		Action:    "chapter_update",
		Title:     "关注作品更新",
		Content:   "您关注的《{{bookTitle}}》更新了第{{chapterNumber}}章：{{chapterTitle}}",
		Variables: []string{"bookTitle", "bookId", "chapterNumber", "chapterTitle", "chapterId"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	// 会员通知模板
	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeMembership,
		Action:    "expiring_soon",
		Title:     "会员即将到期",
		Content:   "您的会员将在{{days}}天后到期，请及时续费以享受会员权益。",
		Variables: []string{"days", "expireDate"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeMembership,
		Action:    "expired",
		Title:     "会员已到期",
		Content:   "您的会员已到期，续费后可继续享受会员权益。",
		Variables: []string{"expireDate"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	templates = append(templates, &notification.NotificationTemplate{
		ID:        primitive.NewObjectID().Hex(),
		Type:      notification.NotificationTypeMembership,
		Action:    "renewed",
		Title:     "会员续费成功",
		Content:   "您的会员已成功续费，有效期至{{expireDate}}。",
		Variables: []string{"expireDate"},
		Language:  "zh-CN",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return templates
}
