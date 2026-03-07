package transaction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMongoRunnerNilClient(t *testing.T) {
	runner := NewMongoRunner(nil)

	err := runner.Run(context.Background(), func(context.Context) error {
		return nil
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "mongo client is nil")
}
