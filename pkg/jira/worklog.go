package jira

import "time"

type Worklog struct {
	Duration  time.Duration `json:"duration,omitempty"`
	Comment   string        `json:"comment,omitempty"`
	Submitted bool          `json:"submitted,omitempty"`
	Author    string        `json:"author,omitempty"`
	Updated   time.Time     `json:"updated,omitempty"`
}

// sorting
type Worklogs []Worklog

func (wls Worklogs) Len() int {
	return len(wls)
}
func (wls Worklogs) Swap(i, j int) { wls[i], wls[j] = wls[j], wls[i] }

type ByUpdated struct{ Worklogs }

func (s ByUpdated) Less(i, j int) bool { return s.Worklogs[i].Updated.Before(s.Worklogs[j].Updated) }
