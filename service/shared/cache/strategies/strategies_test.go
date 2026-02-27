package strategies

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTTLStrategy_DefaultTTL(t *testing.T) {
	s := NewTTLStrategy(0)
	assert.Equal(t, 5*time.Minute, s.GetTTL("k"))
	assert.True(t, s.ShouldCache("k", "v"))
	assert.NoError(t, s.OnMiss("k"))
}

func TestStrategyManager_MatchByPrefix(t *testing.T) {
	manager := NewStrategyManager(2 * time.Minute)
	manager.RegisterStrategy("user:", NewTTLStrategy(30*time.Second))

	assert.Equal(t, 30*time.Second, manager.GetStrategy("user:123").GetTTL("user:123"))
	assert.Equal(t, 2*time.Minute, manager.GetStrategy("book:1").GetTTL("book:1"))
}
