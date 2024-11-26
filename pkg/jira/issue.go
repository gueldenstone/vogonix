package jira

type Issue struct {
	Summary   string    `json:"summary,omitempty"`
	Assignee  string    `json:"assignee,omitempty"`
	Key       string    `json:"key,omitempty"`
	WorkLogs  []Worklog `json:"worklogs,omitempty"`
	TimeSpent int       `json:"time_spent,omitempty"`
}
