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

func requireLocationProjectIDValidation(t *testing.T, err error) {
	t.Helper()
	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorValidation, serviceErr.Type)
	require.Contains(t, serviceErr.Error(), "invalid project ID")
}

func TestLocationService_CreateRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	location, err := service.Create(context.Background(), "invalid-project", "user-1", &interfaces.CreateLocationRequest{Name: "测试地点"})

	requireLocationProjectIDValidation(t, err)
	require.Nil(t, location)
}

func TestLocationService_CreateRejectsUserWithoutProjectEditPermission(t *testing.T) {
	service := &LocationService{
		projectRepo: projectReaderStub{project: &writer.Project{}},
	}

	location, err := service.Create(context.Background(), primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), &interfaces.CreateLocationRequest{Name: "测试地点"})

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
	require.Nil(t, location)
}

func TestLocationService_UpdateRejectsUserWithoutProjectEditPermission(t *testing.T) {
	projectID := primitive.NewObjectID()
	service := &LocationService{
		projectRepo: projectReaderStub{project: &writer.Project{OwnedEntity: base.OwnedEntity{AuthorID: primitive.NewObjectID()}}},
	}

	location, err := service.Update(context.Background(), primitive.NewObjectID().Hex(), projectID.Hex(), primitive.NewObjectID().Hex(), &interfaces.UpdateLocationRequest{})

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
	require.Nil(t, location)
}

func TestLocationService_DeleteRejectsUserWithoutProjectEditPermission(t *testing.T) {
	projectID := primitive.NewObjectID()
	service := &LocationService{
		projectRepo: projectReaderStub{project: &writer.Project{OwnedEntity: base.OwnedEntity{AuthorID: primitive.NewObjectID()}}},
	}

	err := service.Delete(context.Background(), primitive.NewObjectID().Hex(), projectID.Hex(), primitive.NewObjectID().Hex())

	var serviceErr *serviceErrors.ServiceError
	require.ErrorAs(t, err, &serviceErr)
	require.Equal(t, serviceErrors.ServiceErrorForbidden, serviceErr.Type)
}

func TestLocationService_GetByIDRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	location, err := service.GetByID(context.Background(), "location-1", "invalid-project")

	requireLocationProjectIDValidation(t, err)
	require.Nil(t, location)
}

func TestLocationService_CreateRelationRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	relation, err := service.CreateRelation(context.Background(), "invalid-project", "user-1", &interfaces.CreateLocationRelationRequest{FromID: "location-1", ToID: "location-2", Type: "adjacent"})

	requireLocationProjectIDValidation(t, err)
	require.Nil(t, relation)
}

func TestLocationService_DeleteRelationRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	err := service.DeleteRelation(context.Background(), "relation-1", "invalid-project", "user-1")

	requireLocationProjectIDValidation(t, err)
}
