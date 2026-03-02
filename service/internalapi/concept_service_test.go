package internalapi

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/writer"
)

type mockConceptRepo struct {
	concept        *writer.Concept
	batchConcepts  []*writer.Concept
	deleteCalledID string
}

func (m *mockConceptRepo) Create(ctx context.Context, concept *writer.Concept) error { return nil }
func (m *mockConceptRepo) GetByID(ctx context.Context, id string) (*writer.Concept, error) {
	if m.concept == nil {
		return nil, errors.New("not found")
	}
	return m.concept, nil
}
func (m *mockConceptRepo) ListByProject(ctx context.Context, projectID string) ([]*writer.Concept, error) {
	return nil, nil
}
func (m *mockConceptRepo) Search(ctx context.Context, projectID, category, keyword string) ([]*writer.Concept, error) {
	return nil, nil
}
func (m *mockConceptRepo) Update(ctx context.Context, concept *writer.Concept) error { return nil }
func (m *mockConceptRepo) Delete(ctx context.Context, id string) error {
	m.deleteCalledID = id
	return nil
}
func (m *mockConceptRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.Concept, error) {
	return m.batchConcepts, nil
}

func TestConceptService_GetConcept_ProjectMismatch(t *testing.T) {
	repo := &mockConceptRepo{
		concept: &writer.Concept{ProjectID: "project-a"},
	}
	svc := NewConceptService(repo)

	_, err := svc.GetConcept(context.Background(), "u1", "project-b", "concept-1")
	if err == nil {
		t.Fatal("expected error for project mismatch")
	}
}

func TestConceptService_Delete_ProjectMismatch(t *testing.T) {
	repo := &mockConceptRepo{
		concept: &writer.Concept{ProjectID: "project-a"},
	}
	svc := NewConceptService(repo)

	err := svc.Delete(context.Background(), "u1", "project-b", "concept-1")
	if err == nil {
		t.Fatal("expected error for project mismatch")
	}
	if repo.deleteCalledID != "" {
		t.Fatalf("delete should not be called, got %s", repo.deleteCalledID)
	}
}

func TestConceptService_BatchGet_FilterByProject(t *testing.T) {
	repo := &mockConceptRepo{
		batchConcepts: []*writer.Concept{
			{ProjectID: "project-a"},
			{ProjectID: "project-b"},
			{ProjectID: "project-a"},
		},
	}
	svc := NewConceptService(repo)

	concepts, err := svc.BatchGet(context.Background(), "u1", "project-a", []string{"1", "2", "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(concepts) != 2 {
		t.Fatalf("expected 2 concepts, got %d", len(concepts))
	}
}
