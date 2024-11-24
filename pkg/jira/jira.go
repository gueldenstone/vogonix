package jira

import (
	"context"
	"fmt"
	"time"

	v3 "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	str2dur "github.com/xhit/go-str2duration/v2"
)

const (
	jiraTimeLayout = "2006-01-02T15:04:05.000-0700"
)

type JiraInstance struct {
	ctx       context.Context
	atlassian *v3.Client
}

type Issue struct {
	Summary  string
	Assignee string
	Key      string
}

type Worklog struct {
	Duration time.Duration
	Comment  string
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

func (jira *JiraInstance) Startup(ctx context.Context) {
	jira.ctx = ctx
}

func (jira JiraInstance) GetAssignedIssues() ([]Issue, error) {
	jql := "assignee = currentUser() AND status NOT IN ('Done') ORDER BY created DESC"
	fields := []string{"status", "worklog", "assignee", "summary"}
	expand := []string{}

	jira_issues, _, err := jira.atlassian.Issue.Search.Get(context.Background(), jql, fields, expand, 0, 50, "")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve jira issues: %w", err)
	}

	issues := make([]Issue, 0)
	for _, jira_issue := range jira_issues.Issues {
		issues = append(issues, Issue{
			Key:      jira_issue.Key,
			Summary:  jira_issue.Fields.Summary,
			Assignee: jira_issue.Fields.Assignee.DisplayName,
		})
	}
	return issues, nil
}

func (jira JiraInstance) GetBaseUrl() string {
	return jira.atlassian.Site.String()
}

func (jira JiraInstance) GetWorkLogs(issueId string) ([]Worklog, error) {
	jira_worklogs, response, err := jira.atlassian.Issue.Worklog.Issue(jira.ctx, issueId, 0, 1000, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve worklogs from: %s: %w", response.Status, err)
	}

	worklogs := make([]Worklog, 0)
	for _, jira_worklog := range jira_worklogs.Worklogs {
		duration, err := str2dur.ParseDuration(jira_worklog.TimeSpent)
		if err != nil {
			runtime.LogWarning(jira.ctx, "Unable to parse duration for a worklog... skipping this entry.")
			continue
		}
		worklogs = append(worklogs, Worklog{
			Duration: duration,
		})

	}
	return worklogs, nil
}

func (jira JiraInstance) GetTimeSpentOnIssue(issueId string) (string, error) {
	worklogs, err := jira.GetWorkLogs(issueId)
	if err != nil {
		return "", fmt.Errorf("unable to get worklogs for %s: %w", issueId, err)
	}
	summedDuration := 0 * time.Second
	for _, worklog := range worklogs {
		summedDuration += worklog.Duration
	}

	return str2dur.String(summedDuration), nil
}

func (jira JiraInstance) SubmitWorklfowLogs(issueId string, seconds int) error {
	runtime.LogDebugf(jira.ctx, "Received %d seconds for %s", seconds, issueId)
	return nil
}
