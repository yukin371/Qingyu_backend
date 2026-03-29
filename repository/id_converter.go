package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared/types"
)

var (
	ErrEmptyID         = types.ErrEmptyID
	ErrInvalidIDFormat = types.ErrInvalidIDFormat
)

func StringToObjectId(id string) (primitive.ObjectID, error) {
	return types.ParseObjectID(id)
}

func ObjectIdToString(id primitive.ObjectID) string {
	return types.ToHex(id)
}
func StringSliceToObjectIDSlice(ids []string) ([]primitive.ObjectID, error) {
	return types.ParseOptionalObjectIDSlice(ids)
}
func ObjectIDSliceToStringSlice(ids []primitive.ObjectID) []string {
	return types.ToHexSlice(ids)
}
func StringPtrToObjectIdPtr(idPtr *string) (*primitive.ObjectID, error) {
	if idPtr == nil {
		return nil, nil
	}
	return types.ParseOptionalObjectID(*idPtr)
}
func ObjectIdPtrToStringPtr(idPtr *primitive.ObjectID) *string {
	if idPtr == nil {
		return nil
	}
	hex := types.ToHex(*idPtr)
	return &hex
}
func ConvertCategoryIdsToString(categoryIds []primitive.ObjectID) []string {
	return types.ToHexSlice(categoryIds)
}
func ConvertCategoryIdsToObjectID(categoryIds []string) ([]primitive.ObjectID, error) {
	return types.ParseOptionalObjectIDSlice(categoryIds)
}
func ParseID(id string) (primitive.ObjectID, error) {
	return types.ParseObjectID(id)
}
func ParseOptionalID(id string) (*primitive.ObjectID, error) {
	return types.ParseOptionalObjectID(id)
}
func ParseIDs(ids []string) ([]primitive.ObjectID, error) {
	return types.ParseObjectIDSliceStrict(ids)
}
func ParseOptionalIDs(ids []string) ([]primitive.ObjectID, error) {
	return types.ParseOptionalObjectIDSlice(ids)
}
func IsIDError(err error) bool {
	return types.IsIDError(err)
}
