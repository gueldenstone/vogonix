package jira

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	str2dur "github.com/xhit/go-str2duration/v2"
)

func (jira JiraInstance) GetBaseUrl() string {
	return jira.atlassian.Site.String()
}

func (jira JiraInstance) getRemoteIssues() ([]Issue, error) {
	jql := "assignee = currentUser() AND status NOT IN ('Done') ORDER BY created DESC"
	fields := []string{"status", "worklog", "assignee", "summary"}
	expand := []string{}

	ctx, cancel := context.WithTimeout(jira.ctx, 1*time.Second)
	defer cancel()
	jira_issues, _, err := jira.atlassian.Issue.Search.Get(ctx, jql, fields, expand, 0, 50, "")
	if err != nil {
		return nil, err
	}

	remoteIssues := make([]Issue, 0)
	for _, jira_issue := range jira_issues.Issues {
		issueKey := jira_issue.Key
		workLogs, err := jira.getWorkLogs(issueKey)
		if err != nil {
			jira.LogWarningf("error getting worklogs for %s: %s", issueKey, err.Error())
		}
		issue := Issue{
			Key:      issueKey,
			Summary:  jira_issue.Fields.Summary,
			Assignee: jira_issue.Fields.Assignee.DisplayName,
			WorkLogs: workLogs,
		}
		remoteIssues = append(remoteIssues, issue)
	}

	return remoteIssues, nil
}

func (jira JiraInstance) getWorkLogs(issueId string) ([]Worklog, error) {
	var worklogs WorkLogs = make([]Worklog, 0)
	jira_worklogs, response, err := jira.atlassian.Issue.Worklog.Issue(jira.ctx, issueId, 0, 1000, 0, []string{"all"})
	if err != nil {
		return worklogs, fmt.Errorf("unable to retrieve worklogs from: %s: %w", response.Status, err)
	}

	for _, jira_worklog := range jira_worklogs.Worklogs {
		duration, err := str2dur.ParseDuration(jira_worklog.TimeSpent)
		if err != nil {
			jira.LogWarning("Unable to parse duration for a worklog... skipping this entry.")
			continue
		}
		updated, err := time.Parse(jiraTimeLayout, jira_worklog.Updated)
		if err != nil {
			jira.LogWarning("Unable to parse updated for a worklog... skipping this entry.")
			continue
		}
		worklogs = append(worklogs, Worklog{
			Duration: duration,
			Author:   jira_worklog.Author.DisplayName,
			Updated:  updated.Format(time.DateTime),
		})

	}
	sort.Sort(ByUpdated{worklogs})
	return worklogs, nil
}

func (jira JiraInstance) AddWorklog(issueKey string, duration time.Duration) error {
	options := &models.WorklogOptionsScheme{
		Notify:         false,
		AdjustEstimate: "leave",
	}
	payload := &models.WorklogADFPayloadScheme{
		Visibility:       nil,
		Started:          time.Now().Format(jiraTimeLayout),
		TimeSpentSeconds: int(duration.Seconds()),
	}

	p, resp, err := jira.atlassian.Issue.Worklog.Add(jira.ctx, issueKey, payload, options)
	if err != nil {
		jira.LogWarningf("error submitting: %s resp: %+v, payload: %+v", err.Error(), resp.Bytes.String(), p)
		return err
	}
	return nil
}
