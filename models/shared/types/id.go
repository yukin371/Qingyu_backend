package types

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrEmptyID         = errors.New("ID cannot be empty")
)

func ParseObjectID(s string) (primitive.ObjectID, error) {
	if s == "" {
		return primitive.NilObjectID, ErrEmptyID
	}
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, s)
	}
	return oid, nil
}

func MustParseObjectID(s string) primitive.ObjectID {
	oid, err := ParseObjectID(s)
	if err != nil {
		panic(err)
	}
	return oid
}

func ParseOptionalObjectID(s string) (*primitive.ObjectID, error) {
	if s == "" {
		return nil, nil
	}
	oid, err := ParseObjectID(s)
	if err != nil {
		return nil, err
	}
	return &oid, nil
}

func ToHex(id primitive.ObjectID) string {
	return id.Hex()
}

func IsValidObjectID(s string) bool {
	_, err := primitive.ObjectIDFromHex(s)
	return err == nil
}

func ParseObjectIDSlice(ss []string) ([]primitive.ObjectID, map[int]error) {
	oids := make([]primitive.ObjectID, 0, len(ss))
	errMap := make(map[int]error)
	for i, s := range ss {
		oid, err := ParseObjectID(s)
		if err != nil {
			errMap[i] = err
			continue
		}
		oids = append(oids, oid)
	}
	return oids, errMap
}

func ParseOptionalObjectIDSlice(ss []string) ([]primitive.ObjectID, error) {
	if len(ss) == 0 {
		return nil, nil
	}
	oids := make([]primitive.ObjectID, 0, len(ss))
	for _, s := range ss {
		if s == "" {
			continue
		}
		oid, err := ParseObjectID(s)
		if err != nil {
			return nil, err
		}
		oids = append(oids, oid)
	}
	return oids, nil
}

func ParseObjectIDSliceStrict(ss []string) ([]primitive.ObjectID, error) {
	if len(ss) == 0 {
		return nil, nil
	}
	oids := make([]primitive.ObjectID, 0, len(ss))
	for i, s := range ss {
		oid, err := ParseObjectID(s)
		if err != nil {
			return nil, fmt.Errorf("ids[%d]: %w", i, err)
		}
		oids = append(oids, oid)
	}
	return oids, nil
}

func ToHexSlice(ids []primitive.ObjectID) []string {
	if ids == nil {
		return nil
	}
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = ToHex(id)
	}
	return result
}

func GenerateNewObjectID() string {
	return primitive.NewObjectID().Hex()
}

func IsNilObjectID(id primitive.ObjectID) bool {
	return id.IsZero()
}

func IsIDError(err error) bool {
	return errors.Is(err, ErrEmptyID) || errors.Is(err, ErrInvalidIDFormat)
}
