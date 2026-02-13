package storage

import "testing"

func TestBackendFactoryValidateConfigLocal(t *testing.T) {
	factory := NewBackendFactory()

	if err := factory.ValidateConfig(&LocalBackendConfig{BasePath: ""}); err == nil {
		t.Fatalf("expected validation error for empty local base_path")
	}

	if err := factory.ValidateConfig(&LocalBackendConfig{BasePath: "./uploads"}); err != nil {
		t.Fatalf("expected local config to be valid, got: %v", err)
	}
}

func TestBackendFactoryValidateConfigMinIO(t *testing.T) {
	factory := NewBackendFactory()

	tests := []struct {
		name string
		cfg  *MinIOBackendConfig
	}{
		{
			name: "missing endpoint",
			cfg: &MinIOBackendConfig{
				AccessKey:  "ak",
				SecretKey:  "sk",
				BucketName: "bucket",
			},
		},
		{
			name: "missing access key",
			cfg: &MinIOBackendConfig{
				Endpoint:   "127.0.0.1:9000",
				SecretKey:  "sk",
				BucketName: "bucket",
			},
		},
		{
			name: "missing secret key",
			cfg: &MinIOBackendConfig{
				Endpoint:   "127.0.0.1:9000",
				AccessKey:  "ak",
				BucketName: "bucket",
			},
		},
		{
			name: "missing bucket name",
			cfg: &MinIOBackendConfig{
				Endpoint:  "127.0.0.1:9000",
				AccessKey: "ak",
				SecretKey: "sk",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := factory.ValidateConfig(tt.cfg); err == nil {
				t.Fatalf("expected minio config validation to fail")
			}
		})
	}

	validCfg := &MinIOBackendConfig{
		Endpoint:   "127.0.0.1:9000",
		AccessKey:  "ak",
		SecretKey:  "sk",
		BucketName: "bucket",
	}
	if err := factory.ValidateConfig(validCfg); err != nil {
		t.Fatalf("expected minio config to be valid, got: %v", err)
	}
}

func TestBackendFactorySupportedBackends(t *testing.T) {
	factory := NewBackendFactory()
	supported := factory.GetSupportedBackends()

	hasLocal := false
	hasMinIO := false
	for _, backend := range supported {
		if backend == BackendTypeLocal {
			hasLocal = true
		}
		if backend == BackendTypeMinIO {
			hasMinIO = true
		}
	}

	if !hasLocal || !hasMinIO {
		t.Fatalf("supported backends should include local and minio, got: %v", supported)
	}
}

