package writer_test

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTemplateRepo(t *testing.T) (*writerRepo.MongoTemplateRepository, context.Context, func()) {
	t.Helper()
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewMongoTemplateRepository(db).(*writerRepo.MongoTemplateRepository)
	ctx := context.Background()
	return repo, ctx, func() {
		_ = db.Collection("templates").Drop(ctx)
		cleanup()
	}
}

func TestTemplateRepository_CreateSetsID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, ctx, cleanup := setupTemplateRepo(t)
	defer cleanup()

	template := &writer.Template{
		Name:        "测试模板",
		Description: "用于验证 Create 会回设 ID",
		Type:        writer.TemplateTypeChapter,
		Category:    "test",
		Content:     "内容 {{var.title}}",
		CreatedBy:   "tester",
	}

	err := repo.Create(ctx, template)
	require.NoError(t, err)
	assert.False(t, template.ID.IsZero())

	saved, err := repo.GetByID(ctx, template.ID)
	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Equal(t, template.ID, saved.ID)
	assert.Equal(t, template.Name, saved.Name)
}
