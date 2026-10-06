package writer

import (
	"context"
	"testing"

	serviceErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/interfaces"
	"github.com/stretchr/testify/require"
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

func TestLocationService_GetByIDRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	location, err := service.GetByID(context.Background(), "location-1", "invalid-project")

	requireLocationProjectIDValidation(t, err)
	require.Nil(t, location)
}

func TestLocationService_CreateRelationRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	relation, err := service.CreateRelation(context.Background(), "invalid-project", &interfaces.CreateLocationRelationRequest{FromID: "location-1", ToID: "location-2", Type: "adjacent"})

	requireLocationProjectIDValidation(t, err)
	require.Nil(t, relation)
}

func TestLocationService_DeleteRelationRejectsInvalidProjectID(t *testing.T) {
	service := &LocationService{}

	err := service.DeleteRelation(context.Background(), "relation-1", "invalid-project")

	requireLocationProjectIDValidation(t, err)
}
