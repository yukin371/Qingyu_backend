package recommendation

import (
	recModel "Qingyu_backend/models/recommendation"
	"context"
	"testing"
)

type stubRecommendationRepository struct {
	recorded  *recModel.UserBehavior
	behaviors []*recModel.UserBehavior
}

func (s *stubRecommendationRepository) RecordBehavior(_ context.Context, behavior *recModel.UserBehavior) error {
	s.recorded = behavior
	return nil
}

func (s *stubRecommendationRepository) GetUserBehaviors(_ context.Context, _ string, _ int) ([]*recModel.UserBehavior, error) {
	return s.behaviors, nil
}

func (s *stubRecommendationRepository) GetItemBehaviors(_ context.Context, _ string, _ int) ([]*recModel.UserBehavior, error) {
	return nil, nil
}

func (s *stubRecommendationRepository) Health(_ context.Context) error {
	return nil
}

func TestRecordUserBehavior_NormalizesLegacyFavorite(t *testing.T) {
	repo := &stubRecommendationRepository{}
	service := NewRecommendationService(repo, nil)

	err := service.RecordUserBehavior(context.Background(), &RecordBehaviorRequest{
		UserID:     "user-1",
		ItemID:     "book-1",
		ItemType:   "book",
		ActionType: "favorite",
	})
	if err != nil {
		t.Fatalf("expected record behavior to succeed: %v", err)
	}
	if repo.recorded == nil || repo.recorded.ActionType != recModel.ActionTypeCollect {
		t.Fatalf("expected favorite to normalize to %s, got %+v", recModel.ActionTypeCollect, repo.recorded)
	}
}

func TestGetUserBehaviors_NormalizesLegacyFavorite(t *testing.T) {
	repo := &stubRecommendationRepository{
		behaviors: []*recModel.UserBehavior{{
			ID:         "1",
			UserID:     "user-1",
			ItemID:     "book-1",
			ItemType:   "book",
			ActionType: recModel.ActionTypeFavorite,
		}},
	}
	service := NewRecommendationService(repo, nil)

	behaviors, err := service.GetUserBehaviors(context.Background(), "user-1", 10)
	if err != nil {
		t.Fatalf("expected get user behaviors to succeed: %v", err)
	}
	if len(behaviors) != 1 || behaviors[0].ActionType != recModel.ActionTypeCollect {
		t.Fatalf("expected favorite to normalize to %s, got %+v", recModel.ActionTypeCollect, behaviors)
	}
}

func TestRecordUserBehavior_RejectsUnknownBehaviorType(t *testing.T) {
	repo := &stubRecommendationRepository{}
	service := NewRecommendationService(repo, nil)

	err := service.RecordUserBehavior(context.Background(), &RecordBehaviorRequest{
		UserID:     "user-1",
		ItemID:     "book-1",
		ItemType:   "book",
		ActionType: "unknown",
	})
	if err == nil {
		t.Fatal("expected invalid behavior type to fail")
	}
}
