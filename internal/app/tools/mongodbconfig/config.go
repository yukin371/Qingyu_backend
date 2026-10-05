// Package mongodbconfig owns MongoDB environment resolution for diagnostic tools.
package mongodbconfig

import (
	"os"
	"strings"
)

// URI returns the configured MongoDB URI or the local diagnostic default.
func URI() string {
	if uri := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI")); uri != "" {
		return uri
	}
	if uri := strings.TrimSpace(os.Getenv("MONGODB_URI")); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017/?directConnection=true"
}

// Database returns the configured MongoDB database or the local default.
func Database() string {
	if name := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE")); name != "" {
		return name
	}
	if name := strings.TrimSpace(os.Getenv("MONGODB_DATABASE")); name != "" {
		return name
	}
	return "qingyu"
}
