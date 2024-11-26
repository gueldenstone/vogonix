package jira

import "time"

type Worklog struct {
	Duration  time.Duration `json:"duration,omitempty"`
	Comment   string        `json:"comment,omitempty"`
	Submitted bool          `json:"submitted,omitempty"`
	Author    string        `json:"author,omitempty"`
	Updated   string        `json:"updated,omitempty"`
}

// sorting
type WorkLogs []Worklog

func (wls WorkLogs) Len() int {
	return len(wls)
}
func (wls WorkLogs) Swap(i, j int) { wls[i], wls[j] = wls[j], wls[i] }

type ByUpdated struct{ WorkLogs }

func (s ByUpdated) Less(i, j int) bool { return s.WorkLogs[i].Updated < s.WorkLogs[j].Updated }
