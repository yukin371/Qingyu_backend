package recommendation

import (
	"context"
	"testing"

	reco "Qingyu_backend/models/recommendation"

	"github.com/stretchr/testify/assert"
)

type stubTableRepo struct {
	listPage     int
	listPageSize int
	updateFields map[string]interface{}
	tableByID    *reco.RecommendationTable
}

func (s *stubTableRepo) Create(ctx context.Context, table *reco.RecommendationTable) error {
	return nil
}
func (s *stubTableRepo) Delete(ctx context.Context, id string) error { return nil }
func (s *stubTableRepo) GetByID(ctx context.Context, id string) (*reco.RecommendationTable, error) {
	return s.tableByID, nil
}
func (s *stubTableRepo) GetByTypePeriod(ctx context.Context, tableType reco.TableType, period string, source reco.TableSource) (*reco.RecommendationTable, error) {
	return nil, nil
}
func (s *stubTableRepo) List(ctx context.Context, tableType *reco.TableType, source *reco.TableSource, page, pageSize int) ([]*reco.RecommendationTable, int64, error) {
	s.listPage = page
	s.listPageSize = pageSize
	return nil, 0, nil
}
func (s *stubTableRepo) UpsertByTypePeriod(ctx context.Context, table *reco.RecommendationTable) error {
	return nil
}
func (s *stubTableRepo) Health(ctx context.Context) error { return nil }
func (s *stubTableRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	s.updateFields = updates
	return nil
}

func TestRecommendationTableService_ListTables_NormalizesPagination(t *testing.T) {
	repo := &stubTableRepo{}
	svc := NewRecommendationTableService(repo)

	_, _, err := svc.ListTables(context.Background(), nil, nil, 0, 999)
	assert.NoError(t, err)
	assert.Equal(t, 1, repo.listPage)
	assert.Equal(t, 20, repo.listPageSize)
}

func TestRecommendationTableService_UpdateManualTable_DoesNotSetUpdatedAtInService(t *testing.T) {
	repo := &stubTableRepo{
		tableByID: &reco.RecommendationTable{
			Source: reco.TableSourceManual,
			Status: reco.TableStatusActive,
		},
	}
	svc := NewRecommendationTableService(repo)

	err := svc.UpdateManualTable(context.Background(), "65f000000000000000000001", "自定义榜", "custom", []reco.TableItem{
		{BookID: "b1", Rank: 1},
	}, reco.TableStatusActive, "admin-1")
	assert.NoError(t, err)

	_, exists := repo.updateFields["updated_at"]
	assert.False(t, exists)
	assert.Equal(t, "admin-1", repo.updateFields["updated_by"])
}
