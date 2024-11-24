package jira

import (
	"context"
	"fmt"
	"time"

	v3 "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/gueldenstone/vogonix/pkg/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	str2dur "github.com/xhit/go-str2duration/v2"
)

const (
	jiraTimeLayout = "2006-01-02T15:04:05.000-0700"
	bucketName     = "worklogs"
)

type TimerData struct {
	closeChan chan bool
	pauseChan chan bool
	ticker    *time.Ticker
}

type JiraInstance struct {
	ctx       context.Context
	atlassian *v3.Client
	timers    map[string]*TimerData
	store     *storage.Storage
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

func NewJiraInstance(url, user, token string, store *storage.Storage) (*JiraInstance, error) {
	atlassian, err := v3.New(nil, url)
	if err != nil {
		return nil, err
	}
	atlassian.Auth.SetBasicAuth(user, token)
	store.AddBucket(bucketName)
	return &JiraInstance{
		atlassian: atlassian,
		timers:    make(map[string]*TimerData),
		store:     store,
	}, nil
}

func (jira JiraInstance) LogDebugf(fmt string, args ...interface{}) {
	runtime.LogDebugf(jira.ctx, fmt, args...)
}

func (jira JiraInstance) LogWarningf(fmt string, args ...interface{}) {
	runtime.LogWarningf(jira.ctx, fmt, args...)
}

func (jira JiraInstance) LogWarning(msg string) {
	runtime.LogWarning(jira.ctx, msg)
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
		issueKey := jira_issue.Key
		issues = append(issues, Issue{
			Key:      issueKey,
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
			jira.LogWarning("Unable to parse duration for a worklog... skipping this entry.")
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

func (jira *JiraInstance) StartTimer(issueId string) {
	// Is there already a timer?
	if timerData, ok := jira.timers[issueId]; ok && timerData != nil {
		jira.timers[issueId].pauseChan <- false
		return
	}

	// Setup new timer
	timerData := &TimerData{
		ticker:    time.NewTicker(1000 * time.Millisecond),
		closeChan: make(chan bool),
		pauseChan: make(chan bool),
	}

	go func() {
		paused := false
		for {
			select {
			case <-timerData.closeChan:
				jira.LogDebugf("stop %s timer because of channel close", issueId)
				return
			case <-jira.ctx.Done():
				jira.LogDebugf("stop %s timer because of channel context done", issueId)
				return
			case pause := <-timerData.pauseChan:
				jira.LogDebugf("received pause with %v for %s timer", pause, issueId)
				paused = pause
			case <-timerData.ticker.C:
				if !paused {
					worklog, _ := jira.GetTimeFromStore(issueId)
					worklog += 1 * time.Second
					jira.UpdateTrackedTime(issueId, worklog)
					jira.LogDebugf("tick for %s at: %s", issueId, worklog)
				}
			}
		}
	}()

	jira.timers[issueId] = timerData
}
func (jira *JiraInstance) PauseTimer(issueId string) {
	jira.LogDebugf("Pausing timer for %s", issueId)
	jira.timers[issueId].pauseChan <- true
}

func (jira *JiraInstance) ResetTimer(issueId string) {
	if jira.timers[issueId] != nil {
		jira.timers[issueId].closeChan <- true
		delete(jira.timers, issueId)
	}
	jira.UpdateTrackedTime(issueId, 0)
}

func (jira JiraInstance) GetCurrentTimerValue(issueId string) int {
	worklog, err := jira.GetTimeFromStore(issueId)
	if err != nil {
		jira.LogWarningf("unable to get worklog time for %s: %w", issueId, err)
		return 0
	}
	return int(worklog.Seconds())
}

func (jira JiraInstance) SubmitWorklog(issueId string) error {
	seconds := jira.GetCurrentTimerValue(issueId)
	jira.ResetTimer(issueId)
	jira.LogDebugf("Received %d seconds for %s", seconds, issueId)
	return nil
}

func (jira JiraInstance) GetTimeFromStore(issueId string) (time.Duration, error) {
	worklogStr, err := jira.store.GetValue(bucketName, issueId)
	if err != nil {
		return 0 * time.Second, fmt.Errorf("no stored time for %s: %w", issueId, err)
	}
	worklogTime, err := time.ParseDuration(worklogStr)
	if err != nil {
		return 0 * time.Second, fmt.Errorf("unable to parse time '%s' for %s: %w", worklogStr, issueId, err)
	}
	return worklogTime, nil
}

func (jira JiraInstance) WriteWorklogToStore(issueId string, t time.Duration) error {
	return jira.store.UpdateValue(bucketName, issueId, t.String())
}

func (jira JiraInstance) UpdateTrackedTime(issueId string, trackedTime time.Duration) {
	err := jira.WriteWorklogToStore(issueId, trackedTime)
	if err != nil {
		jira.LogWarningf("unable to update tracked time for %s: %w", issueId, err)
		return
	}
	runtime.EventsEmit(jira.ctx, "timer_tick_"+issueId, int(trackedTime.Seconds()))
}
