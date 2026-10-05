package testhierarchy

import "testing"

func TestParseApplyDefaultsToPreview(t *testing.T) {
	apply, err := parseApply(nil)
	if err != nil || apply {
		t.Fatalf("parseApply() = %v, %v", apply, err)
	}
}

func TestResponseCodeRejectsNonNumericCode(t *testing.T) {
	if err := responseCode(map[string]interface{}{"code": "bad"}); err == nil {
		t.Fatal("responseCode() accepted non-numeric code")
	}
}
