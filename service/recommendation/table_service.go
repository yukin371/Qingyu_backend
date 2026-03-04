package recommendation

import (
	reco "Qingyu_backend/models/recommendation"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
	"context"
	"fmt"
	"time"
)

type RecommendationTableService interface {
	GetTable(ctx context.Context, id string) (*reco.RecommendationTable, error)
	ListTables(ctx context.Context, tableType *reco.TableType, source *reco.TableSource, page, pageSize int) ([]*reco.RecommendationTable, int64, error)
	UpsertAutoTable(ctx context.Context, tableType reco.TableType, period string, items []reco.TableItem, updatedBy string) error
	CreateManualTable(ctx context.Context, name, period string, items []reco.TableItem, updatedBy string) error
	UpdateManualTable(ctx context.Context, id, name, period string, items []reco.TableItem, status reco.TableStatus, updatedBy string) error
	DeleteTable(ctx context.Context, id string) error
}

type RecommendationTableServiceImpl struct {
	repo recoRepo.TableRepository
}

func NewRecommendationTableService(repo recoRepo.TableRepository) RecommendationTableService {
	return &RecommendationTableServiceImpl{repo: repo}
}

func (s *RecommendationTableServiceImpl) GetTable(ctx context.Context, id string) (*reco.RecommendationTable, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RecommendationTableServiceImpl) ListTables(ctx context.Context, tableType *reco.TableType, source *reco.TableSource, page, pageSize int) ([]*reco.RecommendationTable, int64, error) {
	page, pageSize = normalizePagination(page, pageSize)
	return s.repo.List(ctx, tableType, source, page, pageSize)
}

func (s *RecommendationTableServiceImpl) UpsertAutoTable(ctx context.Context, tableType reco.TableType, period string, items []reco.TableItem, updatedBy string) error {
	if tableType != reco.TableTypeWeekly && tableType != reco.TableTypeMonthly && tableType != reco.TableTypeMonthlyVote {
		return fmt.Errorf("auto table type not supported: %s", tableType)
	}
	if period == "" {
		return fmt.Errorf("period is required")
	}

	assignRanks(items)
	table := &reco.RecommendationTable{
		Name:      defaultTableName(tableType),
		Type:      tableType,
		Period:    period,
		Source:    reco.TableSourceAuto,
		Status:    reco.TableStatusActive,
		Items:     items,
		UpdatedBy: updatedBy,
		Metadata: map[string]interface{}{
			"refreshMode": "scheduled",
		},
		UpdatedAt: time.Now(),
	}
	return s.repo.UpsertByTypePeriod(ctx, table)
}

func (s *RecommendationTableServiceImpl) CreateManualTable(ctx context.Context, name, period string, items []reco.TableItem, updatedBy string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if period == "" {
		period = "custom"
	}
	assignRanks(items)
	table := &reco.RecommendationTable{
		Name:      name,
		Type:      reco.TableTypeManual,
		Period:    period,
		Source:    reco.TableSourceManual,
		Status:    reco.TableStatusActive,
		Items:     items,
		UpdatedBy: updatedBy,
		Metadata: map[string]interface{}{
			"editable": true,
		},
	}
	return s.repo.Create(ctx, table)
}

func (s *RecommendationTableServiceImpl) UpdateManualTable(ctx context.Context, id, name, period string, items []reco.TableItem, status reco.TableStatus, updatedBy string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("table not found")
	}
	if existing.Source != reco.TableSourceManual {
		return fmt.Errorf("only manual table can be updated")
	}
	if status == "" {
		status = existing.Status
	}
	assignRanks(items)
	updates := map[string]interface{}{
		"name":       name,
		"period":     period,
		"items":      items,
		"status":     status,
		"updated_by": updatedBy,
	}
	return s.repo.Update(ctx, id, updates)
}

func (s *RecommendationTableServiceImpl) DeleteTable(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func defaultTableName(tableType reco.TableType) string {
	switch tableType {
	case reco.TableTypeWeekly:
		return "周榜"
	case reco.TableTypeMonthly:
		return "月榜"
	case reco.TableTypeMonthlyVote:
		return "月票榜"
	default:
		return "推荐榜"
	}
}

func assignRanks(items []reco.TableItem) {
	for i := range items {
		if items[i].Rank <= 0 {
			items[i].Rank = i + 1
		}
		if items[i].AddedAt <= 0 {
			items[i].AddedAt = time.Now().Unix()
		}
	}
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}
