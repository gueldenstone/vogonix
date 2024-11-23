package jira

import (
	"context"
	"fmt"

	v3 "github.com/ctreminiom/go-atlassian/jira/v3"
)

type JiraInstance struct {
	atlassian *v3.Client
}

type Issue struct {
	summary  string
	assignee string
	key      string
}

func NewJiraInstance(url, user, token string) (*JiraInstance, error) {
	atlassian, err := v3.New(nil, url)
	if err != nil {
		return nil, err
	}
	atlassian.Auth.SetBasicAuth(user, token)
	return &JiraInstance{
		atlassian: atlassian,
	}, nil
}

func (jira JiraInstance) GetAssignedIssues() ([]Issue, error) {
	jql := "assignee = currentUser() AND status NOT IN ('Done') ORDER BY created DESC"
	fields := []string{"status", "worklog"}
	expand := []string{}

	jira_issues, _, err := jira.atlassian.Issue.Search.Get(context.Background(), jql, fields, expand, 0, 50, "")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve jira issues: %w", err)
	}

	issues := make([]Issue, jira_issues.Total)
	for _, jira_issue := range jira_issues.Issues {
		issues = append(issues, Issue{
			key:      jira_issue.Key,
			summary:  jira_issue.Fields.Summary,
			assignee: jira_issue.Fields.Assignee.DisplayName,
		})
	}
	return issues, nil
}
