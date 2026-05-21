package config

import "testing"

func TestAIQuotaConfigGetConsistencyThreshold_Defaults(t *testing.T) {
	cfg := &AIQuotaConfig{}

	userThreshold := cfg.GetConsistencyThreshold("user")
	if userThreshold.WarningTokens != 200 || userThreshold.CriticalTokens != 1000 {
		t.Fatalf("unexpected default token thresholds: %+v", userThreshold)
	}
	if userThreshold.WarningRatio != 0.1 || userThreshold.CriticalRatio != 0.2 {
		t.Fatalf("unexpected default ratio thresholds: %+v", userThreshold)
	}
}

func TestAIQuotaConfigGetConsistencyThreshold_MergesConfiguredScope(t *testing.T) {
	cfg := &AIQuotaConfig{
		ConsistencyThresholds: &QuotaConsistencyThresholdsConfig{
			Global: &QuotaConsistencyThresholdConfig{
				WarningTokens:  600,
				CriticalTokens: 2400,
				WarningRatio:   0.3,
				CriticalRatio:  0.55,
			},
		},
	}

	globalThreshold := cfg.GetConsistencyThreshold("global")
	if globalThreshold.WarningTokens != 600 || globalThreshold.CriticalTokens != 2400 {
		t.Fatalf("unexpected configured token thresholds: %+v", globalThreshold)
	}
	if globalThreshold.WarningRatio != 0.3 || globalThreshold.CriticalRatio != 0.55 {
		t.Fatalf("unexpected configured ratio thresholds: %+v", globalThreshold)
	}

	workflowThreshold := cfg.GetConsistencyThreshold("workflow")
	if workflowThreshold.WarningTokens != 200 || workflowThreshold.CriticalTokens != 1000 {
		t.Fatalf("unexpected fallback token thresholds: %+v", workflowThreshold)
	}
	if workflowThreshold.WarningRatio != 0.1 || workflowThreshold.CriticalRatio != 0.2 {
		t.Fatalf("unexpected fallback ratio thresholds: %+v", workflowThreshold)
	}
}
