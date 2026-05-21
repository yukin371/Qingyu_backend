package transaction

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMongoTransactionsDisabled(t *testing.T) {
	t.Setenv(disableMongoTransactionsEnv, "true")
	if !MongoTransactionsDisabled() {
		t.Fatal("expected mongo transactions to be disabled")
	}
}

func TestMongoTransactionsDisabled_DefaultFalse(t *testing.T) {
	t.Setenv(disableMongoTransactionsEnv, "")
	if MongoTransactionsDisabled() {
		t.Fatal("expected mongo transactions to be enabled by default")
	}
}

func TestMongoTransactionsDisabled_FromDotEnv(t *testing.T) {
	t.Setenv(disableMongoTransactionsEnv, "")

	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte(disableMongoTransactionsEnv+"=true\n"), 0o644); err != nil {
		t.Fatalf("failed to write dotenv file: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	if !MongoTransactionsDisabled() {
		t.Fatal("expected mongo transactions to be disabled from .env")
	}
}
