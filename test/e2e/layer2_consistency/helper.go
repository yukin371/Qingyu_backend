//go:build e2e
// +build e2e

package layer2_consistency

import "Qingyu_backend/test/e2e/data"

// filterIssuesByType 按类型过滤问题
func filterIssuesByType(issues []data.ConsistencyIssue, issueType string) []data.ConsistencyIssue {
	var filtered []data.ConsistencyIssue
	for _, issue := range issues {
		if issue.Type == issueType {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

