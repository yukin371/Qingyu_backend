package internalapi

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/writer"
)

type mockWriterDraftRepo struct {
	doc            *writer.WriterDraft
	batchDocs      []*writer.WriterDraft
	deleteCalledID string
}

func (m *mockWriterDraftRepo) Create(ctx context.Context, doc *writer.WriterDraft) error { return nil }
func (m *mockWriterDraftRepo) GetByID(ctx context.Context, id string) (*writer.WriterDraft, error) {
	if m.doc == nil {
		return nil, errors.New("not found")
	}
	return m.doc, nil
}
func (m *mockWriterDraftRepo) GetByProjectAndChapter(ctx context.Context, projectID string, chapterNum int) (*writer.WriterDraft, error) {
	return nil, errors.New("not found")
}
func (m *mockWriterDraftRepo) ListByProject(ctx context.Context, projectID string, limit int) ([]*writer.WriterDraft, error) {
	return nil, nil
}
func (m *mockWriterDraftRepo) Update(ctx context.Context, doc *writer.WriterDraft) error { return nil }
func (m *mockWriterDraftRepo) Delete(ctx context.Context, id string) error {
	m.deleteCalledID = id
	return nil
}
func (m *mockWriterDraftRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.WriterDraft, error) {
	return m.batchDocs, nil
}

func TestWriterDraftService_GetDocument_ProjectMismatch(t *testing.T) {
	repo := &mockWriterDraftRepo{
		doc: &writer.WriterDraft{ProjectID: "project-a"},
	}
	svc := NewWriterDraftService(repo)

	_, err := svc.GetDocument(context.Background(), "u1", "project-b", "doc-1")
	if err == nil {
		t.Fatal("expected error for project mismatch")
	}
}

func TestWriterDraftService_DeleteDocument_ProjectMismatch(t *testing.T) {
	repo := &mockWriterDraftRepo{
		doc: &writer.WriterDraft{ProjectID: "project-a"},
	}
	svc := NewWriterDraftService(repo)

	err := svc.DeleteDocument(context.Background(), "u1", "project-b", "doc-1")
	if err == nil {
		t.Fatal("expected error for project mismatch")
	}
	if repo.deleteCalledID != "" {
		t.Fatalf("delete should not be called, got %s", repo.deleteCalledID)
	}
}

func TestWriterDraftService_BatchGetDocuments_FilterByProject(t *testing.T) {
	repo := &mockWriterDraftRepo{
		batchDocs: []*writer.WriterDraft{
			{ProjectID: "project-a"},
			{ProjectID: "project-b"},
			{ProjectID: "project-a"},
		},
	}
	svc := NewWriterDraftService(repo)

	docs, err := svc.BatchGetDocuments(context.Background(), "u1", "project-a", []string{"1", "2", "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}
}
