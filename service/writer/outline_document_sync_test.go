package writer

import (
	"context"
	"testing"

	writerModels "Qingyu_backend/models/writer"
	"github.com/stretchr/testify/require"
)

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
