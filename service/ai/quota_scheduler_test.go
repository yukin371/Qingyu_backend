package ai

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuotaScheduler_StartStopWithoutOptionalDependencies(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	scheduler := NewQuotaScheduler(nil, nil, nil, nil, logger)

	err := scheduler.Start()
	assert.NoError(t, err)

	scheduler.Stop()

	assert.Contains(t, buffer.String(), "quota scheduler started")
	assert.Contains(t, buffer.String(), "quota scheduler stopped")
}

func TestQuotaScheduler_CheckCrossServiceConsistencySkipsWithoutPhase3Client(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)
	scheduler := NewQuotaScheduler(nil, nil, nil, nil, logger)

	err := scheduler.checkCrossServiceConsistency(context.Background())

	assert.NoError(t, err)
	assert.Contains(t, buffer.String(), "skip quota consistency check: Phase3Client not configured")
}
