package mongodb_test

import (
	"context"
	"strings"
	"testing"

	"Qingyu_backend/models/bookstore"
	mongodb "Qingyu_backend/repository/mongodb/bookstore"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoChapterContentRepository_BatchCreateSetsIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoChapterContentRepository(db)
	ctx := context.Background()

	contents := []*bookstore.ChapterContent{
		{
			ChapterID: primitive.NewObjectID(),
			Content:   "chapter content 1",
		},
		{
			ChapterID: primitive.NewObjectID(),
			Content:   "chapter content 2",
		},
	}

	err := repo.BatchCreate(ctx, contents)
	require.NoError(t, err)
	for _, content := range contents {
		assert.False(t, content.ID.IsZero())
	}
}

func TestMongoRankingRepository_UpdateRankingsSetsIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongodb.NewMongoRankingRepository(db.Client(), db.Name())
	ctx := context.Background()

	items := []*bookstore.RankingItem{
		{
			BookID: primitive.NewObjectID(),
			Type:   bookstore.RankingTypeWeekly,
			Rank:   1,
			Score:  9.9,
			Period: "2026-W10",
		},
		{
			BookID: primitive.NewObjectID(),
			Type:   bookstore.RankingTypeWeekly,
			Rank:   2,
			Score:  8.8,
			Period: "2026-W10",
		},
	}

	err := repo.UpdateRankings(ctx, bookstore.RankingTypeWeekly, "2026-W10", items)
	if err != nil {
		if strings.Contains(err.Error(), "Transaction numbers are only allowed on a replica set member or mongos") {
			t.Skip("mongodb test environment does not support transactions")
		}
		require.NoError(t, err)
	}
	for _, item := range items {
		assert.False(t, item.ID.IsZero())
	}
}
