package writer

import (
	"context"
	"errors"
	"testing"

	writerModels "Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"github.com/stretchr/testify/require"
)

type outlineLookupStub struct {
	writerRepo.OutlineRepository
	node *writerModels.OutlineNode
}

func (s *outlineLookupStub) FindByID(ctx context.Context, id string) (*writerModels.OutlineNode, error) {
	return s.node, nil
}

type documentUpdateErrorStub struct {
	writerRepo.DocumentRepository
}

func (s *documentUpdateErrorStub) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return errors.New("document update failed")
}

func TestOutlineDocumentSyncService_FindOrCreateGlobalOutlineRejectsInvalidProjectID(t *testing.T) {
	service := &OutlineDocumentSyncService{}

	outlineID, err := service.findOrCreateGlobalOutline(context.Background(), "invalid-project")

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid project ID")
	require.Empty(t, outlineID)
}

func TestOutlineDocumentSyncService_SyncFromDocumentCreationRejectsInvalidProjectID(t *testing.T) {
	service := &OutlineDocumentSyncService{}

	outline, err := service.SyncFromDocumentCreation(context.Background(), "invalid-project", &writerModels.Document{})

	require.Error(t, err)
	require.ErrorContains(t, err, "invalid project ID")
	require.Nil(t, outline)
}

func TestOutlineDocumentSyncService_SyncLevelChangePropagatesDetachError(t *testing.T) {
	documentErr := &documentUpdateErrorStub{}
	service := NewOutlineDocumentSyncService(
		&outlineLookupStub{node: &writerModels.OutlineNode{DocumentID: "document-1"}},
		documentErr,
		nil,
		nil,
		nil,
	)

	err := service.SyncLevelChangeToDocument(context.Background(), "outline-1", 0)

	require.Error(t, err)
	require.ErrorContains(t, err, "document update failed")
}
