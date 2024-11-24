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

type TimerData struct {
	currentDuration time.Duration
	closeChan       chan bool
	pauseChan       chan bool
	ticker          *time.Ticker
}

type JiraInstance struct {
	ctx       context.Context
	atlassian *v3.Client
	timers    map[string]*TimerData
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
		timers:    make(map[string]*TimerData),
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
					timerData.currentDuration += 1 * time.Second
					runtime.EventsEmit(jira.ctx, "timer_tick_"+issueId, timerData.currentDuration.Seconds())
					jira.LogDebugf("tick for %s at: %s", issueId, timerData.currentDuration)
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
		runtime.EventsEmit(jira.ctx, "timer_tick_"+issueId, 0)
		delete(jira.timers, issueId)
	}
}

func (jira JiraInstance) GetCurrentTimerValue(issueId string) int {
	if timerData, ok := jira.timers[issueId]; ok && timerData != nil {
		return int(timerData.currentDuration.Seconds())
	}
	return 0
}

func (jira JiraInstance) SubmitWorklog(issueId string, seconds int) error {
	jira.LogDebugf("Received %d seconds for %s", seconds, issueId)
	return nil
}
