package jira

import (
	"context"
	"fmt"
	"sort"
	"time"

	v3 "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/gueldenstone/vogonix/pkg/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	str2dur "github.com/xhit/go-str2duration/v2"
)

const (
	jiraTimeLayout  = "2006-01-02T15:04:05.000-0700"
	timeBucketName  = "worklogs"
	issueBucketName = "issues"
)

type TimerData struct {
	closeChan chan bool
	pauseChan chan bool
	ticker    *time.Ticker
}

type JiraInstance struct {
	ctx         context.Context
	atlassian   *v3.Client
	timers      map[string]*TimerData
	store       *storage.Storage
	OfflineMode bool
}

type Issue struct {
	Summary   string    `json:"summary,omitempty"`
	Assignee  string    `json:"assignee,omitempty"`
	Key       string    `json:"key,omitempty"`
	WorkLogs  []Worklog `json:"worklogs,omitempty"`
	TimeSpent int       `json:"time_spent,omitempty"`
}

type Worklog struct {
	Duration  time.Duration `json:"duration,omitempty"`
	Comment   string        `json:"comment,omitempty"`
	Submitted bool          `json:"submitted,omitempty"`
	Author    string        `json:"author,omitempty"`
	Updated   string        `json:"updated,omitempty"`
}

func NewJiraInstance(url, user, token string, store *storage.Storage) (*JiraInstance, error) {
	atlassian, err := v3.New(nil, url)
	if err != nil {
		return nil, err
	}
	atlassian.Auth.SetBasicAuth(user, token)
	store.AddBucket(timeBucketName)
	store.AddBucket(issueBucketName)
	return &JiraInstance{
		atlassian:   atlassian,
		timers:      make(map[string]*TimerData),
		store:       store,
		OfflineMode: false,
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

func (jira JiraInstance) LogDebug(msg string) {
	runtime.LogDebug(jira.ctx, msg)
}

func (jira *JiraInstance) Startup(ctx context.Context) {
	jira.ctx = ctx
}

func (jira JiraInstance) GetRemoteIssues() ([]Issue, error) {
	jira.LogDebug("GetRemoteIssues")
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
		workLogs, err := jira.GetWorkLogs(issueKey)
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

func (jira JiraInstance) GetAssignedIssues() ([]Issue, error) {

	// retrieval mail fail
	remoteIssues, err := jira.GetRemoteIssues()
	if err != nil {
		jira.OfflineMode = true
	}

	for _, remoteIssue := range remoteIssues {
		jira.store.UpdateStructuredValue(issueBucketName, remoteIssue.Key, remoteIssue)
	}

	// get locally stored issues
	storedIssues, err := jira.GetAllStoredIssues()
	if err != nil {
		return nil, err
	}

	return storedIssues, nil
}

func (jira JiraInstance) GetBaseUrl() string {
	return jira.atlassian.Site.String()
}

type WorkLogs []Worklog

func (wls WorkLogs) Len() int {
	return len(wls)
}
func (wls WorkLogs) Swap(i, j int)      { wls[i], wls[j] = wls[j], wls[i] }
func (wls WorkLogs) Less(i, j int) bool { return wls[i].Updated < wls[j].Updated }

func (jira JiraInstance) GetWorkLogs(issueId string) ([]Worklog, error) {
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
	sort.Sort(worklogs)
	return worklogs, nil
}

func (jira *JiraInstance) StartTimer(issueId string) {
	// Is there already a timer?
	if timerData, ok := jira.timers[issueId]; ok && timerData != nil {
		jira.timers[issueId].pauseChan <- false
		return
	}

	// Setup new timer
	timerData := &TimerData{
		ticker:    time.NewTicker(100 * time.Millisecond),
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

func (jira JiraInstance) GetCurrentTimerValue(issueId string) time.Duration {
	worklog, err := jira.GetTimeFromStore(issueId)
	if err != nil {
		jira.LogWarningf("unable to get worklog time for %s: %s", issueId, err.Error())
		return 0
	}
	return worklog
}

func (jira JiraInstance) SubmitWorklog(issueId string) error {
	timer := jira.GetCurrentTimerValue(issueId)
	options := &models.WorklogOptionsScheme{
		Notify:         false,
		AdjustEstimate: "leave",
	}

	payload := &models.WorklogADFPayloadScheme{
		Visibility: nil,
		// TimeSpent:  "5h",
		Started:          time.Now().Format(jiraTimeLayout),
		TimeSpentSeconds: int(timer.Round(time.Minute).Seconds()),
		// TimeSpent:        str2dur.String(timer),
	}
	jira.LogDebugf("%+v", payload)
	jira.LogDebugf("Received %s for %s", timer.String(), issueId)
	p, resp, err := jira.atlassian.Issue.Worklog.Add(jira.ctx, issueId, payload, options)
	if err != nil {
		jira.LogWarningf("error submitting: %s resp: %+v, payload: %+v", err.Error(), string(resp.Bytes.Bytes()), p)
		return err
	}
	jira.ResetTimer(issueId)
	runtime.EventsEmit(jira.ctx, "update_worklogs")
	return nil
}

func (jira JiraInstance) GetTimeFromStore(issueId string) (time.Duration, error) {
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

func (jira JiraInstance) WriteWorklogToStore(issueId string, t time.Duration) error {
	return jira.store.UpdateStringValue(timeBucketName, issueId, t.String())
}

func (jira JiraInstance) UpdateTrackedTime(issueId string, trackedTime time.Duration) {
	err := jira.WriteWorklogToStore(issueId, trackedTime)
	if err != nil {
		jira.LogWarningf("unable to update tracked time for %s: %s", issueId, err.Error())
		return
	}
	runtime.EventsEmit(jira.ctx, "timer_tick_"+issueId, trackedTime)
}

func (jira JiraInstance) WriteIssueDataToStore(issue Issue) error {
	return jira.store.UpdateStructuredValue(issueBucketName, issue.Key, issue)
}

func (jira JiraInstance) ReadIssueDataFromStore(issueId string) (Issue, error) {
	issue := Issue{}
	err := jira.store.GetStructuredValue(issueBucketName, issueId, &issue)
	return issue, err
}

func (jira JiraInstance) GetAllStoredIssues() ([]Issue, error) {
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
