package jira

import (
	"fmt"
	"time"
)

func (jira JiraInstance) getTimeFromStore(issueId string) (time.Duration, error) {
	worklogStr, err := jira.store.GetStringValue(timeBucketName, issueId)
	if err != nil {
		return 0 * time.Second, fmt.Errorf("no stored time for %s: %w", issueId, err)
	}
	worklogTime, err := time.ParseDuration(worklogStr)
	if err != nil {
		return 0 * time.Second, fmt.Errorf("unable to parse time '%s' for %s: %w", worklogStr, issueId, err)
	}
	return worklogTime, nil
}

func (jira JiraInstance) writeWorklogToStore(issueId string, t time.Duration) error {
	return jira.store.UpdateStringValue(timeBucketName, issueId, t.String())
}

func (jira JiraInstance) writeIssueDataToStore(issue Issue) error {
	return jira.store.UpdateStructuredValue(issueBucketName, issue.Key, issue)
}

func (jira JiraInstance) readIssueDataFromStore(issueId string) (Issue, error) {
	issue := Issue{}
	err := jira.store.GetStructuredValue(issueBucketName, issueId, &issue)
	return issue, err
}

func (jira JiraInstance) getAllStoredIssues() ([]Issue, error) {
	keys, err := jira.store.GetAllKeys(issueBucketName)
	if err != nil {
		return nil, err
	}
	issues := make([]Issue, 0)
	for _, key := range keys {
		issue := Issue{}
		err := jira.store.GetStructuredValue(issueBucketName, key, &issue)
		if err != nil {
			jira.LogWarningf("no data for %s: %s", key, err.Error())
			continue
		}
		issues = append(issues, issue)
	}
	return issues, nil
}
