package finance_test

import (
	"context"
	"testing"

	financeModel "Qingyu_backend/models/finance"
	financeRepo "Qingyu_backend/repository/mongodb/finance"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMembershipRepository_BatchCreateMembershipCardsSetsIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := financeRepo.NewMembershipRepository(db)
	ctx := context.Background()

	planID := primitive.NewObjectID()
	cards := []*financeModel.MembershipCard{
		{
			Code:      "card-test-1",
			PlanID:    planID,
			PlanType:  financeModel.MembershipTypeMonthly,
			Duration:  30,
			BatchID:   "batch-1",
			Status:    financeModel.CardStatusUnused,
			CreatedBy: "tester",
		},
		{
			Code:      "card-test-2",
			PlanID:    planID,
			PlanType:  financeModel.MembershipTypeMonthly,
			Duration:  30,
			BatchID:   "batch-1",
			Status:    financeModel.CardStatusUnused,
			CreatedBy: "tester",
		},
	}

	err := repo.BatchCreateMembershipCards(ctx, cards)
	require.NoError(t, err)
	for _, card := range cards {
		assert.False(t, card.ID.IsZero())
	}
}
