package timer

import (
	"fmt"
	"time"
)

type SectionTimer struct {
	Name     string
	Elapsed  time.Duration
	Sections []Section

	// Internal time tracking of sections
	prev  time.Time
	start time.Time
	end   time.Time
}

type Section struct {
	Name    string
	Start   time.Time
	End     time.Time
	Elapsed time.Duration
}

func (t SectionTimer) String() string {
	return fmt.Sprintf("[Profile Name: %v; Elapsed %v; Sections: %v;]", t.Name, t.Elapsed, t.Sections)
}

func (s Section) String() string {
	return fmt.Sprintf("[Section Name: %v; Start: %v; End: %v; Elapsed %v]", s.Name, s.Start.Unix(), s.End.Unix(), s.Elapsed)
}
