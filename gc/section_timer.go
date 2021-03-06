package gc

import (
	"fmt"
	"runtime/debug"
	"time"
)

func NewSectionTimer(name string) *GCAwareSectionTimer {
	t := time.Now()
	return &GCAwareSectionTimer{
		Name:  name,
		start: t,
		prev:  t,
	}
}

func (g *GCAwareSectionTimer) Mark(label string) {
	t := time.Now()
	mark := mark{name: label, start: g.prev, end: t}
	g.prev = t
	g.marks = append(g.marks, mark)
}

func (g *GCAwareSectionTimer) End(label string) {
	g.Mark(label)
	g.end = time.Now()

	// Pull out info from gc stats
	gcs := debug.GCStats{}
	debug.ReadGCStats(&gcs)
	g.calculateGCTimes(gcs)
}

func (g *GCAwareSectionTimer) calculateGCTimes(gcs debug.GCStats) {
	ends := gcs.PauseEnd
	j := 0
	for i := len(g.marks) - 1; i >= 0; i-- {
		m := g.marks[i]
		section := Section{Name: m.name, Start: m.start, End: m.end}
		for ; j < len(ends) && ends[j].After(m.end); j++ {
			// pass
		}
		for ; j < len(ends) && ends[j].After(m.start) && ends[j].Before(m.end); j++ {
			section.GcCount++
			section.GcTime = section.GcTime + gcs.Pause[j]
		}
		section.Elapsed = m.end.Sub(m.start)
		section.RunTime = section.Elapsed - section.GcTime
		g.Sections = append(g.Sections, section)

		g.GcCount += section.GcCount
		g.GcTime += section.GcTime
	}
	g.Elapsed = g.end.Sub(g.start)
	g.RunTime = g.Elapsed - g.GcTime

	for i, j := 0, len(g.Sections)-1; i < j; i, j = i+1, j-1 {
		g.Sections[i], g.Sections[j] = g.Sections[j], g.Sections[i]
	}
}

func (g GCAwareSectionTimer) CSVHeader() string {
	ret := "elapsed,name,gc_time,run_time,gc_count"
	for _, s := range g.Sections {
		ret = fmt.Sprintf("%s,%s_elapsed,%s_gc_time,%s_run_time,%s_gc_count", ret, s.Name, s.Name, s.Name, s.Name)
	}
	return ret
}

func (g GCAwareSectionTimer) CSVString() string {
	// duration, elapsed, name, gc time, run time, gc count, then the same for each section
	ret := fmt.Sprintf("%v,%v,%v,%v,%v,%v", g.Elapsed, g.Elapsed.Nanoseconds(), g.Name, g.GcTime.Nanoseconds(), g.RunTime.Nanoseconds(), g.GcCount)
	for _, s := range g.Sections {
		ret = fmt.Sprintf("%s,%s", ret, s.CSVString())
	}
	return ret
}

func (s Section) CSVString() string {
	return fmt.Sprintf("%v,%v,%v,%v", s.Elapsed.Nanoseconds(), s.GcTime.Nanoseconds(), s.RunTime.Nanoseconds(), s.GcCount)
}
