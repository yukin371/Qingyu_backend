package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAssertBehavior(t *testing.T) {
	// 测试 assert.NoError
	t.Run("NoError_with_nil", func(t *testing.T) {
		err := error(nil)
		assert.NoError(t, err)
	})
	
	// 测试 assert.Error
	t.Run("Error_with_nil", func(t *testing.T) {
		err := error(nil)
		assert.Error(t, err)
	})
}
