package mongodbconfig

import "testing"

func TestURIAndDatabasePreferPrimaryEnvironment(t *testing.T) {
	t.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI", "mongodb://primary")
	t.Setenv("MONGODB_URI", "mongodb://legacy")
	t.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE", "primary_db")
	t.Setenv("MONGODB_DATABASE", "legacy_db")
	if got := URI(); got != "mongodb://primary" {
		t.Fatalf("URI() = %q", got)
	}
	if got := Database(); got != "primary_db" {
		t.Fatalf("Database() = %q", got)
	}
}

func TestURIAndDatabaseUseDefaults(t *testing.T) {
	for _, key := range []string{
		"QINGYU_DATABASE_PRIMARY_MONGODB_URI",
		"MONGODB_URI",
		"QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE",
		"MONGODB_DATABASE",
	} {
		t.Setenv(key, "")
	}
	if got := URI(); got != "mongodb://localhost:27017/?directConnection=true" {
		t.Fatalf("URI() = %q", got)
	}
	if got := Database(); got != "qingyu" {
		t.Fatalf("Database() = %q", got)
	}
}
