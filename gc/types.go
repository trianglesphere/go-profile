package gc

import (
	"fmt"
	"time"
)

type GCAwareSectionTimer struct {
	Name  string
	start time.Time
	end   time.Time

	// GC Info not filled in until End called
	GcCount  int
	GcTime   time.Duration
	RunTime  time.Duration
	Elapsed  time.Duration
	Sections []Section

	// Internal time tracking of sections
	marks []mark
	prev  time.Time
}

type GCAwareTimer struct {
	Name  string
	start time.Time
	end   time.Time

	// GC Info not filled in until End called
	GcCount int
	GcTime  time.Duration
	RunTime time.Duration
	Elapsed time.Duration
}

type Section struct {
	Name    string
	Start   time.Time
	End     time.Time
	GcCount int
	GcTime  time.Duration
	RunTime time.Duration
	Elapsed time.Duration
}

func (g GCAwareTimer) String() string {
	return fmt.Sprintf("[Profile Name: %v; Elapsed %v; GCTime:  %v; RunTime: %v; GcCount: %v;]", g.Name, g.Elapsed, g.GcTime, g.RunTime, g.GcCount)
}

func (g GCAwareSectionTimer) String() string {
	return fmt.Sprintf("[Profile Name: %v; Elapsed %v; GCTime:  %v; RunTime: %v; GcCount: %v; Sections: %v;]", g.Name, g.Elapsed, g.GcTime, g.RunTime, g.GcCount, g.Sections)
}

func (s Section) String() string {
	return fmt.Sprintf("[Section Name: %v; Start: %v; End: %v; Elapsed %v; GCTime:  %v; RunTime: %v; GcCount: %v]", s.Name, s.Start.Unix(), s.End.Unix(), s.Elapsed, s.GcTime, s.RunTime, s.GcCount)
}

type mark struct {
	name  string
	start time.Time
	end   time.Time
}
