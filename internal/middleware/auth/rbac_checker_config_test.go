package auth

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACChecker_LoadFromFile_AuthorProjectAndDocumentPermissions(t *testing.T) {
	configPath := filepath.Join("..", "..", "..", "configs", "permissions.yaml")

	checker, err := NewRBACChecker(&CheckerConfig{
		Strategy:   "rbac",
		ConfigPath: configPath,
	})
	require.NoError(t, err)

	rbac, ok := checker.(*RBACChecker)
	require.True(t, ok)

	assert.Contains(t, rbac.GetRolePermissions("author"), "project:read")
	assert.Contains(t, rbac.GetRolePermissions("author"), "project:create")
	assert.Contains(t, rbac.GetRolePermissions("author"), "project:update")
	assert.Contains(t, rbac.GetRolePermissions("author"), "project:delete")
	assert.Contains(t, rbac.GetRolePermissions("author"), "document:update")
	assert.Contains(t, rbac.GetRolePermissions("author"), "document:delete")

	allowed, err := checker.Check(context.Background(), "user_003", Permission{
		Resource: "project",
		Action:   "create",
	})
	require.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = checker.Check(context.Background(), "user_003", Permission{
		Resource: "document",
		Action:   "update",
	})
	require.NoError(t, err)
	assert.True(t, allowed)
}
