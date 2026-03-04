package writer

import "errors"

var (
	// ErrOptimisticLockConflict indicates an optimistic lock/version conflict.
	// Repository should return this neutral error and let service map it to business semantics.
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict")
)
