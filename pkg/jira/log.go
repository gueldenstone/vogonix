package jira

import "github.com/wailsapp/wails/v2/pkg/runtime"

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
