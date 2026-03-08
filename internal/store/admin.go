package store

import (
	"context"
	"fmt"

	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type IssuePriorityCount struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type AdminStats struct {
	TotalUsers      int64                `json:"total_users"`
	TotalIssues     int64                `json:"total_issues"`
	IssuesByPriority []IssuePriorityCount `json:"issues_by_priority"`
}

type AdminStore struct {
	queries *dbsqlc.Queries
}

func (a *AdminStore) GetStats(ctx context.Context) (*AdminStats, error) {
	totalUsers, err := a.queries.CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("counting users: %w", err)
	}

	rows, err := a.queries.CountIssuesByPriority(ctx)
	if err != nil {
		return nil, fmt.Errorf("counting issues by priority: %w", err)
	}

	var totalIssues int64
	issuesByPriority := make([]IssuePriorityCount, len(rows))
	for i, row := range rows {
		issuesByPriority[i] = IssuePriorityCount{
			Priority: row.Priority,
			Count:    row.Count,
		}
		totalIssues += row.Count
	}

	return &AdminStats{
		TotalUsers:       totalUsers,
		TotalIssues:      totalIssues,
		IssuesByPriority: issuesByPriority,
	}, nil
}
