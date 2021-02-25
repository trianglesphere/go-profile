package timer

import (
	"fmt"
	"time"
)

func NewSectionTimer(name string) *SectionTimer {
	t := time.Now()
	return &SectionTimer{
		Name:  name,
		start: t,
		prev:  t,
	}
}

func (t *SectionTimer) Mark(label string) {
	now := time.Now()
	sect := Section{Name: label, Start: t.prev, End: now, Elapsed: now.Sub(t.prev)}
	t.prev = now
	t.Sections = append(t.Sections, sect)
}

func (t *SectionTimer) End(label string) {
	t.Mark(label)
	t.end = time.Now()
	t.Elapsed = t.start.Sub(t.end)
}

func (g SectionTimer) CSVHeader() string {
	ret := "duration,elapsed,name"
	for _, s := range g.Sections {
		ret = fmt.Sprintf("%s,%s_elapsed", ret, s.Name)
	}
	return ret
}

func (t SectionTimer) CSVString() string {
	// duration, elapsed, name, then the same for each section
	ret := fmt.Sprintf("%v,%v,%v", t.Elapsed, t.Elapsed.Nanoseconds(), t.Name)
	for _, s := range t.Sections {
		ret = fmt.Sprintf("%s,%s", ret, s.CSVString())
	}
	return ret
}

func (s Section) CSVString() string {
	return fmt.Sprintf("%v", s.Elapsed.Nanoseconds())
}
