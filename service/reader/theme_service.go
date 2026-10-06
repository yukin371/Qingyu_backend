package reader

import (
	"context"
	"fmt"
	"strings"
	"time"

	readerModel "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
)

// ThemeService 是阅读主题的业务 owner，负责主题可见性、owner 校验和激活同步。
type ThemeService interface {
	GetThemes(ctx context.Context, userID string, builtinOnly, publicOnly bool) ([]*readerModel.ReaderTheme, error)
	GetThemeByName(ctx context.Context, userID, name string) (*readerModel.ReaderTheme, error)
	GetActiveTheme(ctx context.Context, userID string) (*readerModel.ReaderTheme, error)
	CreateCustomTheme(ctx context.Context, userID string, req *readerModel.CreateCustomThemeRequest) (*readerModel.ReaderTheme, error)
	UpdateTheme(ctx context.Context, userID, themeID string, req *readerModel.UpdateThemeRequest) error
	DeleteTheme(ctx context.Context, userID, themeID string) error
	ActivateTheme(ctx context.Context, userID, themeName string) error
}

// ReaderThemeService 阅读主题服务实现。
type ReaderThemeService struct {
	themeRepo    readerRepo.ReaderThemeRepository
	settingsRepo readerRepo.ReadingSettingsRepository
}

func NewReaderThemeService(themeRepo readerRepo.ReaderThemeRepository, settingsRepo readerRepo.ReadingSettingsRepository) *ReaderThemeService {
	return &ReaderThemeService{themeRepo: themeRepo, settingsRepo: settingsRepo}
}

func (s *ReaderThemeService) GetThemes(ctx context.Context, userID string, builtinOnly, publicOnly bool) ([]*readerModel.ReaderTheme, error) {
	if s.themeRepo == nil {
		return nil, fmt.Errorf("主题仓储未初始化")
	}
	var themes []*readerModel.ReaderTheme
	if publicOnly {
		publicThemes, _, err := s.themeRepo.GetPublicThemes(ctx, 1, 1000)
		if err != nil {
			return nil, err
		}
		themes = appendUniqueThemes(publicThemes, readerModel.BuiltInThemes...)
	} else {
		builtins, err := s.themeRepo.GetBuiltInThemes(ctx)
		if err != nil {
			return nil, err
		}
		themes = append(themes, builtins...)
		if userID != "" && !builtinOnly {
			userThemes, err := s.themeRepo.GetUserThemes(ctx, userID)
			if err != nil {
				return nil, err
			}
			themes = appendUniqueThemes(themes, userThemes...)
		}
	}
	if builtinOnly {
		filtered := themes[:0]
		for _, theme := range themes {
			if theme != nil && theme.IsBuiltIn {
				filtered = append(filtered, theme)
			}
		}
		themes = filtered
	}
	return themes, nil
}

func (s *ReaderThemeService) GetThemeByName(ctx context.Context, userID, name string) (*readerModel.ReaderTheme, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("主题名称不能为空")
	}
	theme, err := s.themeRepo.GetThemeByName(ctx, name)
	if err == nil && theme != nil {
		if theme.IsBuiltIn || theme.IsPublic || (userID != "" && theme.CreatorID == userID) {
			return theme, nil
		}
		return nil, fmt.Errorf("主题不存在: %s", name)
	}
	for _, builtin := range readerModel.BuiltInThemes {
		if builtin.Name == name {
			return builtin, nil
		}
	}
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("主题不存在: %s", name)
}

func (s *ReaderThemeService) GetActiveTheme(ctx context.Context, userID string) (*readerModel.ReaderTheme, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if s.themeRepo == nil {
		return nil, fmt.Errorf("主题仓储未初始化")
	}
	return s.themeRepo.GetActiveTheme(ctx, userID)
}

func (s *ReaderThemeService) CreateCustomTheme(ctx context.Context, userID string, req *readerModel.CreateCustomThemeRequest) (*readerModel.ReaderTheme, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if req == nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.DisplayName) == "" {
		return nil, fmt.Errorf("主题名称和显示名称不能为空")
	}
	if _, err := s.themeRepo.GetThemeByName(ctx, req.Name); err == nil {
		return nil, fmt.Errorf("主题名称已存在: %s", req.Name)
	} else if !strings.Contains(err.Error(), "不存在") {
		return nil, err
	}
	theme := &readerModel.ReaderTheme{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsBuiltIn:   false,
		IsPublic:    req.IsPublic,
		CreatorID:   userID,
		Colors:      req.Colors,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.themeRepo.CreateTheme(ctx, theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func (s *ReaderThemeService) UpdateTheme(ctx context.Context, userID, themeID string, req *readerModel.UpdateThemeRequest) error {
	theme, err := s.themeRepo.GetTheme(ctx, themeID)
	if err != nil {
		return err
	}
	if err := ensureThemeOwner(userID, theme); err != nil {
		return err
	}
	if req == nil {
		return fmt.Errorf("主题更新内容不能为空")
	}
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.Colors != nil {
		updates["colors"] = *req.Colors
	}
	if len(updates) == 0 {
		return fmt.Errorf("主题更新内容不能为空")
	}
	return s.themeRepo.UpdateTheme(ctx, themeID, updates)
}

func (s *ReaderThemeService) DeleteTheme(ctx context.Context, userID, themeID string) error {
	theme, err := s.themeRepo.GetTheme(ctx, themeID)
	if err != nil {
		return err
	}
	if err := ensureThemeOwner(userID, theme); err != nil {
		return err
	}
	return s.themeRepo.DeleteTheme(ctx, themeID)
}

func (s *ReaderThemeService) ActivateTheme(ctx context.Context, userID, themeName string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	theme, err := s.GetThemeByName(ctx, userID, themeName)
	if err != nil {
		// The frontend uses the canonical name for built-ins and the stable ID
		// for custom themes. Resolve both forms without weakening owner checks.
		theme, err = s.themeRepo.GetTheme(ctx, themeName)
		if err != nil {
			return err
		}
	}
	if !theme.IsBuiltIn && theme.CreatorID != userID {
		return fmt.Errorf("主题不存在")
	}
	// 空 ID 表示仓储回退到了内置常量；此时没有可更新的主题文档，
	// 仍需把用户选择写入阅读设置，待初始化任务补齐内置主题文档后再创建副本。
	if !theme.ID.IsZero() {
		if err := s.themeRepo.SetActiveTheme(ctx, userID, theme.ID.Hex()); err != nil {
			return err
		}
	}
	if s.settingsRepo == nil {
		return fmt.Errorf("阅读设置仓储未初始化")
	}
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if settings == nil {
		settings = &readerModel.ReadingSettings{UserID: userID, FontSize: 16, FontFamily: "serif", LineHeight: 1.5, PageMode: 1}
	}
	settings.Theme = theme.Name
	settings.Background = theme.Colors.Background
	return s.settingsRepo.UpdateByUserID(ctx, userID, settings)
}

func ensureThemeOwner(userID string, theme *readerModel.ReaderTheme) error {
	if theme == nil || theme.IsBuiltIn || theme.CreatorID != userID {
		return fmt.Errorf("主题不存在")
	}
	return nil
}

func appendUniqueThemes(themes []*readerModel.ReaderTheme, additions ...*readerModel.ReaderTheme) []*readerModel.ReaderTheme {
	seen := make(map[string]struct{}, len(themes))
	for _, theme := range themes {
		if theme != nil {
			seen[theme.Name] = struct{}{}
		}
	}
	for _, theme := range additions {
		if theme != nil {
			if _, exists := seen[theme.Name]; exists {
				continue
			}
			seen[theme.Name] = struct{}{}
			themes = append(themes, theme)
		}
	}
	return themes
}
