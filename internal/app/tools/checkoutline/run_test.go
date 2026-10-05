package checkoutline

import "testing"

func TestMongoConfiguration(t *testing.T) {
	t.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI", "mongodb://primary")
	t.Setenv("MONGODB_URI", "mongodb://legacy")
	t.Setenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE", "primary_db")
	t.Setenv("MONGODB_DATABASE", "legacy_db")
	if got := mongoURI(); got != "mongodb://primary" {
		t.Fatalf("mongoURI() = %q", got)
	}
	if got := databaseName(); got != "primary_db" {
		t.Fatalf("databaseName() = %q", got)
	}
}
