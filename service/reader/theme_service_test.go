package reader

import (
	"context"
	"errors"
	"testing"

	readerModel "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fakeThemeRepository struct {
	readerRepo.ReaderThemeRepository
	themes      map[string]*readerModel.ReaderTheme
	updatedID   string
	deletedID   string
	activeUser  string
	activeTheme string
}

func newFakeThemeRepository() *fakeThemeRepository {
	themes := make(map[string]*readerModel.ReaderTheme)
	for _, theme := range readerModel.BuiltInThemes {
		copyTheme := *theme
		copyTheme.ID = primitive.NewObjectID()
		themes[copyTheme.ID.Hex()] = &copyTheme
	}
	return &fakeThemeRepository{themes: themes}
}

func (f *fakeThemeRepository) GetBuiltInThemes(context.Context) ([]*readerModel.ReaderTheme, error) {
	return readerModel.BuiltInThemes, nil
}

func (f *fakeThemeRepository) GetThemeByName(_ context.Context, name string) (*readerModel.ReaderTheme, error) {
	for _, theme := range f.themes {
		if theme.Name == name {
			return theme, nil
		}
	}
	return nil, errors.New("主题不存在")
}

func (f *fakeThemeRepository) CreateTheme(_ context.Context, theme *readerModel.ReaderTheme) error {
	if theme.ID.IsZero() {
		theme.ID = primitive.NewObjectID()
	}
	f.themes[theme.ID.Hex()] = theme
	return nil
}

func (f *fakeThemeRepository) GetTheme(_ context.Context, id string) (*readerModel.ReaderTheme, error) {
	theme, ok := f.themes[id]
	if !ok {
		return nil, errors.New("主题不存在")
	}
	return theme, nil
}

func (f *fakeThemeRepository) UpdateTheme(_ context.Context, id string, _ map[string]interface{}) error {
	f.updatedID = id
	return nil
}

func (f *fakeThemeRepository) DeleteTheme(_ context.Context, id string) error {
	f.deletedID = id
	return nil
}

func (f *fakeThemeRepository) GetUserThemes(_ context.Context, userID string) ([]*readerModel.ReaderTheme, error) {
	var themes []*readerModel.ReaderTheme
	for _, theme := range f.themes {
		if theme.CreatorID == userID || theme.IsBuiltIn {
			themes = append(themes, theme)
		}
	}
	return themes, nil
}

func (f *fakeThemeRepository) SetActiveTheme(_ context.Context, userID, themeID string) error {
	f.activeUser = userID
	f.activeTheme = themeID
	return nil
}

type fakeSettingsRepository struct {
	readerRepo.ReadingSettingsRepository
	settings *readerModel.ReadingSettings
}

func (f *fakeSettingsRepository) GetByUserID(context.Context, string) (*readerModel.ReadingSettings, error) {
	return f.settings, nil
}

func (f *fakeSettingsRepository) UpdateByUserID(_ context.Context, _ string, settings *readerModel.ReadingSettings) error {
	f.settings = settings
	return nil
}

func TestReaderThemeService_UpdateRejectsForeignOwner(t *testing.T) {
	repo := newFakeThemeRepository()
	foreign := &readerModel.ReaderTheme{ID: primitive.NewObjectID(), Name: "foreign", CreatorID: "user-a"}
	repo.themes[foreign.ID.Hex()] = foreign
	service := NewReaderThemeService(repo, &fakeSettingsRepository{})

	err := service.UpdateTheme(context.Background(), "user-b", foreign.ID.Hex(), &readerModel.UpdateThemeRequest{})

	require.Error(t, err)
	require.Empty(t, repo.updatedID)
}

func TestReaderThemeService_DeleteRejectsBuiltIn(t *testing.T) {
	repo := newFakeThemeRepository()
	builtin := readerModel.BuiltInThemes[0]
	builtinCopy := *builtin
	builtinCopy.ID = primitive.NewObjectID()
	repo.themes[builtinCopy.ID.Hex()] = &builtinCopy
	service := NewReaderThemeService(repo, &fakeSettingsRepository{})

	err := service.DeleteTheme(context.Background(), "user-a", builtinCopy.ID.Hex())

	require.Error(t, err)
	require.Empty(t, repo.deletedID)
}

func TestReaderThemeService_ActivateBuiltInPersistsSettings(t *testing.T) {
	repo := newFakeThemeRepository()
	light, err := repo.GetThemeByName(context.Background(), "light")
	require.NoError(t, err)
	settingsRepo := &fakeSettingsRepository{settings: &readerModel.ReadingSettings{UserID: "user-a"}}
	service := NewReaderThemeService(repo, settingsRepo)

	err = service.ActivateTheme(context.Background(), "user-a", "light")

	require.NoError(t, err)
	require.Equal(t, "user-a", repo.activeUser)
	require.Equal(t, light.ID.Hex(), repo.activeTheme)
	require.Equal(t, "light", settingsRepo.settings.Theme)
	require.Equal(t, light.Colors.Background, settingsRepo.settings.Background)
}

func TestReaderThemeService_GetThemeByNameHidesPrivateForeignTheme(t *testing.T) {
	repo := newFakeThemeRepository()
	privateTheme := &readerModel.ReaderTheme{ID: primitive.NewObjectID(), Name: "private", CreatorID: "user-a", IsPublic: false}
	repo.themes[privateTheme.ID.Hex()] = privateTheme
	service := NewReaderThemeService(repo, &fakeSettingsRepository{})

	_, err := service.GetThemeByName(context.Background(), "user-b", privateTheme.Name)

	require.Error(t, err)
	require.Contains(t, err.Error(), "主题不存在")
}

var _ readerRepo.ReaderThemeRepository = (*fakeThemeRepository)(nil)
var _ readerRepo.ReadingSettingsRepository = (*fakeSettingsRepository)(nil)
