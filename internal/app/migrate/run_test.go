package migrateapp

import "testing"

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		migration string
		env       string
		wantErr   bool
	}{
		{name: "valid up", command: "up", migration: "002_create_users_indexes", env: "dev"},
		{name: "missing command", migration: "002_create_users_indexes", env: "dev", wantErr: true},
		{name: "invalid command", command: "status", migration: "002_create_users_indexes", env: "dev", wantErr: true},
		{name: "missing migration", command: "up", env: "dev", wantErr: true},
		{name: "invalid environment", command: "up", migration: "002_create_users_indexes", env: "local", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(tt.command, tt.migration, tt.env)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseName(t *testing.T) {
	tests := map[string]string{
		"dev":        "qingyu_dev",
		"staging":    "qingyu_staging",
		"production": "qingyu",
	}
	for env, want := range tests {
		if got := databaseName(env); got != want {
			t.Fatalf("databaseName(%q) = %q, want %q", env, got, want)
		}
	}
}
