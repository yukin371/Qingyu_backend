//go:build e2e
// +build e2e

package layer3_boundary

import "testing"

func TestConcurrentSocialInteraction(t *testing.T) {
	RunConcurrentSocialInteraction(t)
}

func TestBoundaryDataSizes(t *testing.T) {
	RunBoundaryDataSizes(t)
}
