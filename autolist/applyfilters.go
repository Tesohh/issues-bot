package autolist

import (
	"issues/db"
	"slices"
	"strings"
)

func ApplyFilters(issues []db.Issue, filterMe bool, userId string, showDone bool, priorityFilter string, kindFilter string) []db.Issue {
	if filterMe {
		issues = slices.DeleteFunc(issues, func(issue db.Issue) bool {
			ids := strings.Split(issue.AssigneeIDs, ",")
			return !slices.Contains(ids, userId)
		})
	}

	if !showDone {
		issues = slices.DeleteFunc(issues, func(issue db.Issue) bool {
			return issue.IssueStatus == db.IssueStatusCanceled || issue.IssueStatus == db.IssueStatusDone
		})
	}

	if priorityFilter != "" {
		issues = slices.DeleteFunc(issues, func(issue db.Issue) bool {
			return issue.PriorityRoleID != priorityFilter
		})
	}

	if kindFilter != "" {
		issues = slices.DeleteFunc(issues, func(issue db.Issue) bool {
			return issue.KindRoleID != kindFilter
		})
	}

	return issues
}
