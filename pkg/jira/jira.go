package jira

import (
	"context"
	"sort"
	"time"

	v3 "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/gueldenstone/vogonix/pkg/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	storeDir    string
	OfflineMode bool
}

func NewJiraInstance(url, user, token, storeDir string) (*JiraInstance, error) {
	atlassian, err := v3.New(nil, url)
	if err != nil {
		return nil, err
	}
	atlassian.Auth.SetBasicAuth(user, token)

	return &JiraInstance{
		atlassian:   atlassian,
		timers:      make(map[string]*TimerData),
		storeDir:    storeDir,
		OfflineMode: false,
	}, nil
}

func (jira *JiraInstance) Startup(ctx context.Context) {
	jira.ctx = ctx
	if jira.store == nil {

		store, err := storage.NewStorage(jira.storeDir)
		if err != nil {
			panic(err)
		}
		store.AddBucket(timeBucketName)
		store.AddBucket(issueBucketName)
		jira.store = store
	}
}

func (jira *JiraInstance) Shutdown(ctx context.Context) {
	jira.store.Close()
}

func (jira JiraInstance) GetAssignedIssues() ([]Issue, error) {
	// retrieval mail fail
	remoteIssues, err := jira.getRemoteIssues()
	if err != nil {
		jira.OfflineMode = true
	}

	for _, remoteIssue := range remoteIssues {
		jira.writeIssueDataToStore(remoteIssue)
	}

	// get locally stored issues
	storedIssues, err := jira.getAllStoredIssues()
	if err != nil {
		return nil, err
	}

	return storedIssues, nil
}

func (jira JiraInstance) GetWorkLogs(issueId string) ([]Worklog, error) {
	worklogs, err := jira.getWorkLogs(issueId)
	sort.Sort(sort.Reverse(ByUpdated{worklogs}))
	return worklogs, err
}

func (jira *JiraInstance) StartTimer(issueId string) {
	// Is there already a timer?
	if timerData, ok := jira.timers[issueId]; ok && timerData != nil {
		jira.timers[issueId].pauseChan <- false
		return
	}

	// Setup new timer
	timerData := &TimerData{
		ticker:    time.NewTicker(1 * time.Second),
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
					worklog, _ := jira.getTimeFromStore(issueId)
					worklog += 1 * time.Second
					jira.updateTrackedTime(issueId, worklog)
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
	jira.updateTrackedTime(issueId, 0)
}

func (jira JiraInstance) GetCurrentTimerValue(issueId string) time.Duration {
	worklog, err := jira.getTimeFromStore(issueId)
	if err != nil {
		jira.LogWarningf("unable to get worklog time for %s: %s", issueId, err.Error())
		return 0
	}
	return worklog
}

func (jira JiraInstance) SubmitWorklog(issueId string) error {
	duration := jira.GetCurrentTimerValue(issueId)
	err := jira.AddWorklog(issueId, duration.Round(time.Minute))
	if err != nil {
		return err
	}
	jira.ResetTimer(issueId)
	runtime.EventsEmit(jira.ctx, "update_worklogs_"+issueId)
	return nil
}

func (jira JiraInstance) updateTrackedTime(issueId string, trackedTime time.Duration) {
	err := jira.writeWorklogToStore(issueId, trackedTime)
	if err != nil {
		jira.LogWarningf("unable to update tracked time for %s: %s", issueId, err.Error())
		return
	}
	runtime.EventsEmit(jira.ctx, "timer_tick_"+issueId, trackedTime)
}
