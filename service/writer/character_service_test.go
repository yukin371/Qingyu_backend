package writer

import (
	"context"
	"testing"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/writer/base"
	serviceErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/interfaces"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type projectReaderStub struct {
	project *writer.Project
	err     error
}

func (s projectReaderStub) GetByID(context.Context, string) (*writer.Project, error) {
	return s.project, s.err
}

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

func TestCharacterService_CreateRejectsUserWithoutProjectEditPermission(t *testing.T) {
	projectID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	service := &CharacterService{
		projectRepo: projectReaderStub{
			project: &writer.Project{OwnedEntity: base.OwnedEntity{AuthorID: authorID}},
		},
	}

	character, err := service.Create(context.Background(), projectID.Hex(), primitive.NewObjectID().Hex(), &interfaces.CreateCharacterRequest{Name: "测试角色"})

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
	require.Nil(t, character)
}

func TestCharacterService_UpdateRejectsUserWithoutProjectEditPermission(t *testing.T) {
	projectID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	service := &CharacterService{
		projectRepo: projectReaderStub{
			project: &writer.Project{OwnedEntity: base.OwnedEntity{AuthorID: authorID}},
		},
	}

	character, err := service.Update(context.Background(), primitive.NewObjectID().Hex(), projectID.Hex(), primitive.NewObjectID().Hex(), &interfaces.UpdateCharacterRequest{})

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
	require.Nil(t, character)
}

func TestCharacterService_DeleteRejectsUserWithoutProjectEditPermission(t *testing.T) {
	projectID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	service := &CharacterService{
		projectRepo: projectReaderStub{
			project: &writer.Project{OwnedEntity: base.OwnedEntity{AuthorID: authorID}},
		},
	}

	err := service.Delete(context.Background(), primitive.NewObjectID().Hex(), projectID.Hex(), primitive.NewObjectID().Hex())

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
}

func TestCharacterService_GetByIDRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	character, err := service.GetByID(context.Background(), "507f1f77bcf86cd799439011", "invalid-project")

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, character)
}

func TestCharacterService_CreateRelationRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	relation, err := service.CreateRelation(context.Background(), "invalid-project", "user-1", &interfaces.CreateRelationRequest{FromID: "character-1", ToID: "character-2", Type: "朋友"})

	requireCharacterProjectIDValidation(t, err)
	require.Nil(t, relation)
}

func TestCharacterService_DeleteRelationRejectsInvalidProjectID(t *testing.T) {
	service := &CharacterService{}

	err := service.DeleteRelation(context.Background(), "relation-1", "invalid-project", "user-1")

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
