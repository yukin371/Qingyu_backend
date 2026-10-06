package writer

import (
	"context"
	"testing"

	serviceErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/interfaces"
	"github.com/stretchr/testify/require"
)

func requireCharacterProjectIDValidation(t *testing.T, err error) {
	t.Helper()
	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorValidation, serviceErr.Type)
	require.Contains(t, serviceErr.Error(), "invalid project ID")
}

func TestCharacterService_CreateRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	character, err := service.Create(context.Background(), "invalid-project", "user-1", &interfaces.CreateCharacterRequest{Name: "测试角色"})

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, character)
}

func TestCharacterService_GetByIDRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	character, err := service.GetByID(context.Background(), "507f1f77bcf86cd799439011", "invalid-project")

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, character)
}

func TestCharacterService_CreateRelationRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	relation, err := service.CreateRelation(context.Background(), "invalid-project", &interfaces.CreateRelationRequest{FromID: "character-1", ToID: "character-2", Type: "朋友"})

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, relation)
}

func TestCharacterService_DeleteRelationRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	err := service.DeleteRelation(context.Background(), "relation-1", "invalid-project")

	requireCharacterProjectIDValidation(t, err)
}

func TestCharacterService_CreateRelationTimelineEventRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	event, err := service.CreateRelationTimelineEvent(context.Background(), "invalid-project", &interfaces.CreateRelationTimelineEventRequest{RelationID: "relation-1"})

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, event)
}

func TestCharacterService_GetRelationTimelineRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	events, err := service.GetRelationTimeline(context.Background(), "relation-1", "invalid-project")

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, events)
}

func TestCharacterService_UpdateRelationTimelineEventRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	event, err := service.UpdateRelationTimelineEvent(context.Background(), "event-1", "invalid-project", &interfaces.UpdateRelationTimelineEventRequest{})

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, event)
}

func TestCharacterService_DeleteRelationTimelineEventRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	err := service.DeleteRelationTimelineEvent(context.Background(), "event-1", "invalid-project")

	requireCharacterProjectIDValidation(t, err)
}
